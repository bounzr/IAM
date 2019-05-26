package repository

import (
	"github.com/gofrs/uuid"
)

type Resource interface {
	GetName() string
	SetName(name string)
	GetResourceType() string
	SetResourceType(resourceType string)
	GetUUID() uuid.UUID
}

