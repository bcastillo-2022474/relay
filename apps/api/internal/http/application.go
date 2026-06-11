package http

import (
	"context"
	"net/http"

	appcommand "github.com/bcastillo-2022474/relay/internal/domain/application/command"
	"github.com/danielgtaylor/huma/v2"
)

type applicationResponse struct {
	Body struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}

func RegisterApplicationRoutes(api huma.API, createApp *appcommand.CreateCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-application",
		Method:        http.MethodPost,
		Path:          "/v1/app",
		Summary:       "Create an application",
		Tags:          []string{"Applications"},
		DefaultStatus: http.StatusCreated,
	}, wrap(func(ctx context.Context, input *struct {
		Body struct {
			Name string `json:"name" minLength:"1" doc:"Application name"`
			Slug string `json:"slug" minLength:"1" doc:"URL-safe identifier"`
		}
	}) (*applicationResponse, error) {
		app, err := createApp.Execute(ctx, appcommand.CreateInput{
			OrganizationID: OrgIDFromCtx(ctx),
			Caller:         CallerFromCtx(ctx),
			Name:           input.Body.Name,
			Slug:           input.Body.Slug,
		})
		if err != nil {
			return nil, err
		}
		resp := &applicationResponse{}
		resp.Body.ID = app.ID.String()
		resp.Body.Name = app.Name
		resp.Body.Slug = app.Slug.String()
		return resp, nil
	}))
}
