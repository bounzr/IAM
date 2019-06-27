package repository

import (
	"../config"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ResourceManager interface {
	init()
	close()
	setResourceTag(resource ResourceTagger)
	deleteResourceTag(uuid uuid.UUID)
	loadResourceTag(uuid uuid.UUID, resource ResourceTagger) error
}

func initResources() {
	implementation := config.IAM.Resources.Implementation
	switch implementation {
	case "leveldb":
		resourceManager = &ResourceManagerLeveldb{resourcesPath: "./rep/resource"}
	default:
		resourceManager = &ResourceManagerBasic{}
	}
	resourceManager.init()
}

func AddResource(resource ResourceTagProvider) {
	resourceManager.setResourceTag(resource.GetResourceTag())
}

func GetUserScim(uuid uuid.UUID) (*scim2.User, error) {
	userResource := &ResourceTag{}
	err := resourceManager.loadResourceTag(uuid, userResource)
	if err != nil {
		return nil, err
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

func GetResourceMetadata(uuid uuid.UUID) (tag *ResourceTag, ok bool) {
	rtag := &ResourceTag{}
	err := resourceManager.loadResourceTag(uuid, rtag)
	if err != nil {
		log.Error("can not get resource tag", zap.String("id", uuid.String()), zap.Error(err))
		return nil, false
	}
	return rtag, true
}

//todo GetGroup scim
func GetGroupScim(uuid uuid.UUID) (*scim2.Group, error) {
	return nil, nil
}
