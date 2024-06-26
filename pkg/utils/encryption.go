package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func encrypt(data []byte, key []byte) (string, error) {
	// Ensure the key is 32 bytes long for AES-256
	if len(key) != 32 {
		return "", fmt.Errorf("key must be 32 bytes long")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext, _ := pkcs7Pad(data, block.BlockSize())

	// Generate a random IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...), nil
}

func EncryptResponse() fiber.Handler {
	key := os.Getenv("SECRET_KEY")

	return func(c *fiber.Ctx) error {
		// Proceed with next middleware
		if err := c.Next(); err != nil {
			return err
		}

		// Get response body
		body := c.Response().Body()

		// Encrypt response body
		encrypted, err := encrypt(body, []byte(key))
		if err != nil {
			return err
		}

		// Set encrypted body as response
		c.Response().SetBody([]byte(encrypted))
		return nil
	}
}
