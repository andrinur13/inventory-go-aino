package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

var key = []byte("tourismNtransportationEncryption")

//TrimLeftChars : trimming prefix
func TrimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

//EncryptQR : encript qr
func EncryptQR(text string, qrPrefix string) string {
	qrArr := strings.Split(text, "#")
	qrStan := qrArr[len(qrArr)-1]
	qrMidArr := qrArr[:len(qrArr)-1]
	qrMid := strings.Join(qrMidArr, "#")
	qrEnd := encrypt(key, qrStan)
	encryptedQR := qrPrefix + "#" + qrMid + "#" + qrEnd

	return encryptedQR
}

//DecryptQR : decrypt qr
func DecryptQR(text string) (string, error) {
	qrArr := strings.Split(text, "#")

	if len(qrArr) < 2 {
		return "", fmt.Errorf("Invalid QR Code")
	}

	qrStan := qrArr[len(qrArr)-1]
	qrMidArr := qrArr[:len(qrArr)-1]
	qrMidArr = append(qrMidArr[:0], qrMidArr[1:]...)
	qrMid := strings.Join(qrMidArr, "#")

	qrStanDecode, err := decrypt(key, qrStan)

	if err != nil {
		return "", fmt.Errorf("Invalid QR Code")
	}

	qrCode := qrMid + "#" + qrStanDecode

	return qrCode, nil
}

func encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(key []byte, cryptoText string) (string, error) {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("Invalid QR Code")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
