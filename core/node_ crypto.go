package core

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/dop251/goja"
)

func toBytes(v interface{}) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case *Buffer:
		return v.value
	}
	return nil
}

func toString(v interface{}) string {
	switch v := v.(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case *Buffer:
		return string(v.value)
	}
	return ""
}

type Crypto struct {
	vm *goja.Runtime
}

type Hmac struct {
	// decode    bool
	vm        *goja.Runtime
	Algorithm string
	middle    interface{}
}

func (h *Hmac) Digest(EF string) interface{} {
	var result interface{}
	var SF = ""
	switch h.Algorithm {
	case "md5", "sha1", "sha256", "sha512":
		result = h.middle.(hash.Hash).Sum(nil)
	case "aes-256-cbc":
	}
	if EF == "" {
		return &Buffer{
			value: result.([]byte),
		}
	}
	return Convert(h.vm, result, SF, EF)
}

func (h *Hmac) Update(data interface{}, SF, EF string) interface{} {
	switch h.Algorithm {
	case "md5", "sha1", "sha256", "sha512":
		h.middle.(hash.Hash).Write(toBytes(data))
	case "aes-256-cbc":
		if SF == "hex" {
			data = Convert(h.vm, data, "hex", "")
		}
		return Convert(h.vm, h.middle.(*AesCipher).Update(toBytes(data)), "", EF)
	}
	return nil
}

func (h *Hmac) Final(EF string) interface{} {
	switch h.Algorithm {
	case "md5", "sha1", "sha256", "sha512":
	case "aes-256-cbc":
		return Convert(h.vm, h.middle.(*AesCipher).Final(), "", EF)
	}
	return nil
}

func (c *Crypto) CreateCipheriv(algorithm string, key, iv interface{}) *Hmac {
	var middle interface{}
	var err error
	switch algorithm {
	case "aes-256-cbc":

		middle, err = NewAesCipher(c.vm, toBytes(key), toBytes(iv), false)
		if err != nil {
			panic(Error(c.vm, err))
		}
	}
	return &Hmac{
		vm: c.vm,
		// decode:    true,
		middle:    middle,
		Algorithm: algorithm,
	}
}

func (c *Crypto) CreateDecipheriv(algorithm string, key, iv interface{}) *Hmac {
	var middle interface{}
	var err error
	switch algorithm {
	case "aes-256-cbc":
		middle, err = NewAesCipher(c.vm, toBytes(key), toBytes(iv), true)
		if err != nil {
			panic(Error(c.vm, err))
		}
	}
	return &Hmac{
		// decode:    true,
		vm:        c.vm,
		middle:    middle,
		Algorithm: algorithm,
	}
}

func (c *Crypto) CreateHash(algorithm string) *Hmac {
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
		vm:        c.vm,
		middle:    hashFunc(),
		Algorithm: algorithm,
	}
}

func (c *Crypto) CreateHmac(algorithm string, key interface{}) *Hmac {
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
		vm:        c.vm,
		middle:    hmac.New(hashFunc, toBytes(key)),
		Algorithm: algorithm,
	}
}

type AesCipher struct {
	block  cipher.Block
	iv     []byte
	mode   cipher.BlockMode
	decode bool
	buffer bytes.Buffer
	vm     *goja.Runtime
}

func NewAesCipher(vm *goja.Runtime, key, iv []byte, decode bool) (*AesCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	if !decode {
		mode = cipher.NewCBCEncrypter(block, iv)
	}
	return &AesCipher{
		vm:     vm,
		block:  block,
		iv:     iv,
		mode:   mode,
		decode: decode,
	}, nil
}

func (ac *AesCipher) Update(text []byte) []byte {
	//fmt.Println("Update", string(text))
	if ac.decode {
		// Decrypt
		var err error
		//fmt.Println("Decrypt1", string(text))
		ac.mode.CryptBlocks(text, text)
		//fmt.Println("Decrypt2", string(text))
		text, err = pkcs7Unpadding(text)
		if err != nil {
			panic(Error(ac.vm, err))
		}
		//fmt.Println("Decrypt3", string(text), err)
	} else {
		// Encrypt
		text = pkcs7Padding(text, ac.block.BlockSize())
		ac.mode.CryptBlocks(text, text)
	}
	ac.buffer.Write(text)
	return nil
	// return text
}

func (ac *AesCipher) Final() []byte {
	return ac.buffer.Bytes()
}

func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("pkcs7: invalid padding")
	}
	return data[:(length - unpadding)], nil
}

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

//	func (c *Cipher) Rc4Decrypt(key, ciphertext string) (string, error) {
//		cipher, err := rc4.NewCipher([]byte(key))
//		if err != nil {
//			return "", err
//		}
//		data, err := base64.StdEncoding.DecodeString(ciphertext)
//		if err != nil {
//			return "", err
//		}
//		cipher.XORKeyStream(data, data)
//		return string(data), nil
//	}

func Convert(vm *goja.Runtime, data interface{}, fromFormat string, toFormat string) interface{} {
	//fmt.Println(data, fromFormat, toFormat)
	var dataBytes []byte

	switch input := data.(type) {
	case string:
		dataBytes = []byte(input)
	case []byte:
		dataBytes = input
	case *Buffer:
		dataBytes = input.value
	default:
		panic(Error(vm, fmt.Errorf("invalid data type: %T", data)))
	}

	switch fromFormat {
	case "hex":
		// 将数据从 HEX 编码格式转换为字节数组
		decodedData, err := hex.DecodeString(string(dataBytes))
		if err != nil {
			panic(Error(vm, err))
		}
		dataBytes = decodedData
	case "base64":
		// 将数据从 Base64 编码格式转换为字节数组
		decodedData, err := base64.StdEncoding.DecodeString(string(dataBytes))
		if err != nil {
			panic(Error(vm, err))
		}
		dataBytes = decodedData
	case "bytes", "binary":
		// 不需要进行转换
	case "utf8", "utf-8", "":

	default:
		panic(Error(vm, fmt.Errorf("unsupported input format: %s", fromFormat)))
	}

	//fmt.Println(string(dataBytes), fromFormat, toFormat)

	switch toFormat {
	case "hex":
		// 将数据转换为 HEX 编码格式
		//fmt.Println("hex", hex.EncodeToString(dataBytes), fromFormat, toFormat)
		return hex.EncodeToString(dataBytes)
	case "base64":
		// 将数据转换为 Base64 编码格式
		return base64.StdEncoding.EncodeToString(dataBytes)
	case "bytes", "binary":
		// 不需要进行转换
		return dataBytes
	case "utf8", "utf-8", "":
		// 将数据转换为 UTF-8 编码格式
		return string(dataBytes)
	default:
		panic(Error(vm, fmt.Errorf("unsupported output format: %s", toFormat)))
	}
}

func cryptoModule(vm *goja.Runtime, module *goja.Object) {
	cryto := Crypto{
		vm: vm,
	}
	o := module.Get("exports").(*goja.Object)
	o.Set("createCipheriv", cryto.CreateCipheriv)
	o.Set("createDecipheriv", cryto.CreateDecipheriv)
	o.Set("createHash", cryto.CreateHash)
	o.Set("createHmac", cryto.CreateHmac)
}
