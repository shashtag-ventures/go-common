package crypto_test

import (
	"testing"

	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/stretchr/testify/assert"
)

const testKey = "this-is-a-32-byte-key-1234567890" // Exactly 32 bytes

func TestEncryptDecrypt(t *testing.T) {
	plainText := "Hello, World! This is a secret message."

	// Test Case 1: Successful Encryption and Decryption
	t.Run("Success", func(t *testing.T) {
		cipherText, err := crypto.Encrypt(plainText, testKey)
		assert.NoError(t, err)
		assert.NotEmpty(t, cipherText)
		assert.NotEqual(t, plainText, cipherText)

		decryptedText, err := crypto.Decrypt(cipherText, testKey)
		assert.NoError(t, err)
		assert.Equal(t, plainText, decryptedText)
	})

	// Test Case 2: Encryption with invalid key size
	t.Run("Invalid Key Size", func(t *testing.T) {
		shortKey := "too-short"
		_, err := crypto.Encrypt(plainText, shortKey)
		assert.Error(t, err)
	})

	// Test Case 3: Decryption with wrong key
	t.Run("Wrong Key", func(t *testing.T) {
		cipherText, err := crypto.Encrypt(plainText, testKey)
		assert.NoError(t, err)

		wrongKey := "another-32-byte-key-1234567890-!"
		_, err = crypto.Decrypt(cipherText, wrongKey)
		assert.Error(t, err)
	})

	// Test Case 4: Decryption of invalid base64
	t.Run("Invalid Base64", func(t *testing.T) {
		_, err := crypto.Decrypt("not-base64-!", testKey)
		assert.Error(t, err)
	})

	// Test Case 5: Decryption of short cipher text
	t.Run("Short CipherText", func(t *testing.T) {
		_, err := crypto.Decrypt("YWI=", testKey) // base64 for "ab"
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cipher text too short")
	})
}
