package repository

import (
	"../config"
	"../logger"
	"go.uber.org/zap"
)

var (
	log					*zap.Logger
	clientManager		ClientRepository
	groupManager		GroupRepository
	resourceManager 	resourceRepository
	sessionManager  	SessionRepository
	tokenManager    	TokenRepository
	userRepositories 	map[string]UserRepository
)

//init repositories manager
func Init() {
	log = logger.GetLogger()
	initClients()
	initUsers()
	initTokens()
	initSessions()
	initResources()
	initGroups()

	_, err := GetUser(config.IAM.Users.Admin.Username);
	if  err != nil {
		log.Error("get user", zap.Error(err))
		err = AddTechnicalUser(config.IAM.Users.Admin.Repository, config.IAM.Users.Admin.Username, config.IAM.Users.Admin.Password)
	}
	if err != nil {
		log.Error("add technical user", zap.Error(err))
	}
}

func initClients() {
	implementation := config.IAM.Clients.Implementation
	switch implementation {
	case "leveldb":
		clientManager = NewClientRepositoryLeveldb()
	default:
		clientManager = NewClientRepositoryBasic()
	}
}

func initGroups(){
	implementation := config.IAM.Groups.Implementation
	switch implementation {
	case "leveldb":
		groupManager = NewGroupRepositoryLeveldb()
	default:
		groupManager = NewGroupRepositoryBasic()
	}
}

func initResources() {
	implementation := config.IAM.Resources.Implementation
	switch implementation {
	case "leveldb":
		resourceManager = NewResourceRepositoryLeveldb()
	default:
		resourceManager = NewResourceRepositoryBasic()
	}
}

func initSessions() {
	implementation := config.IAM.Sessions.Implementation
	switch implementation {
	case "leveldb":
		sessionManager = NewSessionRepositoryLeveldb()
	default:
		sessionManager = NewSessionRepository()
	}
}

func initTokens() {
	implementation := config.IAM.Tokens.Implementation
	switch implementation {
	case "leveldb":
		tokenManager = NewTokenRepositoryLeveldb()
	default:
		tokenManager = NewTokenRepositoryBasic()
	}
}

func initUsers() {
	userRepositories = make(map[string]UserRepository)

	//TODO different user repos management
	//todo implementation
	name := "main"
	userRepo := NewUserRepository(name)
	addUserRepository(userRepo, name)
}

func addUserRepository(ur UserRepository, id string) error {
	//verify that repository already exists
	if _, ok := userRepositories[id]; ok {
		log.Error("add user repository", zap.Error(ErrRepositoryNotAvailable))
		return ErrRepositoryNotAvailable
	}
	userRepositories[id] = ur
	return nil
}
