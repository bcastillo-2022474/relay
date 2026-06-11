package event

import (
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
)

type TypeID string

// EventType is a schema'd event definition within an application,
// e.g. "invoice.paid". Payloads published for this type are validated
// against PayloadSchema at ingest time.
type EventType struct {
	ID             TypeID
	Name           string
	ApplicationID  application.ID
	OrganizationID organization.ID
	PayloadSchema  *types.PayloadSchema
}

func NewEventType(
	name string,
	appID application.ID,
	orgID organization.ID,
	schema *types.PayloadSchema,
) (EventType, error) {
	if name == "" {
		return EventType{}, fmt.Errorf("event type name cannot be empty")
	}

	return EventType{
		ID:             TypeID(uuid.NewString()),
		Name:           name,
		ApplicationID:  appID,
		OrganizationID: orgID,
		PayloadSchema:  schema,
	}, nil
}
