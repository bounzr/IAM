package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"strings"
)

type UserManagerLeveldb struct {
	cfgDB    *leveldb.DB
	cfgPath  string
	nameDB   *leveldb.DB
	namePath string
	uuidDB   *leveldb.DB
	uuidPath string
}

func (r *UserManagerLeveldb) init() {
	if len(r.cfgPath) == 0 {
		r.cfgPath = "./rep/user/cfg"
	}
	r.openCfgDB()
	if len(r.namePath) == 0 {
		r.namePath = "./rep/user/name"
	}
	r.openNameDB()
	if len(r.uuidPath) == 0 {
		r.uuidPath = "./rep/user/uuid"
	}
	r.openUUIDDB()
}

func (r *UserManagerLeveldb) close() {
	defer r.cfgDB.Close()
	defer r.nameDB.Close()
	defer r.uuidDB.Close()
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

func (r *UserManagerLeveldb) deleteUser(userID interface{}) {
	username := r.getUsername(userID)
	if len(username) == 0 {
		return
	}
	err := r.nameDB.Delete([]byte(username), nil)
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
	iter := r.nameDB.NewIterator(nil, nil)
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

func (r *UserManagerLeveldb) getUser(userID interface{}) (*User, bool) {
	username := r.getUsername(userID)
	if len(username) == 0 {
		return nil, false
	}
	dataBytes, err := r.nameDB.Get([]byte(username), nil)
	if err != nil {
		log.Debug("can not get user", zap.String("username", username), zap.Error(err))
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

func (r *UserManagerLeveldb) getUsername(userID interface{}) string {
	switch userID.(type) {
	case string:
		log.Debug("userID is a string")
		return userID.(string)
	case uuid.UUID:
		dataBytes, err := r.uuidDB.Get(userID.(uuid.UUID).Bytes(), nil)
		if err != nil {
			log.Error("can not decode uid", zap.ByteString("uuid", userID.(uuid.UUID).Bytes()), zap.Error(err))
			return ""
		}
		return string(dataBytes)
	default:
		return ""
	}

}

func (r *UserManagerLeveldb) openCfgDB() {
	var err error
	r.cfgDB, err = leveldb.OpenFile(r.cfgPath, nil)
	if err != nil {
		log.Error("can not open user config repository", zap.Error(err))
	}
}

func (r *UserManagerLeveldb) openNameDB() {
	var err error
	r.nameDB, err = leveldb.OpenFile(r.namePath, nil)
	if err != nil {
		log.Error("can not open user repository", zap.Error(err))
	}
}

func (r *UserManagerLeveldb) openUUIDDB() {
	var err error
	r.uuidDB, err = leveldb.OpenFile(r.uuidPath, nil)
	if err != nil {
		log.Error("can not open uuid repository", zap.Error(err))
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
	err = r.nameDB.Put([]byte(user.UserName), data.Bytes(), nil)
	if err != nil {
		log.Error("can not save user", zap.String("username", user.UserName), zap.Error(err))
	} else {
		err = r.uuidDB.Put(user.ID.Bytes(), []byte(user.UserName), nil)
		if err != nil {
			log.Error("can not save user", zap.String("username", user.UserName), zap.Error(err))
		} else {
			log.Debug("user updated", zap.String("user ID", user.ID.String()), zap.String("username", user.UserName))
		}
	}
}
