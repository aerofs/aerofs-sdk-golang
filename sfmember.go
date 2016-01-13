package aerofs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func (c *Client) ListSFMember(id string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members"}, "/")
	newHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) AddSFMember(id, email string, permissions []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members"}, "/")
	link := c.getURL(route, "")

	newMember := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) SetSFMemberPermissions(id, email string, permissions, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	newPerms := map[string]interface{}{
		"permissions": permissions,
	}
	data, err := json.Marshal(newPerms)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemoveSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DEL", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}
