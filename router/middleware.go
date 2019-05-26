package router

import (
	"net/http"

	"go.uber.org/zap"

	"../repository"
)

// middleware - a type that describes a middleware, at the core of this
// implementation a middleware is merely a function that takes a handler
// function, and returns a handler function.
type middleware func(http.HandlerFunc) http.HandlerFunc

// middlewareChain - a function that takes a handler function, a list of middlewares
// and creates a new application stack as a single http handler
func middlewareChain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	// if there are no more middlewares, we just return the
	// handlerfunc, as we are done recursing.
	if len(m) == 0 {
		return f
	}
	// otherwise pop the middleware from the list,
	// and call build chain recursively as it's parameter
	return m[0](middlewareChain(f, m[1:cap(m)]...))
}

/**
All requests to the token endpoint must be authenticated.
Either pass client id and secret via Basic Authentication or add client_id and client_secret fields to the POST body.
When providing the client_id and client_secret in the Authorization header it is expected to be:
    client_id:client_secret
    Base64 encoded

The authorization server MUST:
o require client authentication for confidential clients or for any client that was issued client credentials (or with other
authentication requirements),
o authenticate the client if client authentication is included,
o ensure that the authorization code was issued to the authenticated confidential client, or if the client is public, ensure that the
code was issued to "client_id" in the request,
o verify that the authorization code is valid, and
o ensure that the "redirect_uri" parameter is present if the "redirect_uri" parameter was included in the initial authorization
request as described in Section 4.1.1, and if included ensure that their values are identical.
*/
var basicClientAuthSecurity = func(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID, clientSecret, ok := r.BasicAuth()
		if !ok {
			clientID = r.PostForm.Get("client_id")
			clientSecret = r.PostForm.Get("client_secret")
			if len(clientID) == 0 || len(clientSecret) == 0 {
				log.Debug("client id or secret are empty", zap.String("client", clientID))
				http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
				return
			}
		}
		clientCtx, ok := authenticateClient(clientID, clientSecret)
		if !ok {
			log.Debug("invalid client authentication attempt", zap.String("client", clientID))
			http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
			return
		}
		ctx := newContextWithClient(r.Context(), clientCtx)
		f(w, r.WithContext(ctx))
	}
}

//verifies basic authentication and adds the user to the context
var basicUserAuthSecurity = func(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Debug("invalid basic auth login attempt", zap.String("user", username))
			http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
			return
		}
		user, err := authenticateUser(username, password)
		if err != nil {
			log.Error("invalid user authentication attempt", zap.String("user", username))
			http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
			return
		}
		ctx := newContextWithUser(r.Context(), user)
		f(w, r.WithContext(ctx))
	}
}

//verifies session cookie authentication
var sessionCookieSecurity = func(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := validateLoginSession(w, r)
		if err == nil && user != nil {
			ctx := newContextWithUser(r.Context(), user)
			f(w, r.WithContext(ctx))
			return
		}
		if err != nil || user == nil {
			log.Error("can not validate login session", zap.Error(err))
			addTargetURLToSession(w, r)
			http.Redirect(w, r, "/bounzr/login", 302)
			return
		}
	}
}

//select between cookie or basic security
//TODO what if I want to login with basicAuth or Cookie interchangeably
/*
var checkSecurity = func(f http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if _, _, ok := r.BasicAuth(); ok{
			basicUserAuthSecurity(f)
		}
		//default
		sessionCookieSecurity(f)

	}
}
*/
