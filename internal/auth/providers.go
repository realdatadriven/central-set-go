// internal/auth/providers.go
package auth

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

type ProviderConfig struct {
	Name        string
	Config      *oauth2.Config
	ProviderKey string        // "google", "github", "company-okta", ...
	UserInfoURL string        // optional – for fetching /userinfo
}

var (
	providers = make(map[string]*ProviderConfig) // key = providerKey lowercase
)

func LoadProviders() error {
	// You can also read from a config file / DB in larger systems
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "OAUTH_") {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]

		// Example: OAUTH_GITHUB_CLIENT_ID → github
		if !strings.HasPrefix(key, "OAUTH_") || strings.Contains(key, "_") == false {
			continue
		}

		providerPart := strings.SplitN(strings.TrimPrefix(key, "OAUTH_"), "_", 2)[0]
		providerKey := strings.ToLower(providerPart)

		// Initialize if first time seeing this provider
		if _, exists := providers[providerKey]; !exists {
			providers[providerKey] = &ProviderConfig{
				Name:        providerKey,
				ProviderKey: providerKey,
				Config: &oauth2.Config{
					RedirectURL: getOrDefault("OAUTH_"+providerPart+"_REDIRECT_URL", os.Getenv("OAUTH_DEFAULT_REDIRECT_URL")),
					Scopes:      strings.Fields(getOrDefault("OAUTH_"+providerPart+"_SCOPES", os.Getenv("OAUTH_DEFAULT_SCOPES"))),
				},
			}
		}

		p := providers[providerKey]

		switch {
		case strings.HasSuffix(key, "_CLIENT_ID"):
			p.Config.ClientID = value
		case strings.HasSuffix(key, "_CLIENT_SECRET"):
			p.Config.ClientSecret = value
		case strings.HasSuffix(key, "_REDIRECT_URL"):
			p.Config.RedirectURL = value
		case strings.HasSuffix(key, "_AUTH_URL"):
			p.Config.Endpoint.AuthURL = value
		case strings.HasSuffix(key, "_TOKEN_URL"):
			p.Config.Endpoint.TokenURL = value
		case strings.HasSuffix(key, "_USERINFO_URL"):
			p.UserInfoURL = value
		}
	}

	// Apply well-known endpoints for popular providers (optional fallback)
	for key, p := range providers {
		if p.Config.Endpoint.AuthURL == "" {
			switch key {
			case "google":
				p.Config.Endpoint = google.Endpoint
			case "github":
				p.Config.Endpoint = github.Endpoint
				// GitHub userinfo = https://api.github.com/user
				p.UserInfoURL = "https://api.github.com/user"
			// add more well-known providers if you want
			}
		}
	}

	if len(providers) == 0 {
		return fmt.Errorf("no OAuth2 providers configured")
	}

	return nil
}

func GetProvider(providerKey string) (*ProviderConfig, bool) {
	p, ok := providers[strings.ToLower(providerKey)]
	return p, ok
}

func getOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}