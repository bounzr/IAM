package router

import (
	"bounzr/iam/pages"
	"bounzr/iam/repository"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

//newBounzrRouter returns a new router with Bounzr basic endpoints
func newBounzrRouter(router *mux.Router) {
	router.HandleFunc("/", chain(indexPageGetHandler, sessionCookieSecurity)).Methods(http.MethodGet)
	router.HandleFunc("/login", loginPageHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/logout", logoutPageGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/register", registerPageHandler).Methods(http.MethodGet, http.MethodPost)
}

func indexPageGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usr, ok := fromContextGetUser(ctx)
	if !ok {
		//TODO redirect to login instead of error
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		log.Debug("can not get user from context", zap.Error(repository.ErrInvalidLogin))
		return
	}
	user, found := repository.GetUser(usr.GetUserID())
	if !found {
		log.Debug("user not found", zap.String("user id", usr.GetUserID().String()))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	data := user.GetScim()
	err := pages.RenderPage(w, "index", data)
	if err != nil {
		log.Error("can not render index webpage", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginPageGetHandler(w, r)
	}
	if r.Method == http.MethodPost {
		loginPagePostHandler(w, r)
	}
}

/**
if already logged in, user is sent to main page. Else will ask for username password
*/
func loginPageGetHandler(w http.ResponseWriter, r *http.Request) {
	_, err := validateLoginSession(w, r)
	if err == nil {
		landLoginRequest(w, r)
		return
	}
	log.Error("can not validate login session", zap.Error(err))
	err = pages.RenderPage(w, "login", nil)
	if err != nil {
		log.Error("can not render login webpage", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//landLoginRequest redirects to the previous requested url before being forced to login
func landLoginRequest(w http.ResponseWriter, r *http.Request) {
	url, err := getTargetURLFromSession(w, r)
	if err != nil {
		log.Error("can not get target url", zap.String("url", url), zap.Error(err))
	}
	http.Redirect(w, r, url, 302)
}

func loginPagePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	user, valid := authenticateUser(username, password)
	if !valid {
		log.Debug("user authentication not valid", zap.String("username", username))
		pages.RenderPage(w, "login", repository.ErrInvalidLogin.Error())
		return
	}
	sessionToken := repository.NewSessionToken(user)
	session, _ := BounzrCookieStore.Get(r, SessionCookie)
	session.Values[UserSessionToken] = sessionToken
	session.Save(r, w)
	landLoginRequest(w, r)
}

func logoutPageGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := BounzrCookieStore.Get(r, SessionCookie)
	//remove session from session store
	deleteSession(r)
	//remove session from cookie
	delete(session.Values, UserSessionToken)
	session.Save(r, w)
	landLoginRequest(w, r)
}

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		registerPageGetHandler(w, r)
	}
	if r.Method == http.MethodPost {
		registerPagePostHandler(w, r)
	}
}

func registerPageGetHandler(w http.ResponseWriter, r *http.Request) {
	err := pages.RenderPage(w, "signup", nil)
	if err != nil {
		log.Error("can not render signup template", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//todo replace AddAdminUser with addUserScim
func registerPagePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	//todo use scim
	err := repository.AddAdminUser("main", username, password)
	if err != nil {
		log.Error("can not add technical user", zap.String("username", username), zap.Error(err))
		pages.RenderPage(w, "signup", err.Error())
		return
	}
	http.Redirect(w, r, "/bounzr/login", 302)
}
