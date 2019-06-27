package repository

import (
	"../oauth2"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
)

//UserManager contains user and profile information
type UserManager interface {
	close()
	deleteUser(username string)
	findUsers() ([]User, error)
	getRepositoryName() string
	getUser(username string) (*User, bool)
	init()
	setRepositoryName(name string)
	setUser(user *User)
	validateUser(username string, password string) error
}

var userRepositories map[string]UserManager

func initUsers() {
	userRepositories = make(map[string]UserManager)

	//TODO different user repos management
	//todo implementation
	name := "main"
	userRepo := NewUserManager(name)
	addUserRepository(userRepo, name)
}

func addUserRepository(ur UserManager, id string) error {
	//verify that repository already exists
	if _, ok := userRepositories[id]; ok {
		log.Error("add user repository", zap.Error(ErrRepositoryNotAvailable))
		return ErrRepositoryNotAvailable
	}
	userRepositories[id] = ur
	return nil
}

func NewUserManager(name string) UserManager {
	/*
		repo := &UserManagerBasic{
			name: name,
		}*/
	repo := &UserManagerLeveldb{
		cfgPath:  "./rep/user_cfg",
		userPath: "./rep/user",
	}
	repo.init()
	repo.setRepositoryName(name)
	return repo
}

func AddScimUser(scimUser *scim2.User) error {
	repository := "main" //todo config repository
	users, err := getUserRepository(repository)
	if err != nil {
		log.Error("can not get user", zap.String("username", scimUser.UserName), zap.Error(err))
		return err
	}
	username := strings.ToLower(scimUser.UserName)
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not get uuid", zap.String("username", scimUser.UserName), zap.Error(err))
		return err
	}
	password := scimUser.Password
	if len(password) == 0 {
		log.Error("password is empty", zap.String("username", scimUser.UserName), zap.Error(scim2.ErrBadRequestInvalidValue))
		return scim2.ErrBadRequestInvalidValue
	}
	user := NewUser(id, username, password, repository)
	user.setScim(scimUser)
	users.setUser(user)
	AddResource(user)
	return nil
}

//addUser adds new user to defined repository
func AddTechnicalUser(repository string, username string, password string) error {
	users, err := getUserRepository(repository)
	if err != nil {
		log.Error("can not add technical user", zap.String("username", username), zap.Error(err))
		return err
	}
	username = strings.ToLower(username)
	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not get uuid", zap.String("username", username), zap.Error(err))
		return err
	}
	user := NewUser(id, username, password, repository)
	users.setUser(user)
	return nil
}

//todo find users(filter, attributes)
func FindUsers() []scim2.User {
	var users []scim2.User
	for repName, repo := range userRepositories {
		repUsers, err := repo.findUsers()
		if err != nil {
			log.Error("can not add users from repository", zap.String("repository", repName), zap.Error(err))
			continue
		}
		for _, user := range repUsers {
			users = append(users, *user.GetScim())
		}
	}
	return users
}

//GetAuthorizationRequest returns an authorization request from a user for a client by using the corresponding consent token
func GetAuthorizationRequest(user *UserCtx, consentToken *ConsentToken) (authorizationRequest *oauth2.AuthorizationRequest, err error) {
	rep, err := getUserRepository(user.RepositoryName)
	if err != nil {
		log.Error("can not get user repository", zap.String("username", user.UserName), zap.String("repository", user.RepositoryName))
		return nil, err
	}
	usr, ok := rep.getUser(user.UserName)
	if !ok {
		log.Error("user not found", zap.String("username", user.UserName), zap.Error(ErrUsernameNotFound))
		return nil, ErrSessionNotFound
	}
	authorizationRequest, ok = usr.getClientAuthorizationRequest(consentToken)
	if ok {

		return authorizationRequest, nil
	}
	log.Error("client authorization request not found for client", zap.String("client ID", consentToken.ClientID.String()), zap.Error(ErrSessionNotFound))
	return nil, ErrSessionNotFound
}

//getUserRepository returns the repository with the given name
func getUserRepository(repName string) (UserManager, error) {
	repo, ok := userRepositories[repName]
	if !ok {
		log.Error("can not find repository", zap.String("repository", repName), zap.Error(ErrRepositoryNotAvailable))
		return nil, ErrRepositoryNotAvailable
	}
	return repo, nil
}

//getUser searches username in all repositories and return a struct for a matching user
func GetUser(username string) (*User, error) {
	if len(username) == 0 {
		return nil, ErrUsernameNotFound
	}
	for _, store := range userRepositories {
		user, ok := store.getUser(username)
		if ok {
			return user, nil
		}
	}
	return nil, ErrUsernameNotFound
}

//SetAuthorizationRequest saves the client authorization request into the corresponding user
func SetAuthorizationRequest(user *UserCtx, authorizationRequest *oauth2.AuthorizationRequest) (consentToken *ConsentToken, err error) {
	consentToken = newConsentToken(uuid.FromStringOrNil(authorizationRequest.ClientID))
	rep, err := getUserRepository(user.RepositoryName)
	if err != nil {
		return nil, err
	}
	usr, ok := rep.getUser(user.UserName)
	if !ok {
		return nil, ErrUsernameNotFound
	}
	usr.setClientAuthorizationRequest(consentToken, authorizationRequest)
	rep.setUser(usr)
	return consentToken, nil
}

//ValidateUser validates an user in the specified repository
func ValidateUser(username string, password string) (*UserCtx, error) {
	user, err := GetUser(username)
	if err != nil {
		return nil, err
	}
	us, err := getUserRepository(user.RepositoryName)
	if err != nil {
		return nil, err
	}
	err = us.validateUser(user.UserName, password)
	if err != nil {
		return nil, err
	}
	return user.GetUserCtx(), nil
}
