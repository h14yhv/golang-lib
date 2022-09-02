package rest

import (
	"fmt"
	"net/url"
)

type (
	ProxyConfig struct {
		Enable   bool   `json:"enable" yaml:"enable"`
		Scheme   string `json:"scheme" yaml:"scheme"`
		Host     string `json:"host" yaml:"host"`
		Port     int    `json:"port" yaml:"port"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	}

	AuthenticateConfig struct {
		Type       string      `json:"type" yaml:"type"`
		Attributes interface{} `json:"attributes" yaml:"attributes"`
	}

	AuthenticateBasicAuthConfig struct {
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	}

	AuthenticateTokenConfig struct {
		Token string `json:"token" yaml:"token"`
	}
)

func (body *ProxyConfig) String() string {
	host := body.Host
	if body.Port != 0 {
		host = fmt.Sprintf("%s:%d", host, body.Port)
	}
	u := url.URL{
		Scheme: body.Scheme,
		Host:   host,
	}
	if body.Username != "" || body.Password != "" {
		u.User = url.UserPassword(body.Username, body.Password)
	}
	// Success
	return u.String()
}

func (body *AuthenticateConfig) GetBasicAuthConfig() AuthenticateBasicAuthConfig {
	auth := AuthenticateBasicAuthConfig{}
	if body != nil {
		if user, ok := body.Attributes.(map[string]interface{})["username"]; ok {
			auth.Username = user.(string)
		}
		if pass, ok := body.Attributes.(map[string]interface{})["password"]; ok {
			auth.Password = pass.(string)
		}
	}
	// Success
	return auth
}

func (body *AuthenticateConfig) GetTokenConfig() AuthenticateTokenConfig {
	auth := AuthenticateTokenConfig{}
	if body != nil {
		if token, ok := body.Attributes.(map[string]interface{})["token"]; ok {
			auth.Token = token.(string)
		}
	}
	// Success
	return auth
}
