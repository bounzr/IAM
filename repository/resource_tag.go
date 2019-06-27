package repository

import (
	"../config"
	"../scim2"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type ResourceTag struct {
	Created        time.Time
	LastModified   time.Time
	RepositoryName string
	ResourceType   string
	ID             uuid.UUID
	Name           string
}

func (m *ResourceTag) GetScimMetadata() *scim2.Metadata {
	hostname := config.IAM.Server.Hostname
	location := fmt.Sprintf("https://%s/scim2/%s/%s", hostname, m.ResourceType, m.ID.String())
	meta := &scim2.Metadata{
		Created:      m.Created.String(),
		LastModified: m.LastModified.String(),
		Location:     location,
		ResourceType: m.ResourceType,
	}
	return meta
}

func (m *ResourceTag) GetName() string {
	return m.Name
}

func (m *ResourceTag) SetName(name string) {
	m.Name = name
}

func (m *ResourceTag) GetResourceType() string {
	return m.ResourceType
}

func (m *ResourceTag) GetUUID() uuid.UUID {
	return m.ID
}

func (m *ResourceTag) SetResourceType(resourceType string) {
	m.ResourceType = resourceType
}
