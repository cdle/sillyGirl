package core

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

func GetTransport(proxyUrl string, user, password string) (*http.Transport, error) {
	var auth *proxy.Auth
	if user != "" && password != "" {
		auth = &proxy.Auth{User: user, Password: password}
	}
	switch {
	case strings.HasPrefix(proxyUrl, "socks5://"):
		dialer, err := proxy.SOCKS5("tcp", strings.TrimPrefix(proxyUrl, "socks5://"), auth, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return &http.Transport{Dial: dialer.Dial}, nil
	case strings.HasPrefix(proxyUrl, "http://"):
		proxyUrl = strings.TrimPrefix(proxyUrl, "http://")
		url, err := url.Parse("http://" + proxyUrl)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{Proxy: http.ProxyURL(url)}
		if auth != nil {
			// Encode the credentials in Base64
			authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
			transport.ProxyConnectHeader = http.Header{
				"Proxy-Authorization": {authHeader},
			}
		}
		return transport, nil
	case strings.HasPrefix(proxyUrl, "https://"):
		proxyUrl = strings.TrimPrefix(proxyUrl, "https://")
		url, err := url.Parse("https://" + proxyUrl)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{Proxy: http.ProxyURL(url)}
		if auth != nil {
			// Encode the credentials in Base64
			authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password))
			transport.ProxyConnectHeader = http.Header{
				"Proxy-Authorization": {authHeader},
			}
		}
		return transport, nil
	default:
		return nil, errors.New("proxy url schema error")
	}
}
