package util

import (
	"crypto/rand"
	"fmt"
)

// NewRandomID generates a random string ID
func NewRandomID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%8x", b[0:4]), nil
}
