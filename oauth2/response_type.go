package oauth2

import "errors"

type ResponseType int

const(
	nullResponseType	ResponseType = -1
	Code				ResponseType = iota
	Token
)

var responseTypeValueMap = map[string]ResponseType{
	"code": Code,
	"token": Token,
}

var ErrResponseTypeNotFound = errors.New("wrong response type requested. Returning code")

func NewResponseType(s string) (ResponseType, error){
	if val, ok := responseTypeValueMap[s]; ok{
		return val, nil
	}
	return nullResponseType, ErrResponseTypeNotFound
}

func (rt ResponseType)	String() string{
	switch rt{
	case Code:
		return "code"
	case Token:
		return "token"
	}
	return "null"
}


