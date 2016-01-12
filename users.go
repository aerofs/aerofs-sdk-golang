package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// USER API Calls
// The format is to construct a URL, marshal the request
// and unmarshal the parsed HTTP response

// Retrieve array of Appliance users
func (c *Client) ListUsers(limit int, after, before *int) (*[]User, error) {
	query := url.Values{"limit": []string{string(limit)}}
	if before != nil {
		query.Add("before", string(*before))
	}
	if after != nil {
		query.Add("after", string(*after))
	}

	link := url.URL{Scheme: "https",
		Host:     c.Host,
		Path:     strings.Join([]string{API, "users"}, "/"),
		RawQuery: query.Encode(),
	}
	fmt.Println(link.String())
	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	listResponse := ListUserResponse{}
	err = GetEntity(res, &listResponse)
	return &listResponse.Users, err
}

// Retrieve a single user
func (c *Client) GetUser(email string) (*User, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users", email}, "/"),
	}
	fmt.Println(url)

	res, err := c.get(url.String())
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	user := User{}
	err = GetEntity(res, &user)
	return &user, err

}

// TODO : Should Shares, Invitations be ignored?
func (c *Client) CreateUser(user User) (*User, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users"}, "/"),
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	res, err := c.post(url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newUser := User{}
	err = GetEntity(res, &newUser)
	return &newUser, err
}
