package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/domain/application"
	endpoint2 "github.com/bcastillo-2022474/relay/internal/domain/endpoint"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type CreateCommand struct {
	repo    endpoint2.Repository
	appRepo application.Repository
	authz   types.Authorization
	log     *slog.Logger
}

func NewCreateCommand(
	repo endpoint2.Repository,
	appRepo application.Repository,
	authz types.Authorization,
	log *slog.Logger,
) *CreateCommand {
	return &CreateCommand{
		repo:    repo,
		appRepo: appRepo,
		authz:   authz,
		log:     log,
	}
}

type CreateInput struct {
	ApplicationID  types.ApplicationID
	OrganizationID types.OrganizationID
	URL            string
	Description    string
	Caller         types.Caller
}

func (c *CreateCommand) Execute(ctx context.Context, input CreateInput) (endpoint2.Endpoint, error) {
	c.log.Info("creating endpoint",
		"url", input.URL,
		"app_id", input.ApplicationID,
		"org_id", input.OrganizationID)

	if err := c.authz.CanSubscribeToEvent(input.Caller, input.OrganizationID); err != nil {
		return endpoint2.Endpoint{}, fmt.Errorf("checking authorization: %w", err)
	}

	appResult, err := c.appRepo.FindByID(ctx, input.ApplicationID, input.OrganizationID)
	if err != nil {
		return endpoint2.Endpoint{}, fmt.Errorf("checking application: %w", err)
	}
	if !appResult.Ok {
		return endpoint2.Endpoint{}, apperr.NotFound("application %q not found in organization %q",
			input.ApplicationID, input.OrganizationID)
	}

	createdEndpoint, err := endpoint2.New(input.ApplicationID, input.OrganizationID, input.URL, input.Description)
	if err != nil {
		return endpoint2.Endpoint{}, apperr.Invalid(err, "invalid endpoint")
	}

	if err := c.repo.Save(ctx, createdEndpoint); err != nil {
		return endpoint2.Endpoint{}, fmt.Errorf("saving endpoint: %w", err)
	}

	c.log.Info("endpoint created",
		"id", createdEndpoint.ID,
		"url", createdEndpoint.URL,
		"app_id", createdEndpoint.ApplicationID)
	return createdEndpoint, nil
}
