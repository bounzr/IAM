package oauth2

import "strings"

//Code Grant AuthorizationRequest The client constructs the request URI by adding parameters
/**
The client constructs the request URI by adding the following parameters to the query component of the
authorization endpoint URI using the "application/x-www-form-urlencoded" format, per Appendix B:

      response_type
            REQUIRED.  Value MUST be set to "code".
      client_id
            REQUIRED.  The client identifier as described in Section 2.2.
      redirect_uri
            OPTIONAL.  As described in Section 3.1.2. - REQUIRED after security review
      Scope
            OPTIONAL.  The Scope of the access request as described by Section 3.3.
      State
         RECOMMENDED.  An opaque value used by the client to maintain State between the request and callback.
         The authorization server includes this value when redirecting the user-agent back to the client.
         The parameter SHOULD be used for preventing cross-site request forgery as described in Section 10.12.

The client directs the resource owner to the constructed URI using an HTTP redirection response, or by other
means available to it via the user-agent. For example, the client directs the user-agent to make the following
HTTP request using TLS (with extra line breaks for display purposes only):
    GET /authorize?response_type=code&client_id=s6BhdRkqt3&State=xyz
        &redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
	Host: server.example.com

	The authorization server validates the request to ensure that all
   required parameters are present and valid.  If the request is valid,
   the authorization server authenticates the resource owner and obtains
   an authorization decision (by asking the resource owner or by
   establishing approval via other means).
**/
//Implicit Grant AuthorizationRequest
/**
The client constructs the request URI by adding the following parameters to the query component of the authorization
endpoint URI using the "application/x-www-form-urlencoded" format, per Appendix B:
	response_type
		REQUIRED. Value MUST be set to "token".
	client_id
		REQUIRED. The client identifier as described in Section 2.2.
	redirect_uri
		OPTIONAL. As described in Section 3.1.2.
	Scope
		OPTIONAL. The Scope of the access request as described by Section 3.3.
	State
		RECOMMENDED. An opaque value used by the client to maintain State between the request and callback. The
		authorization server includes this value when redirecting the user-agent back to the client. The parameter
		SHOULD be used for preventing cross-site request forgery as described in Section 10.12.
The client directs the resource owner to the constructed URI using an HTTP redirection response, or by other means
available to it via the user-agent.
For example, the client directs the user-agent to make the following HTTP request using TLS (with extra line breaks for
display purposes only):
GET /authorize?response_type=token&client_id=s6BhdRkqt3&State=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
Host: server.example.com
*/

//AuthorizationRequest used for code or implicit grants
type AuthorizationRequest struct {
	ResponseType string `schema:"response_type,required"` //token or code
	ClientID     string `schema:"client_id,required"`
	RedirectURI  string `schema:"redirect_uri,required"`
	Scope        string `schema:"Scope"`
	State        string `schema:"State"`
}

//GetScopesList returns a slice of the AuthorizationRequest scopes
func (ar *AuthorizationRequest) GetScopesList() []string {
	if len(strings.TrimSpace(ar.Scope)) > 0 {
		return strings.Split(strings.TrimSpace(ar.Scope), " ")
	}
	return nil
}

//GetScopesMap returns a map of the AuthorizationRequest scopes
func (ar *AuthorizationRequest) GetScopesMap() map[string]struct{} {
	scopesMap := make(map[string]struct{})
	scopesList := ar.GetScopesList()
	for _, value := range scopesList {
		scopesMap[value] = struct{}{}
	}
	return scopesMap
}

//MatchScopes assign to Scope only the elements that are common to the approved Scopes and the previous authorization request scopes
func (ar *AuthorizationRequest) MatchScopes(approvedScopes []string) {
	availScopes := ar.GetScopesMap()
	var sb strings.Builder
	//todo replace the trim for the if index value == 0 {sb.WriteString(scp)} else {sb.WriteString(" "  + scp)}
	for _, scp := range approvedScopes {
		if _, ok := availScopes[scp]; ok {
			sb.WriteString(scp + " ")
		}
	}
	ar.Scope = strings.Trim(sb.String(), " ")
	sb.Reset()
}
