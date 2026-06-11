package http

import (
	"context"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/domain/endpoint"
	endpointcommand "github.com/bcastillo-2022474/relay/internal/domain/endpoint/command"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type createEndpointRequest struct {
	AppID types.ApplicationID `path:"appId" doc:"Application ID"`
	Body  struct {
		URL         string `json:"url" minLength:"1"`
		Description string `json:"description,omitempty"`
	}
}

type endpointResponse struct {
	Body endpointBody
}

type endpointBody struct {
	ID            types.EndpointID    `json:"id"`
	ApplicationID types.ApplicationID `json:"application_id"`
	URL           string              `json:"url"`
	Description   string              `json:"description,omitempty"`
	SigningSecret string              `json:"signing_secret" doc:"Shown once on creation; used to verify webhook signatures"`
}

func newEndpointResponse(ep endpoint.Endpoint) *endpointResponse {
	return &endpointResponse{Body: endpointBody{
		ID:            ep.ID,
		ApplicationID: ep.ApplicationID,
		URL:           ep.URL,
		Description:   ep.Description,
		SigningSecret: ep.SigningSecret,
	}}
}

func RegisterEndpointRoutes(api huma.API, createEndpoint *endpointcommand.CreateCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-endpoint",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/endpoint",
		Summary:       "Create an endpoint",
		Tags:          []string{"Endpoints"},
		DefaultStatus: http.StatusCreated,
	}, wrap(func(ctx context.Context, req *createEndpointRequest) (*endpointResponse, error) {
		ep, err := createEndpoint.Execute(ctx, endpointcommand.CreateInput{
			ApplicationID:  req.AppID,
			OrganizationID: OrgIDFromCtx(ctx),
			URL:            req.Body.URL,
			Description:    req.Body.Description,
			Caller:         CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		return newEndpointResponse(ep), nil
	}))
}
