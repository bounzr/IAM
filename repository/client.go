package repository

import (
	"../config"
	"../oauth2"
	"../scim2"
	"../utils"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

type Client struct {
	Contacts                []string
	Created                 time.Time
	GrantTypes              map[string]struct{}
	ID                      uuid.UUID
	Jwks                    string //todo jwks object
	JwksURI                 string
	LastModified            time.Time
	LogoURI                 string
	Name                    string
	OwnerID                 uuid.UUID
	PolicyURI               string
	RedirectURIs            map[string]struct{}
	ResponseTypes           map[string]struct{}
	Scope                   string
	Secret                  string
	SecretExpiresAt         time.Time
	SoftwareID              string
	SoftwareVersion         string
	TokenEndpointAuthMethod string
	TosURI                  string
	URI                     string
}

func (c *Client) GetResourceTag() *ResourceTag {
	resourceTag := &ResourceTag{
		Created:        c.Created,
		LastModified:   c.LastModified,
		RepositoryName: "client", //todo repository names
		ResourceType:   "Client",
		ID:             c.ID,
		Name:           c.Name,
	}
	return resourceTag
}

func (c *Client) GetScim() *scim2.Client {
	memberGroupFilter := make(map[string]interface{})
	memberGroupFilter["member"] = c.ID
	groups := FindGroupAssignments(memberGroupFilter)

	client := &scim2.Client{
		Name:                    c.Name,
		PasswordExpiresAt:       c.SecretExpiresAt.Unix(),
		URI:                     c.URI,
		Contacts:                c.Contacts,
		GrantTypes:              c.GetGrantTypes(),
		Groups:                  groups,
		ID:                      c.ID.String(),
		JwksURI:                 c.JwksURI,
		Jwks:                    c.Jwks,
		LogoUri:                 c.LogoURI,
		Metadata:                c.GetResourceTag().GetScimMetadata(),
		PolicyUri:               c.PolicyURI,
		RedirectUris:            c.GetRedirectUris(),
		ResponseTypes:           c.GetResponseTypes(),
		Schemas:                 []string{"org:bounzer:iam:scim2:1.0:Client"},
		Scope:                   c.Scope,
		SoftwareId:              c.SoftwareID,
		SoftwareVersion:         c.SoftwareVersion,
		TokenEndpointAuthMethod: c.TokenEndpointAuthMethod,
		TosUri:                  c.TosURI,
	}
	return client
}

func getCommonScope(scope1 string, scope2 string) (commonScope string) {
	m1 := getScopesMap(scope1)
	s2 := strings.Split(strings.TrimSpace(scope2), " ")
	commonScope = ""
	first := true
	for _, v := range s2 {
		if _, ok := m1[v]; ok {
			if first == true {
				commonScope = v
				first = false
			} else {
				commonScope = commonScope + " " + v
			}
		}
	}
	return
}

func getScopesMap(scope string) map[string]struct{} {
	m := make(map[string]struct{})
	scp := strings.Split(strings.TrimSpace(scope), " ")
	for _, v := range scp {
		m[v] = struct{}{}
	}
	return m
}

func getSliceToMap(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, v := range slice {
		set[strings.TrimSpace(v)] = struct{}{}
	}
	return set
}

func (c *Client) GetClientCtx() *oauth2.ClientCtx {
	return oauth2.NewClientCtx(c.ID, c.Name, c.LogoURI)
}

func (c *Client) GetClientInformationResponse() (*oauth2.ClientInformationResponse, error) {
	//todo verify that the client is not expired!
	cir := &oauth2.ClientInformationResponse{
		ClientId:                c.ID.String(),
		ClientSecret:            c.Secret,
		ClientIdIssuedAt:        c.Created.Unix(),
		ClientSecretExpiresAt:   c.SecretExpiresAt.Unix(), //todo expired client secret?
		RedirectUris:            c.GetRedirectUris(),
		TokenEndpointAuthMethod: c.TokenEndpointAuthMethod,
		GrantTypes:              c.GetGrantTypes(),
		ResponseTypes:           c.GetResponseTypes(),
		ClientName:              c.Name,
		ClientUri:               c.URI,
		LogoUri:                 c.LogoURI,
		Scope:                   c.Scope,
		Contacts:                c.Contacts,
		TosUri:                  c.TosURI,
		PolicyUri:               c.PolicyURI,
		JwksUri:                 c.JwksURI,
		Jwks:                    c.Jwks,
		SoftwareId:              c.SoftwareID,
		SoftwareVersion:         c.SoftwareVersion,
	}
	return cir, nil
}

//GetGrantTypes returns slice of client registered grant types
func (c *Client) GetGrantTypes() []string {
	gts := make([]string, len(c.GrantTypes))
	i := 0
	for k := range c.GrantTypes {
		gts[i] = k
		i++
	}
	return gts
}

func (c *Client) GetRedirectUris() []string {
	rus := make([]string, len(c.RedirectURIs))
	i := 0
	for k := range c.RedirectURIs {
		rus[i] = k
		i++
	}
	return rus
}

//GetResponseTypes returns slice of registered response types
func (c *Client) GetResponseTypes() []string {
	rts := make([]string, len(c.ResponseTypes))
	i := 0
	for k := range c.ResponseTypes {
		rts[i] = k
		i++
	}
	return rts
}

//HasGrantType returns true if the requested Grant Type is registered for the client
func (c *Client) HasGrantType(grantType string) (ok bool) {
	_, ok = c.GrantTypes[grantType]
	return
}

//HasResponseType returns true if the requested Response Type is registered for the client
func (c *Client) HasResponseType(responseType string) (ok bool) {
	_, ok = c.ResponseTypes[responseType]
	return
}

//HasResponseType returns true if the requested Response Type is registered for the client
func (c *Client) HasRedirectURI(redirectURI string) (ok bool) {
	_, ok = c.RedirectURIs[redirectURI]
	return
}

func NewClientFromOauth(request *oauth2.ClientRegistrationRequest) (*Client, error) {
	clientID, err := uuid.NewV4()
	issuedAt := time.Now()
	secretExpiresAt := issuedAt.Add(config.IAM.Clients.GetSecretDuration())

	cli := &Client{
		Contacts:                request.Contacts,
		Created:                 issuedAt,
		GrantTypes:              getSliceToMap(request.GrantTypes),
		ID:                      clientID,
		Jwks:                    request.Jwks,    //todo jwks is an object and not a string. develop object
		JwksURI:                 request.JwksURI, //todo if jwksuri is provided get jwks from uri
		LastModified:            issuedAt,
		LogoURI:                 request.LogoUri,
		Name:                    request.ClientName,
		PolicyURI:               request.PolicyUri,
		RedirectURIs:            getSliceToMap(request.RedirectUris),
		ResponseTypes:           getSliceToMap(request.ResponseTypes),
		Scope:                   strings.TrimSpace(request.Scope),
		Secret:                  utils.GetRandomPassword(16), //todo password generator
		SecretExpiresAt:         secretExpiresAt,
		SoftwareID:              request.SoftwareId,
		SoftwareVersion:         request.SoftwareVersion,
		TokenEndpointAuthMethod: request.TokenEndpointAuthMethod,
		TosURI:                  request.TosUri,
		URI:                     request.ClientUri,
	}
	return cli, err
}

func NewClientFromScim(request *scim2.Client) (*Client, error) {
	clientID, err := uuid.NewV4()
	issuedAt := time.Now()
	secretExpiresAt := issuedAt.Add(config.IAM.Clients.GetSecretDuration())

	cli := &Client{
		Contacts:                request.Contacts,
		Created:                 issuedAt,
		GrantTypes:              getSliceToMap(request.GrantTypes),
		ID:                      clientID,
		Jwks:                    request.Jwks,    //todo jwks is an object and not a string. develop object
		JwksURI:                 request.JwksURI, //todo if jwksuri is provided get jwks from uri
		LastModified:            issuedAt,
		LogoURI:                 request.LogoUri,
		Name:                    request.Name,
		PolicyURI:               request.PolicyUri,
		RedirectURIs:            getSliceToMap(request.RedirectUris),
		ResponseTypes:           getSliceToMap(request.ResponseTypes),
		Scope:                   strings.TrimSpace(request.Scope),
		Secret:                  utils.GetRandomPassword(16),
		SecretExpiresAt:         secretExpiresAt,
		SoftwareID:              request.SoftwareId,
		SoftwareVersion:         request.SoftwareVersion,
		TokenEndpointAuthMethod: request.TokenEndpointAuthMethod,
		TosURI:                  request.TosUri,
		URI:                     request.URI,
	}
	return cli, err
}

//todo Secret must be encrypted + salted
//todo client secret must be modifiable and the secretExpiresAt reset
func (c *Client) ValidateClientSecret(secret string) (ok bool) {
	if !utils.InTimeSpan(c.Created, c.SecretExpiresAt, time.Now()) {
		return false
	}
	ok = false
	if strings.Compare(c.Secret, secret) == 0 {
		ok = true
	}
	return
}

func (c *Client) ValidateScope(requestedScope string) (validatedScope string) {
	return getCommonScope(c.Scope, requestedScope)
}
