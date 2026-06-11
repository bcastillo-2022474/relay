package event_type

import "github.com/bcastillo-2022474/relay/internal/shared/types"

type TypeRepository interface {
	Save(eventType EventType) error
	FindByID(id types.EventTypeID, orgID types.OrganizationID) (types.Result[EventType], error)
	FindByName(name string, appID types.ApplicationID, orgID types.OrganizationID) (types.Result[EventType], error)
	Delete(eventType EventType) error
}
