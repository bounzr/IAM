package repository

import (
	"../scim2"
	"github.com/gofrs/uuid"
	"time"
)

type Group struct {
	Metadata *ResourceTag
	Members  map[uuid.UUID]interface{}
}

func NewGroup(id uuid.UUID, name string) *Group {
	rm := make(map[uuid.UUID]interface{})
	currentTime := time.Now()
	metadata := &ResourceTag{
		Created:      currentTime,
		ID:           id,
		LastModified: currentTime,
		Name:         name,
		ResourceType: "Group",
	}
	group := &Group{
		Metadata: metadata,
		Members:  rm,
	}
	return group
}

func (g *Group) AddResource(resource ResourceTagger) {
	id := resource.GetUUID()
	g.Members[id] = resource
}

func (g *Group) DeleteResource(resource uuid.UUID) {
	delete(g.Members, resource)
}

func (g *Group) GetScim() *scim2.Group {
	var scimMembers []scim2.GroupMember
	//todo ref
	for m, _ := range g.Members {
		sm := scim2.GroupMember{
			Ref:   "https://localhost/scim2/users/" + m.String(),
			Type:  "User",
			Value: m.String(),
		}
		scimMembers = append(scimMembers, sm)
	}

	group := &scim2.Group{
		DisplayName: g.Metadata.Name,
		ID:          g.Metadata.ID.String(),
		Members:     scimMembers,
		Metadata:    g.Metadata.GetScimMetadata(),
		Schemas:     []string{"urn:ietf:params:scim:schemas:core:2.0:Group"},
	}
	return group
}

func (g *Group) GetResourceTag() *ResourceTag {
	return g.Metadata
}
