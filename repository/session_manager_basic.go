package repository

type SessionManagerBasic struct {
	sessionCache map[SessionToken]*UserCtx
}

func (r *SessionManagerBasic) init() {
	r.sessionCache = make(map[SessionToken]*UserCtx)
}

func (r *SessionManagerBasic) close() {
	//nothing
}

func (r *SessionManagerBasic) setSessionContext(token SessionToken, user *UserCtx) {
	r.sessionCache[token] = user
}

func (r *SessionManagerBasic) deleteSessionContext(token SessionToken) error {
	_, ok := r.sessionCache[token]
	if ok {
		delete(r.sessionCache, token)
		return nil
	}
	return ErrSessionInvalid
}

func (r *SessionManagerBasic) getSessionContext(token SessionToken) (*UserCtx, error) {
	user, ok := r.sessionCache[token]
	if ok {
		return user, nil
	}
	return nil, ErrSessionNotFound
}
