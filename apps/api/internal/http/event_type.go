package http

import (
	"context"
	"encoding/json"
	"net/http"

	etcommand "github.com/bcastillo-2022474/relay/internal/domain/event_type/command"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type eventTypeResponse struct {
	Body struct {
		ID            string          `json:"id"`
		Name          string          `json:"name"`
		ApplicationID string          `json:"application_id"`
		PayloadSchema json.RawMessage `json:"payload_schema,omitempty"`
	}
}

func RegisterEventTypeRoutes(api huma.API, createEventType *etcommand.CreateTypeCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-event-type",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/event-type",
		Summary:       "Create an event type",
		Tags:          []string{"Event Types"},
		DefaultStatus: http.StatusCreated,
	}, wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			Name          string          `json:"name" minLength:"1"`
			PayloadSchema json.RawMessage `json:"payload_schema,omitempty"`
		}
	}) (*eventTypeResponse, error) {
		var schema *types.PayloadSchema
		if len(input.Body.PayloadSchema) > 0 {
			ps, err := types.NewPayloadSchema(input.Body.PayloadSchema)
			if err != nil {
				return nil, huma.Error422UnprocessableEntity("invalid payload schema", err)
			}
			schema = &ps
		}
		appID, err := types.ParseApplicationID(input.AppID)
		if err != nil {
			return nil, apperr.Invalid(err, "invalid application id %q", input.AppID)
		}
		et, err := createEventType.Execute(ctx, etcommand.CreateTypeInput{
			Name:           input.Body.Name,
			ApplicationID:  appID,
			OrganizationID: OrgIDFromCtx(ctx),
			PayloadSchema:  schema,
			Caller:         CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		resp := &eventTypeResponse{}
		resp.Body.ID = et.ID.String()
		resp.Body.Name = et.Name
		resp.Body.ApplicationID = et.ApplicationID.String()
		if et.HasPayloadSchema() {
			resp.Body.PayloadSchema = et.PayloadSchema.JSON()
		}
		return resp, nil
	}))
}
