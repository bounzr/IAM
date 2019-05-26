package router

import (
	"github.com/gofrs/uuid"
	"net/http"

	"../oauth2"
	"../repository"
	"go.uber.org/zap"
)

func authenticateClient(clientID, clientSecret string) (clientCtx *oauth2.ClientCtx, ok bool) {
	client, found := repository.GetClient(uuid.FromStringOrNil(clientID))
	if found {
		if ok = client.ValidateClientSecret(clientSecret); ok {
			return client.GetClientCtx(), ok
		}
	}
	log.Debug("client context not found in repository")
	return nil, false
}

func authenticateUser(username, password string) (*repository.UserCtx, error) {
	usrCtx, err := repository.ValidateUser(username, password)
	if err == nil {
		return usrCtx, err
	}
	log.Error("can not validate user", zap.String("username", username), zap.Error(err))
	return nil, err
}

func validateLoginSession(w http.ResponseWriter, r *http.Request) (user *repository.UserCtx, error error) {
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	if err != nil {
		log.Error("can not retrieve session cookie from repository", zap.Error(err))
		return nil, err
	}
	ust := session.Values[UserSessionToken]
	if ust == nil {
		log.Error("can not retrieve user session token", zap.Error(repository.ErrSessionNotFound))
		return nil, repository.ErrSessionNotFound
	}
	var token = &repository.SessionToken{}
	token, ok := ust.(*repository.SessionToken)
	if !ok {
		log.Error("can not cast session token", zap.String("user id", user.GetUserID().String()), zap.Error(repository.ErrSessionNotFound))
		return nil, repository.ErrSessionNotFound
	}
	usr, error := repository.GetSessionUser(*token)
	if error != nil {
		log.Error("can not retrieve user session", zap.Error(error))
		delete(session.Values, UserSessionToken)
		session.Save(r, w)
	}
	return usr, error
}

//TODO remove the session
func deleteSession(r *http.Request) error {
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	if err != nil {
		log.Error("can not retrieve session cookie", zap.Error(err))
		return err
	}
	ust := session.Values[UserSessionToken]
	var token = &repository.SessionToken{}
	token, ok := ust.(*repository.SessionToken)
	if ok {
		log.Debug("can not cast session token")
		err = repository.DeleteSessionUser(*token)
		if err != nil {
			log.Error("can not delete session user", zap.Error(err))
			return err
		}
	}
	return nil
}

func addTargetURLToSession(w http.ResponseWriter, r *http.Request) error {
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	if err != nil && session == nil {
		log.Error("session nil and can not retrieve session cookie", zap.Error(err))
		return err
	}
	//log.Debugf("target url: %s", r.URL.String())
	session.Values[TargetUrl] = r.URL.String()
	session.Save(r, w)
	return nil
}

func getTargetURLFromSession(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := BounzrCookieStore.Get(r, SessionCookie)
	if err != nil {
		log.Error("can not retrieve session cookie", zap.Error(err))
		return "/bounzr", err
	}
	turl := session.Values[TargetUrl]
	var target string
	target, ok := turl.(string)
	if !ok {
		log.Debug("can not cast target url. Default will be returned")
		return "/bounzr", repository.ErrSessionInvalid
	}
	return target, nil
}
