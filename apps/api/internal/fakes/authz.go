package fakes

import "github.com/bcastillo-2022474/relay/internal/shared/types"

type AllowAllAuthorization struct{}

func (AllowAllAuthorization) CanPublishMessage(types.Caller, types.OrganizationID) error { return nil }
func (AllowAllAuthorization) CanSubscribeToEvent(types.Caller, types.OrganizationID) error {
	return nil
}
func (AllowAllAuthorization) CanCreateApplication(types.Caller, types.OrganizationID) error {
	return nil
}
func (AllowAllAuthorization) CanCreateEventType(types.Caller, types.OrganizationID) error { return nil }
