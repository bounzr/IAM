package repository

import (
	"go.uber.org/zap"
	"strings"
)

//UserRepositoryBasic In memory repository. Used for tests
type UserRepositoryBasic struct {
	name  string
	users map[string]*User
}

//init the repository
func (r *UserRepositoryBasic) init() {
	r.users = make(map[string]*User)
}

//setUser register a new user
func (r *UserRepositoryBasic) setUser(user *User) {
	//add the new user to the repo
	r.users[user.UserName] = user
	return
}

//validateUser if username matches password returns nil, otherwise returns error
func (r *UserRepositoryBasic) validateUser(username string, password string) error {
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

func (r *UserRepositoryBasic) deleteUser(username string) {
	delete(r.users, username)
}

func (r *UserRepositoryBasic) getRepositoryName() string {
	return r.name
}

//getUser get user by username
func (r *UserRepositoryBasic) getUser(username string) (*User, bool) {
	//an user always has an user structure
	usr, ok := r.users[username]
	if !ok {
		if strings.Compare(username, "admin") != 0 {
			log.Debug("username is not in repository", zap.String("username", username))
		}
		return nil, false
	}
	return usr, ok
}

func (r *UserRepositoryBasic) setRepositoryName(name string) {
	r.name = name
}
