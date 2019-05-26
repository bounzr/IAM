package repository

import (
	"../config"
	"../oauth2"
	"github.com/gofrs/uuid"
	"strings"


	"go.uber.org/zap"
)

type TokenRepository interface {
	deleteAccessToken(tokenHint *oauth2.AccessTokenHint)
	init()
	setAccessToken(accessToken *oauth2.AccessToken) error
	setAuthorizationCode(code *oauth2.AuthorizationCode) error
	validateAccessToken(tokenHint *oauth2.AccessTokenHint) (token *oauth2.AccessToken, ok bool)
	validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (code *oauth2.AuthorizationCode, ok bool)
}

func getAccessToken(opt *oauth2.AccessTokenOptions) (token *oauth2.AccessToken) {
	accessToken, refreshToken := oauth2.NewAccessToken(opt, config.IAM.Tokens.GetAccessDuration(), config.IAM.Tokens.GetRefreshDuration())
	_, ath, rth := accessToken.GetTokenHints()
	_, athExists := tokenManager.validateAccessToken(ath)
	rthExists := false
	if rth != nil {
		_, rthExists = tokenManager.validateAccessToken(rth)
	}
	if athExists || rthExists {
		return getAccessToken(opt)
	}
	tokenManager.setAccessToken(accessToken)
	if refreshToken != nil {
		tokenManager.setAccessToken(refreshToken)
	}
	return accessToken
}

