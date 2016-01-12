package aerofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

func (c *Client) ListSharedFolders(email string, etags []string) (*[]SharedFolder, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users", email, "shares"}, "/"),
	}

	if len(etags) > 0 {
		v := url.Values{"If-None-Match": etags}
		link.RawQuery = v.Encode()
	}

	res, err := c.get(link.String())
	if err != nil {
		return nil, err
	}

	folderList := []SharedFolder{}
	err = GetEntity(res, &folderList)
	return &folderList, err
}

func (c *Client) ListSharedFolderMetadata(sid string, etags []string) (*SharedFolder, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "shares", sid}, "/"),
	}

	if len(etags) > 0 {
		v := url.Values{"If-None-Match": etags}
		link.RawQuery = v.Encode()
	}

	res, err := c.get(link.String())
	if err != nil {
		return nil, err
	}

	folder := SharedFolder{}
	err = GetEntity(res, &folder)
	return &folder, err
}

func (c *Client) CreatedSharedFolder(name string) (*SharedFolder, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "shares"}, "/"),
	}

	folder := SharedFolder{Name: name}
	data, err := json.Marshal(folder)
	if err != nil {
		return nil, errors.New("Unable to marshal SharedFolder with given string")
	}

	res, err := c.post(link.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	folder = SharedFolder{}
	err = GetEntity(res, &folder)
	return &folder, err
}
