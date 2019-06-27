package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type ResourceManagerLeveldb struct {
	resources     *leveldb.DB
	resourcesPath string
}

func (r *ResourceManagerLeveldb) init() {
	if len(r.resourcesPath) == 0 {
		r.resourcesPath = "./rep/resource"
	}
	db, err := leveldb.OpenFile(r.resourcesPath, nil)
	if err != nil {
		log.Error("can not init resource repository", zap.Error(err))
	}
	r.resources = db
}

func (r *ResourceManagerLeveldb) close() {
	defer r.resources.Close()
}

func (r *ResourceManagerLeveldb) deleteResourceTag(uuid uuid.UUID) {
	err := r.resources.Delete(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not delete resource", zap.String("resource ID", uuid.String()), zap.Error(err))
	} else {
		log.Debug("deleted resource", zap.String("resource ID", uuid.String()))
	}
}

/*
func (r *ResourceManagerLeveldb) getResourceMetadata(uuid uuid.UUID) (resource ResourceTagger, ok bool) {
	dataBytes, err := r.resources.Get(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not get resource", zap.String("resource ID", uuid.String()), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var res ResourceTagger
	err = dec.Decode(res)
	if err != nil {
		log.Error("can not decode resource", zap.String("resource ID", uuid.String()), zap.Error(err))
		return nil, false
	}
	log.Debug("resource retrieved", zap.String("resource ID", res.GetUUID().String()), zap.String("name", res.GetName()), zap.String("type", res.GetResourceType()))
	return res, true
}
*/
func (r *ResourceManagerLeveldb) loadResourceTag(uuid uuid.UUID, resource ResourceTagger) error {
	dataBytes, err := r.resources.Get(uuid.Bytes(), nil)
	if err != nil {
		log.Error("can not load resource", zap.String("resource ID", uuid.String()), zap.Error(err))
		return err
	} else {
		log.Debug("data loaded", zap.Int("bytes size", len(dataBytes)))
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	err = dec.Decode(resource)
	if err != nil {
		log.Error("can not decode resource", zap.String("resource ID", uuid.String()), zap.Error(err))
		return err
	}
	log.Debug("resource loaded", zap.String("resource ID", resource.GetUUID().String()), zap.String("name", resource.GetName()), zap.String("type", resource.GetResourceType()))
	return nil
}

func (r *ResourceManagerLeveldb) setResourceTag(resource ResourceTagger) {
	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	err := enc.Encode(resource)
	if err != nil {
		log.Error("can not encode resource", zap.String("resource ID", resource.GetUUID().String()), zap.Error(err))
	} else {
		log.Debug("data encoded", zap.Int("bytes size", len(data.Bytes())))
	}
	err = r.resources.Put(resource.GetUUID().Bytes(), data.Bytes(), nil)
	if err != nil {
		log.Error("can not update resource", zap.String("resource ID", resource.GetUUID().String()), zap.Error(err))
	} else {
		log.Debug("updated resource", zap.String("resource ID", resource.GetUUID().String()))
	}
}
