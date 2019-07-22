package scim2

type User struct {
	Active            bool                  `json:"active,omitempty"`
	Addresses         []Address             `json:"addresses,omitempty"`
	DisplayName       string                `json:"displayName,omitempty"` //The name of the user, suitable for display to end-users
	Entitlements      []string              `json:"entitlements,omitempty"`
	Emails            []MultiValueAttribute `json:"emails,omitempty"`
	ExternalId        string                `json:"externalId,omitempty"` //A String that is an identifier for the resource as defined by the provisioning client. This identifier MUST be unique across the SCIM service provider’s entire set of resources
	Groups            []GroupAssignment     `json:"groups,omitempty"`
	ID                string                `json:"id"` // A unique identifier for a SCIM resource as defined by the service provider
	Ims               []MultiValueAttribute `json:"ims,omitempty"`
	Locale            string                `json:"locale,omitempty"`
	Metadata          *Metadata             `json:"meta"`
	Name              *Name                 `json:"name,omitempty"`     //The components of the user’s name
	NickName          string                `json:"nickName,omitempty"` //The casual way to address the user in real life
	Password          string                `json:"password,omitempty"`
	PhoneNumbers      []MultiValueAttribute `json:"phoneNumbers,omitempty"`
	Photos            []MultiValueAttribute `json:"photos,omitempty"`
	PreferredLanguage string                `json:"preferredLanguage,omitempty"`
	ProfileURL        string                `json:"profileUrl,omitempty"`
	Roles             []string              `json:"roles,omitempty"`
	Schemas           []string              `json:"schemas,omitempty"`
	Timezone          string                `json:"timezone,omitempty"`
	Title             string                `json:"title,omitempty"`
	UserName          string                `json:"userName"` //A unique identifier for the user, typically used by the user to directly authenticate to the service provider. This identifier MUST be unique across the service provider’s entire set of Users
	UserType          string                `json:"userType,omitempty"`
	X509Certificates  []MultiValueAttribute `json:"x509Certificates,omitempty"`
}

type Address struct {
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Type          string `json:"type,omitempty"`
}

type Name struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

type MultiValueAttribute struct {
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

func (u *User) GetGroups() []GroupAssignment {
	return u.Groups
}

func (u *User) GetID() string {
	return u.ID
}
