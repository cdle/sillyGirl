package core

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
)

type Crypto struct{}

type Hmac struct {
	Algorithm string
	hash      hash.Hash
}

func (h *Hmac) Digest(class string) string {
	value := h.hash.Sum(nil)
	switch class {
	case "hex":
		return fmt.Sprintf("%x", value)
	case "base64":
		return base64.StdEncoding.EncodeToString(value)
	}
	return ""
}

func (h *Hmac) Update(data string) {
	h.hash.Write([]byte(data))
}

func (c *Crypto) CreateHmac(algorithm, key string) *Hmac {
	var hashFunc func() hash.Hash
	switch algorithm {
	case "md5":
		hashFunc = md5.New
	case "sha1":
		hashFunc = sha1.New
	case "sha256":
		hashFunc = sha256.New
	case "sha512":
		hashFunc = sha512.New
	default:
		return nil
	}
	return &Hmac{
		hash:      hmac.New(hashFunc, []byte(key)),
		Algorithm: algorithm,
	}
}

// func (c *Crypto) CreateCipheriv(algorithm, key, data string) *Hmac {
// 	switch algorithm {
// 	case "aes-256-cbc":
// 		block, err := aes.NewCipher([]byte(key))
// 		if err != nil {
// 			return "", err
// 		}
// 		ciphertext := make([]byte, aes.BlockSize+len(data))
// 		iv := ciphertext[:aes.BlockSize]
// 		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 			return "", err
// 		}
// 		stream := cipher.NewCFBEncrypter(block, iv)
// 		stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))
// 		return base64.StdEncoding.EncodeToString(ciphertext), nil
// 	}
// }

// type Cipher struct{}

// func (c *Cipher) AesEncrypt(key, data string) (string, error) {
// 	block, err := aes.NewCipher([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}
// 	ciphertext := make([]byte, aes.BlockSize+len(data))
// 	iv := ciphertext[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return "", err
// 	}
// 	stream := cipher.NewCFBEncrypter(block, iv)
// 	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(data))
// 	return base64.StdEncoding.EncodeToString(ciphertext), nil
// }

// func (c *Cipher) AesDecrypt(key, ciphertext string) (string, error) {
// 	block, err := aes.NewCipher([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}
// 	data, err := base64.StdEncoding.DecodeString(ciphertext)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(data) < aes.BlockSize {
// 		return "", fmt.Errorf("ciphertext too short")
// 	}
// 	iv := data[:aes.BlockSize]
// 	data = data[aes.BlockSize:]
// 	stream := cipher.NewCFBDecrypter(block, iv)
// 	stream.XORKeyStream(data, data)
// 	return string(data), nil
// }

// func (c *Cipher) Rc4Encrypt(key, data string) (string, error) {
// 	cipher, err := rc4.NewCipher([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}
// 	ciphertext := make([]byte, len(data))
// 	cipher.XORKeyStream(ciphertext, []byte(data))
// 	return base64.StdEncoding.EncodeToString(ciphertext), nil
// }

// func (c *Cipher) Rc4Decrypt(key, ciphertext string) (string, error) {
// 	cipher, err := rc4.NewCipher([]byte(key))
// 	if err != nil {
// 		return "", err
// 	}
// 	data, err := base64.StdEncoding.DecodeString(ciphertext)
// 	if err != nil {
// 		return "", err
// 	}
// 	cipher.XORKeyStream(data, data)
// 	return string(data), nil
// }
