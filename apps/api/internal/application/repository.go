package application

import (
	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(app Application) error
	FindBySlug(slug string, organizationID organization.ID) (types.Result[Application], error)
	FindByID(id ID, organizationID organization.ID) (types.Result[Application], error)
	Delete(app Application) error
}
