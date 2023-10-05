package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	// We're using a 32 byte long secret key.
	// This is probably something you generate first
	// then put into and environment variable.
	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

func Encrypt(plaintext string) (string, error) {
	// Convert the secret key to 32-byte array
	key := []byte(secretKey)

	// Create a new AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a new random initialization vector
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Create a new AES cipher block mode using the block and IV
	mode := cipher.NewCBCEncrypter(block, iv)

	// Pad the plaintext to the nearest multiple of the block size
	paddedPlaintext := padPlaintext([]byte(plaintext), block.BlockSize())

	// Create a byte slice for the ciphertext
	ciphertext := make([]byte, len(paddedPlaintext))

	// Encrypt the padded plaintext using the mode
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Concatenate the IV and ciphertext
	encrypted := append(iv, ciphertext...)

	// Encode the encrypted byte slice to base64
	encoded := base64.URLEncoding.EncodeToString(encrypted)

	// Return the encrypted string in UTF-8
	return encoded, nil
}

func Decrypt(encrypted string) (string, error) {
	// Convert the secret key to 32-byte array
	key := []byte(secretKey)

	// Decode the base64 encoded string
	decoded, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// Extract the IV and ciphertext from the decoded byte slice
	iv := decoded[:aes.BlockSize]
	ciphertext := decoded[aes.BlockSize:]

	// Create a new AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block mode using the block and IV
	mode := cipher.NewCBCDecrypter(block, iv)

	// Create a byte slice for the padded plaintext
	paddedPlaintext := make([]byte, len(ciphertext))

	// Decrypt the ciphertext using the mode
	mode.CryptBlocks(paddedPlaintext, ciphertext)

	// Unpad the padded plaintext
	unpaddedPlaintext := unpadPlaintext(paddedPlaintext)

	// Return the decrypted plaintext
	return string(unpaddedPlaintext), nil
}
func unpadPlaintext(paddedPlaintext []byte) []byte {
	padding := int(paddedPlaintext[len(paddedPlaintext)-1])
	unpadded := paddedPlaintext[:len(paddedPlaintext)-padding]
	return unpadded
}

// padPlaintext pads the plaintext to the nearest multiple of the block size.
func padPlaintext(plaintext []byte, blockSize int) []byte {
	padding := blockSize - (len(plaintext) % blockSize)
	padded := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)
	return padded
}
func GenerateString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
