package scim2

type Resource struct {
	ID         string   `json:"id"`
	ExternalID string   `json:"externalId,omitempty"`
	Metadata   Metadata `json:"meta"`
}


//Metadata A complex attribute containing resource metadata.  All "meta" sub-attributes are assigned by the service
// provider (have a "mutability" of "readOnly"), and all of these sub-attributes have a "returned" characteristic of
// "default".  This attribute SHALL be ignored when provided by clients.
type Metadata struct {
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Location     string `json:"location,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Version      string `json:"version,omitempty"`
}
