package repository

import (
	"../scim2"
	"time"
	"github.com/gofrs/uuid"
)

type Metadata struct {
	Created        time.Time
	LastModified   time.Time
	RepositoryName string
	ResourceType   string
	ID             uuid.UUID
	Name           string
}

func (m *Metadata)GetScimMetadata() *scim2.Metadata {
	meta := &scim2.Metadata{
		Created: 		m.Created.String(),
		LastModified:	m.LastModified.String(),
		//todo hostname
		Location:		"https://localhost/scim2/" + m.ResourceType + "/" + m.ID.String(),
		ResourceType:	m.ResourceType,
	}
	return meta;
}


func (m *Metadata) GetName() string {
	return m.Name
}

func (m *Metadata) SetName(name string) {
	m.Name = name
}

func (m *Metadata) GetResourceType() string {
	return m.ResourceType
}

func (m *Metadata) GetUUID() uuid.UUID {
	return m.ID
}

func (m *Metadata) SetResourceType(resourceType string) {
	m.ResourceType = resourceType
}