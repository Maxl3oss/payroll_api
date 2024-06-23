package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/gofiber/fiber/v2"
)

func encryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	// Return as base64 string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func EncryptResponse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Proceed with next middleware
		if err := c.Next(); err != nil {
			return err
		}

		// Get response body
		body := c.Response().Body()

		// Encrypt response body
		encrypted, err := encryptAES(body, []byte("Zi4VwqYgHXNbBQRRETetjPZVRHKibAux"))
		if err != nil {
			return err
		}

		// Set encrypted body as response
		c.Response().SetBody([]byte(encrypted))
		return nil
	}
}
