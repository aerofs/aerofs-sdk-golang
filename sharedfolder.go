package aerofs

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) ListSharedFolders(email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "shares"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) ListSharedFolderMetadata(sid string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", sid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) CreateSharedFolder(name string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares"}, "/")
	link := c.getURL(route, "")
	data := []byte(fmt.Sprintf(`{"name" : %s"}`, name))

	res, err := c.post(link, bytes.NewBuffer(data))

	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}
