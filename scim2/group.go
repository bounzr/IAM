package scim2

type Group struct {
	DisplayName string        `json:"displayName,omitempty"` // A human-readable name for the Group.  REQUIRED.
	ID          string        `json:"id,omitempty"`
	Members     []GroupMember `json:"members,omitempty"`
	Metadata    *Metadata      `json:"meta,omitempty"`
	Schemas     []string      `json:"schemas,omitempty"`
}

type GroupMember struct {
	Ref   string `json:"$ref,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
