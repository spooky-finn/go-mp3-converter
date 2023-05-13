package pkg

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(s string) string {
	hashFunc := sha1.New()
	if _, err := hashFunc.Write([]byte(s)); err != nil {
		panic(err)
	}
	return hex.EncodeToString(hashFunc.Sum(nil)[:12])
}
