package oauth2

import (
	"bounzr/iam/token"
	"bounzr/iam/utils"
	"github.com/gofrs/uuid"
	"time"
)

/**
//code grant access token response
The authorization server issues an access token and optional refresh token, and constructs the response by adding the
following parameters to the entity-body of the HTTP response with a 200 (OK) status code:

access_token
REQUIRED. The access token issued by the authorization server.

token_type
REQUIRED. The type of the token issued ex. Bearer

expires_in
RECOMMENDED. The lifetime in seconds of the access token. For example, the value "3600" denotes that the access token will
expire in one hour from the time the response was generated. If omitted, the authorization server SHOULD provide the
expiration time via other means or document the default value.

refresh_token
OPTIONAL. The refresh token, which can be used to obtain new access token using the same authorization grant

Scope
OPTIONAL, if identical to the Scope requested by the client; otherwise, REQUIRED.

The parameters are included in the entity-body of the HTTP response using the "application/json" media type as defined by [RFC4627]. The
parameters are serialized into a JavaScript Object Notation (JSON) structure by adding each parameter at the highest structure level.
Parameter names and string values are included as JSON strings.
Numerical values are included as JSON numbers. The order of parameters does not matter and can vary.
The authorization server MUST include the HTTP "Cache-Control" response header field [RFC2616] with a value of "no-store" in any
response containing token, credentials, or other sensitive information, as well as the "Pragma" response header field [RFC2616]
with a value of "no-cache".
For example:
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache
{
"access_token":"2YotnFZFEjr1zCsicMWpAA",
"token_type":"example",
"expires_in":3600,
"refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
"example_parameter":"example_value"
}
*/

//implicit grant access token response
/**
If the resource owner grants the access request, the authorization server issues an access token and delivers it to the
client by adding the following parameters to the fragment component of the redirection URI using the
"application/x-www-form-urlencoded" format, per Appendix B:

	access_token
	REQUIRED. The access token issued by the authorization server.

	token_type
	REQUIRED. The type of the token issued as described in Section 7.1. Value is case insensitive.

	expires_in
	RECOMMENDED. The lifetime in seconds of the access token. For example, the value "3600" denotes that the access
		token will expire in one hour from the time the response was generated.
		If omitted, the authorization server SHOULD provide the expiration time via other means or document the default
		value.

	Scope
	OPTIONAL, if identical to the Scope requested by the client; otherwise, REQUIRED. The Scope of the access token as
	described by Section 3.3.

	State
	REQUIRED if the "State" parameter was present in the client authorization request. The exact value received from the
	client.

The authorization server MUST NOT issue a refresh token.
For example, the authorization server redirects the user-agent by sending the following HTTP response (with extra line
breaks for	display purposes only):
HTTP/1.1 302 Found
Location: http://example.com/cb#access_token=2YotnFZFEjr1zCsicMWpAA&state=xyz&token_type=example&expires_in=3600
*/

type TokenUnit struct {
	Active         bool
	ClientID       uuid.UUID //client_id of the Relying Party as an client value
	ExpirationTime time.Time //Expiration time on or after which the ID TokenUnit MUST NOT be accepted for processing
	IssuedAt       time.Time
	Issuer         string    //identifies the principal that issued the TokenUnit
	NotBefore      time.Time //the time before which the TokenUnit MUST NOT be accepted for processing
	Scope          []byte
	State          string
	OwnerID        uuid.UUID //user_id A locally unique and never reassigned identifier within the Issuer for the End-User
	ParentToken    []byte
	Token          []byte        //access_token or refresh_token
	TokenAuthType  TokenAuthType //bearer, mac
	TokenHintType  TokenHintType //access_token, refresh_token
}

type AccessTokenOptions struct {
	AddRefreshToken bool      //include refresh TokenUnit
	ClientID        uuid.UUID //client_id
	Issuer          string    //server host
	Scope           []byte
	State           string    //client State
	OwnerID         uuid.UUID //user_id
}

type AccessTokenHint struct {
	Token string `json:"token"`
	Hint  string `json:"token_type_hint"`
}

