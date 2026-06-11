package apperr

import "fmt"

// Kind classifies a domain error so the transport edge can map it to a
// protocol status. The domain never imports HTTP; the edge never parses
// error strings.
type Kind int

const (
	KindNotFound Kind = iota + 1
	KindConflict
	KindInvalid
	KindForbidden
)

// Error is a domain error with a kind. Wrap it freely with
// fmt.Errorf("context: %w", err) — the kind survives the chain and is
// recovered at the edge via errors.As.
type Error struct {
	Kind Kind
	Msg  string
	Err  error // optional cause
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

func (e *Error) Unwrap() error { return e.Err }

func NotFound(format string, a ...any) *Error {
	return &Error{Kind: KindNotFound, Msg: fmt.Sprintf(format, a...)}
}

func Conflict(format string, a ...any) *Error {
	return &Error{Kind: KindConflict, Msg: fmt.Sprintf(format, a...)}
}

func Forbidden(format string, a ...any) *Error {
	return &Error{Kind: KindForbidden, Msg: fmt.Sprintf(format, a...)}
}

// Invalid carries the validation cause so its detail reaches the client.
func Invalid(err error, format string, a ...any) *Error {
	return &Error{Kind: KindInvalid, Msg: fmt.Sprintf(format, a...), Err: err}
}
