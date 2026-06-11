package event_type

import (
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
)

type EventType struct {
	ID             types.EventTypeID
	Name           string
	ApplicationID  types.ApplicationID
	OrganizationID types.OrganizationID
	PayloadSchema  *types.PayloadSchema
}

func NewEventType(
	name string,
	appID types.ApplicationID,
	orgID types.OrganizationID,
	schema *types.PayloadSchema,
) (EventType, error) {
	if name == "" {
		return EventType{}, fmt.Errorf("event_type type name cannot be empty")
	}

	return EventType{
		ID:             types.EventTypeID(uuid.NewString()),
		Name:           name,
		ApplicationID:  appID,
		OrganizationID: orgID,
		PayloadSchema:  schema,
	}, nil
}

func (e EventType) HasPayloadSchema() bool {
	return e.PayloadSchema != nil
}
