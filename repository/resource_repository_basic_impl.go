package repository

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ResourceRepositoryBasic struct {
	resources map[uuid.UUID]Resource
}

func NewResourceRepositoryBasic() resourceRepository {
	repo := &ResourceRepositoryBasic{}
	repo.init()
	return repo
}

func (r *ResourceRepositoryBasic) init() {
	r.resources = make(map[uuid.UUID]Resource)
}

func (r *ResourceRepositoryBasic) deleteResource(uuid uuid.UUID) {
	delete(r.resources, uuid)
}

func (r *ResourceRepositoryBasic) loadResource(uuid uuid.UUID, resource Resource) (ok bool) {
	res, ok := r.resources[uuid]
	if ok {
		resource.SetName(res.GetName())
		resource.SetResourceType(res.GetResourceType())
	}else{
		log.Error("can not load resource", zap.String("resource id", uuid.String()), zap.Error(ErrResourceNotFound))
	}
	return ok
}

func (r *ResourceRepositoryBasic) saveResource(resource Resource) {
	r.resources[resource.GetUUID()] = resource
}
