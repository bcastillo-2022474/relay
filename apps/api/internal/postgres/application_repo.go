package postgres

import (
	"context"
	"errors"

	"github.com/bcastillo-2022474/relay/internal/domain/application"
	"github.com/bcastillo-2022474/relay/internal/postgres/db"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ApplicationRepo struct {
	q *db.Queries
}

func NewApplicationRepo(conn db.DBTX) *ApplicationRepo {
	return &ApplicationRepo{q: db.New(conn)}
}

func (r *ApplicationRepo) Save(ctx context.Context, app application.Application) error {
	return r.q.InsertApplication(ctx, db.InsertApplicationParams{
		ID:             uuid.UUID(app.ID),
		Name:           app.Name,
		Slug:           app.Slug.String(),
		OrganizationID: uuid.UUID(app.OrganizationID),
	})
}

func (r *ApplicationRepo) FindByID(ctx context.Context, id types.ApplicationID, orgID types.OrganizationID) (types.Result[application.Application], error) {
	row, err := r.q.FindApplicationByID(ctx, db.FindApplicationByIDParams{
		ID:             uuid.UUID(id),
		OrganizationID: uuid.UUID(orgID),
	})
	return appResult(row, err)
}

func (r *ApplicationRepo) FindBySlug(ctx context.Context, slug string, orgID types.OrganizationID) (types.Result[application.Application], error) {
	row, err := r.q.FindApplicationBySlug(ctx, db.FindApplicationBySlugParams{
		Slug:           slug,
		OrganizationID: uuid.UUID(orgID),
	})
	return appResult(row, err)
}

func (r *ApplicationRepo) Delete(ctx context.Context, app application.Application) error {
	return r.q.DeleteApplication(ctx, db.DeleteApplicationParams{
		ID:             uuid.UUID(app.ID),
		OrganizationID: uuid.UUID(app.OrganizationID),
	})
}

func appResult(row db.Application, err error) (types.Result[application.Application], error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Result[application.Application]{}, nil
	}
	if err != nil {
		return types.Result[application.Application]{}, err
	}
	return types.Result[application.Application]{Found: rowToApplication(row), Ok: true}, nil
}

func rowToApplication(row db.Application) application.Application {
	return application.Application{
		ID:             types.ApplicationID(row.ID),
		Name:           row.Name,
		Slug:           types.SlugFromTrusted(row.Slug),
		OrganizationID: types.OrganizationID(row.OrganizationID),
	}
}
