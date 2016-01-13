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
	"strings"
)

// GroupMember Functions
func (c *Client) ListGroupMembers(groupId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) AddGroupMember(groupId, name string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members"}, "/")
	link := c.getURL(route, "")
	newMember := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal provided group member")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetGroupMember(groupId, email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemoveMember(groupId, email string) error {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
