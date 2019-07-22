package repository

import (
	"../oauth2"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"strings"
	"time"
)

//User information
type User struct {
	AccessTokens                       map[uuid.UUID]*oauth2.AccessTokenHint
	Attributes                         *UserAttributes
	AuthorizationRequests              map[uuid.UUID]*oauth2.AuthorizationRequest //[oauth_client} auth request
	AuthorizationRequestsConsentTokens map[uuid.UUID]*ConsentToken
	ID                                 uuid.UUID
	Metadata                           *ResourceTag
	Password                           string
	RefreshTokens                      map[uuid.UUID]*oauth2.AccessTokenHint
	RepositoryName                     string
	UserName                           string
}

type UserAttributes struct {
	Active            bool
	Addresses         []scim2.Address
	DisplayName       string
	Entitlements      []string
	Emails            []scim2.MultiValueAttribute
	ExternalId        string
	Ims               []scim2.MultiValueAttribute
	Locale            string
	Name              *scim2.Name
	NickName          string
	PhoneNumbers      []scim2.MultiValueAttribute
	Photos            []scim2.MultiValueAttribute
	PreferredLanguage string
	ProfileURL        string
	Roles             []string
	Timezone          string
	Title             string
	UserType          string
	X509Certificates  []scim2.MultiValueAttribute
}

type UserCtx struct {
	RepositoryName string
	UserID         uuid.UUID
	UserName       string
}

//NewUser creates new user
func NewUser(username string, password string, repository string) (*User, error) {
	atm := make(map[uuid.UUID]*oauth2.AccessTokenHint)
	arm := make(map[uuid.UUID]*oauth2.AuthorizationRequest)
	arctm := make(map[uuid.UUID]*ConsentToken)
	rtm := make(map[uuid.UUID]*oauth2.AccessTokenHint)
	currentTime := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		log.Error("can not get uuid for new user", zap.String("username", username), zap.Error(err))
		return nil, err
	}

	metadata := &ResourceTag{
		Created:        currentTime,
		ID:             id,
		LastModified:   currentTime,
		Name:           username,
		RepositoryName: repository,
		ResourceType:   "User",
	}

	var attributes *UserAttributes
	attributes = &UserAttributes{}

	user := &User{
		AccessTokens:                       atm,
		Attributes:                         attributes,
		AuthorizationRequests:              arm,
		AuthorizationRequestsConsentTokens: arctm,
		ID:                                 id,
		Metadata:                           metadata,
		Password:                           password,
		RefreshTokens:                      rtm,
		RepositoryName:                     repository,
		UserName:                           username,
	}

	return user, nil
}

func (u *User) DeleteClientAccessToken(clientID uuid.UUID) {
	delete(u.AccessTokens, clientID)
}

func (u *User) DeleteClientRefreshToken(clientID uuid.UUID) {
	delete(u.RefreshTokens, clientID)
}

func (u *User) DeleteClientTokens(clientID uuid.UUID) {
	delete(u.AccessTokens, clientID)
	delete(u.RefreshTokens, clientID)
}

func (u *User) GetClientAccessToken(clientID uuid.UUID) (tokenReference *oauth2.AccessTokenHint, ok bool) {
	if tokenReference, ok := u.AccessTokens[clientID]; ok {
		return tokenReference, true
	}
	return nil, false
}

func (u *User) GetClientRefreshToken(clientID uuid.UUID) (tokenReference *oauth2.AccessTokenHint, ok bool) {
	if len(clientID) == 0 {
		return nil, false
	}
	if refreshToken, ok := u.RefreshTokens[clientID]; ok {
		return refreshToken, true
	}
	return nil, false
}

