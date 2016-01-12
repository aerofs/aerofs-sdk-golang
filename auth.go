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
	// Parameters required for network requests
	Config AuthConfig

	// The URL of an AeroFS Appliance instance
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

	url := fmt.Sprintf("https://%v/%v?%v", auth.AppUrl, route, query)
	return url
}

// Retrieve User OAuth token, scopes given an Authorization code
func (auth *AuthClient) GetAccessToken(code string) (string, []string, error) {
	// Regular API requests are 'application/json' but authorization uses
	// Urlencoding
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("client_id", auth.Config.Id)
	v.Set("client_secret", auth.Config.Secret)
	v.Set("redirect_uri", auth.Config.Redirect)

	link := url.URL{Scheme: "https",
		Host: auth.AppUrl,
		Path: strings.Join([]string{"auth", "token"}, "/"),
	}
	body := bytes.NewBuffer([]byte(v.Encode()))
	encoding := "application/x-www-form-urlencoded"

	res, err := http.Post(link.String(), encoding, body)
	if err != nil {
		return "", []string{}, err
	}

	accessResponse := Access{}
	err = GetEntity(res, &accessResponse)
	grantedScopes := strings.Split(accessResponse.Scopes, ",")
	return accessResponse.Token, grantedScopes, err
}
