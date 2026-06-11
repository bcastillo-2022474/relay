package fakes

import (
	"context"
	"sync"
	"time"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type InMemoryMessageRepository struct {
	mu   sync.Mutex
	rows map[types.MessageID]*messageRow
}

type messageRow struct {
	msg           message.Message
	nextAttemptAt time.Time
}

func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{rows: make(map[types.MessageID]*messageRow)}
}

func (r *InMemoryMessageRepository) Save(ctx context.Context, msg message.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rows[msg.ID] = &messageRow{msg: msg, nextAttemptAt: time.Now()}
	return nil
}

func (r *InMemoryMessageRepository) ClaimBatch(ctx context.Context, limit int) ([]message.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	var claimed []message.Message
	for _, row := range r.rows {
		if len(claimed) >= limit {
			break
		}
		if row.msg.Status == message.StatusPending && !row.nextAttemptAt.After(now) {
			// Mirror the real adapter: claiming increments attempts and
			// returns the post-increment count. The contract suite enforces
			// this parity.
			row.msg.Attempts++
			claimed = append(claimed, row.msg)
		}
	}
	return claimed, nil
}

func (r *InMemoryMessageRepository) MarkPublished(ctx context.Context, id types.MessageID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if row, ok := r.rows[id]; ok {
		row.msg.Status = message.StatusPublished
	}
	return nil
}

func (r *InMemoryMessageRepository) Reschedule(ctx context.Context, id types.MessageID, nextAttempt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if row, ok := r.rows[id]; ok {
		row.nextAttemptAt = nextAttempt
	}
	return nil
}

func (r *InMemoryMessageRepository) MarkFailed(ctx context.Context, id types.MessageID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if row, ok := r.rows[id]; ok {
		row.msg.Status = message.StatusFailed
	}
	return nil
}
