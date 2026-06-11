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
