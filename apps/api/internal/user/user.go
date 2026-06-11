package user

import "github.com/bcastillo-2022474/relay/internal/organization"

type ID string

type User struct {
	ID             ID
	Name           string
	OrganizationID organization.ID
}