type AccessTokenResponse struct {
	AccessToken   string `json:"access_token"`
	TokenAuthType string `json:"token_type"`
	ExpiresIn     int64  `json:"expires_in"`
	RefreshToken  string `json:"refresh_token,omitempty"`
	Scope         string `json:"Scope,omitempty"`
	State         string `json:"State,omitempty"`
}

/**
If the request fails due to a missing, invalid, or mismatching redirection URI, or if the client identifier is missing
or invalid, the authorization server SHOULD inform the resource owner of the error and MUST NOT automatically redirect
the user-agent to the invalid redirection URI.
If the resource owner denies the access request or if the request fails for reasons other than a missing or invalid
redirection URI, the authorization server informs the client by adding the following parameters to the fragment
component of the redirection URI using the "application/x-www-form-urlencoded" format, per Appendix B:

	error
	REQUIRED. A single ASCII [USASCII] error code from the following:
		invalid_request
		The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than
		once, or is otherwise malformed.
		unauthorized_client
		The client is not authorized to request an access token using this method.
		access_denied
		The resource owner or authorization server denied the request.
		unsupported_response_type
		The authorization server does not support obtaining an access token using this method.
		invalid_scope
		The requested Scope is invalid, unknown, or malformed.
		server_error
		The authorization server encountered an unexpected condition that prevented it from fulfilling the request.
		(This error code is needed because a 500 Internal Server Error HTTP status code cannot be returned to the client
		via an HTTP redirect.)
		temporarily_unavailable
		The authorization server is currently unable to handle the request due to a temporary overloading or maintenance
		of the server. (This error code is needed because a 503 Service Unavailable HTTP status code cannot be returned
		to the client via an HTTP redirect.)
		Values for the "error" parameter MUST NOT include characters outside the set %x20-21 / %x23-5B / %x5D-7E.

	error_description
	OPTIONAL. Human-readable ASCII [USASCII] text providing additional information, used to assist the client developer
	in understanding the error that occurred.
	Values for the "error_description" parameter MUST NOT include characters outside the set %x20-21 / %x23-5B / %x5D-7E.

	error_uri
	OPTIONAL. A URI identifying a human-readable web page with information about the error, used to provide the client
	developer with additional information about the error.
	Values for the "error_uri" parameter MUST conform to the URI-reference syntax and thus MUST NOT include characters
	outside the set %x21 / %x23-5B / %x5D-7E.

	State
	REQUIRED if a "State" parameter was present in the client authorization request. The exact value received from the
	client.

For example, the authorization server redirects the user-agent by sending the following HTTP response:
HTTP/1.1 302 Found
Location: https://client.example.com/cb#error=access_denied&state=xyz
*/

type AccessTokenErrorResponse struct {
	Error            string `json:error`
	ErrorDescription string `json:error_description`
	ErrorURI         string `json:error_uri`
	State            string `json:State`
}

func NewTokenSet(opt *AccessTokenOptions, accessDuration time.Duration, refreshDuration time.Duration) (accessToken *TokenUnit, refreshToken *TokenUnit) {
	creationTime := time.Now()
	accessToken = &TokenUnit{
		Active:         true,
		ClientID:       opt.ClientID,
		ExpirationTime: creationTime.Add(accessDuration),
		IssuedAt:       creationTime,
		Issuer:         opt.Issuer,
		NotBefore:      creationTime,
		Scope:          []byte(opt.Scope),
		State:          opt.State,
		OwnerID:        opt.OwnerID,
		ParentToken:    nil,
		Token:          token.GetToken(),
		TokenAuthType:  NewTokenAuthType("Bearer"),
		TokenHintType:  NewTokenHintType("access_token"),
	}
	if !opt.AddRefreshToken {
		return accessToken, nil
	}

	refreshToken = &TokenUnit{
		Active:         true,
		ClientID:       opt.ClientID,
		ExpirationTime: creationTime.Add(refreshDuration),
		IssuedAt:       creationTime,
		Issuer:         opt.Issuer,
		NotBefore:      creationTime,
		Scope:          []byte(opt.Scope),
		State:          opt.State,
		OwnerID:        opt.OwnerID,
		ParentToken:    accessToken.Token,
		Token:          token.GetToken(),
		TokenAuthType:  NewTokenAuthType("Bearer"),
		TokenHintType:  NewTokenHintType("refresh_token"),
	}

	return accessToken, refreshToken
}

