package event

import (
	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type TypeRepository interface {
	Save(eventType EventType) error
	FindByID(id TypeID, orgID organization.ID) (types.Result[EventType], error)
	FindByName(name string, appID application.ID, orgID organization.ID) (types.Result[EventType], error)
	Delete(eventType EventType) error
}
