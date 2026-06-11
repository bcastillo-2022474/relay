package postgres

import (
	"context"
	"errors"

	"github.com/bcastillo-2022474/relay/internal/domain/endpoint"
	"github.com/bcastillo-2022474/relay/internal/postgres/db"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EndpointRepo struct {
	q *db.Queries
}

func NewEndpointRepo(conn db.DBTX) *EndpointRepo {
	return &EndpointRepo{q: db.New(conn)}
}

func (r *EndpointRepo) Save(ctx context.Context, ep endpoint.Endpoint) error {
	return r.q.InsertEndpoint(ctx, db.InsertEndpointParams{
		ID:             uuid.UUID(ep.ID),
		ApplicationID:  uuid.UUID(ep.ApplicationID),
		OrganizationID: uuid.UUID(ep.OrganizationID),
		Url:            ep.URL,
		Description:    ep.Description,
		SigningSecret:  ep.SigningSecret,
	})
}

func (r *EndpointRepo) FindByID(ctx context.Context, id types.EndpointID, orgID types.OrganizationID) (types.Result[endpoint.Endpoint], error) {
	row, err := r.q.FindEndpointByID(ctx, db.FindEndpointByIDParams{
		ID:             uuid.UUID(id),
		OrganizationID: uuid.UUID(orgID),
	})
	return endpointResult(row, err)
}

func (r *EndpointRepo) FindByURL(ctx context.Context, url string, orgID types.OrganizationID) (types.Result[endpoint.Endpoint], error) {
	row, err := r.q.FindEndpointByURL(ctx, db.FindEndpointByURLParams{
		Url:            url,
		OrganizationID: uuid.UUID(orgID),
	})
	return endpointResult(row, err)
}

func (r *EndpointRepo) Delete(ctx context.Context, ep endpoint.Endpoint) error {
	return r.q.DeleteEndpoint(ctx, db.DeleteEndpointParams{
		ID:             uuid.UUID(ep.ID),
		OrganizationID: uuid.UUID(ep.OrganizationID),
	})
}

func endpointResult(row db.Endpoint, err error) (types.Result[endpoint.Endpoint], error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Result[endpoint.Endpoint]{}, nil
	}
	if err != nil {
		return types.Result[endpoint.Endpoint]{}, err
	}
	return types.Result[endpoint.Endpoint]{Found: rowToEndpoint(row), Ok: true}, nil
}

func rowToEndpoint(row db.Endpoint) endpoint.Endpoint {
	return endpoint.Endpoint{
		ID:             types.EndpointID(row.ID),
		ApplicationID:  types.ApplicationID(row.ApplicationID),
		OrganizationID: types.OrganizationID(row.OrganizationID),
		URL:            row.Url,
		Description:    row.Description,
		SigningSecret:  row.SigningSecret,
		Disabled:       row.Disabled,
	}
}
