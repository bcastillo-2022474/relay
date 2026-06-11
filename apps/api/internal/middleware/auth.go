package middleware

import (
	"context"
	"net/http"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

type contextKey string

const (
	callerKey contextKey = "caller"
	orgIDKey  contextKey = "org_id"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: verify JWT, extract claims
		caller := types.Caller{ID: "stub-user-id", Kind: types.User}
		orgID := types.OrganizationID("stub-org-id")

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
