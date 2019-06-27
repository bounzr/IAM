package oauth2

import (
	"github.com/gofrs/uuid"
	"strings"
	"time"

	"../utils"
)

type Client struct {
	Contacts                []string
	GrantTypes              map[GrantType]struct{}
	Groups                  []string //who can admin the client
	ID                      uuid.UUID
	IDIssuedAt              time.Time
	Jwks                    string //todo jwks object
	JwksURI                 string
	LogoURI                 string
	Name                    string
	OwnerID                 uuid.UUID
	PolicyURI               string
	RedirectUris            map[string]struct{}
	ResponseTypes           map[ResponseType]struct{}
	Scope                   string
	Secret                  string
	SecretExpiresAt         time.Time
	SoftwareID              string
	SoftwareVersion         string
	TokenEndpointAuthMethod TokenEndpointAuthMethod
	TosURI                  string
	URI                     string
}

type ClientCtx struct {
	ID      uuid.UUID
	Name    string
	logoURI string
}

type Clients struct {
	RegisteredClients []*ClientSummary
}

type ClientSummary struct {
	ID     string
	Name   string
	Groups []string
	Owner  string
}

//NewOauth2Client creates a new oauth2 client based on the client registration request. Returns pointer to client or error
func NewOauth2Client(request *ClientRegistrationRequest) (*Client, error) {
	clientID, err := uuid.NewV4()
	issuedAt := time.Now().UTC()

	cli := &Client{
		Contacts:                request.Contacts,
		GrantTypes:              getGrantTypes(request.GrantTypes),
		Groups:                  []string{"admin"},
		ID:                      clientID,
		IDIssuedAt:              issuedAt,
		Jwks:                    request.Jwks,    //todo jwks is an object and not a string. develop object
		JwksURI:                 request.JwksURI, //todo if jwksuri is provided get jwks from uri
		LogoURI:                 request.LogoUri,
		Name:                    request.ClientName,
		PolicyURI:               request.PolicyUri,
		RedirectUris:            getRedirectURIs(request.RedirectUris),
		ResponseTypes:           getResponseTypes(request.ResponseTypes),
		Scope:                   strings.TrimSpace(request.Scope),
		Secret:                  utils.GetRandomPassword(16),
		SecretExpiresAt:         issuedAt.AddDate(1, 0, 0),
		SoftwareID:              request.SoftwareId,
		SoftwareVersion:         request.SoftwareVersion,
		TokenEndpointAuthMethod: getTokenEndpointAuthMethod(request.TokenEndpointAuthMethod),
		TosURI:                  request.TosUri,
		URI:                     request.ClientUri,
	}

	return cli, err
}

func NewClientCtx(id uuid.UUID, name string, logoURI string) *ClientCtx {
	return &ClientCtx{
		ID:      id,
		Name:    name,
		logoURI: logoURI,
	}
}

func getGrantTypes(types []string) map[GrantType]struct{} {
	gts := make(map[GrantType]struct{})
	for i := 0; i < len(types); i++ {
		var t GrantType
		t, err := NewGrantType(types[i])
		if err == nil {
			gts[t] = struct{}{}
		}
	}
	return gts
}

func getRedirectURIs(uris []string) map[string]struct{} {
	rus := make(map[string]struct{})
	for _, uri := range uris {
		rus[uri] = struct{}{}
	}
	return rus
}

func getResponseTypes(types []string) map[ResponseType]struct{} {
	rts := make(map[ResponseType]struct{})
	for i := 0; i < len(types); i++ {
		var t ResponseType
		t, err := NewResponseType(types[i])
		if err == nil {
			rts[t] = struct{}{}
		}
	}
	return rts
}

func getTokenEndpointAuthMethod(s string) TokenEndpointAuthMethod {
	var m TokenEndpointAuthMethod
	r, _ := m.Parse(s)
	return r
}

