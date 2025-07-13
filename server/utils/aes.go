package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var key = []byte("thisis32bitlongpassphraseimusing") // 32-byte key for AES-256
var iv = []byte("thisis16bytesiv!")                  // 16-byte IV

func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("input is empty")
	}
	padding := int(src[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}
	return src[:length-padding], nil
}

func EncryptAES(plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedText := pad([]byte(plaintext))
	cipherText := make([]byte, len(paddedText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(cipherText string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	rawData, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(rawData)%aes.BlockSize != 0 {
		return "", errors.New("invalid ciphertext length")
	}

	plainText := make([]byte, len(rawData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, rawData)

	unpadded, err := unpad(plainText)
	if err != nil {
		return "", err
	}
	return string(unpadded), nil
}
