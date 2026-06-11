package http

import (
	"context"
	"net/http"

	endpointcommand "github.com/bcastillo-2022474/relay/internal/domain/endpoint/command"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type endpointResponse struct {
	Body struct {
		ID            string `json:"id"`
		ApplicationID string `json:"application_id"`
		URL           string `json:"url"`
		Description   string `json:"description,omitempty"`
		SigningSecret string `json:"signing_secret" doc:"Shown once on creation"`
	}
}

func RegisterEndpointRoutes(api huma.API, createEndpoint *endpointcommand.CreateCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-endpoint",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/endpoint",
		Summary:       "Create an endpoint",
		Tags:          []string{"Endpoints"},
		DefaultStatus: http.StatusCreated,
	}, wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			URL         string `json:"url" minLength:"1"`
			Description string `json:"description,omitempty"`
		}
	}) (*endpointResponse, error) {
		appID, err := types.ParseApplicationID(input.AppID)
		if err != nil {
			return nil, apperr.Invalid(err, "invalid application id %q", input.AppID)
		}
		ep, err := createEndpoint.Execute(ctx, endpointcommand.CreateInput{
			ApplicationID:  appID,
			OrganizationID: OrgIDFromCtx(ctx),
			URL:            input.Body.URL,
			Description:    input.Body.Description,
			Caller:         CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		resp := &endpointResponse{}
		resp.Body.ID = ep.ID.String()
		resp.Body.ApplicationID = ep.ApplicationID.String()
		resp.Body.URL = ep.URL
		resp.Body.Description = ep.Description
		resp.Body.SigningSecret = ep.SigningSecret
		return resp, nil
	}))
}
