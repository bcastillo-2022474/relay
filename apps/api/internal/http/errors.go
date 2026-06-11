package http

import (
	"context"
	"errors"

	"github.com/bcastillo-2022474/relay/internal/shared/apperr"
	"github.com/danielgtaylor/huma/v2"
)

func mapError(err error) error {
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

func wrap[I, O any](fn func(context.Context, *I) (*O, error)) func(context.Context, *I) (*O, error) {
	return func(ctx context.Context, in *I) (*O, error) {
		out, err := fn(ctx, in)
		return out, mapError(err)
	}
}
