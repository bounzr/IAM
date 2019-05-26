package oauth2

import "errors"

/**
If the resource owner denies the access request or if the request fails for reasons other than a missing
or invalid redirection URI, the authorization server informs the client by adding the following parameters
to the query component of the redirection URI using the "application/x-www-form-urlencoded" format,
per Appendix B:

error
	REQUIRED.  A single ASCII [USASCII] error code from the following:
        invalid_request
                The request is missing a required parameter, includes an invalid parameter value,
                includes a parameter more than once, or is otherwise malformed.
   	unauthorized_client
                The client is not authorized to request an authorization code using this method.
        access_denied
                The resource owner or authorization server denied the request.
        unsupported_response_type
                The authorization server does not support obtaining an authorization code using this method.
        invalid_scope
                The requested Scope is invalid, unknown, or malformed.
        server_error
                The authorization server encountered an unexpected condition that prevented it from
                fulfilling the request. (This error code is needed because a 500 Internal Server Error HTTP
                status code cannot be returned to the client via an HTTP redirect.)
        temporarily_unavailable
                The authorization server is currently unable to handle the request due to a temporary
                overloading or maintenance of the server.  (This error code is needed because a 503
                Service Unavailable HTTP status code cannot be returned to the client via an HTTP redirect.)
      Values for the "error" parameter MUST NOT include characters outside the set %x20-21 / %x23-5B / %x5D-7E.
error_description
        OPTIONAL.  Human-readable ASCII [USASCII] text providing additional information, used to assist the client developer in
        understanding the error that occurred.
        Values for the "error_description" parameter MUST NOT include characters outside the
        set %x20-21 / %x23-5B / %x5D-7E.
error_uri
        OPTIONAL.  A URI identifying a human-readable web page with information about the error,
        used to provide the client developer with additional information about the error.
        Values for the "error_uri" parameter MUST conform to the URI-reference syntax and thus MUST NOT
        include characters outside the set %x21 / %x23-5B / %x5D-7E.
State
        REQUIRED if a "State" parameter was present in the client authorization request.  The exact value
        received from the client.
   For example, the authorization server redirects the user-agent by sending the following HTTP response:
   HTTP/1.1 302 Found
   Location: https://client.example.com/cb?error=access_denied&state=xyz
**/

var (

	//Oauth2 Dynamic Registration Error Responses\\

	//ErrInvalidClientMetadata error response to an inconsistent registration request
	ErrInvalidClientMetadata = errors.New("invalid_client_metadata")

	//OAuth2 Grant Authorization Code Error Responses\\

	//ErrRedirectionURIInfo returns error to be displayed describing that the Redirection URI is wrong
	ErrRedirectionURIInfo = errors.New("missing, invalid, or mismatching redirection URI")
	//ErrClientIdentifierInfo returns error to be displayed describing that the Client Identifier is wrong
	ErrClientIdentifierInfo = errors.New("client identifier is missing or invalid")

	//Token Request error
	ErrAccessDenied       = errors.New("access_denied")
	ErrAccessDeniedInfo   = errors.New("the resource owner or authorization server denied the request")
	ErrInvalidRequest     = errors.New("invalid_request")
	ErrInvalidRequestInfo = errors.New("the request is missing a required parameter, includes an invalid parameter " +
		"value, includes a parameter more than once, or is otherwise malformed")
	ErrInvalidScope               = errors.New("invalid_scope")
	ErrInvalidScopeInfo           = errors.New("the requested Scope is invalid, unknown, or malformed")
	ErrServerError                = errors.New("server_error")
	ErrServerErrorInfo            = errors.New("the authorization server encountered an unexpected condition that prevented it from fulfilling the request")
	ErrTemporarilyUnavailable     = errors.New("temporarily_unavailable")
	ErrTemporarilyUnavailableInfo = errors.New("the authorization server is currently unable to handle the request " +
		"due to a temporary overloading or maintenance of the server")
	ErrUnauthorizedClient          = errors.New("unauthorized_client")
	ErrUnauthorizedClientInfo      = errors.New("the client is not authorized to request an access token using this method")
	ErrUnsupportedResponseType     = errors.New("unsupported_response_type")
	ErrUnsupportedResponseTypeInfo = errors.New("the authorization server does not support obtaining an access token using this method")
	ErrUnsupportedTokenType        = errors.New("unsupported_token_type")
	ErrUnsupportedTokenTypeInfo    = errors.New("the authorization server does not support the revocation of the presented token type")
)
