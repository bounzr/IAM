package oauth2

import (
	"../tokens"
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
OPTIONAL. The refresh token, which can be used to obtain new access tokens using the same authorization grant

Scope
OPTIONAL, if identical to the Scope requested by the client; otherwise, REQUIRED.

The parameters are included in the entity-body of the HTTP response using the "application/json" media type as defined by [RFC4627]. The
parameters are serialized into a JavaScript Object Notation (JSON) structure by adding each parameter at the highest structure level.
Parameter names and string values are included as JSON strings.
Numerical values are included as JSON numbers. The order of parameters does not matter and can vary.
The authorization server MUST include the HTTP "Cache-Control" response header field [RFC2616] with a value of "no-store" in any
response containing tokens, credentials, or other sensitive information, as well as the "Pragma" response header field [RFC2616]
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

type AccessToken struct {
	AccessToken    []byte    //access_token
	ClientID       uuid.UUID //client_id of the Relying Party as an client value
	ExpirationTime time.Time //Expiration time on or after which the ID Token MUST NOT be accepted for processing
	IssuedAt       time.Time
	Issuer         string    //identifies the principal that issued the AccessToken
	NotBefore      time.Time //the time before which the AccessToken MUST NOT be accepted for processing
	RefreshToken   []byte
	Scope          []byte
	State          string
	OwnerID        uuid.UUID     //user_id A locally unique and never reassigned identifier within the Issuer for the End-User
	TokenAuthType  TokenAuthType //bearer, mac
	TokenHintType  TokenHintType //access_token, refresh_token
}

type AccessTokenOptions struct {
	AddRefreshToken bool      //include refresh AccessToken
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

func NewAccessToken(opt *AccessTokenOptions) (accessToken *AccessToken, refreshToken *AccessToken) {
	createTime := time.Now()
	expireTimeAccess := createTime.Add(time.Minute * 60) //todo config an expiration time from properties

	accToken := tokens.GetToken() //todo AccessToken generator

	accessToken = &AccessToken{
		AccessToken:    accToken,
		ClientID:       opt.ClientID,
		ExpirationTime: expireTimeAccess,
		IssuedAt:       createTime,
		Issuer:         opt.Issuer,
		NotBefore:      createTime,
		Scope:          []byte(opt.Scope),
		State:          opt.State,
		OwnerID:        opt.OwnerID,
		RefreshToken:   nil,
		TokenAuthType:  NewTokenAuthType("Bearer"),
		TokenHintType:  NewTokenHintType("access_token"),
	}

	if opt.AddRefreshToken {
		refreshToken = NewRefreshToken(accessToken)
		accessToken.RefreshToken = refreshToken.RefreshToken
	}

	return
}

func NewRefreshToken(accessToken *AccessToken) (refreshToken *AccessToken) {
	createTime := time.Now()
	expireTimeRefresh := accessToken.ExpirationTime.Add(time.Hour * 24) // todo config an expiration time from properties
	token := tokens.GetToken()
	refreshToken = &AccessToken{
		AccessToken:    accessToken.AccessToken,
		ClientID:       accessToken.ClientID,
		ExpirationTime: expireTimeRefresh,
		IssuedAt:       createTime,
		Issuer:         accessToken.Issuer,
		NotBefore:      accessToken.IssuedAt,
		Scope:          accessToken.Scope,
		State:          "",
		OwnerID:        accessToken.OwnerID,
		RefreshToken:   token,
		TokenAuthType:  NewTokenAuthType("Bearer"),
		TokenHintType:  NewTokenHintType("refresh_token"),
	}
	accessToken.RefreshToken = token
	return refreshToken
}

func (t *AccessToken) GetAccessTokenResponse() (response *AccessTokenResponse) {
	response = &AccessTokenResponse{}
	response.AccessToken = string(t.AccessToken)
	response.ExpiresIn = t.GetExpirationTime() - time.Now().Unix()
	response.RefreshToken = string(t.RefreshToken)
	response.Scope = t.GetScope()
	response.State = t.State
	response.TokenAuthType = t.TokenAuthType.String()
	return
}

func (t *AccessToken) GetClient() (clientID uuid.UUID) {
	return t.ClientID
}

func (t *AccessToken) GetExpirationTime() int64 {
	return t.ExpirationTime.Unix()
}

func (t *AccessToken) GetIntrospectionResponse() (response *IntrospectionResponse) {
	ok := t.GetExpirationTime() > time.Now().Unix()
	//if not ok then return false
	if !ok {
		response = &IntrospectionResponse{
			Active: "false",
		}
		return
	}
	response = &IntrospectionResponse{
		Active:        "true",
		ClientID:      t.GetClient().String(),
		Expires:       t.GetExpirationTime(),
		IssuedAt:      t.GetIssuedAt(),
		Issuer:        t.GetIssuer(),
		NotBefore:     t.GetNotBefore(),
		Scope:         t.GetScope(),
		OwnerID:       t.GetResourceOwner().String(),
		TokenAuthType: t.GetTokenAuthType(),
		Username:      t.GetResourceOwner().String(),
	}
	return
}

func (t *AccessToken) GetIssuedAt() int64 {
	return t.IssuedAt.Unix()
}

func (t *AccessToken) GetIssuer() (issuer string) {
	return t.Issuer
}

func (t *AccessToken) GetNotBefore() int64 {
	return t.NotBefore.Unix()
}

func (t *AccessToken) GetScope() (scope string) {
	return string(t.Scope)
}

func (t *AccessToken) GetResourceOwner() (owner uuid.UUID) {
	return t.OwnerID
}

func (t *AccessToken) GetToken() (token []byte) {
	token = nil
	if t.TokenHintType == RefreshTokenHintType {
		token = t.RefreshToken
	} else {
		token = t.AccessToken
	}
	return
}

func (t *AccessToken) GetTokenAuthType() (tokenType string) {
	return t.TokenAuthType.String()
}

func (t *AccessToken) GetTokenHints() (clientID uuid.UUID, accessToken *AccessTokenHint, refreshToken *AccessTokenHint) {
	clientID = t.ClientID
	if len(t.AccessToken) > 0 {
		accessToken = &AccessTokenHint{
			Token: string(t.AccessToken),
			Hint:  AccessTokenHintType.String(),
		}
	}
	if len(t.RefreshToken) > 0 {
		refreshToken = &AccessTokenHint{
			Token: string(t.RefreshToken),
			Hint:  RefreshTokenHintType.String(),
		}
	}
	return
}

func (t *AccessToken) ValidateScope(requestedScope string) (validatedScope string) {
	return getCommonScope(string(t.Scope), requestedScope)
}
