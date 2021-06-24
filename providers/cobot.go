package providers

import (
	"context"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/requests"
	"net/url"
)

// BitbucketProvider represents an Bitbucket based Identity Provider
type CobotProvider struct {
	*ProviderData
	Space       string
}

var _ Provider = (*CobotProvider)(nil)

const (
	cobotProviderName = "Cobot"
	cobotDefaultScope = "email"
)

var (
	// Default Login URL for Cobot
	// Pre-parsed URL of https://www.cobot.me/oauth/authorize
	cobotDefaultLoginURL = &url.URL{
		Scheme: "https",
		Host:   "www.cobot.me",
		Path:   "/oauth2/authorize",
	}

	// Default Redeem URL for Cobot
	// Pre-parsed URL of https://www.cobot.me/oauth/access_token
	cobotDefaultRedeemURL = &url.URL{
		Scheme: "https",
		Host:   "www.cobot.me",
		Path:   "/oauth2/access_token",
	}

	// Default Validation URL for Cobot
	// Pre-parsed URL of https://www.cobot.me/api/user
	cobotDefaultValidateURL = &url.URL{
		Scheme: "https",
		Host:   "www.cobot.me",
		Path:   "/api/user",
	}
)

// NewCobotProvider initiates a new CobotProvider
func NewCobotProvider(p *ProviderData) *CobotProvider {
	p.setProviderDefaults(providerDefaults{
		name:        cobotProviderName,
		loginURL:    cobotDefaultLoginURL,
		redeemURL:   cobotDefaultRedeemURL,
		profileURL:  nil,
		validateURL: cobotDefaultValidateURL,
		scope:       bitbucketDefaultScope,
	})
	return &CobotProvider{ProviderData: p}
}

// GetEmailAddress returns the email of the authenticated user
func (p *CobotProvider) GetEmailAddress(ctx context.Context, s *sessions.SessionState) (string, error) {
	var response struct {
		Email string `json:"email"`
	}

	requestURL := p.ValidateURL.String() + "?access_token=" + s.AccessToken
	err := requests.New(requestURL).
		WithContext(ctx).
		Do().
		UnmarshalInto(&response)
	if err != nil {
		logger.Errorf("failed making request: %v", err)
		return "", err
	}

	return response.Email, nil
}
