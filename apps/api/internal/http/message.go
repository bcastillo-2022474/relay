package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
	msgcommand "github.com/bcastillo-2022474/relay/internal/domain/message/command"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type messageHandlers struct {
	publishMsg *msgcommand.PublishCommand
}

func RegisterMessageRoutes(api huma.API, publishMsg *msgcommand.PublishCommand) {
	h := messageHandlers{publishMsg: publishMsg}

	huma.Register(api, huma.Operation{
		OperationID:   "publish-message",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/msg",
		Summary:       "Publish a message",
		Tags:          []string{"Messages"},
		DefaultStatus: http.StatusAccepted, // durable acceptance, not delivery
	}, wrap(h.publish))
}

type publishMessageRequest struct {
	AppID types.ApplicationID `path:"appId" doc:"Application ID"`
	Body  struct {
		EventType string          `json:"event_type" minLength:"1"`
		Payload   json.RawMessage `json:"payload"`
	}
}

type messageResponse struct {
	Body messageBody
}

type messageBody struct {
	MessageID types.MessageID `json:"message_id"`
	Status    message.Status  `json:"status"`
}

func newMessageResponse(msg message.Message) *messageResponse {
	return &messageResponse{Body: messageBody{
		MessageID: msg.ID,
		Status:    msg.Status,
	}}
}

func (h messageHandlers) publish(ctx context.Context, req *publishMessageRequest) (*messageResponse, error) {
	msg, err := h.publishMsg.Execute(ctx, msgcommand.PublishCommandInput{
		Payload:        req.Body.Payload,
		EventType:      req.Body.EventType,
		ApplicationID:  req.AppID,
		OrganizationID: OrgIDFromCtx(ctx),
		Caller:         CallerFromCtx(ctx),
	})

	if err != nil {
		return nil, err
	}
	return newMessageResponse(msg), nil
}
