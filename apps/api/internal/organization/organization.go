package organization

import (
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Organization struct {
	ID   types.OrganizationID
	Name string
	Slug types.Slug
}
