package endpoint

import (
	"context"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(ctx context.Context, endpoint Endpoint) error
	FindByID(ctx context.Context, id types.EndpointID, orgID types.OrganizationID) (types.Result[Endpoint], error)
	FindByURL(ctx context.Context, url string, orgID types.OrganizationID) (types.Result[Endpoint], error)
	Delete(ctx context.Context, endpoint Endpoint) error
}

type SubscriptionRepository interface {
	Subscribe(ctx context.Context, endpointID types.EndpointID, eventTypeID types.EventTypeID, orgID types.OrganizationID) error
	Unsubscribe(ctx context.Context, endpointID types.EndpointID, eventTypeID types.EventTypeID, orgID types.OrganizationID) error
	FindSubscribedEventTypes(ctx context.Context, endpointID types.EndpointID, orgID types.OrganizationID) ([]types.EventTypeID, error)
}
