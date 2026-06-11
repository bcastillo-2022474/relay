package message

import (
	"encoding/json"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Message struct {
	EventType types.EventTypeID
	Payload   json.RawMessage
}

func NewMessage(eventType types.EventTypeID, payload json.RawMessage) Message {
	return Message{
		EventType: eventType,
		Payload:   payload,
	}
}
