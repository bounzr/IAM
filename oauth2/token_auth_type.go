package oauth2

type TokenAuthType int

const (
	nullTokenAuthType   TokenAuthType = -1
	BearerTokenAuthType TokenAuthType = iota
	MACTokenAuthType
)

var tokenAuthTypeValueMap = map[string]TokenAuthType{
	"Bearer": BearerTokenAuthType,
	"MAC":    MACTokenAuthType,
}

func NewTokenAuthType(t string) TokenAuthType {
	if val, ok := tokenAuthTypeValueMap[t]; ok {
		return val
	}
	return nullTokenAuthType
}

func (t TokenAuthType) String() string {
	switch t {
	case BearerTokenAuthType:
		return "Bearer"
	case MACTokenAuthType:
		return "MAC"
	}
	return "null"

}
