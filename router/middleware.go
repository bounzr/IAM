package router

import (
	"net/http"

	"go.uber.org/zap"

	"../repository"
)

// middleware - a type that describes a middleware, at the core of this implementation a middleware is merely a function
// that takes a handler function, and returns a handler function.
type middleware func(http.HandlerFunc) http.HandlerFunc

// chain - a function that takes a handler function, a list of middlewares and creates a new application stack
// as a single http handler
func chain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	// if there are no more middlewares, we just return the handlerfunc, as we are done recursing.
	if len(m) == 0 {
		return f
	}
	// otherwise pop the middleware from the list, and call build chain recursively as it's parameter
	return m[0](chain(f, m[1:cap(m)]...))
}

/**
Either pass client id and secret via Basic Authentication or add client_id and client_secret fields to the POST body.
When providing the client_id and client_secret in the Authorization header it is expected to be:
    client_id:client_secret
    Base64 encoded
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
		user, valid := authenticateUser(username, password)
		if !valid {
			log.Debug("invalid user authentication attempt", zap.String("user", username))
			http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
			return
		}
		ctx := newContextWithUser(r.Context(), user)
		f(w, r.WithContext(ctx))
	}
}

//method to verify that the logged in user is registered in the given groups
func verifyUserGroups(groups ...string) middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := fromContextGetUser(r.Context())
			if !ok {
				log.Error("user must be authenticated in order to verify groups")
				http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
				return
			}
			if !validateIsUserInGroup(user, groups) {
				log.Error("user does not belong to requested group", zap.String("user", user.UserName))
				http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
				return
			}
			f(w, r)
		}
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
