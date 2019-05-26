package oauth2

import (
	"errors"
)

type TokenEndpointAuthMethod int

const (
	none TokenEndpointAuthMethod = iota
	clientSecretPost
	clientSecretBasic
)

var tokenEndpointAuthMethodValueMap = map[string]TokenEndpointAuthMethod{
	"none":                none,
	"client_secret_post":  clientSecretPost,
	"client_secret_basic": clientSecretBasic,
}

var ErrTokenEndpointAuthMethodNotFound = errors.New("token endpoint authentication method not found. Returning default")

func (m TokenEndpointAuthMethod) String() string {
	switch m {
	case none:
		return "none"
	case clientSecretPost:
		return "client_secret_post"
	case clientSecretBasic:
		return "client_secret_basic"
	default:
		return "none"
	}
}

func (m TokenEndpointAuthMethod) Parse(s string) (TokenEndpointAuthMethod, error) {
	if val, ok := tokenEndpointAuthMethodValueMap[s]; ok {
		return val, nil
	}
	return clientSecretPost, ErrTokenEndpointAuthMethodNotFound
}
