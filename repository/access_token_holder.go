package repository

import (
	"bounzr/iam/oauth2"
	"github.com/gofrs/uuid"
)

type AccessTokenHolder interface {
	DeleteClientAccessToken(clientID uuid.UUID)
	DeleteClientRefreshToken(clientID uuid.UUID)
	DeleteClientTokens(clientID uuid.UUID)
	GetClientAccessToken(clientID uuid.UUID) (tokenReference *oauth2.AccessTokenHint, ok bool)
	GetClientRefreshToken(clientID uuid.UUID) (tokenReference *oauth2.AccessTokenHint, ok bool)
	SetClientTokens(accessToken *oauth2.TokenUnit, refreshToken *oauth2.TokenUnit) *oauth2.AccessTokenResponse
}
