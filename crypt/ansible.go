package crypt

import (
	"bytes"
	"fmt"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/pbkdf2"
)

const (
	KEYLENGTH  = 32
	SALTLENGTH = 32

	MAGIC = "$ANSIBLE_VAULT"
)

func GenerateKeys(password, salt []byte) ([]byte, []byte, []byte) {
	key := pbkdf2.Key(password, salt, 10000, (2*KEYLENGTH)+aes.BlockSize, sha256.New)
	return key[:KEYLENGTH], key[KEYLENGTH:(KEYLENGTH * 2)], key[(KEYLENGTH * 2) : (KEYLENGTH*2)+aes.BlockSize]
}

func pkcs7Pad(src []byte) []byte {
	padding := aes.BlockSize - (len(src) % aes.BlockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs7Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("Invalid padding")
	}
	return src[:(length - unpadding)], nil
}

func AnsibleEncrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, SALTLENGTH)
	if _, err := rand.Read(salt); err != nil {
		return []byte{}, err
	}
	key1, key2, iv := GenerateKeys([]byte(password), salt)

	block, err := aes.NewCipher(key1)
	if err != nil {
		return []byte{}, err
	}
	stream := cipher.NewCTR(block, iv)
	paddedPlaintext := pkcs7Pad(plaintext)
	ciphertext := make([]byte, len(paddedPlaintext))
	stream.XORKeyStream(ciphertext, paddedPlaintext)

	hash := hmac.New(sha256.New, key2)
	hash.Write(ciphertext)
	signature := hash.Sum(nil)

	encryptedSlice := []string{
		hex.EncodeToString(salt),
		hex.EncodeToString(signature),
		hex.EncodeToString(ciphertext),
	}
	encrypted := strings.Join(encryptedSlice, "\n")
	encryptedBytes := make([]byte, hex.EncodedLen(len(encrypted)))
	hex.Encode(encryptedBytes, []byte(encrypted))

	header := [][]byte{
		[]byte(MAGIC),
		[]byte("1.1"),
		[]byte("AES256"),
	}
	body := [][]byte{
		bytes.Join(header, []byte(";")),
	}
	for i := 0; i < len(encryptedBytes); i += 80 {
		body = append(body, encryptedBytes[i:min(len(encryptedBytes), i+80)])
	}
	body = append(body, []byte{})
	return bytes.Join(body, []byte("\n")), nil
}

func AnsibleDecrypt(encrypted []byte, password string) ([]byte, error) {
	body := bytes.Split(encrypted, []byte("\n"))
	if !bytes.HasPrefix(body[0], []byte(MAGIC+";")) {
		return []byte{}, fmt.Errorf("Invalid encrypted data format")
	}
	encryptedHex := bytes.Join(body[1:], []byte{})
	encryptedBytes := make([]byte, hex.DecodedLen(len(encryptedHex)))
	if _, err := hex.Decode(encryptedBytes, encryptedHex); err != nil {
		return []byte{}, err
	}
	encryptedSlice := strings.SplitN(string(encryptedBytes), "\n", 3)
	if len(encryptedSlice) != 3 {
		return []byte{}, fmt.Errorf("Invalid encrypted data format")
	}
	salt, err := hex.DecodeString(encryptedSlice[0])
	if err != nil {
		return []byte{}, err
	}
	signature, err := hex.DecodeString(encryptedSlice[1])
	if err != nil {
		return []byte{}, err
	}
	ciphertext, err := hex.DecodeString(encryptedSlice[2])
	if err != nil {
		return []byte{}, err
	}
	key1, key2, iv := GenerateKeys([]byte(password), salt)
	hash := hmac.New(sha256.New, key2)
	hash.Write(ciphertext)
	check := hash.Sum(nil)
	if bytes.Compare(signature, check) != 0 {
		return []byte{}, fmt.Errorf("Invalid password")
	}
	block, err := aes.NewCipher(key1)
	if err != nil {
		return []byte{}, err
	}
	paddedPlaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(paddedPlaintext, ciphertext)
	plaintext, err := pkcs7Unpad(paddedPlaintext)
	if err != nil {
		return []byte{}, err
	}
	return plaintext, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
