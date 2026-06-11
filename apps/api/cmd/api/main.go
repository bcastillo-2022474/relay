package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	appcommand "github.com/bcastillo-2022474/relay/internal/application/command"
	endpointcommand "github.com/bcastillo-2022474/relay/internal/endpoint/command"
	etcommand "github.com/bcastillo-2022474/relay/internal/event_type/command"
	"github.com/bcastillo-2022474/relay/internal/fakes"
	"github.com/bcastillo-2022474/relay/internal/httpx"
	msgcommand "github.com/bcastillo-2022474/relay/internal/message/command"
	"github.com/bcastillo-2022474/relay/internal/middleware"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type ApplicationResponse struct {
	Body struct {
		ID   string `json:"id" doc:"Application ID"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}

type EventTypeResponse struct {
	Body struct {
		ID            string          `json:"id" doc:"Event type ID"`
		Name          string          `json:"name"`
		ApplicationID string          `json:"application_id"`
		PayloadSchema json.RawMessage `json:"payload_schema,omitempty"`
	}
}

type EndpointResponse struct {
	Body struct {
		ID            string `json:"id" doc:"Endpoint ID"`
		ApplicationID string `json:"application_id"`
		URL           string `json:"url"`
		Description   string `json:"description,omitempty"`
		SigningSecret string `json:"signing_secret" doc:"Shown once; used to verify webhook signatures"`
	}
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	authz := fakes.AllowAllAuthorization{}
	broker := fakes.NewLogBroker(log)

	// In-memory repositories until Postgres lands. These are the fakes the
	// repository contract suites will later run against.
	appRepo := fakes.NewInMemoryApplicationRepo()
	eventTypeRepo := fakes.NewInMemoryEventTypeRepo()
	endpointRepo := fakes.NewInMemoryEndpointRepo()

	createApp := appcommand.NewCreateCommand(appRepo, authz, log)
	createEventType := etcommand.NewCreateTypeCommand(eventTypeRepo, appRepo, authz, log)
	createEndpoint := endpointcommand.NewCreateCommand(endpointRepo, appRepo, authz, log)
	publishMsg := msgcommand.NewPublishCommand(eventTypeRepo, appRepo, broker, authz, log)

	router := chi.NewMux()
	router.Use(middleware.Auth)

	api := humachi.New(router, huma.DefaultConfig("Relay API", "0.1.0"))

	// POST /v1/app
	huma.Register(api, huma.Operation{
		OperationID:   "create-application",
		Method:        http.MethodPost,
		Path:          "/v1/app",
		Summary:       "Create an application",
		Tags:          []string{"Applications"},
		DefaultStatus: http.StatusCreated,
	}, httpx.Wrap(func(ctx context.Context, input *struct {
		Body struct {
			Name string `json:"name" minLength:"1" doc:"Application name"`
			Slug string `json:"slug" minLength:"1" doc:"URL-safe identifier"`
		}
	}) (*ApplicationResponse, error) {
		app, err := createApp.Execute(appcommand.CreateInput{
			OrganizationID: middleware.OrgIDFromCtx(ctx),
			Caller:         middleware.CallerFromCtx(ctx),
			Name:           input.Body.Name,
			Slug:           input.Body.Slug,
		})
		if err != nil {
			return nil, err
		}
		resp := &ApplicationResponse{}
		resp.Body.ID = string(app.ID)
		resp.Body.Name = app.Name
		resp.Body.Slug = app.Slug.String()
		return resp, nil
	}))

	// POST /v1/app/{appId}/event-type
	huma.Register(api, huma.Operation{
		OperationID:   "create-event-type",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/event-type",
		Summary:       "Create an event type",
		Tags:          []string{"Event Types"},
		DefaultStatus: http.StatusCreated,
	}, httpx.Wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			Name          string          `json:"name" minLength:"1"`
			PayloadSchema json.RawMessage `json:"payload_schema,omitempty"`
		}
	}) (*EventTypeResponse, error) {
		var schema *types.PayloadSchema
		if len(input.Body.PayloadSchema) > 0 {
			ps, err := types.NewPayloadSchema(input.Body.PayloadSchema)
			if err != nil {
				return nil, huma.Error422UnprocessableEntity("invalid payload schema", err)
			}
			schema = &ps
		}
		et, err := createEventType.Execute(etcommand.CreateTypeInput{
			Name:           input.Body.Name,
			ApplicationID:  types.ApplicationID(input.AppID),
			OrganizationID: middleware.OrgIDFromCtx(ctx),
			PayloadSchema:  schema,
			Caller:         middleware.CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		resp := &EventTypeResponse{}
		resp.Body.ID = string(et.ID)
		resp.Body.Name = et.Name
		resp.Body.ApplicationID = string(et.ApplicationID)
		if et.HasPayloadSchema() {
			resp.Body.PayloadSchema = et.PayloadSchema.JSON()
		}
		return resp, nil
	}))

	// POST /v1/app/{appId}/endpoint
	huma.Register(api, huma.Operation{
		OperationID:   "create-endpoint",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/endpoint",
		Summary:       "Create an endpoint",
		Tags:          []string{"Endpoints"},
		DefaultStatus: http.StatusCreated,
	}, httpx.Wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			URL         string `json:"url" minLength:"1"`
			Description string `json:"description,omitempty"`
		}
	}) (*EndpointResponse, error) {
		ep, err := createEndpoint.Execute(endpointcommand.CreateInput{
			ApplicationID:  types.ApplicationID(input.AppID),
			OrganizationID: middleware.OrgIDFromCtx(ctx),
			URL:            input.Body.URL,
			Description:    input.Body.Description,
			Caller:         middleware.CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		resp := &EndpointResponse{}
		resp.Body.ID = string(ep.ID)
		resp.Body.ApplicationID = string(ep.ApplicationID)
		resp.Body.URL = ep.URL
		resp.Body.Description = ep.Description
		resp.Body.SigningSecret = ep.SigningSecret
		return resp, nil
	}))

	// POST /v1/app/{appId}/msg
	huma.Register(api, huma.Operation{
		OperationID:   "publish-message",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/msg",
		Summary:       "Publish a message",
		Tags:          []string{"Messages"},
		DefaultStatus: http.StatusAccepted, // ingest is async; this is an ack, not a result
	}, httpx.Wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			EventType string          `json:"event_type" minLength:"1"`
			Payload   json.RawMessage `json:"payload"`
		}
	}) (*struct{}, error) {
		return nil, publishMsg.Execute(msgcommand.PublishCommandInput{
			Payload:        input.Body.Payload,
			EventType:      input.Body.EventType,
			ApplicationID:  types.ApplicationID(input.AppID),
			OrganizationID: middleware.OrgIDFromCtx(ctx),
			Caller:         middleware.CallerFromCtx(ctx),
		})
	}))

	log.Info("api listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("server error", "err", err)
		os.Exit(1)
	}
}
