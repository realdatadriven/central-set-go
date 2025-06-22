package main

import (
	"context"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user *map[string]interface{}) *http.Request {
	//print("contextSetAuthenticatedUser: ", user, "\n")
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *map[string]interface{} {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*map[string]interface{})
	//print("contextGetAuthenticatedUser: ", user, " ", ok, "\n")
	if !ok {
		return nil
	}
	return user
}
