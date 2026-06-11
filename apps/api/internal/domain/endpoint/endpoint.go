package endpoint

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/bcastillo-2022474/relay/internal/shared/types"
)

const signingSecretLength = 32

type Endpoint struct {
	ID             types.EndpointID
	ApplicationID  types.ApplicationID
	OrganizationID types.OrganizationID
	URL            string
	Description    string
	SigningSecret  string
	Disabled       bool
}

func New(
	appID types.ApplicationID,
	orgID types.OrganizationID,
	url string,
	description string,
) (Endpoint, error) {
	if url == "" {
		return Endpoint{}, fmt.Errorf("endpoint URL cannot be empty")
	}

	secret, err := generateSigningSecret()
	if err != nil {
		return Endpoint{}, fmt.Errorf("generating signing secret: %w", err)
	}

	return Endpoint{
		ID:             types.NewEndpointID(),
		ApplicationID:  appID,
		OrganizationID: orgID,
		URL:            url,
		Description:    description,
		SigningSecret:  secret,
		Disabled:       false,
	}, nil
}

func generateSigningSecret() (string, error) {
	bytes := make([]byte, signingSecretLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (e *Endpoint) Disable() { e.Disabled = true }
func (e *Endpoint) Enable()  { e.Disabled = false }
