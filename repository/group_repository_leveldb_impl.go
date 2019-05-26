package repository

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/gofrs/uuid"
	"bytes"
	"encoding/gob"
	"go.uber.org/zap"
)

type GroupRepositoryLeveldb struct {
	groups *leveldb.DB
}

func NewGroupRepositoryLeveldb() GroupRepository {
	gr := &GroupRepositoryLeveldb{}
	gr.init()
	return gr
}

func (gr *GroupRepositoryLeveldb) init() {
	db, err := leveldb.OpenFile("./rep/group", nil)
	if err != nil {
		log.Error("can not init group repository", zap.Error(err))
	}
	gr.groups = db
	//todo defer db.Close()
}

func (gr *GroupRepositoryLeveldb) addGroupResource(groupID uuid.UUID, resource Resource) (ok bool) {
	res, ok := gr.getGroup(groupID)
	if ok {
		res.Members[resource.GetUUID()] = resource
		gr.setGroup(res)
		log.Debug("resource added to group", zap.String("group id", groupID.String()), zap.String("resource id", resource.GetUUID().String()))
	}else{
		log.Error("group not found", zap.String("group id", groupID.String()), zap.Error(ErrGroupNotFound))
	}
	return
}

func (gr *GroupRepositoryLeveldb) deleteGroup(groupID uuid.UUID) {
	err := gr.groups.Delete(groupID.Bytes(), nil)
	if err != nil {
		log.Error("can not delete group", zap.String("group id", groupID.String()), zap.Error(err))
	} else {
		log.Debug("deleted group", zap.String("group id", groupID.String()))
	}
}

func (gr *GroupRepositoryLeveldb) deleteGroupResource(groupID uuid.UUID, resourceID uuid.UUID) {
	res, ok := gr.getGroup(groupID)
	if ok {
		delete(res.Members, resourceID)
		gr.setGroup(res)
		log.Debug("resource deleted from group", zap.String("group id", groupID.String()), zap.String("resource id", resourceID.String()))
	}else{
		log.Error("group not found", zap.String("group id", groupID.String()), zap.Error(ErrGroupNotFound))
	}
}

func (gr *GroupRepositoryLeveldb) getGroup(groupID uuid.UUID) (*Group, bool) {
	dataBytes, err := gr.groups.Get(groupID.Bytes(), nil)
	if err != nil {
		log.Error("can not get group", zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var group Group
	err = dec.Decode(&group)
	if err != nil {
		log.Error("can not decode group", zap.Error(err))
		return nil, false
	}
	return &group, true
}

func (gr *GroupRepositoryLeveldb) setGroup(group *Group) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(group)
	if err != nil {
		log.Error("can not encode group", zap.Error(err))
	}
	err = gr.groups.Put(group.Metadata.ID.Bytes(), data.Bytes(), nil)
	if err != nil {
		log.Error("can not update group", zap.Error(err))
	} else {
		log.Debug("updated group", zap.String("group id", group.Metadata.ID.String()))
	}
}

