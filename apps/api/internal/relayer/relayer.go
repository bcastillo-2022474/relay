package relayer

import (
	"context"
	"log/slog"
	"math"
	"math/rand/v2"
	"time"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
)

const (
	batchSize    = 100
	maxAttempts  = 10
	pollInterval = time.Second
)

type Relayer struct {
	repo   message.Repository
	broker message.Broker
	log    *slog.Logger
}

func New(repo message.Repository, broker message.Broker, log *slog.Logger) *Relayer {
	return &Relayer{repo: repo, broker: broker, log: log}
}

func (r *Relayer) Run(ctx context.Context) {
	r.log.Info("relayer started")
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Don't make the boot backlog wait for the first tick.
	r.drain(ctx)

	for {
		select {
		case <-ctx.Done():
			r.log.Info("relayer stopping")
			return
		case <-ticker.C:
			r.drain(ctx)
		}
	}
}

// drain claims and publishes until the backlog is exhausted (a short batch),
// so throughput is bounded by the broker, not by the poll interval. The poll
// interval only decides how quickly we notice NEW work.
func (r *Relayer) drain(ctx context.Context) {
	for ctx.Err() == nil {
		if r.tick(ctx) < batchSize {
			return
		}
	}
}

func (r *Relayer) tick(ctx context.Context) int {
	msgs, err := r.repo.ClaimBatch(ctx, batchSize)
	if err != nil {
		r.log.Error("relayer: claim batch failed", "err", err)
		return 0
	}
	for _, msg := range msgs {
		r.process(ctx, msg)
	}
	return len(msgs)
}

func (r *Relayer) process(ctx context.Context, msg message.Message) {
	err := r.broker.Publish(ctx, msg)
	if err == nil {
		// Known trade-off: if this mark fails persistently the row stays
		// pending, gets republished (consumers dedupe on message ID), and can
		// eventually be marked failed despite having been delivered. Accepted
		// until it shows up in practice.
		if markErr := r.repo.MarkPublished(ctx, msg.ID); markErr != nil {
			r.log.Error("relayer: mark published failed", "message_id", msg.ID, "err", markErr)
		}
		return
	}

	r.log.Warn("relayer: publish failed", "message_id", msg.ID, "err", err)

	if msg.Attempts >= maxAttempts {
		if markErr := r.repo.MarkFailed(ctx, msg.ID); markErr != nil {
			r.log.Error("relayer: mark failed error", "message_id", msg.ID, "err", markErr)
		}
		r.log.Error("relayer: message failed permanently", "message_id", msg.ID, "attempts", msg.Attempts)
		return
	}

	next := time.Now().Add(backoff(msg.Attempts))
	if rescheduleErr := r.repo.Reschedule(ctx, msg.ID, next); rescheduleErr != nil {
		r.log.Error("relayer: reschedule failed", "message_id", msg.ID, "err", rescheduleErr)
	}
}

// backoff returns exponential delay with ±20% jitter: min(2^attempts, 300)s.
func backoff(attempts int) time.Duration {
	base := math.Min(math.Pow(2, float64(attempts)), 300)
	jitter := base * 0.2 * (rand.Float64()*2 - 1)
	return time.Duration((base + jitter) * float64(time.Second))
}
