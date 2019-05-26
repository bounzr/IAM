package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type ResourceRepositoryLeveldb struct {
	db *leveldb.DB
}

func NewResourceRepositoryLeveldb() resourceRepository {
	repo := &ResourceRepositoryLeveldb{}
	repo.init()
	return repo
}

func (r *ResourceRepositoryLeveldb) init() {
	db, err := leveldb.OpenFile("./rep/resource", nil)
	if err != nil {
		log.Error("can not init resource repository", zap.Error(err))
	}
	r.db = db
	//todo defer db.Close()
}

func (r *ResourceRepositoryLeveldb) deleteResource(uuid uuid.UUID) {
	err := r.db.Delete(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not delete resource", zap.String("resource id", uuid.String()),zap.Error(err))
	} else {
		log.Debug("deleted resource", zap.String("resource id", uuid.String()))
	}
}

func (r *ResourceRepositoryLeveldb) getResource(uuid uuid.UUID) (resource Resource, ok bool) {
	dataBytes, err := r.db.Get(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not get resource", zap.String("resource id", uuid.String()), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var res Resource
	err = dec.Decode(res)
	if err != nil {
		log.Error("can not decode resource", zap.String("resource id", uuid.String()), zap.Error(err))
		return nil, false
	}
	log.Debug("resource retrieved", zap.String("resource id", res.GetUUID().String()), zap.String("name", res.GetName()), zap.String("type", res.GetResourceType()))
	return res, true
}

func (r *ResourceRepositoryLeveldb) loadResource(uuid uuid.UUID, resource Resource) (ok bool) {
	dataBytes, err := r.db.Get(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not load resource", zap.String("resource id", uuid.String()), zap.Error(err))
		return false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	err = dec.Decode(resource)
	if err != nil {
		log.Error("can not decode resource", zap.String("resource id", uuid.String()), zap.Error(err))
		return false
	}
	log.Debug("resource loaded", zap.String("resource id", resource.GetUUID().String()), zap.String("name", resource.GetName()), zap.String("type", resource.GetResourceType()))
	return true
}

func (r *ResourceRepositoryLeveldb) saveResource(resource Resource) {
	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	err := enc.Encode(resource)
	if err != nil {
		log.Error("can not encode resource", zap.String("resource id", resource.GetUUID().String()), zap.Error(err))
	}
	err = r.db.Put(resource.GetUUID().Bytes(), data.Bytes(), nil)
	if err != nil {
		log.Error("can not update resource", zap.String("resource id", resource.GetUUID().String()), zap.Error(err))
	} else {
		log.Debug("updated resource", zap.String("resource id",resource.GetUUID().String()))
	}
}