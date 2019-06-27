package oauth2

import (
	"github.com/gofrs/uuid"
	"strings"
	"time"

	"../utils"
)

/*
If the resource owner grants the access request, the authorization server issues an authorization code and delivers it to the client by
adding the following parameters to the query component of the redirection URI using the "application/x-www-form-urlencoded":

code
REQUIRED. The authorization code generated by the authorization server. The authorization code MUST expire
shortly after it is issued to mitigate the risk of leaks. A maximum authorization code lifetime of 10 minutes is
RECOMMENDED. The client MUST NOT use the authorization code more than once. If an authorization code is used more than
once, the authorization server MUST deny the request and SHOULD revoke (when possible) all token previously issued based on
that authorization code. The authorization code is bound to the client identifier and redirection URI.

State
REQUIRED if the "State" parameter was present in the client authorization request. The exact value received from the
client.

For example, the authorization server redirects the user-agent by sending the following HTTP response:
HTTP/1.1 302 Found Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz
The client MUST ignore unrecognized response parameters. The authorization code string size is left undefined by this
specification. The client should avoid making assumptions about code value sizes. The authorization server SHOULD document the size of
any value it issues.

Error Response
If the request fails due to a missing, invalid, or mismatching redirection URI, or if the client identifier is missing or invalid,
the authorization server SHOULD inform the resource owner of the error and MUST NOT automatically redirect the user-agent to the
invalid redirection URI.
If the resource owner denies the access request or if the request fails for reasons other than a missing or invalid redirection URI,
the authorization server informs the client by adding the following parameters to the query component of the redirection URI using the
"application/x-www-form-urlencoded" format, per Appendix B:

error
REQUIRED. A single ASCII [USASCII] error code from the following:
--invalid_request
The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than
once, or is otherwise malformed.
--unauthorized_client
The client is not authorized to request an authorization code using this method.
--access_denied
The resource owner or authorization server denied the request.
--unsupported_response_type
The authorization server does not support obtaining an authorization code using this method.
--invalid_scope
The requested Scope is invalid, unknown, or malformed.
--server_error
The authorization server encountered an unexpected condition that prevented it from fulfilling the request.
(This error code is needed because a 500 Internal Server Error HTTP status code cannot be returned to the client
via an HTTP redirect.)
--temporarily_unavailable
The authorization server is currently unable to handle the request due to a temporary overloading or maintenance
of the server. (This error code is needed because a 503 Service Unavailable HTTP status code cannot be returned
to the client via an HTTP redirect.)
Values for the "error" parameter MUST NOT include characters outside the set %x20-21 / %x23-5B / %x5D-7E.

error_description
OPTIONAL. Human-readable ASCII [USASCII] text providing additional information, used to assist the client developer in
understanding the error that occurred.
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
Location: https://client.example.com/cb?error=access_denied&state=xyz
*/
type AuthorizationCode struct {
	Code           string
	ClientID       uuid.UUID
	ExpirationTime time.Time //Expires in 10 mins after returned
	RedirectionURI string
	Scope          []byte
	State          string
	OwnerID        uuid.UUID //resources owner
}

type AuthorizationCodeResponse struct {
	Code      string
	State     string
	ClientURI string
}

/*
type AuthorizationError struct{
	Error	error
	State   string
}
*/

//todo in case of wrong value throw an AuthorizationError
//NewAuthorizationCode returns an new authorization code and its authorization code response
func NewAuthorizationCode(ownerID uuid.UUID, authReq *AuthorizationRequest) *AuthorizationCode {

	code := utils.GetRandomString(22)
	expiration := time.Now().Add(time.Minute * 10)
	uuid := uuid.FromStringOrNil(authReq.ClientID)

	ac := &AuthorizationCode{
		Code:           code,
		ClientID:       uuid,
		ExpirationTime: expiration,
		OwnerID:        ownerID,
		RedirectionURI: authReq.RedirectURI,
		//Scopes:         authReq.GetScopesMap(),
		Scope: []byte(authReq.Scope),
		State: authReq.State,
	}
	return ac
}

func (ac *AuthorizationCode) GetAuthorizationCodeResponse() *AuthorizationCodeResponse {
	ar := &AuthorizationCodeResponse{
		Code:      ac.Code,
		State:     ac.State,
		ClientURI: ac.RedirectionURI,
	}
	return ar
}

func (ac *AuthorizationCode) ValidateAccessTokenRequest(accTokenReq *AuthorizationCodeAccessTokenRequest) (ok bool) {
	if strings.Compare(accTokenReq.Code, ac.Code) != 0 {
		return false
	}
	if strings.Compare(accTokenReq.ClientID, ac.ClientID.String()) != 0 {
		return false
	}
	if time.Now().After(ac.ExpirationTime) {
		return false
	}
	if strings.Compare(accTokenReq.RedirectURI, ac.RedirectionURI) != 0 {
		return false
	}
	return true
}
