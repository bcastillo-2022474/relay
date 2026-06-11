package application

import (
	"context"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(ctx context.Context, app Application) error
	FindBySlug(ctx context.Context, slug string, organizationID types.OrganizationID) (types.Result[Application], error)
	FindByID(ctx context.Context, id types.ApplicationID, organizationID types.OrganizationID) (types.Result[Application], error)
	Delete(ctx context.Context, app Application) error
}
