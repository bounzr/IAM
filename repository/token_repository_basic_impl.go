package repository

import (
	"../oauth2"
	"time"
	"go.uber.org/zap"
)

type TokenRepositoryBasic struct {
	//string is the authorization code
	authorizationCodes map[string]*oauth2.AuthorizationCode

	//string is the token
	accessTokens      map[string]*oauth2.AccessToken
	refreshTokens     map[string]*oauth2.AccessToken
	usedTokens        map[string]struct{}
	blackListedTokens map[string]struct{}
}

func NewTokenRepositoryBasic() TokenRepository {
	mtr := &TokenRepositoryBasic{}
	mtr.init()
	return mtr
}

//init the repository
func (r *TokenRepositoryBasic) init() {
	r.authorizationCodes = make(map[string]*oauth2.AuthorizationCode)
	r.accessTokens = make(map[string]*oauth2.AccessToken)
	r.refreshTokens = make(map[string]*oauth2.AccessToken)
	r.usedTokens = make(map[string]struct{})
	r.blackListedTokens = make(map[string]struct{})
}

func (r *TokenRepositoryBasic) deleteAccessToken(tokenHint *oauth2.AccessTokenHint) {
	r.blackListedTokens[string(tokenHint.Token)] = struct{}{}
	delete(r.accessTokens, string(tokenHint.Token))
	delete(r.refreshTokens, string(tokenHint.Token))
}

func (r *TokenRepositoryBasic) setAuthorizationCode(code *oauth2.AuthorizationCode) error {
	r.authorizationCodes[code.Code] = code
	return nil
}

func (r *TokenRepositoryBasic) setAccessToken(accessToken *oauth2.AccessToken) error {
	if accessToken.TokenHintType == oauth2.AccessTokenHintType {
		at := string(accessToken.GetToken())
		r.accessTokens[at] = accessToken
	}
	if accessToken.TokenHintType == oauth2.RefreshTokenHintType {
		at := string(accessToken.GetToken())
		r.refreshTokens[at] = accessToken
	}
	return nil
}

func (r *TokenRepositoryBasic) validateAccessToken(hint *oauth2.AccessTokenHint) (*oauth2.AccessToken, bool) {
	var token *oauth2.AccessToken
	var ok = false
	if oauth2.NewTokenHintType(hint.Hint) != oauth2.RefreshTokenHintType {
		token, ok = r.accessTokens[string(hint.Token)]
	}
	if !ok && oauth2.NewTokenHintType(hint.Hint) != oauth2.AccessTokenHintType {
		token, ok = r.refreshTokens[string(hint.Token)]
	}
	//verify that it is not expired
	if ok {
		ok = token.GetExpirationTime() > time.Now().Unix()
	}
	//verify that the token is not blacklisted
	if ok {
		_, isBlack := r.blackListedTokens[string(hint.Token)]
		ok = !isBlack
		if !ok {
			r.deleteAccessToken(hint)
			return nil, false
		}
	}
	//token not expired and not blacklisted is ok
	if ok {
		return token, true
	}
	return nil, false
}

func (r *TokenRepositoryBasic) validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (authCode *oauth2.AuthorizationCode, ok bool) {
	code := request.Code
	authCode, ok = r.authorizationCodes[code]
	delete(r.authorizationCodes, code)
	if !ok {
		log.Debug("authorization code not found", zap.String("client id", request.ClientID))
		return nil, false
	}
	ok = authCode.ValidateAccessTokenRequest(request)
	if !ok {
		log.Debug("access token request invalid", zap.String("client id", request.ClientID))
		return nil, false
	}
	return authCode, ok
}
