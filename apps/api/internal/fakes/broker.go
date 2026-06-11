package fakes

import (
	"log/slog"

	"github.com/bcastillo-2022474/relay/internal/message"
)

type LogBroker struct{ log *slog.Logger }

func NewLogBroker(log *slog.Logger) *LogBroker { return &LogBroker{log: log} }

func (b *LogBroker) Publish(msg message.Message) error {
	b.log.Info("broker: publish", "event_type_id", msg.EventType, "payload", string(msg.Payload))
	return nil
}
