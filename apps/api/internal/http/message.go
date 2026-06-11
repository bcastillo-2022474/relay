package http

import (
	"context"
	"encoding/json"
	"net/http"

	msgcommand "github.com/bcastillo-2022474/relay/internal/domain/message/command"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
)

type messageResponse struct {
	Body struct {
		MessageID string `json:"message_id"`
		Status    string `json:"status"`
	}
}

func RegisterMessageRoutes(api huma.API, publishMsg *msgcommand.PublishCommand) {
	huma.Register(api, huma.Operation{
		OperationID:   "publish-message",
		Method:        http.MethodPost,
		Path:          "/v1/app/{appId}/msg",
		Summary:       "Publish a message",
		Tags:          []string{"Messages"},
		DefaultStatus: http.StatusAccepted,
	}, wrap(func(ctx context.Context, input *struct {
		AppID string `path:"appId"`
		Body  struct {
			EventType string          `json:"event_type" minLength:"1"`
			Payload   json.RawMessage `json:"payload"`
		}
	}) (*messageResponse, error) {
		appID, err := types.ParseApplicationID(input.AppID)
		if err != nil {
			return nil, apperr.Invalid(err, "invalid application id %q", input.AppID)
		}
		msg, err := publishMsg.Execute(ctx, msgcommand.PublishCommandInput{
			Payload:        input.Body.Payload,
			EventType:      input.Body.EventType,
			ApplicationID:  appID,
			OrganizationID: OrgIDFromCtx(ctx),
			Caller:         CallerFromCtx(ctx),
		})
		if err != nil {
			return nil, err
		}
		resp := &messageResponse{}
		resp.Body.MessageID = msg.ID.String()
		resp.Body.Status = string(msg.Status)
		return resp, nil
	}))
}