func validateClientAllowsCodeRequest(cli *oauth2.Client, request *oauth2.AuthorizationCodeAccessTokenRequest) (ok bool) {
	ok = true
	if strings.Compare(cli.GetClientID().String(), request.ClientID) != 0 {
		log.Debug("client id does not match", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.ClientID))
		ok = false
		return
	}
	if !cli.HasGrantType(request.GetGrantType()) {
		log.Debug("grant type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token) {
		log.Debug("response type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	if !cli.HasRedirectUri(request.RedirectURI) {
		log.Debug("redirect_uri not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.RedirectURI))
		ok = false
		return
	}
	return
}

func validateClientAllowsImplicitRequest(cli *oauth2.Client, request *oauth2.AuthorizationRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(oauth2.ImplicitGrantType) {
		log.Debug("grant type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", oauth2.ImplicitGrantType.String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token) {
		log.Debug("response type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	if !cli.HasRedirectUri(request.RedirectURI) {
		log.Debug("redirect_uri not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.RedirectURI))
		ok = false
		return
	}
	return
}

func validateClientAllowsPasswordRequest(cli *oauth2.Client, request *oauth2.OwnerPasswordAccessTokenRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(request.GetGrantType()) {
		log.Debug("grant type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token) {
		log.Debug("response type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	return
}

func validateClientAllowsRefreshRequest(cli *oauth2.Client, request *oauth2.RefreshAccessTokenRequest) (ok bool) {
	ok = true
	if !cli.HasGrantType(request.GetGrantType()) {
		log.Debug("grant type not found for client", zap.String("client id", cli.GetClientID().String()), zap.String("requested", request.GetGrantType().String()))
		ok = false
		return
	}
	token, _ := oauth2.NewResponseType("token")
	if !cli.HasResponseType(token) {
		log.Debug("response type not found for client: %s", zap.String("client id", cli.GetClientID().String()), zap.String("requested", token.String()))
		ok = false
		return
	}
	return
}

func AuthorizationCodeGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.AuthorizationCodeAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsCodeRequest(client, request) {
		log.Error("client did not accept the code request", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	authCode, ok := tokenManager.validateAuthorizationCode(request)
	if !ok {
		log.Error("invalid authorization code", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
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

func IntrospectOauth2AccessToken(hint *oauth2.AccessTokenHint) (response *oauth2.IntrospectionResponse) {
	accessToken, ok := tokenManager.validateAccessToken(hint)
	if !ok {
		response = &oauth2.IntrospectionResponse{
			Active: "false",
		}
	} else {
		response = accessToken.GetIntrospectionResponse()
	}
	return
}

func ImplicitGrantOptions(userCtx *UserCtx, request *oauth2.AuthorizationRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(uuid.FromStringOrNil(request.ClientID))
	if !found {
		log.Error("request not valid", zap.String("client id", request.ClientID), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsImplicitRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client id", request.ClientID), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	options := &oauth2.AccessTokenOptions{
		ClientID:        client.GetClientID(),
		AddRefreshToken: false,
		Scope:           []byte(validScope),
		OwnerID:         userCtx.UserID,
		State:           request.State,
	}
	return options
}

func OwnerPasswordGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.OwnerPasswordAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsPasswordRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	userCtx, err := ValidateUser(request.Username, request.Password)
	if err != nil {
		log.Error("invalid user credentials", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrAccessDeniedInfo))
		return nil
	}
	validScope := client.ValidateScope(request.Scope)
	options := &oauth2.AccessTokenOptions{
		ClientID:        client.GetClientID(),
		AddRefreshToken: true,
		Scope:           []byte(validScope),
		OwnerID:         userCtx.UserID,
	}
	return options
}

func RefreshTokenGrantOptions(cliCtx *oauth2.ClientCtx, request *oauth2.RefreshAccessTokenRequest) *oauth2.AccessTokenOptions {
	client, found := GetClient(cliCtx.GetClientID())
	if !found {
		log.Error("request not valid", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	//validate that logged in client matches the access token request attributes
	if !validateClientAllowsRefreshRequest(client, request) {
		log.Error("client did not accept the request", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return nil
	}
	refreshToken, ok := tokenManager.validateAccessToken(request.GetAccessTokenHint())
	if !ok {
		log.Error("invalid refresh token", zap.String("client id", cliCtx.GetClientID().String()), zap.Error(oauth2.ErrAccessDeniedInfo))
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

func RequestOauth2AccessToken(opt *oauth2.AccessTokenOptions) (response *oauth2.AccessTokenResponse, err error) {
	//initial validation. Get owner and client
	if opt == nil {
		return nil, oauth2.ErrInvalidRequest
	}
	client, found := GetClient(opt.ClientID)
	if !found {
		return nil, oauth2.ErrUnauthorizedClient
	}
	userResource := &Metadata{}
	ok := resourceManager.loadResource(opt.OwnerID, userResource)
	//res, ok := resourceManager.getResource(opt.OwnerID)
	if !ok {
		return nil, oauth2.ErrUnauthorizedClient
	}
	/*
		if strings.Compare(res.GetResourceType(), "user") != 0{
			return nil, oauth2.ErrUnauthorizedClient
		}
		userResource, ok := res.(*UserMetadata)
		if !ok{
			return nil, oauth2.ErrUnauthorizedClient
		}
	*/
	users, err := getUserRepository(userResource.RepositoryName)
	if err != nil {
		return nil, err
	}
	owner, ok := users.getUser(userResource.Name)
	if !ok {
		return nil, oauth2.ErrUnauthorizedClient
	}

	//review if owner previously approved a token for client
	tokenHint, ok := owner.GetClientAccessToken(client.GetClientID())
	var token *oauth2.AccessToken
	if ok {
		//review if token is valid in the repository
		token, ok = tokenManager.validateAccessToken(tokenHint)
		//return if token still exists and valid or remove invalid token from user
		if ok {
			return token.GetAccessTokenResponse(), nil
		} else {
			owner.DeleteClientTokens(client.GetClientID())
		}
	}

	//generate all tokens from scratch since we couldnt find an old token
	token = getAccessToken(opt)
	err = tokenManager.setAccessToken(token)
	if err != nil {
		return nil, err
	}
	owner.SetClientTokens(token.GetTokenHints())
	users.setUser(owner)
	response = token.GetAccessTokenResponse()
	return response, nil
}

//RequestOauth2AuthorizationCode returns authorization code response or error
func RequestOauth2AuthorizationCode(context *UserCtx, authorizationRequest *oauth2.AuthorizationRequest) (response *oauth2.AuthorizationCodeResponse, err error) {
	rep, err := getUserRepository(context.RepositoryName)
	if err != nil {
		return nil, err
	}
	userResource := &Metadata{}
	ok := resourceManager.loadResource(context.GetUserID(), userResource)
	if !ok {
		return nil, oauth2.ErrUnauthorizedClient
	}
	user, ok := rep.getUser(userResource.GetName())
	if !ok {
		return nil, oauth2.ErrUnauthorizedClient
	}
	code := oauth2.NewAuthorizationCode(user.ID, authorizationRequest)
	err = tokenManager.setAuthorizationCode(code)
	if err != nil{
		return nil, err
	}
	return code.GetAuthorizationCodeResponse(), nil
}

//ValidateOauth2AuthorizationRequest verifies that the authorization request is compliant
func ValidateOauth2AuthorizationRequest(authorizationRequest *oauth2.AuthorizationRequest) error {
	//1.Immediate display errors
	//Verification of clientID
	clientID := uuid.FromStringOrNil(authorizationRequest.ClientID)
	if len(clientID) == 0 {
		log.Error("client id is empty", zap.Error(oauth2.ErrClientIdentifierInfo))
		return oauth2.ErrClientIdentifierInfo
	}
	client, found := GetClient(clientID)
	if !found {
		log.Error("client id not found", zap.Error(oauth2.ErrClientIdentifierInfo))
		return oauth2.ErrClientIdentifierInfo
	}

	//Verification of redirection URI
	redirectURI := authorizationRequest.RedirectURI
	if len(redirectURI) == 0 {
		log.Error("redirect url is empty", zap.Error(oauth2.ErrRedirectionURIInfo))
		return oauth2.ErrRedirectionURIInfo
	}
	if !client.HasRedirectUri(redirectURI) {
		log.Error("client does not have redirect url", zap.Error(oauth2.ErrRedirectionURIInfo))
		return oauth2.ErrRedirectionURIInfo
	}

	//2. Redirect errors
	//Verification of response type
	responseType := authorizationRequest.ResponseType
	if len(responseType) == 0 {
		log.Error("response type is empty for client", zap.String("client id", clientID.String()), zap.Error(oauth2.ErrInvalidRequestInfo))
		return oauth2.ErrInvalidRequestInfo
	}
	//compare if the grant is code
	if strings.Compare(responseType, "code") == 0 {
		if !(client.HasGrantType(oauth2.AuthorizationCodeGrantType) && client.HasResponseType(oauth2.Code)) {
			log.Error("authorization code must have response type code", zap.String("client id", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
			return oauth2.ErrUnauthorizedClientInfo
		}
		//compare if the grant is token (implicit)
	} else if strings.Compare(responseType, "token") == 0 {
		if !(client.HasGrantType(oauth2.ImplicitGrantType) && client.HasResponseType(oauth2.Token)) {
			log.Error("Implicit grant type must have response type token", zap.String("client id", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
			return oauth2.ErrUnauthorizedClientInfo
		}
		//no grant type or response type
	} else {
		log.Error("no grant type nor response type", zap.String("client id", clientID.String()), zap.Error(oauth2.ErrUnauthorizedClientInfo))
		return oauth2.ErrUnauthorizedClientInfo
	}

	//3. State warning
	//Verification of state otherwise warn the console
	if len(authorizationRequest.State) == 0 {
		log.Warn("client is not providing a state. RFC6749 recommends the use of a state", zap.String("client id", clientID.String()))
	}

	//4. Modification of the request
	//Verification and modification of client allowed scopes
	reqScopes := authorizationRequest.Scope
	authorizationRequest.Scope = client.ValidateScope(reqScopes)

	return nil
}
