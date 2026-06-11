package command

import (
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/event_type"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type CreateTypeCommand struct {
	typeRepo event_type.TypeRepository
	appRepo  application.Repository
	authz    types.Authorization
	log      *slog.Logger
}

func NewCreateTypeCommand(
	typeRepo event_type.TypeRepository,
	appRepo application.Repository,
	authz types.Authorization,
	log *slog.Logger,
) *CreateTypeCommand {
	return &CreateTypeCommand{
		typeRepo: typeRepo,
		appRepo:  appRepo,
		authz:    authz,
		log:      log,
	}
}

type CreateTypeInput struct {
	Name           string
	ApplicationID  types.ApplicationID
	OrganizationID types.OrganizationID
	PayloadSchema  *types.PayloadSchema
	Caller         types.Caller
}

func (c *CreateTypeCommand) Execute(input CreateTypeInput) (event_type.EventType, error) {
	c.log.Info("creating event type",
		"name", input.Name,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	if err := c.authz.CanCreateEventType(input.Caller, input.OrganizationID); err != nil {
		return event_type.EventType{}, fmt.Errorf("checking authorization: %w", err)
	}

	appResult, err := c.appRepo.FindByID(input.ApplicationID, input.OrganizationID)
	if err != nil {
		return event_type.EventType{}, fmt.Errorf("checking application: %w", err)
	}
	if !appResult.Ok {
		return event_type.EventType{}, apperr.NotFound("application %q not found in organization %q",
			input.ApplicationID, input.OrganizationID)
	}

	existing, err := c.typeRepo.FindByName(input.Name, input.ApplicationID, input.OrganizationID)
	if err != nil {
		return event_type.EventType{}, fmt.Errorf("checking event type name uniqueness: %w", err)
	}
	if existing.Ok {
		return event_type.EventType{}, apperr.Conflict("event type %q already exists for application %q",
			input.Name, input.ApplicationID)
	}

	eventType, err := event_type.NewEventType(
		input.Name,
		input.ApplicationID,
		input.OrganizationID,
		input.PayloadSchema,
	)
	if err != nil {
		return event_type.EventType{}, apperr.Invalid(err, "invalid event type")
	}

	if err := c.typeRepo.Save(eventType); err != nil {
		return event_type.EventType{}, fmt.Errorf("saving event type: %w", err)
	}

	c.log.Info("event type created",
		"name", eventType.Name,
		"id", eventType.ID,
		"app_id", eventType.ApplicationID)
	return eventType, nil
}