func (u *User) GetScim() *scim2.User {

	memberGroupFilter := make(map[string]interface{})
	memberGroupFilter["member"] = u.ID
	groups := FindGroupAssignments(memberGroupFilter)

	user := &scim2.User{
		Active:            u.Attributes.Active,
		Addresses:         u.Attributes.Addresses,
		DisplayName:       u.Attributes.DisplayName,
		Entitlements:      u.Attributes.Entitlements,
		Emails:            u.Attributes.Emails,
		ExternalId:        u.Attributes.ExternalId,
		Groups:            groups,
		ID:                u.Metadata.ID.String(),
		Ims:               u.Attributes.Ims,
		Locale:            u.Attributes.Locale,
		Metadata:          u.Metadata.GetScimMetadata(),
		Name:              u.Attributes.Name,
		NickName:          u.Attributes.NickName,
		PhoneNumbers:      u.Attributes.PhoneNumbers,
		Photos:            u.Attributes.Photos,
		PreferredLanguage: u.Attributes.PreferredLanguage,
		ProfileURL:        u.Attributes.ProfileURL,
		Roles:             u.Attributes.Roles,
		Schemas:           []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		Timezone:          u.Attributes.Timezone,
		Title:             u.Attributes.Title,
		UserName:          u.Metadata.Name,
		UserType:          u.Attributes.UserType,
		X509Certificates:  u.Attributes.X509Certificates,
	}
	return user
}

func (u *User) GetUserCtx() *UserCtx {
	return &UserCtx{
		RepositoryName: u.RepositoryName,
		UserID:         u.ID,
		UserName:       u.UserName,
	}
}

func (u *User) GetResourceTag() *ResourceTag {
	return u.Metadata
}

func (u *User) SetClientTokens(accessToken *oauth2.TokenUnit, refreshToken *oauth2.TokenUnit) *oauth2.AccessTokenResponse {
	log.Debug("adding token for user and client", zap.String("user ID", u.ID.String()), zap.String("client ID", accessToken.ClientID.String()))
	if accessToken != nil {
		u.AccessTokens[accessToken.ClientID] = accessToken.GetTokenHint()
	}
	if refreshToken != nil {
		u.RefreshTokens[refreshToken.ClientID] = refreshToken.GetTokenHint()
	}
	return oauth2.GetAccessTokenResponse(accessToken, refreshToken)
}

func (u *User) setClientAuthorizationRequest(token *ConsentToken, request *oauth2.AuthorizationRequest) {
	uuid := uuid.FromStringOrNil(request.ClientID)
	log.Debug("client authorization request set for client", zap.String("user ID", u.ID.String()), zap.String("client ID", uuid.String()))
	u.AuthorizationRequests[uuid] = request
	u.AuthorizationRequestsConsentTokens[uuid] = token
}

func (u *User) getClientAuthorizationRequest(consentToken *ConsentToken) (autReq *oauth2.AuthorizationRequest, ok bool) {
	log.Debug("client authorization request get for client", zap.String("user ID", u.ID.String()), zap.String("client ID", consentToken.ClientID.String()))
	originalToken, tokenOK := u.AuthorizationRequestsConsentTokens[consentToken.ClientID]
	authReq, reqOK := u.AuthorizationRequests[consentToken.ClientID]
	delete(u.AuthorizationRequestsConsentTokens, consentToken.ClientID)
	delete(u.AuthorizationRequests, consentToken.ClientID)
	if tokenOK && reqOK {
		if strings.Compare(consentToken.Token, originalToken.Token) == 0 {
			return authReq, true
		} else {
			log.Debug("consent token not matching")
		}
	} else {
		log.Debug("token or authorization request was not found for the user", zap.String("user ID", u.ID.String()))
	}
	return nil, false
}

func (u *UserCtx) GetUserID() uuid.UUID {
	return u.UserID
}

func (u *UserCtx) GetRepositoryName() string {
	return u.RepositoryName
}

func (u *User) setScim(scim *scim2.User) {
	if scim != nil {
		u.Attributes = &UserAttributes{
			Active:            scim.Active,
			Addresses:         scim.Addresses,
			DisplayName:       scim.DisplayName,
			Entitlements:      scim.Entitlements,
			Emails:            scim.Emails,
			ExternalId:        scim.ExternalId,
			Ims:               scim.Ims,
			Locale:            scim.Locale,
			Name:              scim.Name,
			NickName:          scim.NickName,
			PhoneNumbers:      scim.PhoneNumbers,
			Photos:            scim.Photos,
			PreferredLanguage: scim.PreferredLanguage,
			ProfileURL:        scim.ProfileURL,
			Roles:             scim.Roles,
			Timezone:          scim.Timezone,
			Title:             scim.Title,
			UserType:          scim.UserType,
			X509Certificates:  scim.X509Certificates,
		}
	}
}
