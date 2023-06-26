package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func StrEncrypt(message string, key string) string {
	messageBytes := []byte(message)
	keyBytes := []byte(key)

	ciphertext, err := aesEncrypt(messageBytes, keyBytes)
	if err != nil {
		return ""
	}

	result := base64.StdEncoding.EncodeToString(ciphertext)
	return result
}

func StrDecrypt(message string, key string) string {
	messageBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return ""
	}

	keyBytes := []byte(key)

	result, err := aesDecrypt(messageBytes, keyBytes)
	if err != nil {
		return ""
	}

	return string(result)
}

func aesEncrypt(message []byte, key []byte) ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padLen := aes.BlockSize - (len(message) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	paddedMessage := append(message, padText...)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedMessage))
	mode.CryptBlocks(ciphertext, paddedMessage)

	return append(iv, ciphertext...), nil
}

func aesDecrypt(ciphertext []byte, key []byte) ([]byte, error) {
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedMessage := make([]byte, len(ciphertext))
	mode.CryptBlocks(decryptedMessage, ciphertext)

	padLen := int(decryptedMessage[len(decryptedMessage)-1])
	unpaddedMessage := decryptedMessage[:len(decryptedMessage)-padLen]

	return unpaddedMessage, nil
}
