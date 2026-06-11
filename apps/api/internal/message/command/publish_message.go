package command

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/event_type"
	"github.com/bcastillo-2022474/relay/internal/message"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type PublishCommand struct {
	eventTypeRepository   event_type.TypeRepository
	applicationRepository application.Repository
	messageRepository     message.Repository
	authz                 types.Authorization
	log                   *slog.Logger
}

type PublishCommandInput struct {
	Payload        json.RawMessage
	EventType      string
	ApplicationID  types.ApplicationID
	OrganizationID types.OrganizationID
	Caller         types.Caller
}

func NewPublishCommand(
	eventTypeRepository event_type.TypeRepository,
	applicationRepository application.Repository,
	messageRepository message.Repository,
	authz types.Authorization,
	log *slog.Logger,
) *PublishCommand {
	return &PublishCommand{
		eventTypeRepository:   eventTypeRepository,
		applicationRepository: applicationRepository,
		messageRepository:     messageRepository,
		authz:                 authz,
		log:                   log,
	}
}

func (c *PublishCommand) Execute(input PublishCommandInput) (message.Message, error) {
	c.log.Info("publishing message",
		"event_type", input.EventType,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	if err := c.authz.CanPublishMessage(input.Caller, input.OrganizationID); err != nil {
		return message.Message{}, fmt.Errorf("checking authorization: %w", err)
	}

	app, err := c.applicationRepository.FindByID(input.ApplicationID, input.OrganizationID)
	if err != nil {
		return message.Message{}, fmt.Errorf("finding application: %w", err)
	}
	if !app.Ok {
		return message.Message{}, apperr.NotFound("application %q not found", input.ApplicationID)
	}

	eventType, err := c.eventTypeRepository.FindByName(input.EventType, input.ApplicationID, input.OrganizationID)
	if err != nil {
		return message.Message{}, fmt.Errorf("finding event type: %w", err)
	}
	if !eventType.Ok {
		return message.Message{}, apperr.NotFound("event type %q not found", input.EventType)
	}

	if eventType.Found.HasPayloadSchema() {
		if err := eventType.Found.PayloadSchema.Validate(input.Payload); err != nil {
			return message.Message{}, apperr.Invalid(err, "payload does not conform to event type %q schema", input.EventType)
		}
	}

	msg := message.New(input.OrganizationID, input.ApplicationID, eventType.Found.ID, input.Payload)

	if err := c.messageRepository.Save(msg); err != nil {
		return message.Message{}, fmt.Errorf("saving message: %w", err)
	}

	c.log.Info("message accepted",
		"message_id", msg.ID,
		"event_type_id", eventType.Found.ID,
		"app_id", input.ApplicationID)

	return msg, nil
}
