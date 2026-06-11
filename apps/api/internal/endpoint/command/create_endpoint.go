package command

import (
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/endpoint"
	"github.com/bcastillo-2022474/relay/internal/organization"
)

type CreateCommand struct {
	repo    endpoint.Repository
	appRepo application.Repository
	log     *slog.Logger
}

func NewCreateCommand(
	repo endpoint.Repository,
	appRepo application.Repository,
	log *slog.Logger,
) *CreateCommand {
	return &CreateCommand{
		repo:    repo,
		appRepo: appRepo,
		log:     log,
	}
}

type CreateInput struct {
	ApplicationID  application.ID
	OrganizationID organization.ID
	URL            string
	Description    string
}

func (c *CreateCommand) Execute(input CreateInput) (endpoint.Endpoint, error) {
	c.log.Info("creating endpoint",
		"url", input.URL,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	// Verify the application exists and belongs to this org.
	appResult, err := c.appRepo.FindByID(input.ApplicationID, input.OrganizationID)
	if err != nil {
		return endpoint.Endpoint{}, fmt.Errorf("checking application: %w", err)
	}
	if !appResult.Ok {
		return endpoint.Endpoint{}, fmt.Errorf("application %q not found in organization %q",
			input.ApplicationID, input.OrganizationID)
	}

	createdEndpoint, err := endpoint.New(input.ApplicationID, input.OrganizationID, input.URL, input.Description)
	if err != nil {
		return endpoint.Endpoint{}, fmt.Errorf("creating endpoint: %w", err)
	}

	if err := c.repo.Save(createdEndpoint); err != nil {
		return endpoint.Endpoint{}, fmt.Errorf("saving endpoint: %w", err)
	}

	c.log.Info("endpoint created",
		"id", createdEndpoint.ID,
		"url", createdEndpoint.URL,
		"app_id", createdEndpoint.ApplicationID)
	return createdEndpoint, nil
}
