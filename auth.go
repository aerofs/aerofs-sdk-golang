package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type AuthClient struct {
	Config AppConfig
	AppUrl string
}

// Return the destination URL for user 3rd-party app authorization
func (auth *AuthClient) GetAuthCode() string {
	scopes := strings.Join(auth.Config.Scopes, ",")
	params := fmt.Sprintf(
		"?response_type=code&"+
			"client_id=%v&"+
			"redirect_uri=%v"+
			"scope=%v", auth.Config.Id, auth.Config.Redirect, scopes)
	route := "authorize"
	url := fmt.Sprintf("https://%v/%v/%v", auth.AppUrl, route, params)
	return url
}

// Retrieve User Authorization token
func (auth *AuthClient) GetAcessToken(code string) (string, error) {
	params := fmt.Sprintf(
		"client_id=%v&"+
			"client_secret=%v&"+
			"redirect_uri=%v&"+
			"grant_type=authorization_code&"+
			"code=%v", auth.Config.Id, auth.Config.Secret, auth.Config.Redirect, code)
	route := "auth/token"
	url := fmt.Sprintf("https://%v/%v", auth.AppUrl, route)
	body := bytes.NewBuffer([]byte(params))

	res, err := http.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return "", err
	}

	accessResponse := Access{}
	err = GetEntity(res, &accessResponse)

	return accessResponse.Token, err
}
