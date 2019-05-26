package repository

import (
	"github.com/gofrs/uuid"
	"../scim2"
)

type resourceRepository interface {
	init()
	saveResource(resource Resource)
	deleteResource(uuid uuid.UUID)
	loadResource(uuid uuid.UUID, resource Resource) (ok bool)
}

func GetUserScim(uuid uuid.UUID) (*scim2.User, error) {
	userResource := &Metadata{}
	ok := resourceManager.loadResource(uuid, userResource)
	if !ok {
		return nil, scim2.ErrNotFound
	}
	users, err := getUserRepository(userResource.RepositoryName)
	if err != nil {
		return nil, err
	}
	user, ok := users.getUser(userResource.GetName())
	if !ok {
		return nil, scim2.ErrNotFound
	}
	return user.GetScim(), nil
}

func GetUserResourceMetadata(uuid uuid.UUID) (*Metadata, error) {
	userResource := &Metadata{}
	ok := resourceManager.loadResource(uuid, userResource)
	if !ok {
		return nil, scim2.ErrNotFound
	}
	return userResource, nil
}


//todo GetGroup scim
func GetGroupScim(uuid uuid.UUID) (*scim2.Group, error) {
	return nil, nil
}

