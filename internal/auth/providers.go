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
	Name             string
	Config           *oauth2.Config
	ProviderKey      string   // "google", "github", "company-okta", ...
	UserInfoURL      string   // optional – for fetching /userinfo
	ExpectedIssuer   string   // iss claim must match exactly (or one of)
	ExpectedAudience []string // aud claim must contain at least one of these
	AllowAnyAudience bool     // rare – only for very specific cases
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

// Example inside the provider init loop or a separate setup function
func configureExpectedValues(p *ProviderConfig) {
	provider := strings.ToLower(p.ProviderKey)

	switch provider {
	case "google":
		p.ExpectedIssuer = "https://accounts.google.com"
		// Audience for id_token is usually your CLIENT_ID
		p.ExpectedAudience = []string{os.Getenv("OAUTH_GOOGLE_CLIENT_ID")}
		// Google sometimes also accepts the project number variant, but stick to client_id

	case "github":
		// GitHub does NOT use JWT access_tokens by default (opaque tokens)
		// If you're using GitHub App / OIDC → adjust accordingly
		p.ExpectedIssuer = "" // ← often not needed for GitHub
		p.ExpectedAudience = nil
		// GitHub userinfo uses access_token in header, no JWT validation usually

	case "auth0":
		domain := os.Getenv("AUTH0_DOMAIN") // e.g. "your-tenant.eu.auth0.com"
		p.ExpectedIssuer = "https://" + domain + "/"
		// Audience = your API identifier (not client_id!)
		// → set in Auth0 dashboard → APIs → Identifier
		p.ExpectedAudience = []string{os.Getenv("AUTH0_API_IDENTIFIER")}

	case "keycloak":
		// Issuer = full realm URL
		// Example: https://auth.example.com/realms/my-realm
		issuer := os.Getenv("OAUTH_KEYCLOAK_" + strings.ToUpper(provider) + "_ISSUER")
		if issuer == "" {
			issuer = "https://your-keycloak.com/realms/" + provider // fallback pattern
		}
		p.ExpectedIssuer = issuer
		// Audience = usually client_id (or audience mapper in Keycloak)
		p.ExpectedAudience = []string{os.Getenv("OAUTH_" + strings.ToUpper(provider) + "_CLIENT_ID")}

	default:
		// For custom / unknown providers → load from env or fail
		p.ExpectedIssuer = os.Getenv("OAUTH_" + strings.ToUpper(provider) + "_ISSUER")
		audEnv := os.Getenv("OAUTH_" + strings.ToUpper(provider) + "_AUDIENCE")
		if audEnv != "" {
			p.ExpectedAudience = []string{audEnv}
		}
	}

	// Optional fallback / warning
	if p.ExpectedIssuer == "" && p.UserInfoURL != "" {
		fmt.Printf("Warning: no ExpectedIssuer set for %s – OIDC validation will be weak\n", provider)
	}
}

/**/
