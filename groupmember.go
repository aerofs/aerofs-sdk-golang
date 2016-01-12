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
	"net/url"
	"strings"
)

// GroupMember Functions
func (c *Client) ListGroupMembers(groupID string) (*[]GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	groupMembers := []GroupMember{}
	err = GetEntity(res, &groupMembers)
	return &groupMembers, err
}

func (c *Client) AddGroupMember(groupID string, groupMember GroupMember) (*GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members"}, "/"),
	}

	data, err := json.Marshal(groupMember)
	if err != nil {
		return nil, errors.New("Unable to marshal provided group member")
	}

	res, err := c.post(link.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newMember := GroupMember{}
	err = GetEntity(res, &newMember)
	return &newMember, err
}

func (c *Client) GetGroupMember(groupID, email string) (*GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members", email}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	groupMember := GroupMember{}
	err = GetEntity(res, &groupMember)
	return &groupMember, err
}

func (c *Client) RemoveMember(groupID, email string) error {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members", email}, "/"),
	}

	res, err := c.del(link.String())
	defer res.Body.Close()
	return err
}
