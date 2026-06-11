package http

import (
	"reflect"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// NewAPI builds the huma API and teaches its OpenAPI generator that our
// uuid-backed ID types serialize as strings (reflection alone would see a
// [16]byte array). The types themselves stay framework-free — they only
// implement the stdlib encoding.Text(Un)Marshaler interfaces.
func NewAPI(router chi.Router, title, version string) huma.API {
	cfg := huma.DefaultConfig(title, version)
	registry := cfg.Components.Schemas

	str := reflect.TypeOf("")
	for _, t := range []any{
		types.OrganizationID{},
		types.ApplicationID{},
		types.EndpointID{},
		types.EventTypeID{},
		types.UserID{},
		types.MessageID{},
		types.Slug{},
	} {
		registry.RegisterTypeAlias(reflect.TypeOf(t), str)
	}

	return humachi.New(router, cfg)
}
