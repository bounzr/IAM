package repository

import (
	"../config"
	"../logger"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

var (
	log             *zap.Logger
	clientManager   ClientManager
	groupManager    GroupManager
	resourceManager ResourceManager
	sessionManager  SessionManager
	tokenManager    TokenManager
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
	initResources()
	initGroups()

	adminGroupFilter := make(map[string]interface{})
	adminGroupFilter["name"] = "admin"
	groups, err := groupManager.findGroups(adminGroupFilter)
	if err != nil {
		log.Error("can not get groups", zap.Error(err))
	}
	var adminGroupID uuid.UUID
	if len(groups) == 0 {
		adminGroupID, err = AddGroup("admin")
		if err != nil {
			log.Error("can not get admin group", zap.Error(err))
			panic("admin group is required to run")
		}
	} else {
		adminGroupID = groups[0].Metadata.ID
	}
	adminGroup, err := GetGroup(adminGroupID)
	if err != nil {
		log.Error("can not get admin group", zap.Error(err))
		panic("admin group is required to run")
	}
	if len(adminGroup.Members) == 0 {
		var adminUser *User
		adminUser, err = GetUser(config.IAM.Users.Admin.Username)
		if err != nil {
			log.Debug("no admin user found, new one will be created", zap.Error(err))
			err = AddTechnicalUser(config.IAM.Users.Admin.Repository, config.IAM.Users.Admin.Username, config.IAM.Users.Admin.Password)
			if err != nil {
				log.Error("can not add admin user", zap.String("username", config.IAM.Users.Admin.Username), zap.Error(err))
			}
			adminUser, err = GetUser(config.IAM.Users.Admin.Username)
			if err != nil {
				log.Error("no admin user found", zap.String("username", config.IAM.Users.Admin.Username), zap.Error(err))
			}
		}
		AddGroupResource(adminGroupID, adminUser.Metadata)
	}
}
