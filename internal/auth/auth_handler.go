package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const (
	cookiePrefix = "oauth_state_"
	cookiePath   = "/"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider is required", http.StatusBadRequest)
		return
	}

	cfg, ok := GetProvider(provider)
	if !ok {
		http.Error(w, "unknown OAuth provider", http.StatusBadRequest)
		return
	}

	state := generateSecureState()

	// ─── PKCE ───────────────────────────────────────────────
	verifier, err := GeneratePKCEVerifier() // or oauth2.GenerateVerifier()
	if err != nil {
		http.Error(w, "failed to generate PKCE verifier", http.StatusInternalServerError)
		return
	}

	challenge := S256CodeChallenge(verifier) // or oauth2.S256ChallengeFromVerifier(verifier)

	// Store BOTH state and verifier in cookie (or session/redis)
	// Important: bind them together to prevent mix-ups
	cookieValue := state + "|" + verifier // simple delimiter
	if err := setPKCECookie(w, cookieValue, provider); err != nil {
		http.Error(w, "failed to set PKCE cookie", http.StatusInternalServerError)
		return
	}

	// Also keep state cookie for CSRF (optional double protection)
	setStateCookie(w, state, provider)

	// Build authorization URL with PKCE parameters
	authURL := cfg.Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline, // or Online
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// CallbackHandler exchanges code for token and processes login
func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	if provider == "" {
		http.Error(w, "provider missing", http.StatusBadRequest)
		return
	}

	cfg, ok := GetProvider(provider)
	if !ok {
		http.Error(w, "unknown provider", http.StatusBadRequest)
		return
	}

	// 1. Get state & verifier from cookie
	cookieValue, err := getPKCECookie(r, provider)
	if err != nil {
		http.Error(w, "PKCE data missing", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(cookieValue, "|", 2)
	if len(parts) != 2 {
		http.Error(w, "invalid PKCE cookie format", http.StatusBadRequest)
		return
	}

	storedState, storedVerifier := parts[0], parts[1]

	queryState := r.URL.Query().Get("state")
	if queryState == "" || storedState != queryState {
		http.Error(w, "state mismatch (possible CSRF)", http.StatusBadRequest)
		return
	}

	// Clear cookies after use
	clearPKCECookie(w, provider)
	clearStateCookie(w, provider)

	// 2. Exchange code using the original verifier
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "authorization code missing", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	token, err := cfg.Config.Exchange(
		ctx,
		code,
		oauth2.SetAuthURLParam("code_verifier", storedVerifier), // ← PKCE magic
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("token exchange failed: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(token)
	// Proceed with user info, session/JWT creation, etc.
	// ...
}

// Helpers (move to utils.go if you prefer)
func generateSecureState() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		panic(err) // in real code → return "", err
	}
	return base64.URLEncoding.EncodeToString(b)
}

func setStateCookie(w http.ResponseWriter, state, provider string) error {
	cookieName := cookiePrefix + strings.ToLower(provider)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    state,
		Path:     cookiePath,
		HttpOnly: true,
		Secure:   true, // enforce in production
		SameSite: http.SameSiteStrictMode,
		MaxAge:   10 * 60, // 10 minutes is plenty for login
	})
	return nil
}

func getStateFromCookie(r *http.Request, provider string) (string, error) {
	cookieName := cookiePrefix + strings.ToLower(provider)
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func clearStateCookie(w http.ResponseWriter, provider string) {
	cookieName := cookiePrefix + strings.ToLower(provider)
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   cookiePath,
		MaxAge: -1,
	})
}

const pkceCookiePrefix = "oauth_pkce_"

func setPKCECookie(w http.ResponseWriter, value, provider string) error {
	name := pkceCookiePrefix + strings.ToLower(provider)
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // HTTPS only in production
		SameSite: http.SameSiteStrictMode,
		MaxAge:   10 * 60, // 10 minutes
	})
	return nil
}

func getPKCECookie(r *http.Request, provider string) (string, error) {
	name := pkceCookiePrefix + strings.ToLower(provider)
	c, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func clearPKCECookie(w http.ResponseWriter, provider string) {
	name := pkceCookiePrefix + strings.ToLower(provider)
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
