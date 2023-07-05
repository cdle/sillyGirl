package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func getJS(user, pass, title string) string {
	baseURL := "http://aut.zhelee.cn/plugin/download"
	queryParams := url.Values{}
	queryParams.Set("title", title)
	queryParams.Set("username", user)
	queryParams.Set("password", pass)
	queryParams.Set("version", "1")

	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	u.RawQuery = queryParams.Encode()

	urlStr := u.String()
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, urlStr, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// fmt.Println(string(body))
	return string(body)
}

func main() {
	user := "test1234"
	pass := "test1234"
	title := "hook"
	plaintext := getJS(user, pass, title)
	key := getKey("test1234")
	panic(key)
	decrypted, err := decrypt(string(plaintext), []byte(key))
	// 打印解密后的数据
	fmt.Println(err)
	fmt.Printf("解密后的数据：%s\n", decrypted)
}

func getKey(user string) string {
	t := user + "killsillygirltodeath109times"
	if len(t) > 32 {
		return t[:32]
	}
	// if len(t) > 24 {
	// 	return t[:24]
	// }
	// return t[:16]
	return t[:24]
}

// 使用AES-CBC模式进行解密
func decrypt(encrypted string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	decrypter.CryptBlocks(decrypted, ciphertext)
	return string(decrypted), nil
}
