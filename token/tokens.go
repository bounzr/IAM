package token

var (
	tokenGenerator TokenGenerator
)

func Init() {
	tokenGenerator = NewTokenGeneratorBasic()
}
