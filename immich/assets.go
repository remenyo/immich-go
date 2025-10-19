package immich

import (
	"context"
)

type UpdateAssetsVisibilityRequest struct {
	IDs      []string `json:"ids"`
	Visibility string   `json:"visibility"`
}

func (ic *ImmichClient) UpdateAssetsVisibility(ctx context.Context, ids []string, visibility string) (bool, error) {
	req := UpdateAssetsVisibilityRequest{
		IDs:      ids,
		Visibility: visibility,
	}

	err := ic.newServerCall(ctx, "UpdateAssetsVisibility").
		do(
			putRequest("/assets", setAcceptJSON(), setJSONBody(req)),
		)

	return err == nil, err
}
