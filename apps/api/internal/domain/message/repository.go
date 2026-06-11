package message

import (
	"context"
	"time"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(ctx context.Context, msg Message) error
	ClaimBatch(ctx context.Context, limit int) ([]Message, error)
	MarkPublished(ctx context.Context, id types.MessageID) error
	Reschedule(ctx context.Context, id types.MessageID, nextAttempt time.Time) error
	MarkFailed(ctx context.Context, id types.MessageID) error
}
