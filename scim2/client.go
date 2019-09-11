package scim2

type Client struct {
	Name                    string            `json:"clientName,omitempty"`
	PasswordExpiresAt       int64             `json:"clientSecretExpiresAt,omitempty"`
	URI                     string            `json:"clientUri,omitempty"`
	Contacts                []string          `json:"contacts,omitempty"`
	GrantTypes              []string          `json:"grantTypes,omitempty"`
	Groups                  []GroupAssignment `json:"groups,omitempty"`
	ID                      string            `json:"id"`
	JwksURI                 string            `json:"jwksUri,omitempty"`
	Jwks                    string            `json:"jwks,omitempty"` //todo jwks object
	LogoUri                 string            `json:"logoUri,omitempty"`
	Metadata                *Metadata         `json:"meta"`
	Password                string            `json:"password,omitempty"`
	PolicyUri               string            `json:"policyUri,omitempty"`
	RedirectUris            []string          `json:"redirectUris"`
	ResponseTypes           []string          `json:"responseTypes,omitempty"`
	Schemas                 []string          `json:"schemas,omitempty"`
	Scope                   string            `json:"scope,omitempty"`
	SoftwareId              string            `json:"softwareId,omitempty"`
	SoftwareVersion         string            `json:"softwareVersion,omitempty"`
	TokenEndpointAuthMethod string            `json:"tokenEndpointAuthMethod"`
	TosUri                  string            `json:"tosUri,omitempty"`
}

func (c *Client) GetGroups() []GroupAssignment {
	return c.Groups
}

func (c *Client) GetID() string {
	return c.ID
}
