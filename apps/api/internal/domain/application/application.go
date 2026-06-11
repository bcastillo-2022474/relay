package application

import (
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Application struct {
	ID             types.ApplicationID
	Name           string
	OrganizationID types.OrganizationID
	Slug           types.Slug
}

func New(organizationID types.OrganizationID, name, slug string) (Application, error) {
	slugVO, err := types.NewSlug(slug)
	if err != nil {
		return Application{}, fmt.Errorf("invalid application slug %q: %w", slug, err)
	}

	return Application{
		ID:             types.NewApplicationID(),
		Name:           name,
		OrganizationID: organizationID,
		Slug:           slugVO,
	}, nil
}
