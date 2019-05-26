package tokens

var (
	tokenGenerator TokenGenerator
)

func Init() {
	initTokenGenerator()
}

func initTokenGenerator() {
	tokenGenerator = NewTokenGeneratorBasic()
}
