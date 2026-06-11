package types

type Result[T any] struct {
	Found T
	Ok    bool
}
