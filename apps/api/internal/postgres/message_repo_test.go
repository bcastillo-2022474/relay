package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/bcastillo-2022474/relay/internal/contracts"
	"github.com/bcastillo-2022474/relay/internal/domain/message"
	"github.com/bcastillo-2022474/relay/internal/postgres"
	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TestMessageRepoContract runs the same suite the in-memory fake must pass,
// against real Postgres. Each NewRepo call gets its own transaction that is
// rolled back on cleanup, with the messages table emptied inside it, so the
// suite is isolated from dev data and leaves no residue.
func TestMessageRepoContract(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://relay:relay@localhost:5432/relay"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("creating pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(context.Background()); err != nil {
		t.Skipf("postgres not reachable at %s: %v", dbURL, err)
	}

	contracts.MessageRepositoryContract{
		NewRepo: func(t *testing.T) (message.Repository, types.OrganizationID, types.ApplicationID, types.EventTypeID) {
			ctx := context.Background()
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatalf("begin tx: %v", err)
			}
			t.Cleanup(func() { _ = tx.Rollback(context.Background()) })

			if _, err := tx.Exec(ctx, `DELETE FROM messages`); err != nil {
				t.Fatalf("emptying messages: %v", err)
			}

			orgID, appID, etID := uuid.New(), uuid.New(), uuid.New()
			if _, err := tx.Exec(ctx,
				`INSERT INTO organizations (id, name, slug) VALUES ($1, 'contract test', $2)`,
				orgID, "contract-"+orgID.String()[:8]); err != nil {
				t.Fatalf("seeding organization: %v", err)
			}
			if _, err := tx.Exec(ctx,
				`INSERT INTO applications (id, name, slug, organization_id) VALUES ($1, 'contract test', $2, $3)`,
				appID, "contract-"+appID.String()[:8], orgID); err != nil {
				t.Fatalf("seeding application: %v", err)
			}
			if _, err := tx.Exec(ctx,
				`INSERT INTO event_types (id, name, application_id, organization_id) VALUES ($1, 'contract.test', $2, $3)`,
				etID, appID, orgID); err != nil {
				t.Fatalf("seeding event type: %v", err)
			}

			return postgres.NewMessageRepo(tx),
				types.OrganizationID(orgID),
				types.ApplicationID(appID),
				types.EventTypeID(etID)
		},
	}.Test(t)
}
