package utils

import (
	"math/rand"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes   = "1234567890"
	specialBytes  = "~=+%^*/()[]{}/!@#$?|"
)

var (
	src = rand.NewSource(time.Now().UnixNano())
)

func GetRandomBytes(length int) []byte {
	bytesSet := letterBytes + numberBytes
	buf := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(bytesSet) {
			buf[i] = bytesSet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return buf
}

func GetRandomPassword(length int) string {
	bytes := letterBytes + numberBytes + specialBytes
	randPass := rand.New(src)
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = bytes[randPass.Intn(len(bytes))]
	}
	return string(buf)
}

func GetRandomString(length int) string {
	bytesSet := letterBytes + numberBytes
	buf := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(bytesSet) {
			buf[i] = bytesSet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(buf)
}

func GetRandom32Token() [32]byte {
	var	randomBytes [32]byte
	byteSet := letterBytes + numberBytes
	for i, cache, remain := 31, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(byteSet) {
			randomBytes[i] = byteSet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return randomBytes
}
