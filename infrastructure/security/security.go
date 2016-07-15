package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
)

//var iv = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
const (
	AESKeyLength = 32
)

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncryptAES(key, text string) (string, error) {
	if len(key) != AESKeyLength {
		return "", errors.Errorf("Invalid key length %d should be %d", len(key), AESKeyLength)
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "error creating encryption cipher")
	}

	cipherText := make([]byte, aes.BlockSize+len(text))

	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "error creating iv")
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(text))

	return encodeBase64(cipherText), nil
}

func DecryptAES(key, text string) (string, error) {
	if len(key) != AESKeyLength {
		return "", errors.Errorf("Invalid key length %d should be %d", len(key), AESKeyLength)
	}

	cipherText, err := decodeBase64(text)
	if err != nil {
		return "", errors.Wrap(err, "error unbase64")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrap(err, "error creating decryption cipher")
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

func HmacSha1(secret, message string) (string, error) {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", errors.Wrap(err, "hmacsha1")
	}
	return encodeBase64(h.Sum(nil)), nil
}

func HmacSha256(secret, message string) (string, error) {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	_, err := h.Write([]byte(message))
	if err != nil {
		return "", errors.Wrap(err, "hmacsha256")
	}
	return encodeBase64(h.Sum(nil)), nil
}
