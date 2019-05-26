package repository

import "../utils"

type SessionToken [32]byte

type SessionRepository interface {
	addSessionContext(token SessionToken, user *UserCtx)
	deleteSessionContext(token SessionToken) error
	getSessionContext(token SessionToken) (*UserCtx, error)
	init()
}

//DeleteSessionUser Delete the caches related to the given session token if the token is correct
func DeleteSessionUser(token SessionToken) error {
	return sessionManager.deleteSessionContext(token)
}

//GetSessionUser Returns the user model stored in the user cache for the given id.
func GetSessionUser(token SessionToken) (*UserCtx, error) {
	return sessionManager.getSessionContext(token)
}

//NewSessionToken generate session token for user ctx
func NewSessionToken(user *UserCtx) SessionToken {
	token := utils.GetRandom32Token()
	sessionManager.addSessionContext(token, user)
	return token
}
