package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication

import (
	/*	"encoding/json"
		"io/ioutil" */
	"net/http"
)

const (
	API = "v1.3"
)

//
type Client struct {
	AppUrl string
	Token  string //Oauth Token
}

// Create a client
func NewClient(token, appUrl string) (*Client, error) {
	c := Client{appUrl, token}
	return &c, nil
}

// Retrieve array of Appliance users
func (c *Client) ListUsers(limit int) ([]User, error) {
	route := "users"
	params := "limit=" + string(limit)
	url := c.getPrefix() + route + params
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return []User{}, err
	}
	/*
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []User{}, err
		}
	*/
	user := []User{}
	err = GetEntity(res, &user)
	//err = json.Unmarshal(data, &user)
	return user, err
}

// Return the base URL used for all API routes
func (c *Client) getPrefix() string {
	return "https://" + c.AppUrl + "/api/" + API + "/"
}
