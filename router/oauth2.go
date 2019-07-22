package router

import (
	"../pages"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"

	"../oauth2"
	"../repository"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//newOauth2Router returns router with Oauth2 related routes
func newOauth2Router(router *mux.Router) {
	router.HandleFunc("/authorize", chain(
		oauth2AuthorizeHandler,
		sessionCookieSecurity),
	).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/introspect", chain(
		oauth2IntrospectHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins", "ProtectedResources")),
	).Methods("POST")
	//TODO according to rfc anonymous registration is allowed, token may be allowed.
	router.HandleFunc("/register", chain(oauth2RegisterHandlerPost, basicUserAuthSecurity)).Methods("POST")
	//TODO according to rfc authorization must be token and not basic. Replace basicUserAuthSecurity
	router.HandleFunc("/register/{id:[-a-zA-Z0-9]+}", chain(
		oauth2RegisterHandlerGet,
		basicUserAuthSecurity)).Methods("GET")
	router.HandleFunc("/revoke", chain(oauth2RevokeHandlerPost, basicClientAuthSecurity)).Methods("POST")
	router.HandleFunc("/token", chain(oauth2TokenHandlerPost, basicClientAuthSecurity)).Methods("POST")
}

func oauth2AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		oauth2AuthorizeGetHandler(w, r)
		return
	}
	if r.Method == http.MethodPost {
		oauth2AuthorizePostHandler(w, r)
		return
	}
}

func authorizationRequestErrorRedirect(w http.ResponseWriter, r *http.Request, authorizationRequest *oauth2.AuthorizationRequest, err string) {
	redirectURI := authorizationRequest.RedirectURI
	state := authorizationRequest.State
	errRedirectURI := fmt.Sprintf("%s?error=%s", redirectURI, err)
	if len(state) > 0 {
		errRedirectURI = fmt.Sprintf("%s?error=%s&state=%s", redirectURI, err, state)
	}
	log.Debug("redirecting authorize error to uri ", zap.String("error", err), zap.String("uri", redirectURI))
	http.Redirect(w, r, errRedirectURI, http.StatusFound)
}

