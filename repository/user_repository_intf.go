package repository

import (
	"../oauth2"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
)

//UserRepository contains user and profile information
type UserRepository interface {
	init()
	setUser(user *User)
	validateUser(username string, password string) error
	deleteUser(username string)
	getRepositoryName() string
	getUser(username string) (*User, bool)
	setRepositoryName(name string)
}

func NewUserRepository(name string) UserRepository {
	/*
		repo := &UserRepositoryBasic{
			name: name,
		}*/
	repo := NewUserRepositoryLevelDB()
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
	user := NewUser(id, username, password, repository, scimUser)
	users.setUser(user)
	resourceManager.saveResource(user.GetResourceMetadata())
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
	user := NewUser(id, username, password, repository, nil)
	users.setUser(user)
	//save the unique resource
	resourceManager.saveResource(user.GetResourceMetadata())
	return nil
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
	log.Error("client authorization request not found for client", zap.String("client id", consentToken.ClientID.String()), zap.Error(ErrSessionNotFound))
	return nil, ErrSessionNotFound
}

//getUserRepository returns the repository with the given name
func getUserRepository(repName string) (UserRepository, error) {
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
