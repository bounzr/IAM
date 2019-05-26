package repository

import (

	"github.com/syndtr/goleveldb/leveldb"
	"bytes"
	"encoding/gob"
	"go.uber.org/zap"
)

type SessionRepositoryLeveldb struct {
	sessionDB *leveldb.DB
}

func NewSessionRepositoryLeveldb() SessionRepository {
	rep := &SessionRepositoryLeveldb{}
	rep.init()
	return rep
}

func (r *SessionRepositoryLeveldb) init() {
	var err error
	r.sessionDB, err = leveldb.OpenFile("./rep/session_user", nil)
	if err != nil {
		log.Error("can not init session repository", zap.Error(err))
	}
}

func (r *SessionRepositoryLeveldb) addSessionContext(token SessionToken, user *UserCtx) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(user)
	if err != nil {
		log.Error("can not encode user context", zap.String("user id", user.UserID.String()), zap.Error(err))
	}
	err = r.sessionDB.Put(token[:], data.Bytes(), nil)
	if err != nil {
		log.Error("can not update session context", zap.String("user id", user.UserID.String()), zap.Error(err))
	} else {
		log.Debug("session context udpated for user", zap.String("user id", user.UserID.String()))
	}
}

func (r *SessionRepositoryLeveldb) deleteSessionContext(token SessionToken) error {
	err := r.sessionDB.Delete(token[:], nil)
	if err != nil {
		log.Error("can not delete session context", zap.Error(err))
		return err
	} else {
		log.Debug("session context was deleted")
	}
	return nil
}

func (r *SessionRepositoryLeveldb) getSessionContext(token SessionToken) (*UserCtx, error) {
	dataBytes, err := r.sessionDB.Get(token[:], nil)
	if err != nil {
		log.Error("can not get session context", zap.Error(err))
		return nil, err
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var user UserCtx
	err = dec.Decode(&user)
	if err != nil {
		log.Error("can not decode user context", zap.Error(err))
		return nil, err
	}
	log.Debug("user context retrieved", zap.String("user id", user.UserID.String()))
	return &user, nil
}