/**
4.1.1 Authorization Request

GET /authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
Host: server.example.com
*/
func oauth2AuthorizeGetHandler(w http.ResponseWriter, r *http.Request) {
	var authReq oauth2.AuthorizationRequest
	authorizationRequest := &authReq
	err := decoder.Decode(authorizationRequest, r.URL.Query())
	if err != nil {
		log.Error("can not decode authorization request", zap.String("client id", authorizationRequest.ClientID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//start authorization request validation
	err = repository.ValidateAuthorizationRequest(authorizationRequest)
	if err != nil {
		log.Error("not valid authorization request", zap.String("client id", authorizationRequest.ClientID), zap.Error(err))
	}
	//immediate display errors
	if err == oauth2.ErrClientIdentifierInfo {
		http.Error(w, oauth2.ErrClientIdentifierInfo.Error(), http.StatusBadRequest)
		return
	}
	if err == oauth2.ErrRedirectionURIInfo {
		http.Error(w, oauth2.ErrClientIdentifierInfo.Error(), http.StatusBadRequest)
		return
	}
	//redirect errors
	if err == oauth2.ErrInvalidRequestInfo {
		authorizationRequestErrorRedirect(w, r, authorizationRequest, oauth2.ErrInvalidRequest.Error())
		return
	}
	if err == oauth2.ErrUnauthorizedClientInfo {
		authorizationRequestErrorRedirect(w, r, authorizationRequest, oauth2.ErrUnauthorizedClient.Error())
		return
	}
	//end authorization request validation

	//get user from context as it has logged in by the middleware
	ctx := r.Context()
	usr, ok := fromContextGetUser(ctx)
	if !ok {
		//TODO redirect to login instead of error
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		log.Debug("can not get user from context", zap.Error(repository.ErrInvalidLogin))
		return
	}

	//generate request code for consents page
	consentToken, err := repository.SetAuthorizationRequest(usr, authorizationRequest)
	//TODO on session error must be handled better than a 500
	if err != nil {
		log.Error("can not set authorization request", zap.String("client id", authorizationRequest.ClientID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//todo this store must be available in a multiple nodes environment
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	//TODO on session error must be handled better than a 500
	if err != nil {
		log.Error("can not get session cookie", zap.String("client id", authorizationRequest.ClientID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values[ConsentsToken] = consentToken
	session.Save(r, w)

	//Render consents page
	var p *pages.AuthorizePage

	if len(authorizationRequest.GetScopesList()) == 0 {
		p = &pages.AuthorizePage{
			ClientID:   authorizationRequest.ClientID,
			ClientName: authorizationRequest.ClientID,    //todo get client name from repository
			ClientURI:  authorizationRequest.RedirectURI, //todo get client URI from repository
		}
	} else {
		p = &pages.AuthorizePage{
			ClientID:   authorizationRequest.ClientID,
			ClientName: authorizationRequest.ClientID,    //todo get client name from repository
			ClientURI:  authorizationRequest.RedirectURI, //todo get client URI from repository
			ScopesList: authorizationRequest.GetScopesList(),
		}
	}

	err = pages.RenderPage(w, "authorize.html", p)
	if err != nil {
		log.Error("can not render authorize.html", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return

}

func oauth2AuthorizePostHandler(w http.ResponseWriter, r *http.Request) {
	//get user from context as it has logged in by the middleware
	ctx := r.Context()
	usrCtx, ok := fromContextGetUser(ctx)
	if !ok {
		//TODO redirect to login instead of error
		log.Error("can not get user from context", zap.Error(repository.ErrInvalidLogin))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	//get consents token from the session
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	if err != nil {
		//TODO on session error must be handled better than a 500
		log.Error("can not get session cookie", zap.Error(repository.ErrSessionInvalid))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	val, ok := session.Values[ConsentsToken]
	if !ok {
		if val == nil {
			log.Error("can not get consents token from session cookie", zap.Error(repository.ErrSessionNotFound))
		}
		//without consents token there is no idea where the call is coming from and where is it going.
		http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusForbidden)
		return
	}

	var ct = &repository.ConsentToken{}
	ct, ok = val.(*repository.ConsentToken)
	if !ok {
		log.Error("can not cast consent token from session value", zap.Error(repository.ErrSessionInvalid))
		//TODO on session error must be handled better. Return to login
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	authReq, err := repository.GetAuthorizationRequest(usrCtx, ct)
	if err != nil {
		log.Error("can not get authorization request", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName), zap.Error(err))
		//TODO on session error must be handled better. Return to login
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Error("can not parse authorization form data", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName), zap.Error(err))
		//TODO on session error must be handled better than a 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if deny := r.FormValue("deny"); len(deny) > 0 {
		//TODO redirect error message
		log.Debug("authorization form denied", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if approve := r.FormValue("approve"); len(approve) == 0 {
		//TODO redirect error message
		log.Error("authorization form denied", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var approvedScopes []string
	for key, value := range r.Form {
		if strings.HasPrefix(key, "scope_") && strings.Compare(value[0], "on") == 0 {
			scp := strings.TrimPrefix(key, "scope_")
			approvedScopes = append(approvedScopes, scp)
		}
	}
	authReq.MatchScopes(approvedScopes)

	//authorization code grant response
	if strings.Compare(authReq.ResponseType, "code") == 0 {
		authResponse, err := repository.RequestAuthorizationCode(usrCtx, authReq)
		if err != nil {
			log.Error("can not get authorization request", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName), zap.Error(err))
			//TODO on session error must be handled better than a 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//todo encode to URL. Gorilla maybe?
		//HTTP/1.1 302 Found
		//Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz
		clientURL := fmt.Sprintf("%s?code=%s&state=%s", authResponse.ClientURI, authResponse.Code, authResponse.State)
		//todo in case of error
		//HTTP/1.1 302 Found
		//Location: https://client.example.com/cb?error=access_denied&state=xyz
		http.Redirect(w, r, clientURL, http.StatusFound)
		return
	}

	//implicit code grant response
	if strings.Compare(authReq.ResponseType, "token") == 0 {
		options := repository.ImplicitGrantOptions(usrCtx, authReq)
		tokenResponse, err := repository.RequestAccessToken(options)
		if err != nil {
			//TODO on session error must be handled better than a 500
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		//todo in case of error
		//HTTP/1.1 302 Found
		//Location: https://client.example.com/cb?error=access_denied&state=xyz

		clientURL := ""
		if len(strings.TrimSpace(tokenResponse.Scope)) == 0 {
			//HTTP/1.1 302 Found
			//Location: http://example.com/cb#access_token=2YotnFZFEjr1zCsicMWpAA&state=xyz&token_type=example&expires_in=3600
			clientURL = fmt.Sprintf("%s#access_token=%s&state=%s&token_type=%s&expires_in=%d", authReq.RedirectURI,
				tokenResponse.AccessToken, tokenResponse.State, tokenResponse.TokenAuthType, tokenResponse.ExpiresIn)

		} else {
			clientURL = fmt.Sprintf("%s#access_token=%s&state=%s&token_type=%s&expires_in=%d&scope=%s", authReq.RedirectURI,
				tokenResponse.AccessToken, tokenResponse.State, tokenResponse.TokenAuthType, tokenResponse.ExpiresIn, tokenResponse.Scope)
		}
		http.Redirect(w, r, clientURL, http.StatusFound)
		return
	}
	log.Error("request could not be handled by any grant", zap.String("user id", usrCtx.UserID.String()), zap.String("username", usrCtx.UserName), zap.Error(err))
	//TODO on session error must be handled better than a 500
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

func oauth2IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//token := r.FormValue("token")
	//tokenType := r.FormValue("token_type_hint")
	hint := &oauth2.AccessTokenHint{}
	err = decoder.Decode(hint, r.PostForm)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	//verify that token hint is valid if exist
	if len(hint.Hint) > 0 {
		if oauth2.NewTokenHintType(hint.Hint) == oauth2.NullTokenHintType {
			log.Error("token hint not supported", zap.String("hint", hint.Hint), zap.Error(oauth2.ErrUnsupportedResponseTypeInfo))
			http.Error(w, oauth2.ErrUnsupportedTokenType.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	introspection := repository.IntrospectAccessToken(hint)
	//Marshal clientInfResp to json and write to response
	intJson, err := json.Marshal(introspection)
	if err != nil {
		log.Error("can not marshal introspection token", zap.Error(err))
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(intJson)
	return
}

/**
     GET /register/s6BhdRkqt3 HTTP/1.1
     Accept: application/json
     Host: server.example.com
     Authorization: Bearer reg-23410913-abewfq.123483
**/
func oauth2RegisterHandlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	client, ok := repository.GetClient(uuid.FromStringOrNil(id))
	if !ok {
		log.Error("can not get client", zap.String("client id", id))
		//TODO remove token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		//TODO marshal the error request
		w.Write([]byte("{\"request\":\"not valid\"}"))
		return
	}
	owner := client.OwnerID

	//TODO clean context() after it was used leaving no trace
	var userCtx *repository.UserCtx
	if loggedUser := r.Context().Value(userCtxKey); loggedUser != nil {
		//todo add or handle anonymous scope to user and client
		userCtx = loggedUser.(*repository.UserCtx)
		log.Debug("found logged user context", zap.String("user id", userCtx.GetUserID().String()))
	} else {
		log.Debug("no logged user context key found", zap.String("user id", userCtx.GetUserID().String()))
		//TODO remove token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"request\":\"not valid\"}"))
		return
	}
	//TODO also check if the owner userCtx is in the user scopes
	if userCtx.GetUserID() != owner {
		log.Debug("user from context is not the owner", zap.String("user id", id))
		//TODO remove token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{\"request\":\"not valid\"}"))
		return
	}
	clientInfResp, err := client.GetClientInformationResponse()
	if err != nil {
		log.Error("can not get client information response", zap.String("client id", client.ID.String()), zap.Error(err))
	}
	//Marshal clientInfResp to json and write to response
	cirJSON, err := json.Marshal(clientInfResp)
	if err != nil {
		log.Error("can not marshal the client information response", zap.String("client id", client.ID.String()), zap.Error(err))
	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(cirJSON)
	return
}

//oauth2RegisterHandlerPost is used to register a new Oauth2 client
func oauth2RegisterHandlerPost(w http.ResponseWriter, r *http.Request) {
	//initialize empty request
	clientReq := &oauth2.ClientRegistrationRequest{}
	//Parse json request body and use it to set fields
	err := json.NewDecoder(r.Body).Decode(clientReq)
	if err != nil {
		log.Error("can not decode the client registration request", zap.String("client name", clientReq.ClientName), zap.Error(err))
	}

	//todo process client registration
	client, err := repository.NewClientFromOauth(clientReq)
	if err != nil {
		log.Error("can not create new client", zap.String("client name", clientReq.ClientName), zap.Error(err))
		//todo client.GetClientRegistrationError
	}
	//TODO clean context() after it was used leaving no trace
	//TODO use better func fromContextGetUser(ctx context.Context) (*repository.User, bool)
	if loggedUser := r.Context().Value(userCtxKey); loggedUser != nil {
		//todo add or handle anonymous scope to user and client
		profile := loggedUser.(*repository.UserCtx)
		log.Debug("found logged user context", zap.String("user id", profile.GetUserID().String()))
		client.OwnerID = profile.UserID
	}
	repository.AddClient(client)
	//todo if AddClient() err != nil
	var clientInfResp *oauth2.ClientInformationResponse
	clientInfResp, err = client.GetClientInformationResponse()
	if err != nil {
		log.Error("can not get client information response", zap.String("client id", client.ID.String()), zap.Error(err))
	}
	//Marshal clientInfResp to json and write to response
	cirJSON, err := json.Marshal(clientInfResp)
	if err != nil {
		log.Error("can not marhsal client information response", zap.String("client id", client.ID.String()), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(cirJSON)
	return
}

func oauth2RevokeHandlerPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error("can not parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//token := r.FormValue("token")
	//tokenType := r.FormValue("token_type_hint")
	hint := &oauth2.AccessTokenHint{}
	err = decoder.Decode(hint, r.PostForm)
	if err != nil {
		log.Error("can not decode access token hint", zap.String("hint", hint.Hint), zap.Error(err))
		http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	if len(hint.Hint) > 0 {
		if oauth2.NewTokenHintType(hint.Hint) == oauth2.NullTokenHintType {
			log.Error("token hint not supported", zap.String("hint", hint.Hint), zap.Error(oauth2.ErrUnsupportedResponseTypeInfo))
			http.Error(w, oauth2.ErrUnsupportedTokenType.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	repository.DeleteOauth2AccessToken(hint)

	w.WriteHeader(http.StatusOK)
	return
}

func oauth2TokenHandlerPost(w http.ResponseWriter, r *http.Request) {
	//get client from context as it has logged in by the middleware
	ctx := r.Context()
	cliCtx, ok := fromContextGetClient(ctx)
	if !ok {
		log.Debug("can not get client from context")
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Error("can not parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	grant := r.FormValue("grant_type")
	var options *oauth2.AccessTokenOptions
	//authorization code grant
	if strings.Compare(grant, "authorization_code") == 0 {
		var tokenReq oauth2.AuthorizationCodeAccessTokenRequest
		err = decoder.Decode(&tokenReq, r.PostForm)
		if err != nil {
			log.Error("can not decode authorization code access token request", zap.Error(err))
			http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		options = repository.AuthorizationCodeGrantOptions(cliCtx, &tokenReq)
		//password grant
	} else if strings.Compare(grant, "client_credentials") == 0 {
		var tokenReq oauth2.ClientCredentialsAccessTokenRequest
		err = decoder.Decode(&tokenReq, r.PostForm)
		if err != nil {
			log.Error("can not decode client credentials access token request", zap.Error(err))
			http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		//todo options = repository.ClientCredentialsGrantOptions(cliCtx, &tokenReq)
	} else if strings.Compare(grant, "password") == 0 {
		var tokenReq oauth2.OwnerPasswordAccessTokenRequest
		err = decoder.Decode(&tokenReq, r.PostForm)
		if err != nil {
			log.Error("can not decode owner password access token request", zap.Error(err))
			http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		options = repository.OwnerPasswordGrantOptions(cliCtx, &tokenReq)
		//client credentials grant
	} else if strings.Compare(grant, "refresh_token") == 0 {
		var tokenReq oauth2.RefreshAccessTokenRequest
		err = decoder.Decode(&tokenReq, r.PostForm)
		if err != nil {
			log.Error("can not decode refresh access token request", zap.Error(err))
			http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
			return
		}
		options = repository.RefreshTokenGrantOptions(cliCtx, &tokenReq)
	} else {
		log.Error("can not process access token request", zap.Error(oauth2.ErrInvalidRequestInfo))
		http.Error(w, oauth2.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	tokenResponse, err := repository.RequestAccessToken(options)
	if err != nil {
		log.Error("can not process access token request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	atrJSON, err := json.Marshal(tokenResponse)
	if err != nil {
		log.Error("can not marshal token response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	//todo cache and other security headers. See standard for access token response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(atrJSON)
}
