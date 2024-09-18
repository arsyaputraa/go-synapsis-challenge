package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate random characters
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
