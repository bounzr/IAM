package oauth2

/**
Client Registration Request
-----------------------------
The following client metadata fields are defined by this
specification. The implementation and use of all client metadata
fields is OPTIONAL, unless stated otherwise. All data member types
(strings, arrays, numbers) are defined in terms of their JSON
[RFC7159] representations.

redirect_uris
Array of redirection URI strings for use in redirect-based flows
such as the authorization code and implicit flows. As required by
Section 2 of OAuth 2.0 [RFC6749], clients using flows with
redirection MUST register their redirection URI values.

token_endpoint_auth_method
String indicator of the requested authentication method for the
token endpoint. Values defined by this specification are:
* "none": The client is a public client as defined in OAuth 2.0,
and does not have a client secret.
* "client_secret_post": The client uses the HTTP POST parameters
as defined in OAuth 2.0
* "client_secret_basic": The client uses HTTP Basic as defined in
OAuth 2.0

grant_types
Array of OAuth 2.0 grant type strings that the client can use at
the token endpoint. These grant types are defined as follows:
* "authorization_code": The authorization code grant type defined
in OAuth 2.0.
* "implicit": The implicit grant type.
* "password": The resource owner password credentials grant type.
* "client_credentials": The client credentials grant type.
* "refresh_token": The refresh token grant type.
* "urn:ietf:params:oauth:grant-type:jwt-bearer": The JWT Bearer
TokenUnit Grant Type defined in OAuth JWT Bearer TokenUnit Profiles
[RFC7523].
* "urn:ietf:params:oauth:grant-type:saml2-bearer": The SAML 2.0
Bearer Assertion Grant defined in OAuth SAML 2 Bearer TokenUnit
Profiles [RFC7522].
If the token endpoint is used in the grant type, the value of this
parameter MUST be the same as the value of the "grant_type"
parameter passed to the token endpoint defined in the grant type
definition.

response_types
Array of the OAuth 2.0 response type strings that the client can
use at the authorization endpoint. These response types are
defined as follows:
* "code": The authorization code response type
* "token": The implicit response type
If the authorization endpoint is used by the grant type, the value
of this parameter MUST be the same as the value of the
"response_type" parameter passed to the authorization endpoint
defined in the grant type definition. If omitted, the
default is that the client will use only the "code" response type.

client_name
Human-readable string name of the client to be presented to the
end-user during authorization. If omitted, the authorization
server MAY display the raw "client_id" value to the end-user
instead.

client_uri
URL string of a web page providing information about the client.
If present, the server SHOULD display this URL to the end-user in
a clickable fashion.

logo_uri
URL string that references a logo for the client. If present, the
server SHOULD display this image to the end-user during approval.

Scope
String containing a space-separated list of Scope values (as
described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client
can use when requesting access token. If omitted, an authorization
server MAY register a client with a default set of scopes.

contacts
Array of strings representing ways to contact people responsible
for this client, typically email addresses.

tos_uri
URL string that points to a human-readable terms of service
document for the client that describes a contractual relationship
between the end-user and the client that the end-user accepts when
authorizing the client. The authorization server SHOULD display
this URL to the end-user if it is provided.

policy_uri
URL string that points to a human-readable privacy policy document
that describes how the deployment organization collects, uses,
retains, and discloses personal data. The authorization server
SHOULD display this URL to the end-user if it is provided.

jwks_uri
URL string referencing the client’s JSON Web Key (JWK) Set
[RFC7517] document, which contains the client’s public keys. The
value of this field MUST point to a valid JWK Set document. These
keys can be used by higher-level protocols that use signing or
encryption. For instance, these keys might be used by some
applications for validating signed requests made to the token
endpoint when using JWTs for client authentication [RFC7523]. Use
of this parameter is preferred over the "jwks" parameter, as it
allows for easier key rotation. The "jwks_uri" and "jwks"
parameters MUST NOT both be present in the same request or
response.

jwks
Client’s JSON Web Key Set [RFC7517] document value, which contains
the client’s public keys. The value of this field MUST be a JSON
object containing a valid JWK Set. These keys can be used by
higher-level protocols that use signing or encryption. This
parameter is intended to be used by clients that cannot use the
"jwks_uri" parameter, such as native clients that cannot host
public URLs. The "jwks_uri" and "jwks" parameters MUST NOT both
be present in the same request or response.

software_id
A unique identifier string (e.g., a Universally Unique Identifier
(UUID)) assigned by the client developer or software publisher
used by registration endpoints to identify the client software to
be dynamically registered.

software_version
A version identifier string for the client software identified by
"software_id".

+-----------------------------------------------+-------------------+
| grant_types value includes: | response_types |
| | value includes: |
+-----------------------------------------------+-------------------+
| authorization_code | code |
| implicit | token |
| password | (none) |
| client_credentials | (none) |
| refresh_token | (none) |
| urn:ietf:params:oauth:grant-type:jwt-bearer | (none) |
| urn:ietf:params:oauth:grant-type:saml2-bearer | (none) |
+-----------------------------------------------+-------------------+


The following is a non-normative example request not using an initial
access token:
POST /register HTTP/1.1
Content-Type: application/json
Accept: application/json
Host: server.example.com
{
"redirect_uris": [
"https://client.example.org/callback",
"https://client.example.org/callback2"],
"client_name": "My Example Client",
"client_name#ja-Jpan-JP":
"\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u540D",
"token_endpoint_auth_method": "client_secret_basic",
"logo_uri": "https://client.example.org/logo.png",
"jwks_uri": "https://client.example.org/my_public_keys.jwks",
"example_extension_parameter": "example_value"
}

The following is a non-normative example request using an initial
access token and registering a JWK Set by value (with line breaks
within values for display purposes only):
POST /register HTTP/1.1
Content-Type: application/json
Accept: application/json
Authorization: Bearer ey23f2.adfj230.af32-developer321
Host: server.example.com
{
"redirect_uris": ["https://client.example.org/callback",
"https://client.example.org/callback2"],
"client_name": "My Example Client",
"client_name#ja-Jpan-JP":
"\u30AF\u30E9\u30A4\u30A2\u30F3\u30C8\u540D",
"token_endpoint_auth_method": "client_secret_basic",
"policy_uri": "https://client.example.org/policy.html",
"jwks": {"keys": [{
"e": "AQAB",
"n": "nj3YJwsLUFl9BmpAbkOswCNVx17Eh9wMO-_AReZwBqfaWFcfG
HrZXsIV2VMCNVNU8Tpb4obUaSXcRcQ-VMsfQPJm9IzgtRdAY8NN8Xb7PEcYyk
lBjvTtuPbpzIaqyiUepzUXNDFuAOOkrIol3WmflPUUgMKULBN0EUd1fpOD70p
RM0rlp_gg_WNUKoW1V-3keYUJoXH9NztEDm_D2MQXj9eGOJJ8yPgGL8PAZMLe
2R7jb9TxOCPDED7tY_TU4nFPlxptw59A42mldEmViXsKQt60s1SLboazxFKve
qXC_jpLUt22OC6GUG63p-REw-ZOr3r845z50wMuzifQrMI9bQ",
"kty": "RSA"
}]},
"example_extension_parameter": "example_value"
}

*/

type ClientRegistrationRequest struct {
	RedirectUris            []string `json:"redirect_uris"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	ClientName              string   `json:"client_name"`
	ClientUri               string   `json:"client_uri"`
	LogoUri                 string   `json:"logo_uri"`
	Scope                   string   `json:"Scope"`
	Contacts                []string `json:"contacts"`
	TosUri                  string   `json:"tos_uri"`
	PolicyUri               string   `json:"policy_uri"`
	JwksURI                 string   `json:"jwks_uri"`
	//todo jwks object
	Jwks            string `json:"jwks"`
	SoftwareId      string `json:"software_id"`
	SoftwareVersion string `json:"software_version"`
}
