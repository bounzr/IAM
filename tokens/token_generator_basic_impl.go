package tokens

import (
	"math/rand"
	"time"
)

const (
	idxBits   = 6              // 6 bits to represent a index
	idxMask   = 1<<idxBits - 1 // All 1-bits, as many as IdxBits
	idxMax    = 63 / idxBits   // # of indices fitting in 63 bits
	bytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890."
	tokenSize = 32
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

type TokenGeneratorBasic struct {
	//usedTokens  map[string]struct{}
}

func NewTokenGeneratorBasic() TokenGenerator {
	tgb := &TokenGeneratorBasic{}
	tgb.Init()
	return tgb
}

func (tg *TokenGeneratorBasic) Init() {
	//tg.usedTokens = make(map[string]struct{})
}

func (tg *TokenGeneratorBasic) getToken() []byte {
	buf := make([]byte, tokenSize)
	for i, cache, remain := tokenSize-1, src.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(bytes) {
			buf[i] = bytes[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}
	return buf
}
