package application

import (
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(app Application) error
	FindBySlug(slug string, organizationID types.OrganizationID) (types.Result[Application], error)
	FindByID(id types.ApplicationID, organizationID types.OrganizationID) (types.Result[Application], error)
	Delete(app Application) error
}
