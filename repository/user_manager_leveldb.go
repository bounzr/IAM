package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"strings"
)

type UserManagerLeveldb struct {
	cfgDB    *leveldb.DB
	cfgPath  string
	userDB   *leveldb.DB
	userPath string
}

func (r *UserManagerLeveldb) init() {
	if len(r.cfgPath) == 0 {
		r.cfgPath = "./rep/user_cfg"
	}
	r.openCfgDB()
	if len(r.userPath) == 0 {
		r.userPath = "./rep/user"
	}
	r.openUserDB()
}

func (r *UserManagerLeveldb) close() {
	defer r.cfgDB.Close()
	defer r.userDB.Close()
}

func (r *UserManagerLeveldb) validateUser(username string, password string) error {
	user, ok := r.getUser(username)
	if !ok {
		return ErrInvalidLogin
	}
	if strings.Compare(password, user.Password) == 0 {
		return nil
	}
	return ErrInvalidLogin
}

func (r *UserManagerLeveldb) deleteUser(username string) {
	err := r.userDB.Delete([]byte(username), nil)
	if err != nil {
		log.Error("can not delete user", zap.String("username", username), zap.Error(err))
	} else {
		log.Debug("user deleted", zap.String("username", username))
	}
}

func (r *UserManagerLeveldb) getRepositoryName() string {
	nameByte, err := r.cfgDB.Get([]byte("name"), nil)
	if err != nil {
		log.Error("can not get repository name", zap.Error(err))
		return ""
	}
	return string(nameByte)
}

//todo findUsers(search parameters)
func (r *UserManagerLeveldb) findUsers() ([]User, error) {
	var users []User
	iter := r.userDB.NewIterator(nil, nil)
	for iter.Next() {
		dataBytes := iter.Value()
		data := bytes.NewBuffer(dataBytes)
		dec := gob.NewDecoder(data)
		var user User
		err := dec.Decode(&user)
		if err != nil {
			log.Error("can not decode user", zap.ByteString("id", iter.Key()), zap.Error(err))
		}
		log.Debug("user decoded", zap.ByteString("id", iter.Key()))
		users = append(users, user)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserManagerLeveldb) getUser(username string) (*User, bool) {
	dataBytes, err := r.userDB.Get([]byte(username), nil)
	if err != nil {
		log.Error("can not get user", zap.String("username", username), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var user User
	err = dec.Decode(&user)
	if err != nil {
		log.Error("can not decode user", zap.String("username", username), zap.Error(err))
		return nil, false
	}
	return &user, true
}

func (r *UserManagerLeveldb) openCfgDB() {
	var err error
	r.cfgDB, err = leveldb.OpenFile(r.cfgPath, nil)
	if err != nil {
		log.Error("can not open user config repository", zap.Error(err))
	}
}

func (r *UserManagerLeveldb) openUserDB() {
	var err error
	r.userDB, err = leveldb.OpenFile(r.userPath, nil)
	if err != nil {
		log.Error("can not open user repository", zap.Error(err))
	}
}

func (r *UserManagerLeveldb) setRepositoryName(name string) {
	err := r.cfgDB.Put([]byte("name"), []byte(name), nil)
	if err != nil {
		log.Error("can not set repository name", zap.String("repository", name), zap.Error(err))
	} else {
		log.Debug("repository name set", zap.String("repository", name))
	}
}

func (r *UserManagerLeveldb) setUser(user *User) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(user)
	if err != nil {
		log.Error("can not encode user", zap.String("username", user.UserName), zap.Error(err))
	}
	err = r.userDB.Put([]byte(user.UserName), data.Bytes(), nil)
	if err != nil {
		log.Error("can not save user", zap.String("username", user.UserName), zap.Error(err))
	} else {
		log.Debug("user updated", zap.String("user ID", user.ID.String()), zap.String("username", user.UserName))
	}
}
