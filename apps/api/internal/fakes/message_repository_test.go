package fakes_test

import (
	"testing"

	"github.com/bcastillo-2022474/relay/internal/contracts"
	"github.com/bcastillo-2022474/relay/internal/domain/message"
	"github.com/bcastillo-2022474/relay/internal/fakes"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
)

func TestInMemoryMessageRepositoryContract(t *testing.T) {
	contracts.MessageRepositoryContract{
		NewRepo: func(t *testing.T) (message.Repository, types.OrganizationID, types.ApplicationID, types.EventTypeID) {
			// The fake enforces no foreign keys; any IDs are valid.
			return fakes.NewInMemoryMessageRepository(),
				types.OrganizationID(uuid.New()),
				types.ApplicationID(uuid.New()),
				types.EventTypeID(uuid.New())
		},
	}.Test(t)
}
