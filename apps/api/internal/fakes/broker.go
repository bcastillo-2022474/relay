package fakes

import (
	"context"
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
)

type LogBroker struct{ log *slog.Logger }

func NewLogBroker(log *slog.Logger) *LogBroker { return &LogBroker{log: log} }

func (b *LogBroker) Publish(ctx context.Context, msg message.Message) error {
	b.log.Info("broker: publish", "event_type_id", msg.EventTypeID, "payload", string(msg.Payload))
	return nil
}
