package router

import (
	"net/http"

	"../pages"
	"../repository"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//newBounzrRouter returns a new router with Bounzr basic endpoints
func newBounzrRouter(router *mux.Router) {
	router.HandleFunc("/", middlewareChain(indexPageGetHandler, sessionCookieSecurity)).Methods("GET")
	router.HandleFunc("/login", loginPageGetHandler).Methods("GET")
	router.HandleFunc("/login", loginPagePostHandler).Methods("POST")
	router.HandleFunc("/logout", logoutPageGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPageGetHandler).Methods("GET")
	router.HandleFunc("/register", registerPagePostHandler).Methods("POST")
}

func indexPageGetHandler(w http.ResponseWriter, r *http.Request) {
	err := pages.RenderPage(w, "index.html", nil)
	if err != nil {
		log.Error("can not render index.html", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	err = pages.RenderPage(w, "login.html", nil)
	if err != nil {
		log.Error("can not render login.html", zap.Error(err))
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
	user, err := authenticateUser(username, password)
	if err != nil {
		log.Error("can not authenticate user", zap.String("username", username), zap.Error(err))
		pages.RenderPage(w, "login.html", err.Error())
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

func registerPageGetHandler(w http.ResponseWriter, r *http.Request) {
	err := pages.RenderPage(w, "register.html", nil)
	if err != nil {
		log.Error("can not render register.html", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerPagePostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	//todo use scim
	err := repository.AddTechnicalUser("main", username, password)
	if err != nil {
		log.Error("can not add technical user", zap.String("username", username), zap.Error(err))
		pages.RenderPage(w, "register.html", err.Error())
		return
	}
	http.Redirect(w, r, "/bounzr/login", 302)
}
