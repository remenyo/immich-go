package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/simulot/immich-go/app/cmd"
)

// immich-go entry point
func main() {
	ctx := context.Background()
	err := immichGoMain(ctx)
	if err != nil {
		if e := context.Cause(ctx); e != nil {
			err = e
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// makes immich-go breakable with ^C and run it
func immichGoMain(ctx context.Context) error {
	// Create a context with cancel function to gracefully handle Ctrl+C events
	ctx, cancel := context.WithCancelCause(ctx)

	// Handle Ctrl+C signal (SIGINT)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	// Watch for ^C to be pressed
	go func() {
		<-signalChannel
		fmt.Println("\nCtrl+C received. Shutting down...")
		cancel(errors.New("Ctrl+C received")) // Cancel the context when Ctrl+C is received
	}()

	c, a := cmd.RootImmichGoCommand(ctx)
	err := c.ExecuteContext(ctx)
	if err == nil {
		return nil
	}

	retryMax, _ := c.Flags().GetInt("retry-max")
	retryDelay, _ := c.Flags().GetDuration("retry-delay")

	for i := 0; i < retryMax; i++ {
		a.Log().Info(fmt.Sprintf("command failed, wait for %s, and retry", retryDelay.String()))
		select {
		case <-time.After(retryDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
		c, a = cmd.RootImmichGoCommand(ctx)
		c.SetArgs(os.Args[1:])
		err = c.ExecuteContext(ctx)
		if err == nil {
			return nil
		}
	}

	if err != nil && a.Log().GetSLog() != nil {
		a.Log().Error(err.Error())
	}
	return err
}
