package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AuthClient struct {
	Config AppConfig
	AppUrl string
}

// Return the destination URL for user 3rd-party app authorization
func (auth *AuthClient) GetAuthCode() string {
	scopes := strings.Join(auth.Config.Scopes, ",")
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", auth.Config.Id)
	v.Set("redirect_uri", auth.Config.Redirect)
	v.Set("scope", scopes)
	query := v.Encode()
	route := "authorize"

	url := fmt.Sprintf("https://%v/%v/?%v", auth.AppUrl, route, query)
	return url
}

// Retrieve User Authorization token
func (auth *AuthClient) GetAccessToken(code string) (string, error) {
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("client_id", auth.Config.Id)
	v.Set("client_secret", auth.Config.Secret)
	v.Set("redirect_uri", auth.Config.Redirect)
	data := v.Encode()
	body := bytes.NewBuffer([]byte(data))

	route := "auth/token"
	url := fmt.Sprintf("https://%v/%v", auth.AppUrl, route)

	res, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", err
	}

	accessResponse := Access{}
	err = GetEntity(res, &accessResponse)
	return accessResponse.Token, err
}
