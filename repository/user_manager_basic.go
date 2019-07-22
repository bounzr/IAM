package repository

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
)

//UserManagerBasic In memory repository. Used for tests
type UserManagerBasic struct {
	name  string
	users map[string]*User
	uids  map[uuid.UUID]string
}

//init the repository
func (r *UserManagerBasic) init() {
	r.users = make(map[string]*User)
	r.uids = make(map[uuid.UUID]string)
}

func (r *UserManagerBasic) close() {
	//nothing
}

//setUser register a new user
func (r *UserManagerBasic) setUser(user *User) {
	//add the new user to the repo
	r.users[user.UserName] = user
	r.uids[user.ID] = user.UserName
	return
}

//validateUser if username matches password returns nil, otherwise returns error
func (r *UserManagerBasic) validateUser(username string, password string) error {
	//all usernames handled in lowercase
	username = strings.ToLower(username)
	user := r.users[username]
	stdpwd := user.Password
	if strings.Compare(stdpwd, password) == 0 {
		return nil
	}
	return ErrInvalidLogin

	//TODO hash storage of password
	/*
		hash := s.Store[username]
		if len(hash) == 0{
			return ErrInvalidLogin
		}
		err := bcrypt.CompareHashAndPassword(hash, []byte(password))
		if err != nil{
			return ErrInvalidLogin
		}
		return nil
	*/
	//TODO still not hashing the password in memory.
}

func (r *UserManagerBasic) deleteUser(userID interface{}) {
	username := r.getUsername(userID)
	delete(r.users, username)
}

func (r *UserManagerBasic) getRepositoryName() string {
	return r.name
}

func (r *UserManagerBasic) findUsers() ([]User, error) {
	users := make([]User, len(r.users))
	idx := 0
	for _, user := range r.users {
		users[idx] = *user
		idx++
	}
	return users, nil
}

//getUser get user by username
func (r *UserManagerBasic) getUser(userID interface{}) (*User, bool) {
	//an user always has an user structure
	username := r.getUsername(userID)
	usr, ok := r.users[username]
	if !ok {
		if strings.Compare(username, "admin") != 0 {
			log.Debug("username is not in repository", zap.String("username", username))
		}
		return nil, false
	}
	return usr, ok
}

func (r *UserManagerBasic) getUsername(userID interface{}) string {
	switch userID.(type) {
	case string:
		log.Debug("userID is a string")
		return userID.(string)
	case uuid.UUID:
		return r.uids[userID.(uuid.UUID)]
	default:
		return ""
	}

}

func (r *UserManagerBasic) setRepositoryName(name string) {
	r.name = name
}
