package tokens

type TokenGenerator interface {
	getToken() []byte
}

func GetToken() []byte {
	return tokenGenerator.getToken()
}
