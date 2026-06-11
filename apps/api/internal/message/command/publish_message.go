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
	broker                message.Broker
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
	broker message.Broker,
	authz types.Authorization,
	log *slog.Logger,
) *PublishCommand {
	return &PublishCommand{
		eventTypeRepository:   eventTypeRepository,
		applicationRepository: applicationRepository,
		broker:                broker,
		authz:                 authz,
		log:                   log,
	}
}

func (c *PublishCommand) Execute(input PublishCommandInput) error {
	c.log.Info("publishing message",
		"event_type", input.EventType,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	err := c.authz.CanPublishMessage(input.Caller, input.OrganizationID)
	if err != nil {
		return fmt.Errorf("checking authorization: %w", err)
	}

	app, err := c.applicationRepository.FindByID(input.ApplicationID, input.OrganizationID)
	if err != nil {
		return fmt.Errorf("finding application: %w", err)
	}

	if !app.Ok {
		return apperr.NotFound("application %q not found", input.ApplicationID)
	}

	eventType, err := c.eventTypeRepository.FindByName(input.EventType, input.ApplicationID, input.OrganizationID)
	if err != nil {
		return fmt.Errorf("finding event_type: %w", err)
	}

	if !eventType.Ok {
		return apperr.NotFound("event_type %q not found", input.EventType)
	}

	// PayloadSchema is optional; event types without one accept any payload.
	if eventType.Found.HasPayloadSchema() {
		if err := eventType.Found.PayloadSchema.Validate(input.Payload); err != nil {
			return apperr.Invalid(err, "payload does not conform to event type %q schema", input.EventType)
		}
	}

	msg := message.NewMessage(eventType.Found.ID, input.Payload)

	err = c.broker.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}

	c.log.Info("message published successfully",
		"event_type", input.EventType,
		"event_type_id", eventType.Found.ID,
		"app_id", input.ApplicationID)

	return nil
}
