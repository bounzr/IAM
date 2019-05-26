package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"../repository"
	"../scim2"
	"github.com/gofrs/uuid"
	"strings"
	"encoding/json"
	"go.uber.org/zap"
)

func newScim2Router(router *mux.Router) {
	router.HandleFunc("/users", middlewareChain(scim2UserDeleteHandler, basicUserAuthSecurity)).Methods("DELETE")
	router.HandleFunc("/users/{id:[-a-zA-Z0-9]+}", middlewareChain(scim2UserGetHandler, basicUserAuthSecurity)).Methods("GET")
	router.HandleFunc("/users", middlewareChain(scim2UserPatchHandler, basicUserAuthSecurity)).Methods("PATCH")
	router.HandleFunc("/users", middlewareChain(scim2UserPostHandler, basicUserAuthSecurity)).Methods("POST")
	router.HandleFunc("/users", middlewareChain(scim2UserPutHandler, basicUserAuthSecurity)).Methods("PUT")
	router.HandleFunc("/groups", middlewareChain(scim2GroupDeleteHandler, basicUserAuthSecurity)).Methods("DELETE")
	router.HandleFunc("/groups/{id:[-a-zA-Z0-9]+}", middlewareChain(scim2GroupGetHandler, basicUserAuthSecurity)).Methods("GET")
	router.HandleFunc("/groups", middlewareChain(scim2GroupPatchHandler, basicUserAuthSecurity)).Methods("PATCH")
	router.HandleFunc("/groups", middlewareChain(scim2GroupPostHandler, basicUserAuthSecurity)).Methods("POST")
	router.HandleFunc("/groups", middlewareChain(scim2GroupPutHandler, basicUserAuthSecurity)).Methods("PUT")

}

func scim2UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

//To retrieve a known resource, clients send GET requests to the resource endpoint, e.g., "/Users/{id}", "/Groups/{id}", or
//"/Schemas/{id}", where "{id}" is a resource identifier (for example, the value of the "id" attribute).
//If the resource exists, the server responds with HTTP status code 200 (OK) and includes the result in the body of the response.
func scim2UserGetHandler(w http.ResponseWriter, r *http.Request) {
	//get user from context as it has logged in by the middleware
	//todo validate if user is admin
	ctx := r.Context()
	usr, ok := fromContextGetUser(ctx)
	if !ok {
		log.Error("can not get user from context", zap.Error(scim2.ErrUnauthorized))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	//todo replace GetResourceMetadata with get user and user.getScim to match post request
	resource, err := repository.GetUserScim(uuid.FromStringOrNil(id))
	if err != nil {
		log.Error("can not get scim user resource", zap.String("user id", id), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	//todo an admin should be able to read the resource
	if strings.Compare(usr.UserID.String(), resource.ID) != 0 {
		log.Debug("user id does not match resource id", zap.String("user id", usr.UserID.String()), zap.String("resource id", resource.ID))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	//Marshal clientInfResp to json and write to response
	resourceJson, err := json.Marshal(resource)
	if err != nil {
		log.Error("can not marshal resource json", zap.String("resource id", resource.ID), zap.Error(scim2.ErrInternalError))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(resourceJson)
	return
}

func scim2UserPatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func scim2UserPostHandler(w http.ResponseWriter, r *http.Request) {
	userReq := &scim2.User{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		log.Error("can not decode user", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	err = repository.AddScimUser(userReq)
	if err != nil {
		log.Error("can not add user from scim", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	user, err := repository.GetUser(userReq.UserName)
	if err != nil {
		log.Error("can not get user", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	scim := user.GetScim()
	//Marshal clientInfResp to json and write to response
	scimJson, err := json.Marshal(scim)
	if err != nil {
		log.Error("can not marshal scim json", zap.String("user id",scim.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(scimJson)
	return
}

func scim2UserPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func scim2GroupDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func scim2GroupGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func scim2GroupPatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func scim2GroupPostHandler(w http.ResponseWriter, r *http.Request) {
	groupReq := &scim2.Group{}
	err := json.NewDecoder(r.Body).Decode(groupReq)
	if err != nil {
		log.Error("can not decode group json", zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	id, err := repository.AddScimGroup(groupReq)
	if err != nil {
		log.Error("can not add group from scim", zap.String("group id", groupReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	group, err := repository.GetGroup(id)
	if err != nil {
		log.Error("can not get group", zap.String("group id", id.String()), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	scim := group.GetScim()
	scimJson, err := json.Marshal(scim)
	if err != nil {
		log.Error("can not marshal group to scim json", zap.String("group id",scim.ID),zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(scimJson)
	return
}

func scim2GroupPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}