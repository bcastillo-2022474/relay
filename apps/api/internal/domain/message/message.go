package message

import (
	"encoding/json"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPublished Status = "published"
	StatusFailed    Status = "failed"
)

type Message struct {
	ID             types.MessageID
	OrganizationID types.OrganizationID
	ApplicationID  types.ApplicationID
	EventTypeID    types.EventTypeID
	Payload        json.RawMessage
	Status         Status
}

func New(
	orgID types.OrganizationID,
	appID types.ApplicationID,
	eventTypeID types.EventTypeID,
	payload json.RawMessage,
) Message {
	return Message{
		ID:             types.NewMessageID(),
		OrganizationID: orgID,
		ApplicationID:  appID,
		EventTypeID:    eventTypeID,
		Payload:        payload,
		Status:         StatusPending,
	}
}
