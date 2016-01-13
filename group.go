package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TODO : offset and results should be optional ( via pointers??)
func (c *Client) ListGroups(offset, results int) (*[]byte, *http.Header, error) {
	route := "groups"
	query := url.Values{}
	query.Set("offset", strconv.Itoa(offset))
	query.Set("results", strconv.Itoa(results))
	link := c.getURL(route, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Create a new user group given a new groupname
func (c *Client) CreateGroup(groupName string) (*[]byte, *http.Header, error) {
	route := "groups"
	link := c.getURL(route, "")
	// TODO : Is this preferred to constructing a map, then marshalling?
	// robust vs. bootstrap
	newGroup := []byte(fmt.Sprintf(`{"name" : %s}`, groupName))

	res, err := c.post(link, bytes.NewBuffer(newGroup))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Retrieve a group given a group identifier
func (c *Client) GetGroup(groupId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"request", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err

}

// Remove an existing group
func (c *Client) DeleteGroup(groupId string) error {
	route := strings.Join([]string{API, "groups", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
