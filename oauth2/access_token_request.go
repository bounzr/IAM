package oauth2

/**
Code Grant Access TokenUnit Request

The client makes a request to the token endpoint by sending the following parameters using the "application/x-www-form-urlencoded"
format per Appendix B with a character encoding of UTF-8 in the HTTP request entity-body:

grant_type: REQUIRED. Value MUST be set to "authorization_code".

code: REQUIRED. The authorization code received from the authorization server.

redirect_uri: REQUIRED, if the "redirect_uri" parameter was included in the authorization request as described in
Section 4.1.1, and their values MUST be identical.

client_id: REQUIRED, if the client is not authenticating with the authorization server as described in Section 3.2.1.

For example, the client makes the following HTTP request using TLS (with extra line breaks for display purposes only):
POST /token HTTP/1.1
Host: server.example.com
Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
Content-Type: application/x-www-form-urlencoded
grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb
*/
type AuthorizationCodeAccessTokenRequest struct {
	ClientID    string `schema:"client_id,required"`
	Code        string `schema:"code,required"`
	GrantType   string `schema:"grant_type,required"`
	RedirectURI string `schema:"redirect_uri,required"`
}

/**
Client Credentials Grant Access TokenUnit Request

The client makes a request to the token endpoint by adding the following parameters using the
"application/x-www-form-urlencoded" format per Appendix B with a character encoding of UTF-8 in the HTTP request
entity-body:

   grant_type
         REQUIRED.  Value MUST be set to "client_credentials".
   Scope
         OPTIONAL.  The Scope of the access request as described by Section 3.3.

The client MUST authenticate with the authorization server as described in Section 3.2.1.
For example, the client makes the following HTTP request using transport-layer security (with extra line breaks for
display purposes only):
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded
     grant_type=client_credentials
The authorization server MUST authenticate the client.
*/
type ClientCredentialsAccessTokenRequest struct {
	GrantType string `schema:"grant_type,required"`
	Scope     string `schema:"Scope"`
}

/**
Owner Password Grant Access TokenUnit Request

The client makes a request to the token endpoint by adding the following parameters using the
"application/x-www-form-urlencoded" format per Appendix B with a character encoding of UTF-8 in the HTTP request
entity-body:

   grant_type
         REQUIRED.  Value MUST be set to "password".
   username
         REQUIRED.  The resource owner username.
   password
         REQUIRED.  The resource owner password.
   Scope
         OPTIONAL.  The Scope of the access request as described by Section 3.3.

If the client type is confidential or the client was issued client credentials (or assigned other authentication
requirements), the client MUST authenticate with the authorization server as described in Section 3.2.1.
For example, the client makes the following HTTP request using transport-layer security (with extra line breaks for
display purposes only):
     POST /token HTTP/1.1
     Host: server.example.com
     Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
     Content-Type: application/x-www-form-urlencoded
     grant_type=password&username=johndoe&password=A3ddj3w

*/
type OwnerPasswordAccessTokenRequest struct {
	GrantType string `schema:"grant_type,required"`
	Username  string `schema:"username,required"`
	Password  string `schema:"password,required"`
	Scope     string `schema:"Scope"`
}

/*
If the authorization server issued a refresh token to the client, the client makes a refresh request to the token endpoint
by adding the following parameters using the "application/x-www-form-urlencoded" format per Appendix B with a character
encoding of UTF-8 in the HTTP request entity-body:

	grant_type
		REQUIRED.  Value MUST be set to "refresh_token".
	refresh_token
		REQUIRED.  The refresh token issued to the client.
	Scope
		OPTIONAL.  The Scope of the access request as described by Section 3.3.  The requested Scope MUST NOT include any
		Scope not originally granted by the resource owner, and if omitted is treated as equal to the Scope originally
		granted by the resource owner.

Because refresh token are typically long-lasting credentials used to request additional access token, the refresh token
is bound to the client to which it was issued.  If the client type is confidential or the client was issued client credentials
(or assigned other authentication requirements), the client MUST authenticate with the authorization server as described
in Section 3.2.1.

For example, the client makes the following HTTP request using transport-layer security (with extra line breaks for display purposes
only):
	POST /token HTTP/1.1
	Host: server.example.com
	Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
	Content-Type: application/x-www-form-urlencoded
	grant_type=refresh_token&refresh_token=tGzv3JOkF0XG5Qx2TlKWIA

The authorization server MUST:
   o  require client authentication for confidential clients or for any client that was issued client credentials (or
		with other authentication requirements),
   o  authenticate the client if client authentication is included and ensure that the refresh token was issued to the
		authenticated client, and
   o  validate the refresh token.
If valid and authorized, the authorization server issues an access token as described in Section 5.1.  If the request
failed verification or is invalid, the authorization server returns an error response as described in Section 5.2.
The authorization server MAY issue a new refresh token, in which case the client MUST discard the old refresh token and
replace it with the new refresh token.  The authorization server MAY revoke the old refresh token after issuing a new
refresh token to the client.  If a new refresh token is issued, the refresh token Scope MUST be identical to that of the
refresh token included by the client in the request.
*/
type RefreshAccessTokenRequest struct {
	GrantType    string `schema:"grant_type,required"`
	RefreshToken string `schema:"refresh_token,required"`
	Scope        string `schema:"Scope"`
}

func (atr *AuthorizationCodeAccessTokenRequest) GetGrantType() GrantType {
	gt, _ := NewGrantType(atr.GrantType)
	return gt
}

func (atr *ClientCredentialsAccessTokenRequest) GetGrantType() GrantType {
	gt, _ := NewGrantType(atr.GrantType)
	return gt
}

func (atr *OwnerPasswordAccessTokenRequest) GetGrantType() GrantType {
	gt, _ := NewGrantType(atr.GrantType)
	return gt
}

func (atr *RefreshAccessTokenRequest) GetGrantType() GrantType {
	gt, _ := NewGrantType(atr.GrantType)
	return gt
}

func (atr *RefreshAccessTokenRequest) GetAccessTokenHint() (hint *AccessTokenHint) {
	hint = &AccessTokenHint{
		Token: atr.RefreshToken,
		Hint:  RefreshTokenHintType.String(),
	}
	return
}
