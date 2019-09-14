package repository

import (
	"bounzr/iam/config"
	"bounzr/iam/utils"
)

type SessionToken [32]byte

type SessionManager interface {
	setSessionContext(token SessionToken, user *UserCtx)
	deleteSessionContext(token SessionToken) error
	getSessionContext(token SessionToken) (*UserCtx, error)
	init()
	close()
}

func initSessions() {
	implementation := config.IAM.Sessions.Implementation
	switch implementation {
	case "leveldb":
		sessionManager = &SessionManagerLeveldb{sessionsPath: "./rep/session"}
	default:
		sessionManager = &SessionManagerBasic{}
	}
	sessionManager.init()
}

//DeleteSessionUser Delete the caches related to the given session token if the token is correct
func DeleteSessionUser(token SessionToken) error {
	return sessionManager.deleteSessionContext(token)
}

//GetSessionUser Returns the user model stored in the user cache for the given ID.
func GetSessionUser(token SessionToken) (*UserCtx, error) {
	return sessionManager.getSessionContext(token)
}

//NewSessionToken generate session token for user ctx
func NewSessionToken(user *UserCtx) SessionToken {
	token := utils.GetRandom32Token()
	sessionManager.setSessionContext(token, user)
	return token
}
