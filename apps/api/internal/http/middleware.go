package http

import (
	"context"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
	"github.com/google/uuid"
)

type contextKey string

const (
	callerKey contextKey = "caller"
	orgIDKey  contextKey = "org_id"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: verify JWT, extract claims.
		// Stub values must satisfy real invariants — the DB stores org IDs as
		// uuid — so the stub org is a fixed UUID, not a fake-looking string.
		caller := types.Caller{ID: "stub-user-id", Kind: types.User}
		orgID := types.OrganizationID(uuid.MustParse("00000000-0000-0000-0000-000000000001"))

		ctx := context.WithValue(r.Context(), callerKey, caller)
		ctx = context.WithValue(ctx, orgIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CallerFromCtx(ctx context.Context) types.Caller {
	return ctx.Value(callerKey).(types.Caller)
}

func OrgIDFromCtx(ctx context.Context) types.OrganizationID {
	return ctx.Value(orgIDKey).(types.OrganizationID)
}
