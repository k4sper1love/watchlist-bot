package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts a plain text string using AES-GCM encryption.
// It retrieves the master key, generates a random nonce, and seals the data with the key and nonce.
// The resulting cipher text is encoded in base64 for safe storage or transmission.
// Returns the encrypted text as a base64-encoded string or an error if encryption fails.
func Encrypt(plainText string) (string, error) {
	// Retrieve the master key from the environment variable.
	key, err := getMasterKey()
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block using the master key.
	// This block is used to initialize the AES-GCM mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil
	}

	// Initialize AES-GCM mode, which provides authenticated encryption.
	// AES-GCM combines encryption and integrity verification in a single step.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random nonce (number used once) for the encryption process.
	// The nonce size is determined by the AES-GCM implementation.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plain text using the AES-GCM seal function.
	// The nonce is prepended to the cipher text for use during decryption.
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	// Encode the cipher text (including the nonce) in base64 for safe storage or transmission.
	// Base64 encoding ensures that the encrypted data can be safely stored as a string.
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts a base64-encoded cipher text using AES-GCM decryption.
// It retrieves the master key, decodes the base64 input, and extracts the nonce and cipher text.
// The plain text is then decrypted and returned as a string.
// Returns an error if the decryption process fails or if the input is invalid.
func Decrypt(encryptedText string) (string, error) {
	// Retrieve the master key from the environment variable.
	key, err := getMasterKey()
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded cipher text into raw bytes.
	// This step reverses the encoding performed during encryption.
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block using the master key.
	// This block is used to initialize the AES-GCM mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Initialize AES-GCM mode, which provides authenticated decryption.
	// AES-GCM ensures both confidentiality and integrity of the encrypted data.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the cipher text.
	// The nonce is prepended to the cipher text during encryption.
	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		// If the cipher text is too short to contain a valid nonce, return an error.
		return "", errors.New("cipher text too short")
	}

	// Split the cipher text into the nonce and the actual encrypted data.
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	// Decrypt the cipher text using the AES-GCM open function.
	// This step verifies the integrity of the data and recovers the original plain text.
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	// Return the decrypted plain text as a string.
	return string(plainText), nil
}
