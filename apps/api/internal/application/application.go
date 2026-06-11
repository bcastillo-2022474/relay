package application

import (
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
)

type ID string

// Application is a logical event source within an organization,
// e.g. "billing-service".
type Application struct {
	ID             ID
	Name           string
	OrganizationID organization.ID
	Slug           types.Slug
}

func New(organizationID organization.ID, name, slug string) (Application, error) {
	slugVO, err := types.NewSlug(slug)
	if err != nil {
		return Application{}, fmt.Errorf("invalid application slug %q: %w", slug, err)
	}

	return Application{
		ID:             ID(uuid.NewString()),
		Name:           name,
		OrganizationID: organizationID,
		Slug:           slugVO,
	}, nil
}
