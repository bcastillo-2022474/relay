package endpoint

import "github.com/bcastillo-2022474/relay/internal/shared/types"

type Repository interface {
	Save(endpoint Endpoint) error
	FindByID(id types.EndpointID, orgID types.OrganizationID) (types.Result[Endpoint], error)
	FindByURL(url string, orgID types.OrganizationID) (types.Result[Endpoint], error)
	Delete(endpoint Endpoint) error
}

type SubscriptionRepository interface {
	Subscribe(endpointID types.EndpointID, eventTypeID types.EventTypeID, orgID types.OrganizationID) error
	Unsubscribe(endpointID types.EndpointID, eventTypeID types.EventTypeID, orgID types.OrganizationID) error
	FindSubscribedEventTypes(endpointID types.EndpointID, orgID types.OrganizationID) ([]types.EventTypeID, error)
}
