package command

import (
	"context"
	"fmt"
	"log/slog"

	application2 "github.com/bcastillo-2022474/relay/internal/domain/application"
	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type CreateCommand struct {
	repo  application2.Repository
	authz types.Authorization
	log   *slog.Logger
}

func NewCreateCommand(repo application2.Repository, authz types.Authorization, log *slog.Logger) *CreateCommand {
	return &CreateCommand{repo: repo, authz: authz, log: log}
}

type CreateInput struct {
	OrganizationID types.OrganizationID
	Name           string
	Slug           string
	Caller         types.Caller
}

func (c *CreateCommand) Execute(ctx context.Context, input CreateInput) (application2.Application, error) {
	c.log.Info("creating application", "name", input.Name, "org_id", input.OrganizationID)

	if err := c.authz.CanCreateApplication(input.Caller, input.OrganizationID); err != nil {
		return application2.Application{}, fmt.Errorf("checking authorization: %w", err)
	}

	existing, err := c.repo.FindBySlug(ctx, input.Slug, input.OrganizationID)
	if err != nil {
		return application2.Application{}, fmt.Errorf("checking slug uniqueness: %w", err)
	}
	if existing.Ok {
		return application2.Application{}, apperr.Conflict("application slug %q already in use", input.Slug)
	}

	app, err := application2.New(input.OrganizationID, input.Name, input.Slug)
	if err != nil {
		return application2.Application{}, apperr.Invalid(err, "invalid application")
	}

	if err := c.repo.Save(ctx, app); err != nil {
		return application2.Application{}, fmt.Errorf("saving application: %w", err)
	}

	c.log.Info("application created",
		"name", app.Name, "id", app.ID, "org_id", app.OrganizationID)
	return app, nil
}
