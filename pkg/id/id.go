package id

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

const LENGTH = 24

func Ascending() (string, error) {
	return generateID(false)
}

func Descending() (string, error) {
	return generateID(true)
}

func generateID(descending bool) (string, error) {
	now := time.Now().UnixMilli()
	if descending {
		now = ^now
	}

	timeBytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		timeBytes[i] = byte(now >> (40 - 8*i))
	}

	randomBytes := make([]byte, (LENGTH-12)/2)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	result := make([]byte, LENGTH)
	hex.Encode(result[:12], timeBytes)
	hex.Encode(result[12:], randomBytes)

	return string(result), nil
}
