package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/domain/event_type"
	"github.com/bcastillo-2022474/relay/internal/postgres/db"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventTypeRepo struct {
	q *db.Queries
}

func NewEventTypeRepo(conn db.DBTX) *EventTypeRepo {
	return &EventTypeRepo{q: db.New(conn)}
}

func (r *EventTypeRepo) Save(ctx context.Context, et event_type.EventType) error {
	var schema []byte
	if et.HasPayloadSchema() {
		schema = et.PayloadSchema.JSON()
	}
	return r.q.InsertEventType(ctx, db.InsertEventTypeParams{
		ID:             uuid.UUID(et.ID),
		Name:           et.Name,
		ApplicationID:  uuid.UUID(et.ApplicationID),
		OrganizationID: uuid.UUID(et.OrganizationID),
		PayloadSchema:  schema,
	})
}

func (r *EventTypeRepo) FindByID(ctx context.Context, id types.EventTypeID, orgID types.OrganizationID) (types.Result[event_type.EventType], error) {
	row, err := r.q.FindEventTypeByID(ctx, db.FindEventTypeByIDParams{
		ID:             uuid.UUID(id),
		OrganizationID: uuid.UUID(orgID),
	})
	return eventTypeResult(row, err)
}

func (r *EventTypeRepo) FindByName(ctx context.Context, name string, appID types.ApplicationID, orgID types.OrganizationID) (types.Result[event_type.EventType], error) {
	row, err := r.q.FindEventTypeByName(ctx, db.FindEventTypeByNameParams{
		Name:           name,
		ApplicationID:  uuid.UUID(appID),
		OrganizationID: uuid.UUID(orgID),
	})
	return eventTypeResult(row, err)
}

func (r *EventTypeRepo) Delete(ctx context.Context, et event_type.EventType) error {
	return r.q.DeleteEventType(ctx, db.DeleteEventTypeParams{
		ID:             uuid.UUID(et.ID),
		OrganizationID: uuid.UUID(et.OrganizationID),
	})
}

func eventTypeResult(row db.EventType, err error) (types.Result[event_type.EventType], error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return types.Result[event_type.EventType]{}, nil
	}
	if err != nil {
		return types.Result[event_type.EventType]{}, err
	}
	et, err := rowToEventType(row)
	if err != nil {
		return types.Result[event_type.EventType]{}, err
	}
	return types.Result[event_type.EventType]{Found: et, Ok: true}, nil
}

func rowToEventType(row db.EventType) (event_type.EventType, error) {
	et := event_type.EventType{
		ID:             types.EventTypeID(row.ID),
		Name:           row.Name,
		ApplicationID:  types.ApplicationID(row.ApplicationID),
		OrganizationID: types.OrganizationID(row.OrganizationID),
	}
	// The schema was validated at write time, but the compiled form cannot be
	// stored — rebuilding it is construction, not re-validation. Failure here
	// means corrupt data and deserves a loud error.
	if len(row.PayloadSchema) > 0 {
		ps, err := types.NewPayloadSchema(row.PayloadSchema)
		if err != nil {
			return event_type.EventType{}, fmt.Errorf("compiling stored payload schema for event type %s: %w", et.ID, err)
		}
		et.PayloadSchema = &ps
	}
	return et, nil
}
