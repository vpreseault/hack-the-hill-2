package database

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateID() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}