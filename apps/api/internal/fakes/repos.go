package fakes

import (
	"sync"

	"github.com/bcastillo-2022474/relay/internal/application"
	"github.com/bcastillo-2022474/relay/internal/endpoint"
	"github.com/bcastillo-2022474/relay/internal/event_type"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

// In-memory repositories for local development and tests.
// They enforce the same org-scoping the real implementations will,
// so lookups across organizations come back not-found, not leaked.

type InMemoryApplicationRepo struct {
	mu   sync.RWMutex
	apps map[types.ApplicationID]application.Application
}

func NewInMemoryApplicationRepo() *InMemoryApplicationRepo {
	return &InMemoryApplicationRepo{apps: make(map[types.ApplicationID]application.Application)}
}

func (r *InMemoryApplicationRepo) Save(app application.Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apps[app.ID] = app
	return nil
}

func (r *InMemoryApplicationRepo) FindBySlug(slug string, orgID types.OrganizationID) (types.Result[application.Application], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, app := range r.apps {
		if app.Slug.String() == slug && app.OrganizationID == orgID {
			return types.Result[application.Application]{Found: app, Ok: true}, nil
		}
	}
	return types.Result[application.Application]{}, nil
}

func (r *InMemoryApplicationRepo) FindByID(id types.ApplicationID, orgID types.OrganizationID) (types.Result[application.Application], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	app, ok := r.apps[id]
	if !ok || app.OrganizationID != orgID {
		return types.Result[application.Application]{}, nil
	}
	return types.Result[application.Application]{Found: app, Ok: true}, nil
}

func (r *InMemoryApplicationRepo) Delete(app application.Application) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.apps, app.ID)
	return nil
}

type InMemoryEventTypeRepo struct {
	mu    sync.RWMutex
	types map[types.EventTypeID]event_type.EventType
}

func NewInMemoryEventTypeRepo() *InMemoryEventTypeRepo {
	return &InMemoryEventTypeRepo{types: make(map[types.EventTypeID]event_type.EventType)}
}

func (r *InMemoryEventTypeRepo) Save(et event_type.EventType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[et.ID] = et
	return nil
}

func (r *InMemoryEventTypeRepo) FindByID(id types.EventTypeID, orgID types.OrganizationID) (types.Result[event_type.EventType], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	et, ok := r.types[id]
	if !ok || et.OrganizationID != orgID {
		return types.Result[event_type.EventType]{}, nil
	}
	return types.Result[event_type.EventType]{Found: et, Ok: true}, nil
}

func (r *InMemoryEventTypeRepo) FindByName(name string, appID types.ApplicationID, orgID types.OrganizationID) (types.Result[event_type.EventType], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, et := range r.types {
		if et.Name == name && et.ApplicationID == appID && et.OrganizationID == orgID {
			return types.Result[event_type.EventType]{Found: et, Ok: true}, nil
		}
	}
	return types.Result[event_type.EventType]{}, nil
}

func (r *InMemoryEventTypeRepo) Delete(et event_type.EventType) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.types, et.ID)
	return nil
}

type InMemoryEndpointRepo struct {
	mu        sync.RWMutex
	endpoints map[types.EndpointID]endpoint.Endpoint
}

func NewInMemoryEndpointRepo() *InMemoryEndpointRepo {
	return &InMemoryEndpointRepo{endpoints: make(map[types.EndpointID]endpoint.Endpoint)}
}

func (r *InMemoryEndpointRepo) Save(ep endpoint.Endpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.endpoints[ep.ID] = ep
	return nil
}

func (r *InMemoryEndpointRepo) FindByID(id types.EndpointID, orgID types.OrganizationID) (types.Result[endpoint.Endpoint], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ep, ok := r.endpoints[id]
	if !ok || ep.OrganizationID != orgID {
		return types.Result[endpoint.Endpoint]{}, nil
	}
	return types.Result[endpoint.Endpoint]{Found: ep, Ok: true}, nil
}

func (r *InMemoryEndpointRepo) FindByURL(url string, orgID types.OrganizationID) (types.Result[endpoint.Endpoint], error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ep := range r.endpoints {
		if ep.URL == url && ep.OrganizationID == orgID {
			return types.Result[endpoint.Endpoint]{Found: ep, Ok: true}, nil
		}
	}
	return types.Result[endpoint.Endpoint]{}, nil
}

func (r *InMemoryEndpointRepo) Delete(ep endpoint.Endpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.endpoints, ep.ID)
	return nil
}
