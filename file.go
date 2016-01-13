package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// File Related Operations
// TODO : include on-demand fields
func (c *Client) GetFileMetadata(fileId string, fields []string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"files", fileId}, "/")
	query := url.Values{"fields": fields}
	link := c.getURL(route, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetFilePath(fileId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId, "path"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

// Instantiate a newfile
func (c *Client) CreateFile(file File) (*File, error) {
	route := "files"
	link := c.getURL(route, "")

	data, err := json.Marshal(file)
	if err != nil {
		return nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newFile := File{}
	err = GetEntity(res, &newFile)
	newFile.Etag = res.Header.Get("ETag")

	return &newFile, err
}

func (c *Client) DeleteFile(fileid string, etags []string) error {
	route := strings.Join([]string{"files", fileid}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
