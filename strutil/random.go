package strutil

import (
	"crypto/rand"

	"github.com/google/uuid"
)

// AlphanumericChars contains standard alphanumeric characters (a-z, 0-9).
const AlphanumericChars = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandomString generates a random string of a given length using a specific character set.
// If charSet is empty, it defaults to AlphanumericChars.
// It uses crypto/rand for security, with a fallback to UUID-based randomness.
func RandomString(length int, charSet string) string {
	if charSet == "" {
		charSet = AlphanumericChars
	}

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback if crypto/rand fails (extremely rare)
		for i := range b {
			b[i] = charSet[uuid.New().ID()%uint32(len(charSet))]
		}
		return string(b)
	}

	for i := range b {
		b[i] = charSet[int(b[i])%len(charSet)]
	}
	return string(b)
}
