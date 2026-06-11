package command

import (
	"fmt"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/organization"
)

type CreateCommand struct {
	repo application.Repository
	log  *slog.Logger
}

func NewCreateCommand(repo application.Repository, log *slog.Logger) *CreateCommand {
	return &CreateCommand{repo: repo, log: log}
}

type CreateInput struct {
	OrganizationID organization.ID
	Name           string
	Slug           string
}

func (c *CreateCommand) Execute(input CreateInput) error {
	c.log.Info("creating application", "name", input.Name, "org_id", input.OrganizationID)

	app, err := application.New(input.OrganizationID, input.Name, input.Slug)
	if err != nil {
		return fmt.Errorf("creating application: %w", err)
	}

	if err := c.repo.Save(app); err != nil {
		return fmt.Errorf("saving application: %w", err)
	}

	c.log.Info("application created",
		"name", app.Name, "id", app.ID, "org_id", app.OrganizationID)
	return nil
}
