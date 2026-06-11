package http

import (
	"context"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/domain/application"
	appcommand "github.com/bcastillo-2022474/relay/internal/domain/application/command"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type createApplicationRequest struct {
	Body struct {
		Name string `json:"name" minLength:"1" doc:"Application name"`
		Slug string `json:"slug" minLength:"1" doc:"URL-safe identifier"`
	}
}

type applicationResponse struct {
	Body applicationBody
}

type applicationBody struct {
	ID   types.ApplicationID `json:"id"`
	Name string              `json:"name"`
	Slug types.Slug          `json:"slug"`
}

func newApplicationResponse(app application.Application) *applicationResponse {
	return &applicationResponse{Body: applicationBody{
		ID:   app.ID,
		Name: app.Name,
		Slug: app.Slug,
	}}
}

func RegisterApplicationRoutes(api huma.API, createApp *appcommand.CreateCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-application",
		Method:        http.MethodPost,
		Path:          "/v1/app",
		Summary:       "Create an application",
		Tags:          []string{"Applications"},
		DefaultStatus: http.StatusCreated,
	}, wrap(func(ctx context.Context, req *createApplicationRequest) (*applicationResponse, error) {
		app, err := createApp.Execute(ctx, appcommand.CreateInput{
			OrganizationID: OrgIDFromCtx(ctx),
			Caller:         CallerFromCtx(ctx),
			Name:           req.Body.Name,
			Slug:           req.Body.Slug,
		})

		if err != nil {
			return nil, err
		}
		return newApplicationResponse(app), nil
	}))
}
