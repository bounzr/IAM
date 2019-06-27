package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"reflect"
	"strings"
)

type GroupManagerLeveldb struct {
	groups     *leveldb.DB
	groupsPath string
}

func (gr *GroupManagerLeveldb) close() {
	defer gr.groups.Close()
}

func (gr *GroupManagerLeveldb) findGroups(conditions map[string]interface{}) ([]Group, error) {
	var groups []Group
	iter := gr.groups.NewIterator(nil, nil)
	for iter.Next() {
		dataBytes := iter.Value()
		data := bytes.NewBuffer(dataBytes)
		dec := gob.NewDecoder(data)
		var group Group
		err := dec.Decode(&group)
		if err != nil {
			log.Error("can not decode group", zap.Error(err))
		}
		log.Debug("group decoded", zap.String("id", group.Metadata.ID.String()))
		nameCondition := conditions["name"]
		var nameToFind string
		if nameCondition != nil {
			nameToFind = nameCondition.(string)
		}
		if len(nameToFind) > 0 {
			log.Debug("finding groups with name condition", zap.String("name", conditions["name"].(string)))
			if strings.Compare(nameToFind, group.Metadata.Name) != 0 {
				continue
			}
		}
		memberCondition := conditions["member"]
		var memberToFind uuid.UUID
		if memberCondition != nil && reflect.TypeOf(memberToFind) == reflect.TypeOf(uuid.UUID{}) {
			memberToFind = memberCondition.(uuid.UUID)
			log.Debug("finding groups with member condition", zap.String("member", conditions["member"].(uuid.UUID).String()))
			_, ok := group.Members[memberToFind]
			if !ok {
				continue
			}
		}
		//todo control other conditions or wrong filters
		groups = append(groups, group)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return groups, nil
}

func (gr *GroupManagerLeveldb) init() {
	if len(gr.groupsPath) == 0 {
		gr.groupsPath = "./rep/group"
	}
	db, err := leveldb.OpenFile(gr.groupsPath, nil)
	if err != nil {
		log.Error("can not init group repository", zap.Error(err))
	}
	gr.groups = db
}

func (gr *GroupManagerLeveldb) setGroupResource(groupID uuid.UUID, resource ResourceTagger) (ok bool) {
	res, ok := gr.getGroup(groupID)
	if ok {
		res.Members[resource.GetUUID()] = resource
		gr.setGroup(res)
		log.Debug("resource added to group", zap.String("group ID", groupID.String()), zap.String("resource ID", resource.GetUUID().String()))
	} else {
		log.Error("group not found", zap.String("group ID", groupID.String()), zap.Error(ErrGroupNotFound))
	}
	return
}

func (gr *GroupManagerLeveldb) deleteGroup(groupID uuid.UUID) {
	err := gr.groups.Delete(groupID.Bytes(), nil)
	if err != nil {
		log.Error("can not delete group", zap.String("group ID", groupID.String()), zap.Error(err))
	} else {
		log.Debug("deleted group", zap.String("group ID", groupID.String()))
	}
}

func (gr *GroupManagerLeveldb) deleteGroupResource(groupID uuid.UUID, resourceID uuid.UUID) {
	res, ok := gr.getGroup(groupID)
	if ok {
		delete(res.Members, resourceID)
		gr.setGroup(res)
		log.Debug("resource deleted from group", zap.String("group ID", groupID.String()), zap.String("resource ID", resourceID.String()))
	} else {
		log.Error("group not found", zap.String("group ID", groupID.String()), zap.Error(ErrGroupNotFound))
	}
}

func (gr *GroupManagerLeveldb) getGroup(groupID uuid.UUID) (*Group, bool) {
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

func (gr *GroupManagerLeveldb) setGroup(group *Group) {
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
		log.Debug("updated group", zap.String("group ID", group.Metadata.ID.String()))
	}
}
