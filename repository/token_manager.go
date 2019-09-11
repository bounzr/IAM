package repository

import (
	"../config"
	"../oauth2"
	"../utils"
	"bytes"
	"github.com/gofrs/uuid"
	"strings"
	"time"

	"go.uber.org/zap"
)

type TokenManager interface {
	deleteAccessToken(tokenHint *oauth2.AccessTokenHint)
	init()
	close()
	setTokenUnit(token *oauth2.TokenUnit) error
	setAuthorizationCode(code *oauth2.AuthorizationCode) error
	getTokenUnit(tokenHint *oauth2.AccessTokenHint) (token *oauth2.TokenUnit, ok bool)
	validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (code *oauth2.AuthorizationCode, ok bool)
}

func AuthorizationCodeGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.AuthorizationCodeAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsCodeRequest(client, request) {
		log.Error("client did not accept the code request", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	authCode, ok := tokenManager.validateAuthorizationCode(request)
	if !ok {
		log.Error("invalid authorization code", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	options := &oauth2.AccessTokenOptions{
		ClientID:        authCode.ClientID,
		AddRefreshToken: true,
		Scope:           authCode.Scope,
		OwnerID:         authCode.OwnerID,
	}
	return options
}

func DeleteOauth2AccessToken(tokenHint *oauth2.AccessTokenHint) {
	tokenManager.deleteAccessToken(tokenHint)
}

func getAccessTokens(opt *oauth2.AccessTokenOptions) (accessToken, refreshToken *oauth2.TokenUnit) {
	accessToken, refreshToken = oauth2.NewTokenSet(opt, config.IAM.Tokens.GetAccessDuration(), config.IAM.Tokens.GetRefreshDuration())
	err := tokenManager.setTokenUnit(accessToken)
	if err != nil {
		log.Error("could not get access token", zap.String("client ID", opt.ClientID.String()), zap.String("owner ID", opt.OwnerID.String()))
		return nil, nil
	}
	if refreshToken != nil {
		err = tokenManager.setTokenUnit(refreshToken)
		if err != nil {
			log.Error("could not get refresh token", zap.String("client ID", opt.ClientID.String()), zap.String("owner ID", opt.OwnerID.String()))
			return accessToken, nil
		}
	}
	return accessToken, refreshToken
}

func getRefreshTokens(token *oauth2.TokenUnit) (accessToken, refreshToken *oauth2.TokenUnit) {
	refreshToken = oauth2.NewRefreshToken(token, config.IAM.Tokens.GetRefreshDuration())
	err := tokenManager.setTokenUnit(refreshToken)
	if err != nil {
		log.Error("could not get refresh token", zap.String("client ID", token.ClientID.String()), zap.String("owner ID", token.OwnerID.String()))
		return token, nil
	}
	return token, refreshToken
}

func initTokens() {
	implementation := config.IAM.Tokens.Implementation
	switch implementation {
	case "leveldb":
		tokenManager = &TokenManagerLeveldb{tokensPath: "./rep/token", codesPath: "./rep/code"}
	default:
		tokenManager = &TokenManagerBasic{}
	}
	tokenManager.init()
}

func IntrospectAccessToken(hint *oauth2.AccessTokenHint) (response *oauth2.IntrospectionResponse) {
	token, ok := ValidateAccessToken(hint)
	if !ok {
		response = &oauth2.IntrospectionResponse{
			Active: false,
		}
		return
	}
	response = token.GetIntrospectionResponse()
	client, ok := GetClient(uuid.FromStringOrNil(response.ClientID))
	//todo audience is the protected resource
	if ok {
		log.Debug("client URI found", zap.String("id", response.ClientID), zap.String("URI", client.URI))
		response.Audience = client.URI
	} else {
		log.Debug("client URI not found", zap.String("id", response.ClientID))
	}

	if response.OwnerID == response.ClientID {
		response.OwnerName = client.Name
	} else {
		owner, ok := GetUser(uuid.FromStringOrNil(response.OwnerID))
		if ok {
			log.Debug("owner username found", zap.String("id", response.OwnerID), zap.String("username", owner.UserName))
			response.OwnerName = owner.UserName
		} else {
			log.Debug("owner username not found", zap.String("id", response.OwnerID))
		}
	}

	response.Issuer = "https://" + config.IAM.Server.Hostname + ":" + config.IAM.Server.Port
	response.TokenID = string(token.Token)
	return
}

func ValidateAccessToken(hint *oauth2.AccessTokenHint) (token *oauth2.TokenUnit, ok bool) {
	token, ok = tokenManager.getTokenUnit(hint)
	if !ok {
		return nil, false
	}
	if !token.Active {
		return nil, false
	}
	if !utils.InTimeSpan(token.IssuedAt, token.ExpirationTime, time.Now()) {
		token.Active = false
		tokenManager.setTokenUnit(token)
		return nil, false
	}
	return token, ok
}

func ImplicitGrantOptions(userCtx *UserCtx, request *oauth2.AuthorizationRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(uuid.FromStringOrNil(request.ClientID))
	if !found {
		log.Error("request not valid", zap.String("client ID", request.ClientID), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsImplicitRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client ID", request.ClientID), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	options := &oauth2.AccessTokenOptions{
		ClientID:        client.ID,
		AddRefreshToken: false,
		Scope:           []byte(validScope),
		OwnerID:         userCtx.UserID,
		State:           request.State,
	}
	return options
}

func ClientCredentialsGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.ClientCredentialsAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsClientCredentialsRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	options := &oauth2.AccessTokenOptions{
		ClientID:        client.ID,
		AddRefreshToken: false,
		Scope:           []byte(validScope),
		OwnerID:         client.ID,
	}
	return options
}

func OwnerPasswordGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.OwnerPasswordAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsPasswordRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	userCtx, valid := ValidateUser(request.Username, request.Password)
	if !valid {
		log.Error("invalid user credentials", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrAccessDeniedInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	options := &oauth2.AccessTokenOptions{
		ClientID:        client.ID,
		AddRefreshToken: true,
		Scope:           []byte(validScope),
		OwnerID:         userCtx.UserID,
	}
	return options
}

func RefreshTokenGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.RefreshAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsRefreshRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	refreshToken, ok := tokenManager.getTokenUnit(request.GetAccessTokenHint())
	if !ok {
		log.Error("invalid refresh token", zap.String("client ID", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrAccessDeniedInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	validScope = refreshToken.ValidateScope(validScope)

	options := &oauth2.AccessTokenOptions{
		ClientID:        refreshToken.GetClient(),
		AddRefreshToken: true,
		Scope:           []byte(validScope),
		OwnerID:         refreshToken.GetResourceOwner(),
	}
	return options
}

func RequestAccessToken(opt *oauth2.AccessTokenOptions) (response *oauth2.AccessTokenResponse, err error) {
	//initial validation. Get owner and client
	if opt == nil {
		return nil, oauth2.ErrInvalidRequest
	}
	client, found := GetClient(opt.ClientID)
	if !found {
		return nil, oauth2.ErrUnauthorizedClient
	}
	//ownerID can be the client self in client credentials grant
	var owner AccessTokenHolder
	if opt.ClientID == opt.OwnerID {
		log.Debug("ownerID equals clientID")
		owner = client
	} else {
		user, ok := GetUser(opt.OwnerID)
		if !ok {
			return nil, oauth2.ErrUnauthorizedClient
		}
		owner = user
	}

	//refresh token always removed and replaced if requested
	owner.DeleteClientRefreshToken(client.ID)

	//review if owner previously approved an access token for client
	accessTokenHint, ok := owner.GetClientAccessToken(client.ID)
	var accessToken *oauth2.TokenUnit
	if ok {
		//review if token is valid in the repository
		accessToken, ok = ValidateAccessToken(accessTokenHint)
		//return if token still exists and valid or remove invalid token from user
		if !ok {
			owner.DeleteClientAccessToken(client.ID)
		}
	}

	//compare scopes requested with access token
	if accessToken != nil {
		space := []byte{' '}
		oldScope := accessToken.Scope
		newScope := bytes.Split(opt.Scope, space)

		for _, element := range newScope {
			if !bytes.Contains(oldScope, element) {
				owner.DeleteClientAccessToken(client.ID)
				accessToken = nil
				break
			}
		}
	}

	//access token exists but regenerate refresh token required
	if accessToken != nil {
		//is refresh token requested?
		if opt.AddRefreshToken {
			response = owner.SetClientTokens(getRefreshTokens(accessToken))
		} else {
			owner.DeleteClientRefreshToken(opt.ClientID)
			response = owner.SetClientTokens(accessToken, nil)
		}
	} else {
		//generate all token from scratch since we couldnt find an old token
		response = owner.SetClientTokens(getAccessTokens(opt))
	}
	switch owner.(type) {
	case *User:
		user := owner.(*User)
		users, err := getUserRepository(user.RepositoryName)
		if err != nil {
			return nil, err
		}
		users.setUser(user)
	case *Client:
		client := owner.(*Client)
		clientManager.setClient(client)
	default:
		log.Error("token holder type not identifed", zap.String("resource id", opt.ClientID.String()))
	}

	return response, nil
}

//RequestAuthorizationCode returns authorization code response or error
func RequestAuthorizationCode(context *UserCtx, authorizationRequest *oauth2.AuthorizationRequest) (response *oauth2.AuthorizationCodeResponse, err error) {
	rep, err := getUserRepository(context.RepositoryName)
	if err != nil {
		return nil, err
	}
	user, ok := rep.getUser(context.GetUserID())
	if !ok {
		return nil, oauth2.ErrUnauthorizedClient
	}
	code := oauth2.NewAuthorizationCode(user.ID, authorizationRequest)
	err = tokenManager.setAuthorizationCode(code)
	if err != nil {
		return nil, err
	}
	return code.GetAuthorizationCodeResponse(), nil
}

func validateClientAllowsCodeRequest(cli *Client, request *oauth2.AuthorizationCodeAccessTokenRequest) (ok bool) {
	ok = true
	if strings.Compare(cli.ID.String(), request.ClientID) != 0 {
		log.Debug("client ID does not match", zap.String("client ID", cli.ID.String()), zap.String("requested", request.ClientID))
		ok = false
		return
	}
	if !cli.HasGrantType(request.GrantType) {
		log.Debug("grant type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token.String()) {
		log.Debug("response type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	if !cli.HasRedirectURI(request.RedirectURI) {
		log.Debug("redirect_uri not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.RedirectURI))
		ok = false
		return
	}
	return
}

func validateClientAllowsImplicitRequest(cli *Client, request *oauth2.AuthorizationRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(oauth2.ImplicitGrantType.String()) {
		log.Debug("grant type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", oauth2.ImplicitGrantType.String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token.String()) {
		log.Debug("response type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	if !cli.HasRedirectURI(request.RedirectURI) {
		log.Debug("redirect_uri not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.RedirectURI))
		ok = false
		return
	}
	return
}

func validateClientAllowsClientCredentialsRequest(cli *Client, request *oauth2.ClientCredentialsAccessTokenRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(request.GrantType) {
		log.Debug("grant type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token.String()) {
		log.Debug("response type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	return
}

func validateClientAllowsPasswordRequest(cli *Client, request *oauth2.OwnerPasswordAccessTokenRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(request.GrantType) {
		log.Debug("grant type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token.String()) {
		log.Debug("response type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	return
}

func validateClientAllowsRefreshRequest(cli *Client, request *oauth2.RefreshAccessTokenRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(request.GrantType) {
		log.Debug("grant type not found for client", zap.String("client ID", cli.ID.String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token.String()) {
		log.Debug("response type not found for client: %s", zap.String("client ID", cli.ID.String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	return
}

//ValidateAuthorizationRequest verifies that the authorization request is compliant
func ValidateAuthorizationRequest(authorizationRequest *oauth2.AuthorizationRequest) error {
	//1.Immediate display errors
	//Verification of clientID
	clientID := uuid.FromStringOrNil(authorizationRequest.ClientID)
	if len(clientID) == 0 {
		log.Error("client ID is empty", zap.Error(oauth2.ErrClientIdentifierInfo))
		return oauth2.ErrClientIdentifierInfo
	}
	client, found := GetClient(clientID)
	if !found {
		log.Error("client ID not found", zap.Error(oauth2.ErrClientIdentifierInfo))
		return oauth2.ErrClientIdentifierInfo
	}

	//Verification of redirection URI
	redirectURI := authorizationRequest.RedirectURI
	if len(redirectURI) == 0 {
		log.Error("redirect url is empty", zap.Error(oauth2.ErrRedirectionURIInfo))
		return oauth2.ErrRedirectionURIInfo
	}
	if !client.HasRedirectURI(redirectURI) {
		log.Error("client does not have redirect url", zap.Error(oauth2.ErrRedirectionURIInfo))
		return oauth2.ErrRedirectionURIInfo
	}

	//2. Redirect errors
	//Verification of response type
	responseType := authorizationRequest.ResponseType
	if len(responseType) == 0 {
		log.Error("response type is empty for client", zap.String("client ID", clientID.String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return oauth2.ErrInvalidRequestInfo
	}
	//compare if the grant is code
	if strings.Compare(responseType, "code") == 0 {
		if !(client.HasGrantType(oauth2.AuthorizationCodeGrantType.String()) && client.HasResponseType(oauth2.Code.String())) {
			log.Error("authorization code must have response type code", zap.String("client ID", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
			return oauth2.ErrUnauthorizedClientInfo
		}
		//compare if the grant is token (implicit)
	} else if strings.Compare(responseType, "token") == 0 {
		if !(client.HasGrantType(oauth2.ImplicitGrantType.String()) && client.HasResponseType(oauth2.Token.String())) {
			log.Error("Implicit grant type must have response type token", zap.String("client ID", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
			return oauth2.ErrUnauthorizedClientInfo
		}
		//no grant type or response type
	} else {
		log.Error("no grant type nor response type", zap.String("client ID", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
		return oauth2.ErrUnauthorizedClientInfo
	}

	//3. State warning
	//Verification of state otherwise warn the console
	if len(authorizationRequest.State) == 0 {
		log.Warn("client is not providing a state. RFC6749 recommends the use of a state", zap.String("client ID", clientID.String()))
	}

	//4. Modification of the request
	//Verification and modification of client allowed scopes
	reqScopes := authorizationRequest.Scope
	authorizationRequest.Scope = client.ValidateScope(reqScopes)

	return nil
}
