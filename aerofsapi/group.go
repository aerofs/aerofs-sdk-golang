package aerofsapi

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// This file maps all routes exposed on the AeroFS API

// Group related calls

func (c *Client) ListGroups(offset, results int) ([]byte, *http.Header, error) {
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

	return unpackageResponse(res)
}

func (c *Client) CreateGroup(groupName string) ([]byte, *http.Header, error) {
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

	return unpackageResponse(res)
}

func (c *Client) GetGroup(groupId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"request", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) DeleteGroup(groupId string) error {
	route := strings.Join([]string{API, "groups", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
