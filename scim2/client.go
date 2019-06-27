package scim2

type Client struct {
	Name                    string            `json:"client_name,omitempty"`
	PasswordExpiresAt       int64             `json:"client_secret_expires_at,omitempty"`
	URI                     string            `json:"client_uri,omitempty"`
	Contacts                []string          `json:"contacts,omitempty"`
	GrantTypes              []string          `json:"grant_types,omitempty"`
	Groups                  []GroupAssignment `json:"groups,omitempty"`
	ID                      string            `json:"id"`
	JwksURI                 string            `json:"jwks_uri,omitempty"`
	Jwks                    string            `json:"jwks,omitempty"` //todo jwks object
	LogoUri                 string            `json:"logo_uri,omitempty"`
	Metadata                *Metadata         `json:"meta"`
	Password                string            `json:"password,omitempty"`
	PolicyUri               string            `json:"policy_uri,omitempty"`
	RedirectUris            []string          `json:"redirect_uris"`
	ResponseTypes           []string          `json:"response_types,omitempty"`
	Schemas                 []string          `json:"schemas,omitempty"`
	Scope                   string            `json:"Scope,omitempty"`
	SoftwareId              string            `json:"software_id,omitempty"`
	SoftwareVersion         string            `json:"software_version,omitempty"`
	TokenEndpointAuthMethod string            `json:"token_endpoint_auth_method"`
	TosUri                  string            `json:"tos_uri,omitempty"`
}