func NewRefreshToken(accessToken *TokenUnit, duration time.Duration) (refreshToken *TokenUnit) {
	creationTime := time.Now()
	refreshToken = &TokenUnit{
		Active:         true,
		ClientID:       accessToken.ClientID,
		ExpirationTime: creationTime.Add(duration),
		IssuedAt:       creationTime,
		Issuer:         accessToken.Issuer,
		NotBefore:      creationTime,
		Scope:          []byte(accessToken.Scope),
		State:          accessToken.State,
		OwnerID:        accessToken.OwnerID,
		ParentToken:    accessToken.Token,
		Token:          token.GetToken(),
		TokenAuthType:  NewTokenAuthType("Bearer"),
		TokenHintType:  NewTokenHintType("refresh_token"),
	}
	return
}

func GetAccessTokenResponse(accessToken, refreshToken *TokenUnit) (response *AccessTokenResponse) {
	expiresIn := accessToken.GetExpirationTime() - time.Now().Unix()
	if expiresIn < 0 {
		expiresIn = 0
	}
	if refreshToken != nil {
		response = &AccessTokenResponse{
			AccessToken:   string(accessToken.Token),
			ExpiresIn:     expiresIn,
			RefreshToken:  string(refreshToken.Token),
			Scope:         accessToken.GetScope(),
			State:         accessToken.State,
			TokenAuthType: accessToken.TokenAuthType.String(),
		}
	} else {
		response = &AccessTokenResponse{
			AccessToken:   string(accessToken.Token),
			ExpiresIn:     expiresIn,
			Scope:         accessToken.GetScope(),
			State:         accessToken.State,
			TokenAuthType: accessToken.TokenAuthType.String(),
		}
	}
	return
}

func (t *TokenUnit) GetClient() (clientID uuid.UUID) {
	return t.ClientID
}

func (t *TokenUnit) GetExpirationTime() int64 {
	return t.ExpirationTime.Unix()
}

func (t *TokenUnit) GetIntrospectionResponse() (response *IntrospectionResponse) {
	response = &IntrospectionResponse{
		Active: false,
	}
	if !t.Active {
		return
	}
	if !utils.InTimeSpan(t.NotBefore, t.ExpirationTime, time.Now()) {
		t.Active = false
		return
	}
	response = &IntrospectionResponse{
		Active:        t.Active,
		ClientID:      t.GetClient().String(),
		Expires:       t.GetExpirationTime(),
		IssuedAt:      t.GetIssuedAt(),
		Issuer:        t.GetIssuer(),
		NotBefore:     t.GetNotBefore(),
		Scope:         t.GetScope(),
		OwnerID:       t.GetResourceOwner().String(),
		TokenAuthType: t.GetTokenAuthType(),
	}
	return
}

func (t *TokenUnit) GetIssuedAt() int64 {
	return t.IssuedAt.Unix()
}

func (t *TokenUnit) GetIssuer() (issuer string) {
	return t.Issuer
}

func (t *TokenUnit) GetNotBefore() int64 {
	return t.NotBefore.Unix()
}

func (t *TokenUnit) GetScope() (scope string) {
	return string(t.Scope)
}

func (t *TokenUnit) GetResourceOwner() (owner uuid.UUID) {
	return t.OwnerID
}

func (t *TokenUnit) GetToken() (token []byte) {
	token = t.Token
	return
}

func (t *TokenUnit) GetTokenAuthType() (tokenType string) {
	return t.TokenAuthType.String()
}

func (t *TokenUnit) GetTokenHint() *AccessTokenHint {
	hint := &AccessTokenHint{
		Token: string(t.Token),
		Hint:  AccessTokenHintType.String(),
	}
	return hint
}

func (t *TokenUnit) ValidateScope(requestedScope string) (validatedScope string) {
	return getCommonScope(string(t.Scope), requestedScope)
}
