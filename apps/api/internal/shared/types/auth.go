package types

type Authorization interface {
	CanPublishMessage(caller Caller, organizationId OrganizationID) error
	CanSubscribeToEvent(caller Caller, organizationId OrganizationID) error
	CanCreateApplication(caller Caller, organizationId OrganizationID) error
	CanCreateEventType(caller Caller, organizationId OrganizationID) error
}
