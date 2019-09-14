package router

import (
	"bounzr/iam/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

const (
	ConsentsToken    = "consents_token"
	SessionCookie    = "session"
	TargetUrl        = "_target_url"
	UserSessionToken = "bounzr_token"
)

//TODO authKey and encryptKey in the shared registry available for other nodes?
var (
	log               *zap.Logger
	authKey           = securecookie.GenerateRandomKey(64)
	encryptKey        = securecookie.GenerateRandomKey(32)
	BounzrCookieStore = sessions.NewCookieStore([]byte(authKey), []byte(encryptKey))
	decoder           = schema.NewDecoder()
)

func Init() {
	log = logger.GetLogger()
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	//flexible slash path
	r.StrictSlash(true)
	//static files
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	//bounzr subrouter
	bounzrSubRouter := r.PathPrefix("/bounzr").Subrouter()
	newBounzrRouter(bounzrSubRouter)
	//oauth2 subrouter
	oauth2SubRouter := r.PathPrefix("/oauth2").Subrouter()
	newOauth2Router(oauth2SubRouter)
	//scim2 subrouter
	scim2SubRouter := r.PathPrefix("/scim2").Subrouter()
	newScim2Router(scim2SubRouter)

	return r
}
