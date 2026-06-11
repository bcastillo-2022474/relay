package organization

import "github.com/bcastillo-2022474/relay/internal/shared/types"

type ID string

type Organization struct {
	ID   ID
	Name string
	Slug types.Slug
}
