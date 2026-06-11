package types

import (
	"fmt"
	"regexp"
)

// Must stay in sync with the is_valid_slug() function in
// migrations/001_init_write_side.sql.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const maxSlugLength = 63

type Slug struct {
	value string
}

func (s Slug) String() string {
	return s.value
}

// MarshalText lets encoding/json write a Slug as its plain string value.
// There is deliberately no UnmarshalText: slugs enter the system through
// NewSlug so validation cannot be bypassed by deserialization.
func (s Slug) MarshalText() ([]byte, error) {
	return []byte(s.value), nil
}

// SlugFromTrusted rehydrates a slug that was already validated at write time
// (the DB enforces is_valid_slug). Never use it for user input — tightening
// the validation rules must not make existing rows unreadable.
func SlugFromTrusted(value string) Slug { return Slug{value: value} }

func NewSlug(value string) (Slug, error) {
	err := validate(value)
	if err != nil {
		return Slug{}, err
	}

	return Slug{value: value}, nil
}

func validate(value string) error {
	if value == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(value) > maxSlugLength {
		return fmt.Errorf("slug cannot be longer than %d characters", maxSlugLength)
	}

	if !slugPattern.MatchString(value) {
		return fmt.Errorf("slug must be lowercase letters or digits separated by single hyphens (e.g. %q)", "billing-service")
	}

	return nil
}
