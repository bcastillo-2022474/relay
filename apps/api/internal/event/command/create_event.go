package command

import (
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/event"
	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type CreateTypeCommand struct {
	typeRepo event.TypeRepository
	appRepo  application.Repository
	log      *slog.Logger
}

func NewCreateTypeCommand(
	typeRepo event.TypeRepository,
	appRepo application.Repository,
	log *slog.Logger,
) *CreateTypeCommand {
	return &CreateTypeCommand{
		typeRepo: typeRepo,
		appRepo:  appRepo,
		log:      log,
	}
}

type CreateTypeInput struct {
	Name           string
	ApplicationID  application.ID
	OrganizationID organization.ID
	PayloadSchema  *types.PayloadSchema
}

func (c *CreateTypeCommand) Execute(input CreateTypeInput) error {
	c.log.Info("creating event type",
		"name", input.Name,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	// Verify the application exists and belongs to this org.
	appResult, err := c.appRepo.FindByID(input.ApplicationID, input.OrganizationID)
	if err != nil {
		return fmt.Errorf("checking application: %w", err)
	}
	if !appResult.Ok {
		return fmt.Errorf("application %q not found in organization %q",
			input.ApplicationID, input.OrganizationID)
	}

	eventType, err := event.NewEventType(
		input.Name,
		input.ApplicationID,
		input.OrganizationID,
		input.PayloadSchema,
	)
	if err != nil {
		return fmt.Errorf("creating event type: %w", err)
	}

	if err := c.typeRepo.Save(eventType); err != nil {
		return fmt.Errorf("saving event type: %w", err)
	}

	c.log.Info("event type created",
		"name", eventType.Name,
		"id", eventType.ID,
		"app_id", eventType.ApplicationID)
	return nil
}
