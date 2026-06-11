package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/domain/event_type"
	etcommand "github.com/bcastillo-2022474/relay/internal/domain/event_type/command"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type eventTypeHandlers struct {
	createEventType *etcommand.CreateTypeCommand
}

func RegisterEventTypeRoutes(api huma.API, createEventType *etcommand.CreateTypeCommand) {
	h := eventTypeHandlers{createEventType: createEventType}

	huma.Register(api, huma.Operation{
		OperationID:   "create-event-type",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/event-type",
		Summary:       "Create an event type",
		Tags:          []string{"Event Types"},
		DefaultStatus: http.StatusCreated,
	}, wrap(h.create))
}

type createEventTypeRequest struct {
	AppID types.ApplicationID `path:"appId" doc:"Application ID"`
	Body  struct {
		Name          string          `json:"name" minLength:"1"`
		PayloadSchema json.RawMessage `json:"payload_schema,omitempty" doc:"Optional JSON Schema; payloads are validated against it at publish time"`
	}
}

type eventTypeResponse struct {
	Body eventTypeBody
}

type eventTypeBody struct {
	ID            types.EventTypeID   `json:"id"`
	Name          string              `json:"name"`
	ApplicationID types.ApplicationID `json:"application_id"`
	PayloadSchema json.RawMessage     `json:"payload_schema,omitempty"`
}

func newEventTypeResponse(et event_type.EventType) *eventTypeResponse {
	body := eventTypeBody{
		ID:            et.ID,
		Name:          et.Name,
		ApplicationID: et.ApplicationID,
	}
	if et.HasPayloadSchema() {
		body.PayloadSchema = et.PayloadSchema.JSON()
	}
	return &eventTypeResponse{Body: body}
}

func (h eventTypeHandlers) create(ctx context.Context, req *createEventTypeRequest) (*eventTypeResponse, error) {
	et, err := h.createEventType.Execute(ctx, etcommand.CreateTypeInput{
		Name:           req.Body.Name,
		ApplicationID:  req.AppID,
		OrganizationID: OrgIDFromCtx(ctx),
		PayloadSchema:  req.Body.PayloadSchema,
		Caller:         CallerFromCtx(ctx),
	})

	if err != nil {
		return nil, err
	}
	return newEventTypeResponse(et), nil
}
