package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	alias := make([]byte, length)
	for i := range alias {
		alias[i] = chars[rnd.Intn(len(chars))]
	}

	return string(alias)
}
