package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Retrieve array of Appliance users
// limit : The maximum number of entries returned
// after : An index to the first entry to be retrieved
// before: An index to the last possible entry to be retrieved
func (c *Client) ListUsers(limit int, after, before *int) (*[]byte, *http.Header,
	error) {
	route := "users"
	query := url.Values{}
	query.Set("limit", strconv.Itoa(limit))
	if before != nil {
		query.Set("before", strconv.Itoa(*before))
	}
	if after != nil {
		query.Set("after", strconv.Itoa(*after))
	}
	link := c.getURL(route, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Retrieve a single user
func (c *Client) GetUser(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Create a user with given email and name
func (c *Client) CreateUser(email, firstName, lastName string) (*[]byte,
	*http.Header, error) {
	route := "users"
	link := c.getURL(route, "")

	user := map[string]string{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
	}
	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) UpdateUser(email, firstName, lastName string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	user := map[string]string{}
	user["email"] = email
	if firstName != "" {
		user["first_name"] = firstName
	}
	if lastName != "" {
		user["last_name"] = lastName
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) DeleteUser(email string) error {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

func (c *Client) ChangePassword(email, password string) error {
	route := strings.Join([]string{"users", email, "password"}, "/")
	link := c.getURL(route, "")
	data := []byte(`"` + password + `"`)

	_, err := c.put(link, bytes.NewBuffer(data))
	return err
}

func (c *Client) DisablePassword(email string) error {
	route := strings.Join([]string{"users", email, "password"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

func (c *Client) CheckTwoFactorAuth(email string) error {
	route := strings.Join([]string{"users", email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.get(link)
	return err
}

func (c *Client) DisableTwoFactorAuth(email string) error {
	route := strings.Join([]string{"users", email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}
