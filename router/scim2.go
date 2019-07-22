package router

import (
	"../repository"
	"../scim2"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func newScim2Router(router *mux.Router) {
	router.HandleFunc("/clients", chain(
		clientsHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins", "Clients"))).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users", chain(
		usersHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins"))).Methods(
		http.MethodDelete,
		http.MethodGet,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut)
	router.HandleFunc("/users/{id:[-a-zA-Z0-9]+}", chain(userGetHandler, basicUserAuthSecurity)).Methods(http.MethodGet)
	router.HandleFunc("/groups", chain(
		groupsHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins"))).Methods(
		http.MethodDelete,
		http.MethodGet,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut)
	router.HandleFunc("/groups/{id:[-a-zA-Z0-9]+}", chain(groupGetHandler, basicUserAuthSecurity)).Methods(http.MethodGet)
}

func clientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		clientsHandlerGet(w, r)
	}
	if r.Method == http.MethodPost {
		clientsHandlerPost(w, r)
	}
}

func clientsHandlerGet(w http.ResponseWriter, r *http.Request) {
	clients := repository.FindClients()
	schema := []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}
	totalResults := len(clients)
	listResponse := scim2.ResourceQueryResponse{
		Schemas:      schema,
		TotalResults: totalResults,
	}
	if totalResults > 0 {
		listResponse.Resources = clients
	}
	//Marshal to json and write to response
	jsonResponse, err := json.Marshal(listResponse)
	if err != nil {
		log.Error("can not marshal users json", zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(jsonResponse)
	return
}

func clientsHandlerPost(w http.ResponseWriter, r *http.Request) {
	clientReq := &scim2.Client{}
	err := json.NewDecoder(r.Body).Decode(clientReq)
	if err != nil {
		log.Error("can not decode client", zap.String("client id", clientReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	client, err := repository.NewClientFromScim(clientReq)
	if err != nil {
		log.Error("can not add user from scim", zap.String("user id", clientReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusBadRequest)
		return
	}
	//TODO clean context() after it was used leaving no trace
	//TODO use better func fromContextGetUser(ctx context.Context) (*repository.User, bool)
	if loggedUser := r.Context().Value(userCtxKey); loggedUser != nil {
		//todo add or handle anonymous scope to user and client
		profile := loggedUser.(*repository.UserCtx)
		log.Debug("found logged user context", zap.String("user id", profile.GetUserID().String()))
		client.OwnerID = profile.UserID
	}
	//todo if AddClient() err != nil
	repository.AddClient(client)
	repository.SetResourceGroups(clientReq, client.GetResourceTag())
	scim := client.GetScim()
	//Marshal clientInfResp to json and write to response
	scimJson, err := json.Marshal(scim)
	if err != nil {
		log.Error("can not marshal scim json", zap.String("user id", scim.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	w.Write(scimJson)
	return
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		usersDeleteHandler(w, r)
	}
	if r.Method == http.MethodGet {
		usersGetHandler(w, r)
	}
	if r.Method == http.MethodPatch {
		usersPatchHandler(w, r)
	}
	if r.Method == http.MethodPost {
		usersPostHandler(w, r)
	}
	if r.Method == http.MethodPut {
		usersPutHandler(w, r)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		groupsDeleteHandler(w, r)
	}
	if r.Method == http.MethodGet {
		groupsGetHandler(w, r)
	}
	if r.Method == http.MethodPatch {
		groupsPatchHandler(w, r)
	}
	if r.Method == http.MethodPost {
		groupsPostHandler(w, r)
	}
	if r.Method == http.MethodPut {
		groupsPutHandler(w, r)
	}
}

func usersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func usersGetHandler(w http.ResponseWriter, r *http.Request) {
	users := repository.FindUsers()
	schema := []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}
	totalResults := len(users)
	listResponse := scim2.ResourceQueryResponse{
		Schemas:      schema,
		TotalResults: totalResults,
	}
	if totalResults > 0 {
		listResponse.Resources = users
	}
	//Marshal to json and write to response
	jsonResponse, err := json.Marshal(listResponse)
	if err != nil {
		log.Error("can not marshal users json", zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	w.Write(jsonResponse)
	return
}

//To retrieve a known resource, clients send GET requests to the resource endpoint, e.g., "/Users/{id}", "/Groups/{id}", or
//"/Schemas/{id}", where "{id}" is a resource identifier (for example, the value of the "id" attribute).
//If the resource exists, the server responds with HTTP status code 200 (OK) and includes the result in the body of the response.
func userGetHandler(w http.ResponseWriter, r *http.Request) {
	//get user from context as it has logged in by the middleware
	ctx := r.Context()
	usr, ok := fromContextGetUser(ctx)
	if !ok {
		log.Error("can not get user from context", zap.Error(scim2.ErrUnauthorized))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	userIsAdmin := repository.ValidateResourceInGroup(usr.GetUserID(), "Admins")
	vars := mux.Vars(r)
	id := vars["id"]
	if !userIsAdmin && strings.Compare(usr.UserID.String(), strings.TrimSpace(id)) != 0 {
		log.Debug("user not allowed to request the given resource id", zap.String("user id", usr.UserID.String()), zap.String("resource id", id))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	user, found := repository.GetUser(uuid.FromStringOrNil(id))
	if !found {
		log.Debug("user not found", zap.String("user id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	resource := user.GetScim()
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
	w.Write(resourceJson)
	return
}

func usersPatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func usersPostHandler(w http.ResponseWriter, r *http.Request) {
	userReq := &scim2.User{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		log.Error("can not decode user", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	err = repository.AddScimUser(userReq)
	if err != nil {
		log.Error("can not add user from scim", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	user, found := repository.GetUser(userReq.UserName)
	if !found {
		log.Error("can not get user", zap.String("user id", userReq.ID), zap.Error(repository.ErrUsernameNotFound))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	scim := user.GetScim()
	//Marshal clientInfResp to json and write to response
	scimJson, err := json.Marshal(scim)
	if err != nil {
		log.Error("can not marshal scim json", zap.String("user id", scim.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	w.Write(scimJson)
	return
}

func groupGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func usersPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsGetHandler(w http.ResponseWriter, r *http.Request) {
	filter := make(map[string]interface{})
	groups := repository.FindGroups(filter)
	schema := []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"}
	totalResults := len(groups)
	listResponse := scim2.ResourceQueryResponse{
		Schemas:      schema,
		TotalResults: totalResults,
	}
	if totalResults > 0 {
		listResponse.Resources = groups
	}
	//Marshal to json and write to response
	jsonResponse, err := json.Marshal(listResponse)
	if err != nil {
		log.Error("can not marshal users json", zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	return
}

func groupsPatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsPostHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Error("can not marshal group to scim json", zap.String("group id", scim.ID), zap.Error(err))
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

func groupsPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}
