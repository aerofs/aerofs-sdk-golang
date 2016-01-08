package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	API = "v1.3"
)

//
type Client struct {
	AppUrl string
	Prefix string
	Token  string //Oauth Token
}

// Wrapper around HTTP Delete

// Create a client
func NewClient(token, appUrl string) (*Client, error) {
	prefix := "https://" + appUrl + "/api/" + API
	c := Client{appUrl, prefix, token}
	return &c, nil
}

// Retrieve array of Appliance users
func (c *Client) ListUsers(limit int) (*[]User, error) {
	route := "users"
	params := "limit=" + string(limit)
	url := strings.Join([]string{c.Prefix, route, params}, "/")
	listResponse := ListUserResponse{}
	users := &listResponse.Users

	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return users, err
	}

	err = GetEntity(res, &users)
	return users, err
}

// Retrieve a single user
func (c *Client) GetUser(email string) (*User, error) {
	route := "users"
	params := email
	url := strings.Join([]string{c.Prefix, route, params}, "/")
	user := new(User)

	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return user, err
	}

	err = GetEntity(res, &user)
	return user, err

}

// TODO : Should Shares, Invitations be ignored?
func (c *Client) CreateUser(user User) (*User, error) {
	route := "users"
	url := strings.Join([]string{c.Prefix, route}, "/")
	newUser := new(User)

	data, err := json.Marshal(newUser)
	if err != nil {
		return newUser, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return newUser, err
	}

	err = GetEntity(res, newUser)
	return newUser, err
}

func (c *Client) GetInvitee(email string) (*Invitee, error) {
	route := "invitees"
	params := email
	url := strings.Join([]string{c.Prefix, route, params}, "/")
	invitee := new(Invitee)

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return invitee, err
	}

	err = GetEntity(res, invitee)
	return invitee, err
}

func (c *Client) CreateInvitee(email_to, email_from string) (*Invitee, error) {
	route := "invitees"
	url := strings.Join([]string{c.Prefix, route}, "/")
	invitee := new(Invitee)
	invitee.EmailTo = email_to
	invitee.EmailFrom = email_from
	newInvitee := new(Invitee)

	data, err := json.Marshal(invitee)
	if err != nil {
		return newInvitee, err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return newInvitee, err
	}

	err = GetEntity(res, newInvitee)
	return newInvitee, err
}
