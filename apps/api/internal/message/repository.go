package message

import (
	"time"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Repository interface {
	Save(msg Message) error
	ClaimBatch(limit int) ([]Message, error)
	MarkPublished(id types.MessageID) error
	Reschedule(id types.MessageID, nextAttempt time.Time) error
	MarkFailed(id types.MessageID) error
}
