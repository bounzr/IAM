package repository

import (
	"../config"
	"../logger"
	"encoding/gob"
	"go.uber.org/zap"
)

var (
	log            *zap.Logger
	clientManager  ClientManager
	groupManager   GroupManager
	sessionManager SessionManager
	tokenManager   TokenManager
)

//init repositories Manager
func Init() {
	log = logger.GetLogger()
	gob.Register(&SessionToken{})
	gob.Register(&ConsentToken{})
	gob.Register(&ResourceTag{})

	initClients()
	initUsers()
	initTokens()
	initSessions()
	initGroups()

	adminGroup, err := GetGroup(privateGroups["Admins"])
	if err != nil {
		log.Error("can not get admin group", zap.Error(err))
		panic("admin group is required to run")
	}
	if len(adminGroup.Members) == 0 {
		_, found := GetUser(config.IAM.Users.Admin.Username)
		if !found {
			log.Debug("no admin user found, new one will be created")
			err = AddAdminUser(config.IAM.Users.Admin.Repository, config.IAM.Users.Admin.Username, config.IAM.Users.Admin.Password)
			if err != nil {
				log.Error("can not add admin user", zap.String("username", config.IAM.Users.Admin.Username), zap.Error(err))
			}
			_, found = GetUser(config.IAM.Users.Admin.Username)
			if !found {
				log.Error("no admin user found", zap.String("username", config.IAM.Users.Admin.Username), zap.Error(ErrUsernameNotFound))
			}
		}
	}
}
