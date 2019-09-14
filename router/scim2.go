package router

import (
	"bounzr/iam/repository"
	"bounzr/iam/scim2"
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
		verifyUserGroups("Admins", "Clients"))).Methods(
		http.MethodGet,
		http.MethodPost,
	)
	router.HandleFunc("/clients/{id:[-a-zA-Z0-9]+", chain(
		clientHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins", "Clients"))).Methods(
		http.MethodDelete,
		http.MethodGet,
		http.MethodPatch,
		http.MethodPut,
	)
	router.HandleFunc("/users", chain(
		usersHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins"))).Methods(
		http.MethodGet,
		http.MethodPost,
	)
	router.HandleFunc("/users/{id:[-a-zA-Z0-9]+}", chain(
		userHandler,
		basicUserAuthSecurity)).Methods(
		http.MethodGet,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodPut,
	)
	router.HandleFunc("/groups", chain(
		groupsHandler,
		basicUserAuthSecurity,
		verifyUserGroups("Admins"))).Methods(
		http.MethodDelete,
		http.MethodGet,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
	)
	router.HandleFunc("/groups/{id:[-a-zA-Z0-9]+}", chain(groupGet, basicUserAuthSecurity)).Methods(http.MethodGet)
}

func clientDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Debug("wrong user id", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	repository.DeleteClient(uid)
	w.WriteHeader(http.StatusNoContent)
	return
}

func clientHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		clientDelete(w, r)
	case http.MethodGet:
		clientGet(w, r)
	case http.MethodPut:
		clientPut(w, r)
	default:
		notImplemented(w, r)
	}
}

func clientGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	client, found := repository.GetClient(uuid.FromStringOrNil(id))
	if !found {
		log.Debug("client not found", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	resource := client.GetScim()
	//Marshal clientInfResp to json and write to response
	resourceJson, err := json.Marshal(resource)
	if err != nil {
		log.Error("can not marshal resource json", zap.String("id", resource.ID), zap.Error(scim2.ErrInternalError))
		http.Error(w, repository.ErrResourceNotAvailable.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	w.Write(resourceJson)
	return
}

func clientPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Debug("wrong user id", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	cscim := &scim2.Client{}
	err = json.NewDecoder(r.Body).Decode(cscim)
	if err != nil {
		log.Error("can not decode client", zap.String("client id", cscim.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	repository.ReplaceClientByScim(uid, cscim)

	client, found := repository.GetClient(uid)
	if !found {
		log.Debug("client not found", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
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

func clientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		clientsGet(w, r)
	}
	if r.Method == http.MethodPost {
		clientsPost(w, r)
	}
}

func clientsGet(w http.ResponseWriter, r *http.Request) {
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

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func clientsPost(w http.ResponseWriter, r *http.Request) {
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
	//todo if SetClient() err != nil
	repository.SetClient(client)
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

func groupGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsGet(w http.ResponseWriter, r *http.Request) {
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

func groupsPatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsPost(w http.ResponseWriter, r *http.Request) {
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

func groupsPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		groupsDelete(w, r)
	}
	if r.Method == http.MethodGet {
		groupsGet(w, r)
	}
	if r.Method == http.MethodPatch {
		groupsPatch(w, r)
	}
	if r.Method == http.MethodPost {
		groupsPost(w, r)
	}
	if r.Method == http.MethodPut {
		groupsPut(w, r)
	}
}

//To retrieve a known resource, clients send GET requests to the resource endpoint, e.g., "/Users/{id}", "/Groups/{id}", or
//"/Schemas/{id}", where "{id}" is a resource identifier (for example, the value of the "id" attribute).
//If the resource exists, the server responds with HTTP status code 200 (OK) and includes the result in the body of the response.
func userGet(w http.ResponseWriter, r *http.Request) {
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
		//HTTP/1.1 404 Not Found
		//{
		//"schemas": ["urn:ietf:params:scim:api:messages:2.0:Error"],
		//"detail":"Resource 2819c223-7f76-453a-919d-413861904646 not found",
		//"status": "404"
		//}
		//Set Content-Type header so that clients will know how to read response
		w.Header().Set("Content-Type", "application/scim+json")
		w.WriteHeader(http.StatusNotFound)
		//todo write correct json error
		w.Write([]byte("{\"detail\":\"Resource not found\"}"))
		return
	}
	resource := user.GetScim()
	//Marshal clientInfResp to json and write to response
	resourceJson, err := json.Marshal(resource)
	if err != nil {
		log.Error("can not marshal resource json", zap.String("id", resource.ID), zap.Error(scim2.ErrInternalError))
		http.Error(w, repository.ErrResourceNotAvailable.Error(), http.StatusInternalServerError)
		return
	}
	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/scim+json")
	w.WriteHeader(http.StatusOK)
	w.Write(resourceJson)
	return
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	//get user from context as it has logged in by the middleware
	ctx := r.Context()
	usr, ok := fromContextGetUser(ctx)
	if !ok {
		log.Error("can not get user from context", zap.Error(scim2.ErrUnauthorized))
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]

	if strings.Compare(usr.GetUserID().String(), strings.TrimSpace(id)) == 0 {
		log.Debug("this method can not be used for self delete")
		http.Error(w, repository.ErrInvalidLogin.Error(), http.StatusForbidden)
		return
	}
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Debug("wrong user id", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	repository.DeleteUser(uid)
	w.WriteHeader(http.StatusNoContent)
	return
}

func usersGet(w http.ResponseWriter, r *http.Request) {
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

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		userDelete(w, r)
	}
	if r.Method == http.MethodGet {
		userGet(w, r)
	}
	if r.Method == http.MethodPatch {
		userPatch(w, r)
	}
	if r.Method == http.MethodPut {
		userPut(w, r)
	}
	return
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		usersGet(w, r)
	}
	if r.Method == http.MethodPost {
		usersPost(w, r)
	}
	return
}

func userPatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"request\":\"not implemented\"}"))
	w.WriteHeader(http.StatusNotImplemented)
}

func usersPost(w http.ResponseWriter, r *http.Request) {
	userReq := &scim2.User{}
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		log.Error("can not decode user", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	userID, err := repository.AddScimUser(userReq)
	if err != nil {
		log.Error("can not add user from scim", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	user, found := repository.GetUser(userID)
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

func userPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uid, err := uuid.FromString(id)
	if err != nil {
		log.Debug("wrong user id", zap.String("id", id))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	userReq := &scim2.User{}
	err = json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		log.Error("can not decode user", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	err = repository.ReplaceUserByScim(uid, userReq)
	if err != nil {
		log.Error("can not add user from scim", zap.String("user id", userReq.ID), zap.Error(err))
		http.Error(w, repository.ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}
	user, found := repository.GetUser(uid)
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
