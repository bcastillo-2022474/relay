package endpoint

import (
	"github.com/bcastillo-2022474/relay/internal/organization"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(endpoint Endpoint) error
	FindByID(id ID, orgID organization.ID) (types.Result[Endpoint], error)
	FindByURL(url string, orgID organization.ID) (types.Result[Endpoint], error)
	Delete(endpoint Endpoint) error
}

type SubscriptionRepository interface {
	Subscribe(endpointID ID, eventTypeID string, orgID organization.ID) error
	Unsubscribe(endpointID ID, eventTypeID string, orgID organization.ID) error
	FindSubscribedEventTypes(endpointID ID, orgID organization.ID) ([]string, error)
}
