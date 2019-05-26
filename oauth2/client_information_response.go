package oauth2

/**
Client information Response
---------------------------

client_id
REQUIRED. OAuth 2.0 client identifier string. It SHOULD NOT be
currently valid for any other registered client, though an
authorization server MAY issue the same client identifier to
multiple instances of a registered client at its discretion.

client_secret
OPTIONAL. OAuth 2.0 client secret string. If issued, this MUST
be unique for each "client_id" and SHOULD be unique for multiple
instances of a client using the same "client_id". This value is
used by confidential clients to authenticate to the token
endpoint, as described in OAuth 2.0 [RFC6749], Section 2.3.1.

client_id_issued_at
OPTIONAL. Time at which the client identifier was issued. The
time is represented as the number of seconds from
1970-01-01T00:00:00Z as measured in UTC until the date/time of
issuance.

client_secret_expires_at
REQUIRED if "client_secret" is issued. Time at which the client
secret will expire or 0 if it will not expire. The time is
represented as the number of seconds from 1970-01-01T00:00:00Z as
measured in UTC until the date/time of expiration.

Additionally, the authorization server MUST return all registered
metadata about this client, including any fields provisioned by the
authorization server itself. The authorization server MAY reject or
replace any of the client’s requested metadata values submitted
during the registration and substitute them with suitable values.

The following is a non-normative example response of a successful
registration:
HTTP/1.1 201 Created
Content-Type: application/json
Cache-Control: no-store
Pragma: no-cache
{
"client_id": "s6BhdRkqt3",
"client_secret": "cf136dc3c1fc93f31185e5885805d",
"client_id_issued_at": 2893256800,
"client_secret_expires_at": 2893276800,
"redirect_uris": [
"https://client.example.org/callback",
"https://client.example.org/callback2"],
"grant_types": ["authorization_code", "refresh_token"],
"client_name": "My Example Client",
"client_name#ja-Jpan-JP":
"\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u540D",
"token_endpoint_auth_method": "client_secret_basic",
"logo_uri": "https://client.example.org/logo.png",
"jwks_uri": "https://client.example.org/my_public_keys.jwks",
"example_extension_parameter": "example_value"
}
*/

type ClientInformationResponse struct {
	ClientId string `json:"client_id"`
	//todo client_secret is currently always returning "2PHC()~REute/egh"
	ClientSecret          string `json:"client_secret,omitempty"`
	ClientIdIssuedAt      int64  `json:"client_id_issued_at,omitempty"`
	ClientSecretExpiresAt int64  `json:"client_secret_expires_at,omitempty"`

	//authorization server MUST return all registered metadata about this client
	RedirectUris *[]string `json:"redirect_uris"`
	//todo string enum
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	//todo array string enum
	//todo can this string slices be returned as pointers instead of objects?
	GrantTypes []string `json:"grant_types"`
	//todo array string enum
	ResponseTypes []string `json:"response_types"`
	ClientName    string   `json:"client_name,omitempty"`
	ClientUri     string   `json:"client_uri,omitempty"`
	LogoUri       string   `json:"logo_uri,omitempty"`
	Scope         string   `json:"Scope,omitempty"`
	Contacts      []string `json:"contacts,omitempty"`
	TosUri        string   `json:"tos_uri,omitempty"`
	PolicyUri     string   `json:"policy_uri,omitempty"`
	JwksUri       string   `json:"jwks_uri,omitempty"`
	//todo jwks object
	Jwks            string `json:"jwks,omitempty"`
	SoftwareId      string `json:"software_id,omitempty"`
	SoftwareVersion string `json:"software_version,omitempty"`
}

/*
Client Registration Error Response
----------------------------------
the authorization server returns an HTTP 400 status code (unless otherwise specified) with
content type "application/json" consisting of a JSON object [RFC7159]
describing the error in the response body.
Two members are defined for inclusion in the JSON object:

error
REQUIRED. Single ASCII error code string.

error_description
OPTIONAL. Human-readable ASCII text description of the error used
for debugging.

This specification defines the following error codes:

invalid_redirect_uri
The value of one or more redirection URIs is invalid.

invalid_client_metadata
The value of one of the client metadata fields is invalid and the
server has rejected this request. Note that an authorization
server MAY choose to substitute a valid value for any requested
parameter of a client’s metadata.

invalid_software_statement
The software statement presented is invalid.

unapproved_software_statement
The software statement presented is not approved for use by this
authorization server.

The following is a non-normative example of an error response
resulting from a redirection URI that has been blacklisted by the
authorization server (with line breaks within values for display
purposes only):

HTTP/1.1 400 Bad Request
Content-Type: application/json
Cache-Control: no-store
Pragma: no-cache
{
"error": "invalid_redirect_uri",
"error_description": "The redirection URI
http://sketchy.example.com is not allowed by this server."
}

The following is a non-normative example of an error response
resulting from an inconsistent combination of "response_types" and
"grant_types" values (with line breaks within values for display
purposes only):

HTTP/1.1 400 Bad Request
Content-Type: application/json
Cache-Control: no-store
Pragma: no-cache
{
"error": "invalid_client_metadata",
"error_description": "The grant type ’authorization_code’ must be
registered along with the response type ’code’ but found only
*/

type ClientRegistrationError struct {
	Error string `json:"error"`
	//todo This specification defines error codes
	ErrorDescription string `json:"error_description"`
}
