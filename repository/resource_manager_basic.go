package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ResourceManagerBasic struct {
	resources map[uuid.UUID][]byte
}

func (r *ResourceManagerBasic) init() {
	r.resources = make(map[uuid.UUID][]byte)
}

func (r *ResourceManagerBasic) close() {
	//nothing
}

func (r *ResourceManagerBasic) deleteResourceTag(uuid uuid.UUID) {
	delete(r.resources, uuid)
}

func (r *ResourceManagerBasic) loadResourceTag(uuid uuid.UUID, resource ResourceTagger) error {
	resBytes, ok := r.resources[uuid]
	if ok {
		data := bytes.NewBuffer(resBytes)
		dec := gob.NewDecoder(data)
		err := dec.Decode(resource)
		if err != nil {
			log.Error("can not decode resource", zap.String("resource ID", uuid.String()), zap.Error(err))
			return err
		}
		log.Debug("resource loaded", zap.String("resource ID", resource.GetUUID().String()), zap.String("name", resource.GetName()), zap.String("type", resource.GetResourceType()))
	} else {
		log.Error("can not load resource", zap.String("resource ID", uuid.String()), zap.Error(ErrResourceNotFound))
		return ErrResourceNotFound
	}
	return nil
}

func (r *ResourceManagerBasic) setResourceTag(resource ResourceTagger) {
	data := new(bytes.Buffer)
	enc := gob.NewEncoder(data)
	err := enc.Encode(resource)
	if err != nil {
		log.Error("can not encode resource", zap.String("resource ID", resource.GetUUID().String()), zap.Error(err))
	} else {
		log.Debug("data encoded", zap.Int("bytes size", len(data.Bytes())))
	}
	r.resources[resource.GetUUID()] = data.Bytes()
}
