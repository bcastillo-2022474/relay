package user

import "github.com/bcastillo-2022474/relay/internal/shared/types"

type User struct {
	ID             types.UserID
	Name           string
	OrganizationID types.OrganizationID
}
