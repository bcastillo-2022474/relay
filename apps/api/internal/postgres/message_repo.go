package postgres

import (
	"context"
	"time"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
	"github.com/bcastillo-2022474/relay/internal/postgres/db"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageRepo struct {
	q *db.Queries
}

func NewMessageRepo(pool *pgxpool.Pool) *MessageRepo {
	return &MessageRepo{q: db.New(pool)}
}

func (r *MessageRepo) Save(ctx context.Context, msg message.Message) error {
	return r.q.InsertMessage(ctx, db.InsertMessageParams{
		ID:             uuid.UUID(msg.ID),
		OrganizationID: uuid.UUID(msg.OrganizationID),
		ApplicationID:  uuid.UUID(msg.ApplicationID),
		EventTypeID:    uuid.UUID(msg.EventTypeID),
		Payload:        msg.Payload,
		Status:         string(msg.Status),
	})
}

func (r *MessageRepo) ClaimBatch(ctx context.Context, limit int) ([]message.Message, error) {
	rows, err := r.q.ClaimMessageBatch(ctx, int32(limit))
	if err != nil {
		return nil, err
	}
	msgs := make([]message.Message, 0, len(rows))
	for _, row := range rows {
		msgs = append(msgs, message.Message{
			ID:             types.MessageID(row.ID),
			OrganizationID: types.OrganizationID(row.OrganizationID),
			ApplicationID:  types.ApplicationID(row.ApplicationID),
			EventTypeID:    types.EventTypeID(row.EventTypeID),
			Payload:        row.Payload,
			Status:         message.Status(row.Status),
		})
	}
	return msgs, nil
}

func (r *MessageRepo) MarkPublished(ctx context.Context, id types.MessageID) error {
	return r.q.MarkMessagePublished(ctx, uuid.UUID(id))
}

func (r *MessageRepo) Reschedule(ctx context.Context, id types.MessageID, nextAttempt time.Time) error {
	return r.q.RescheduleMessage(ctx, db.RescheduleMessageParams{
		ID:            uuid.UUID(id),
		NextAttemptAt: nextAttempt,
	})
}

func (r *MessageRepo) MarkFailed(ctx context.Context, id types.MessageID) error {
	return r.q.MarkMessageFailed(ctx, uuid.UUID(id))
}
