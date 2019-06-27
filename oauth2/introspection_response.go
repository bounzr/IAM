package oauth2

/**
The server responds with a JSON object [RFC7159] in "application/json" format with the following top-level members.
   	active
      REQUIRED.  Boolean indicator of whether or not the presented token is currently active.  The specifics of a token’s
		"active" State will vary depending on the implementation of the authorization server and the information it keeps
		about its token, but a "true" value return for the "active" property will generally indicate that a given token
		has been issued by this authorization server, has not been revoked by the resource owner, and is within its given
		time window of validity (e.g., after its issuance time and before its expiration time).  See Section 4 for information
		on implementation of such checks.
   	Scope
      OPTIONAL.  A JSON string containing a space-separated list of scopes associated with this token, in the format described
		in Section 3.3 of OAuth 2.0 [RFC6749].
   	client_id
      OPTIONAL.  Client identifier for the OAuth 2.0 client that requested this token.
   	username
      OPTIONAL.  Human-readable identifier for the resource owner who authorized this token.
   	token_type
      OPTIONAL.  Type of the token as defined in Section 5.1 of OAuth 2.0 [RFC6749].
   	exp
      OPTIONAL.  Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token
		will expire, as defined in JWT [RFC7519].
   	iat
      OPTIONAL.  Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token
		was originally issued, as defined in JWT [RFC7519].
   	nbf
      OPTIONAL.  Integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token
		is not to be used before, as defined in JWT [RFC7519].
	sub
	  OPTIONAL.  Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource
		owner who authorized this token.
	aud
   	  OPTIONAL.  Service-specific string identifier or list of string identifiers representing the intended audience for
		this token, as defined in JWT [RFC7519].
	iss
	  OPTIONAL.  String representing the Issuer of this token, as defined in JWT [RFC7519].
	jti
   	  OPTIONAL.  String identifier for the token, as defined in JWT [RFC7519].
*/
/*
The following is a non-normative example response:
HTTP/1.1 200 OK
Content-Type: application/json
{
"active": true,
"client_id": "l238j323ds-23ij4",
"username": "jdoe",
"Scope": "read write dolphin",
"sub": "Z5O3upPC88QrAjx00dis",
"aud": "https://protected.example.net/resource",
"iss": "https://server.example.com/",
"exp": 1419356238,
"iat": 1419350238,
"extension_field": "twenty-seven"
}
*/
type IntrospectionResponse struct {
	Active        string `json:"active,required"`
	Audience      string `json:"aud,omitempty"`
	ClientID      string `json:"client_id,omitempty"`
	Expires       int64  `json:"exp,omitempty"`
	IssuedAt      int64  `json:"iat,omitempty"`
	Issuer        string `json:"iss,omitempty"`
	Jti           string `json:"jti,omitempty"`
	NotBefore     int64  `json:"nbf,omitempty"`
	Scope         string `json:"Scope,omitempty"`
	OwnerID       string `json:"sub,omitempty"`
	TokenAuthType string `json:"token_type,omitempty"`
	Username      string `json:"username,omitempty"`
}

func (ir *IntrospectionResponse) ValidateScope(requestedScope string) (validatedScope string) {
	return getCommonScope(ir.Scope, requestedScope)
}
