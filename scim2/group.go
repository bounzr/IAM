package scim2

import (
	"../config"
	"fmt"
)

type Group struct {
	DisplayName string        `json:"displayName,omitempty"` // A human-readable name for the Group.  REQUIRED.
	ID          string        `json:"id,omitempty"`
	Members     []GroupMember `json:"members,omitempty"`
	Metadata    *Metadata     `json:"meta,omitempty"`
	Schemas     []string      `json:"schemas,omitempty"`
}

type GroupAssigner interface {
	GetGroups() []GroupAssignment
}

type GroupAssignment struct {
	Display string `json:"display,omitempty"`
	Ref     string `json:"$ref,omitempty"`
	Value   string `json:"value,omitempty"`
}

type GroupMember struct {
	Ref   string `json:"$ref,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

func (g *Group) GetGroupAssignment() *GroupAssignment {
	hostname := config.IAM.Server.Hostname
	location := fmt.Sprintf("https://%s/scim2/%s/%s", hostname, g.Metadata.ResourceType, g.ID)
	assignment := &GroupAssignment{
		Display: g.DisplayName,
		Ref:     location,
		Value:   g.ID,
	}
	return assignment
}
