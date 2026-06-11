package event_type

import (
	"context"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type TypeRepository interface {
	Save(ctx context.Context, eventType EventType) error
	FindByID(ctx context.Context, id types.EventTypeID, orgID types.OrganizationID) (types.Result[EventType], error)
	FindByName(ctx context.Context, name string, appID types.ApplicationID, orgID types.OrganizationID) (types.Result[EventType], error)
	Delete(ctx context.Context, eventType EventType) error
}
