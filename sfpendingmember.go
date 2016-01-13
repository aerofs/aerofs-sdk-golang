package aerofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (c *Client) ListPendingMembers(sid string, etags []string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"shares", sid, "pending"}, "/")
	reqHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &reqHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetPendingMember(id, email string) (*[]byte, *http.Header,
	error) {
	route := strings.Join([]string{"shares", id, "pending", email}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) InviteToSharedFolder(sid, email string, permissions []string,
	note string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", sid, "pending"}, "/")
	link := c.getURL(route, "")

	newInvitee := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
		"note":        note,
	}
	data, err := json.Marshal(newInvitee)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal new invitee")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemovePendingMember(sid, email string) error {
	route := strings.Join([]string{"shares", sid, "pending", email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	_, _, err = unpackageResponse(res)
	return err
}
