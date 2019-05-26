package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"strings"
)

type UserRepositoryLeveldb struct {
	cfgDB  *leveldb.DB
	userDB *leveldb.DB
}

func NewUserRepositoryLevelDB() UserRepository {
	rep := &UserRepositoryLeveldb{}
	rep.init()
	return rep
}

func (r *UserRepositoryLeveldb) init() {
	r.openCfgDB()
	//defer r.cfgDB.Close()
	r.openUserDB()
	//defer r.userDB.Close()
}

func (r *UserRepositoryLeveldb) validateUser(username string, password string) error {
	//r.openUserDB()
	//defer r.userDB.Close()
	user, ok := r.getUser(username)
	if !ok {
		return ErrInvalidLogin
	}
	if strings.Compare(password, user.Password) == 0 {
		return nil
	}
	return ErrInvalidLogin
}

func (r *UserRepositoryLeveldb) deleteUser(username string) {
	//r.openUserDB()
	//defer r.userDB.Close()
	err := r.userDB.Delete([]byte(username), nil)
	if err != nil {
		log.Error("can not delete user", zap.String("username", username), zap.Error(err))
	} else {
		log.Debug("user deleted", zap.String("username", username))
	}
}

func (r *UserRepositoryLeveldb) getRepositoryName() string {
	//r.openCfgDB()
	//defer r.cfgDB.Close()
	nameByte, err := r.cfgDB.Get([]byte("name"), nil)
	if err != nil {
		log.Error("can not get repository name", zap.Error(err))
		return ""
	}
	return string(nameByte)
}

func (r *UserRepositoryLeveldb) getUser(username string) (*User, bool) {
	//r.openUserDB()
	//defer r.userDB.Close()
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

func (r *UserRepositoryLeveldb) openCfgDB() {
	var err error
	r.cfgDB, err = leveldb.OpenFile("./rep/user_cfg", nil)
	if err != nil {
		log.Error("can not open user config repository", zap.Error(err))
	}
}

func (r *UserRepositoryLeveldb) openUserDB() {
	var err error
	r.userDB, err = leveldb.OpenFile("./rep/user", nil)
	if err != nil {
		log.Error("can not open user repository", zap.Error(err))
	}
}

func (r *UserRepositoryLeveldb) setRepositoryName(name string) {
	//r.openCfgDB()
	//defer r.cfgDB.Close()
	err := r.cfgDB.Put([]byte("name"), []byte(name), nil)
	if err != nil {
		log.Error("can not set repository name", zap.String("repository", name), zap.Error(err))
	} else {
		log.Debug("repository name set", zap.String("repository", name))
	}
}

func (r *UserRepositoryLeveldb) setUser(user *User) {
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
		log.Debug("user updated", zap.String("user id", user.ID.String()), zap.String("username", user.UserName))
	}
}
