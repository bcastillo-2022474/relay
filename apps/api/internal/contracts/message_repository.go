// Package contracts holds behavior suites that every implementation of a
// port must pass — the real adapter and its fake run the SAME tests, so a
// fake can never silently drift from the thing it imitates.
//
// This package is only imported from _test files; it must never appear in a
// production import path.
package contracts

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bcastillo-2022474/relay/internal/domain/message"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

// MessageRepositoryContract verifies the claim/mark/reschedule semantics the
// relayer depends on. NewRepo must return a FRESH, EMPTY repository plus a
// valid (org, app, event type) ID triple that Save may reference — the real
// implementation has foreign keys, so it must seed those rows.
type MessageRepositoryContract struct {
	NewRepo func(t *testing.T) (message.Repository, types.OrganizationID, types.ApplicationID, types.EventTypeID)
}

func (c MessageRepositoryContract) Test(t *testing.T) {
	ctx := context.Background()

	save := func(t *testing.T, repo message.Repository, org types.OrganizationID, app types.ApplicationID, et types.EventTypeID) message.Message {
		t.Helper()
		msg := message.New(org, app, et, json.RawMessage(`{"k":"v"}`))
		if err := repo.Save(ctx, msg); err != nil {
			t.Fatalf("Save: %v", err)
		}
		return msg
	}

	claimAll := func(t *testing.T, repo message.Repository) map[types.MessageID]message.Message {
		t.Helper()
		msgs, err := repo.ClaimBatch(ctx, 1000)
		if err != nil {
			t.Fatalf("ClaimBatch: %v", err)
		}
		byID := make(map[types.MessageID]message.Message, len(msgs))
		for _, m := range msgs {
			byID[m.ID] = m
		}
		return byID
	}

	t.Run("claim returns saved messages and increments attempts each time", func(t *testing.T) {
		repo, org, app, et := c.NewRepo(t)
		msg := save(t, repo, org, app, et)

		got, ok := claimAll(t, repo)[msg.ID]
		if !ok {
			t.Fatal("saved pending message was not claimed")
		}
		if got.Attempts != 1 {
			t.Fatalf("first claim: attempts = %d, want 1 (post-increment count)", got.Attempts)
		}
		if got.Status != message.StatusPending {
			t.Fatalf("claimed status = %q, want pending", got.Status)
		}
		if string(got.Payload) == "" {
			t.Fatal("claimed message lost its payload")
		}

		// Not rescheduled and not marked: still due, so claimable again,
		// with the attempt count advancing.
		got, ok = claimAll(t, repo)[msg.ID]
		if !ok {
			t.Fatal("unmarked message was not re-claimed")
		}
		if got.Attempts != 2 {
			t.Fatalf("second claim: attempts = %d, want 2", got.Attempts)
		}
	})

	t.Run("published messages are never claimed again", func(t *testing.T) {
		repo, org, app, et := c.NewRepo(t)
		msg := save(t, repo, org, app, et)

		if err := repo.MarkPublished(ctx, msg.ID); err != nil {
			t.Fatalf("MarkPublished: %v", err)
		}
		if _, ok := claimAll(t, repo)[msg.ID]; ok {
			t.Fatal("published message was claimed")
		}
	})

	t.Run("failed messages are never claimed again", func(t *testing.T) {
		repo, org, app, et := c.NewRepo(t)
		msg := save(t, repo, org, app, et)

		if err := repo.MarkFailed(ctx, msg.ID); err != nil {
			t.Fatalf("MarkFailed: %v", err)
		}
		if _, ok := claimAll(t, repo)[msg.ID]; ok {
			t.Fatal("failed message was claimed")
		}
	})

	t.Run("rescheduled messages are claimable only once due", func(t *testing.T) {
		repo, org, app, et := c.NewRepo(t)
		msg := save(t, repo, org, app, et)

		if err := repo.Reschedule(ctx, msg.ID, time.Now().Add(time.Hour)); err != nil {
			t.Fatalf("Reschedule(future): %v", err)
		}
		if _, ok := claimAll(t, repo)[msg.ID]; ok {
			t.Fatal("message rescheduled into the future was claimed early")
		}

		if err := repo.Reschedule(ctx, msg.ID, time.Now().Add(-time.Second)); err != nil {
			t.Fatalf("Reschedule(past): %v", err)
		}
		if _, ok := claimAll(t, repo)[msg.ID]; !ok {
			t.Fatal("due rescheduled message was not claimed")
		}
	})

	t.Run("claim respects the batch limit", func(t *testing.T) {
		repo, org, app, et := c.NewRepo(t)
		for range 3 {
			save(t, repo, org, app, et)
		}

		msgs, err := repo.ClaimBatch(ctx, 2)
		if err != nil {
			t.Fatalf("ClaimBatch: %v", err)
		}
		if len(msgs) != 2 {
			t.Fatalf("claimed %d messages with limit 2, want 2", len(msgs))
		}
	})
}
