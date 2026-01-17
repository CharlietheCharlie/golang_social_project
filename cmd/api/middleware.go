package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthroizedBasicError(w, r, fmt.Errorf("unauthorization header not found"))
				return
			}
			// parse and get base64 encoded username and password
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthroizedBasicError(w, r, fmt.Errorf("unauthorization header is invalid"))
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthroizedBasicError(w, r, err)
				return
			}

			username := app.config.auth.basic.username
			password := app.config.auth.basic.password

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthroizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			// check the credentials
			next.ServeHTTP(w, r)
		})
	}
}
