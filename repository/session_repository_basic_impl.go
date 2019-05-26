package repository

type SessionRepositoryBasic struct {
	sessionCache map[SessionToken]*UserCtx
}

func NewSessionRepository() SessionRepository {
	nsr := &SessionRepositoryBasic{}
	nsr.init()
	return nsr
}

func (r *SessionRepositoryBasic) init() {
	r.sessionCache = make(map[SessionToken]*UserCtx)
}

func (r *SessionRepositoryBasic) addSessionContext(token SessionToken, user *UserCtx) {
	r.sessionCache[token] = user
}

func (r *SessionRepositoryBasic) deleteSessionContext(token SessionToken) error {
	_, ok := r.sessionCache[token]
	if ok {
		delete(r.sessionCache, token)
		return nil
	}
	return ErrSessionInvalid
}

func (r *SessionRepositoryBasic) getSessionContext(token SessionToken) (*UserCtx, error) {
	user, ok := r.sessionCache[token]
	if ok {
		return user, nil
	}
	return nil, ErrSessionNotFound
}
