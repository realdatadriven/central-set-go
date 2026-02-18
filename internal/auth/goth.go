package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
)

// generateState – simple secure random state (you can also add PKCE here)
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// InitGoth loads all configured providers from env vars
func InitGoth() error {
	var providers []goth.Provider
	// Example: load Google if env vars present
	if cid := os.Getenv("OAUTH_GOOGLE_CLIENT_ID"); cid != "" {
		providers = append(providers,
			google.New(
				cid,
				os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
				os.Getenv("OAUTH_GOOGLE_CALLBACK_URL"), // e.g. http://localhost:8080/auth/google/callback
				"email", "profile", "openid",           // scopes
			),
		)
	}

	// GitHub
	if cid := os.Getenv("OAUTH_GITHUB_CLIENT_ID"); cid != "" {
		providers = append(providers,
			github.New(
				cid,
				os.Getenv("OAUTH_GITHUB_CLIENT_SECRET"),
				os.Getenv("OAUTH_GITHUB_CALLBACK_URL"),
				"user:email", "read:user",
			),
		)
	}

	// Generic OIDC (Auth0, Keycloak, Okta, custom)
	// Use separate env blocks or prefix logic as before
	if oidcClientID := os.Getenv("OAUTH_OIDC_CLIENT_ID"); oidcClientID != "" {
		// Discovery URL example for Keycloak: https://.../realms/{realm}/.well-known/openid-configuration
		discoveryURL := os.Getenv("OAUTH_OIDC_DISCOVERY_URL")
		if discoveryURL == "" {
			return fmt.Errorf("missing OAUTH_OIDC_DISCOVERY_URL for OIDC provider")
		}

		p, err := openidConnect.New(
			oidcClientID,
			os.Getenv("OAUTH_OIDC_CLIENT_SECRET"),
			os.Getenv("OAUTH_OIDC_CALLBACK_URL"),
			discoveryURL,
		)
		if err != nil {
			return err
		}
		// Optional: override name if you want "keycloak-customerA" instead of "openid-connect"
		// p.SetName("keycloak")
		providers = append(providers, p)
	}

	if len(providers) == 0 {
		return fmt.Errorf("no OAuth providers configured")
	}

	goth.UseProviders(providers...)

	// Session store setup (required!)
	// Use secure random key in production (32+ bytes)
	fall_back_secret, _ := generateState()
	// println("fall_back_secret:", fall_back_secret)
	store := sessions.NewCookieStore([]byte(os.Getenv("OAUTH_SESSION_SECRET") + fall_back_secret))
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("ENV") == "production" // only HTTPS in prod
	//store.Options.Domain = http.SameSiteLaxMode
	store.Options.Path = "/"
	gothic.Store = store // gothic is goth's session helper*/

	return nil
}
