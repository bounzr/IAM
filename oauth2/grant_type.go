package oauth2

import "errors"

type GrantType int

const (
	NullGrantType              GrantType = -1
	AuthorizationCodeGrantType GrantType = iota // 0
	ClientCredentialsGrantType
	ImplicitGrantType
	PasswordGrantType
	RefreshTokenGrantType
	JwtBearerGrantType
	Saml2BearerGrantType
)

var grantTypeValueMap = map[string]GrantType{
	"authorization_code": AuthorizationCodeGrantType,
	"client_credentials": ClientCredentialsGrantType,
	"implicit":           ImplicitGrantType,
	"password":           PasswordGrantType,
	"refresh_token":      RefreshTokenGrantType,
	"urn:ietf:params:oauth:grant-type:jwt-bearer":   JwtBearerGrantType,
	"urn:ietf:params:oauth:grant-type:saml2-bearer": Saml2BearerGrantType,
}

var ErrGrantTypeNotFound = errors.New("wrong grant type requested. Returning default")

func NewGrantType(s string) (GrantType, error) {
	if val, ok := grantTypeValueMap[s]; ok {
		return val, nil
	}
	return NullGrantType, ErrGrantTypeNotFound
}

func (t GrantType) String() string {
	switch t {
	case AuthorizationCodeGrantType:
		return "authorization_code"
	case ClientCredentialsGrantType:
		return "client_credentials"
	case ImplicitGrantType:
		return "implicit"
	case PasswordGrantType:
		return "password"
	case RefreshTokenGrantType:
		return "refresh_token"
	case JwtBearerGrantType:
		return "urn:ietf:params:oauth:grant-type:jwt-bearer"
	case Saml2BearerGrantType:
		return "urn:ietf:params:oauth:grant-type:saml2-bearer"
	}
	return "null"
}
