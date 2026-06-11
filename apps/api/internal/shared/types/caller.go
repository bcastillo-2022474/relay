package types

type CallerKind string

const (
	User   CallerKind = "user"
	System CallerKind = "system"
)

type Caller struct {
	ID   string
	Kind CallerKind
}
