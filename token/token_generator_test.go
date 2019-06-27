package token

import "testing"

func TestGetToken(t *testing.T) {
	var tg TokenGenerator
	tg = NewTokenGeneratorBasic()
	token := tg.getToken()
	if len(token) == 0 {
		t.Errorf("Expected token got empty value")
	}
}