func (cli *Client) GetClientInformationResponse() (*ClientInformationResponse, error) {
	cir := &ClientInformationResponse{
		ClientId:                cli.ID.String(),
		ClientSecret:            cli.Secret,
		ClientIdIssuedAt:        cli.IDIssuedAt.Unix(),
		ClientSecretExpiresAt:   cli.SecretExpiresAt.Unix(),
		RedirectUris:            cli.GetRedirectUris(),
		TokenEndpointAuthMethod: cli.TokenEndpointAuthMethod.String(),
		GrantTypes:              cli.GetGrantTypes(),
		ResponseTypes:           cli.GetResponseTypes(),
		ClientName:              cli.Name,
		ClientUri:               cli.URI,
		LogoUri:                 cli.LogoURI,
		Scope:                   cli.Scope,
		Contacts:                cli.Contacts,
		TosUri:                  cli.TosURI,
		PolicyUri:               cli.PolicyURI,
		JwksUri:                 cli.JwksURI,
		Jwks:                    cli.Jwks,
		SoftwareId:              cli.SoftwareID,
		SoftwareVersion:         cli.SoftwareVersion,
	}
	return cir, nil
}

func (cli *Client) GetClientCtx() *ClientCtx {
	return NewClientCtx(cli.ID, cli.Name, cli.LogoURI)
}

//GetClientID returns the client ID
func (cli *Client) GetClientID() uuid.UUID {
	return cli.ID
}

//GetClientName returns the client name
func (cli *Client) GetClientName() string {
	return cli.Name
}

//GetClientURI return the client URI
func (cli *Client) GetClientURI() string {
	return cli.URI
}

//GetGrantTypes returns slice of client registered grant types
func (cli *Client) GetGrantTypes() []string {
	gts := make([]string, len(cli.GrantTypes))
	i := 0
	for k := range cli.GrantTypes {
		gts[i] = k.String()
		i++
	}
	return gts
}

//GetGroups returns the name of the profile owning the oauth client
func (cli *Client) GetGroups() []string {
	return cli.Groups
}

func (cli *Client) GetOwnerID() uuid.UUID {
	return cli.OwnerID
}

func (cli *Client) GetRedirectUris() []string {
	rus := make([]string, len(cli.RedirectUris))
	i := 0
	for k := range cli.RedirectUris {
		rus[i] = k
		i++
	}
	return rus
}

//GetResponseTypes returns slice of registered response types
func (cli *Client) GetResponseTypes() []string {
	rts := make([]string, len(cli.ResponseTypes))
	i := 0
	for k := range cli.ResponseTypes {
		rts[i] = k.String()
		i++
	}
	return rts
}

//HasGrantType returns true if the requested Grant Type is registered for the client
func (cli *Client) HasGrantType(reqGrantType GrantType) (ok bool) {
	_, ok = cli.GrantTypes[reqGrantType]
	return
}

//HasRedirectUri returns the requested redirection URI and true if it is registered for the client
func (cli *Client) HasRedirectUri(requestedURI string) (ok bool) {
	ok = false
	_, ok = cli.RedirectUris[requestedURI]
	return
}

func (cli *Client) HasResponseType(reqResponseType ResponseType) (ok bool) {
	_, ok = cli.ResponseTypes[reqResponseType]
	return
}

//SetOwner sets the new owner of the client. ID of the profile
func (cli *Client) SetOwner(newOwner uuid.UUID) {
	cli.OwnerID = newOwner
}

//todo secret expires. Secret must be encrypted + salted
func (cli *Client) ValidateClientSecret(secret string) (ok bool) {
	ok = false
	if strings.Compare(cli.Secret, secret) == 0 {
		ok = true
		//todo if client expired validation
	}
	return
}

/**
String containing a space-separated list of Scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the
client can use when requesting access token. The semantics of values in this list are service specific. If omitted, an
authorization server MAY register a client with a default set of scopes.
**/

func (cli *Client) ValidateScope(requestedScope string) (validatedScope string) {
	return getCommonScope(cli.Scope, requestedScope)
}

func (cliCtx *ClientCtx) GetClientID() uuid.UUID {
	return cliCtx.ID
}

func (cliCtx *ClientCtx) GetClientName() string {
	return cliCtx.Name
}

func (cliCtx *ClientCtx) GetLogoURI() string {
	return cliCtx.logoURI
}
