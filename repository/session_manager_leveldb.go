package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type SessionManagerLeveldb struct {
	sessions     *leveldb.DB
	sessionsPath string
}

func (r *SessionManagerLeveldb) init() {
	var err error
	if len(r.sessionsPath) == 0 {
		r.sessionsPath = "./rep/session"
	}
	r.sessions, err = leveldb.OpenFile(r.sessionsPath, nil)
	if err != nil {
		log.Error("can not init session repository", zap.Error(err))
	}
}

func (r *SessionManagerLeveldb) close() {
	defer r.sessions.Close()
}

func (r *SessionManagerLeveldb) setSessionContext(token SessionToken, user *UserCtx) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(user)
	if err != nil {
		log.Error("can not encode user context", zap.String("user id", user.UserID.String()), zap.Error(err))
	}
	err = r.sessions.Put(token[:], data.Bytes(), nil)
	if err != nil {
		log.Error("can not update session context", zap.String("user id", user.UserID.String()), zap.Error(err))
	} else {
		log.Debug("session context udpated for user", zap.String("user id", user.UserID.String()))
	}
}

func (r *SessionManagerLeveldb) deleteSessionContext(token SessionToken) error {
	err := r.sessions.Delete(token[:], nil)
	if err != nil {
		log.Error("can not delete session context", zap.Error(err))
		return err
	} else {
		log.Debug("session context was deleted")
	}
	return nil
}

func (r *SessionManagerLeveldb) getSessionContext(token SessionToken) (*UserCtx, error) {
	dataBytes, err := r.sessions.Get(token[:], nil)
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
