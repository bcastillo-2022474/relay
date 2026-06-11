package types

import "github.com/google/uuid"

// IDs are uuid-backed: converting to and from the database is an infallible
// cast (uuid.UUID(id)), never a parse. Strings exist only at the transport
// edge — Parse* on the way in, String() on the way out.

type OrganizationID uuid.UUID
type ApplicationID uuid.UUID
type EndpointID uuid.UUID
type EventTypeID uuid.UUID
type UserID uuid.UUID
type MessageID uuid.UUID

func (id OrganizationID) String() string { return uuid.UUID(id).String() }
func (id ApplicationID) String() string  { return uuid.UUID(id).String() }
func (id EndpointID) String() string     { return uuid.UUID(id).String() }
func (id EventTypeID) String() string    { return uuid.UUID(id).String() }
func (id UserID) String() string         { return uuid.UUID(id).String() }
func (id MessageID) String() string      { return uuid.UUID(id).String() }

func NewApplicationID() ApplicationID { return ApplicationID(uuid.New()) }
func NewEndpointID() EndpointID       { return EndpointID(uuid.New()) }
func NewEventTypeID() EventTypeID     { return EventTypeID(uuid.New()) }
func NewMessageID() MessageID         { return MessageID(uuid.New()) }

func ParseOrganizationID(s string) (OrganizationID, error) {
	u, err := uuid.Parse(s)
	return OrganizationID(u), err
}

func ParseApplicationID(s string) (ApplicationID, error) {
	u, err := uuid.Parse(s)
	return ApplicationID(u), err
}

func ParseEndpointID(s string) (EndpointID, error) {
	u, err := uuid.Parse(s)
	return EndpointID(u), err
}

func ParseEventTypeID(s string) (EventTypeID, error) {
	u, err := uuid.Parse(s)
	return EventTypeID(u), err
}

func ParseMessageID(s string) (MessageID, error) {
	u, err := uuid.Parse(s)
	return MessageID(u), err
}
