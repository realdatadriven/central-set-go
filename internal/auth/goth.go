package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
)

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
	/*store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET") + "fallback-secret-very-long"))
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("ENV") == "production" // only HTTPS in prod
	store.Options.SameSite = http.SameSiteLaxMode
	store.Options.Path = "/"
	//gothic.Store = store // gothic is goth's session helper*/

	return nil
}

func GothLoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider required", http.StatusBadRequest)
		return
	}

	// Goth handles state, redirect, PKCE (for supported providers), etc.
	// You can pass extra oauth2 options if needed
	q := r.URL.Query()
	if q.Get("prompt") != "" {
		// example: force consent screen
		//r = gothic.SetState(r, "state-with-prompt") // optional custom state
	}

	url, err := gothic.GetAuthURL(w, r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GothCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider missing from path", http.StatusBadRequest)
		return
	}

	// Completes the flow: exchanges code, fetches user info
	_, err := gothic.CompleteUserAuth(w, r.WithContext(context.WithValue(r.Context(), gothic.ProviderParamKey, provider)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// gu is goth.User – rich struct with:
	//   - Provider
	//   - UserID / Email / Name / NickName
	//   - AccessToken, RefreshToken, ExpiresAt
	//   - RawData (map[string]interface{} for extra claims)

	// Now: create/link your user in DB, issue your own JWT/session
	// Example:
	/*sessionToken, err := yourCreateJWTOrSession(gu)
	if err != nil {
		// ...
	}*/

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "", //sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
