package repository

import (
	"../oauth2"
	"go.uber.org/zap"
)

type TokenManagerBasic struct {
	//string is the authorization code
	authorizationCodes map[string]*oauth2.AuthorizationCode

	//string is the token
	accessTokens      map[string]*oauth2.TokenUnit
	refreshTokens     map[string]*oauth2.TokenUnit
	usedTokens        map[string]struct{}
	blackListedTokens map[string]struct{}
}

//init the repository
func (r *TokenManagerBasic) init() {
	r.authorizationCodes = make(map[string]*oauth2.AuthorizationCode)
	r.accessTokens = make(map[string]*oauth2.TokenUnit)
	r.refreshTokens = make(map[string]*oauth2.TokenUnit)
	r.usedTokens = make(map[string]struct{})
	r.blackListedTokens = make(map[string]struct{})
}

func (r *TokenManagerBasic) close() {
	//nothing
}

func (r *TokenManagerBasic) deleteAccessToken(tokenHint *oauth2.AccessTokenHint) {
	r.blackListedTokens[string(tokenHint.Token)] = struct{}{}
	delete(r.accessTokens, string(tokenHint.Token))
	delete(r.refreshTokens, string(tokenHint.Token))
}

func (r *TokenManagerBasic) setAuthorizationCode(code *oauth2.AuthorizationCode) error {
	r.authorizationCodes[code.Code] = code
	return nil
}

func (r *TokenManagerBasic) setTokenUnit(accessToken *oauth2.TokenUnit) error {
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

func (r *TokenManagerBasic) getTokenUnit(hint *oauth2.AccessTokenHint) (*oauth2.TokenUnit, bool) {
	var token *oauth2.TokenUnit
	var ok = false
	if oauth2.NewTokenHintType(hint.Hint) != oauth2.RefreshTokenHintType {
		token, ok = r.accessTokens[string(hint.Token)]
	}
	if !ok && oauth2.NewTokenHintType(hint.Hint) != oauth2.AccessTokenHintType {
		token, ok = r.refreshTokens[string(hint.Token)]
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

func (r *TokenManagerBasic) validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (authCode *oauth2.AuthorizationCode, ok bool) {
	code := request.Code
	authCode, ok = r.authorizationCodes[code]
	delete(r.authorizationCodes, code)
	if !ok {
		log.Debug("authorization code not found", zap.String("client ID", request.ClientID))
		return nil, false
	}
	ok = authCode.ValidateAccessTokenRequest(request)
	if !ok {
		log.Debug("access token request invalid", zap.String("client ID", request.ClientID))
		return nil, false
	}
	return authCode, ok
}
