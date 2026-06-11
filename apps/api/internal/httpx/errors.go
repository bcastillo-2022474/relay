package httpx

import (
	"context"
	"errors"

	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/danielgtaylor/huma/v2"
)

// MapError translates domain error kinds into HTTP statuses. Errors without
// a kind pass through unchanged and surface as 500 — genuinely unexpected.
func MapError(err error) error {
	var ae *apperr.Error
	if err == nil || !errors.As(err, &ae) {
		return err
	}
	switch ae.Kind {
	case apperr.KindNotFound:
		return huma.Error404NotFound(ae.Error())
	case apperr.KindConflict:
		return huma.Error409Conflict(ae.Error())
	case apperr.KindInvalid:
		return huma.Error422UnprocessableEntity(ae.Error())
	case apperr.KindForbidden:
		return huma.Error403Forbidden(ae.Error())
	}
	return err
}

// Wrap applies MapError to a huma handler so every operation gets the
// mapping without remembering to call it.
func Wrap[I, O any](fn func(context.Context, *I) (*O, error)) func(context.Context, *I) (*O, error) {
	return func(ctx context.Context, in *I) (*O, error) {
		out, err := fn(ctx, in)
		return out, MapError(err)
	}
}
