package oauth2

type TokenHintType int

const (
	NullTokenHintType   TokenHintType = -1
	AccessTokenHintType TokenHintType = iota
	RefreshTokenHintType
)

var tokenHintTypeValueMap = map[string]TokenHintType{
	"access_token":  AccessTokenHintType,
	"refresh_token": RefreshTokenHintType,
}

func NewTokenHintType(t string) TokenHintType {
	if val, ok := tokenHintTypeValueMap[t]; ok {
		return val
	}
	return NullTokenHintType
}

func (t TokenHintType) String() string {
	switch t {
	case AccessTokenHintType:
		return "access_token"
	case RefreshTokenHintType:
		return "refresh_token"
	}
	return "nil"

}
