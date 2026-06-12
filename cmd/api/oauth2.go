package main

import (
	"net/http"

	"github.com/realdatadriven/central-set-go/internal/request"
	"github.com/realdatadriven/central-set-go/internal/response"
)

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// oauthTokenHandler implements a minimal OAuth2 token endpoint
// that delegates user authentication to the existing _login logic
// and returns the existing JWT as an access token. This is a
// scaffold to integrate a full Fosite-based server later.
func (app *application) oauthTokenHandler(w http.ResponseWriter, r *http.Request) {
	params := Dict{}
	if err := request.DecodeJSON(w, r, &params); err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Reuse existing login flow. Expect payload structure similar to /dyn_api/login/login
	data := app._login(params)
	success, _ := data["success"].(bool)
	if !success {
		_ = response.JSON(w, http.StatusUnauthorized, data)
		return
	}

	tokenStr, _ := data["token"].(string)

	expiresIn := int64(app.config.jwt.tokenExpireHours * 3600)
	resp := oauthTokenResponse{
		AccessToken: tokenStr,
		TokenType:   "bearer",
		ExpiresIn:   expiresIn,
	}
	_ = response.JSON(w, http.StatusOK, resp)
}

// oauthMiddleware protects endpoints by validating the existing JWT
// token (issued by the existing login flow). Future work: replace
// this with Fosite introspection and a full OAuth2/OIDC flow.
func (app *application) oauthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := app.verifyToken(r)
		ok, _ := token["success"].(bool)
		if !ok {
			_ = response.JSON(w, http.StatusUnauthorized, token)
			return
		}
		next(w, r)
	}
}
