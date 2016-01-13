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
	"net/url"
	"strings"
)

// Retrieve the metadata of a specified folder
// Path and children are on demand fields
func (c *Client) GetFolderMetadata(folderId string, fields []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId}, "/")
	query := ""
	if len(fields) > 0 {
		v := url.Values{"fields": fields}
		query = v.Encode()
	}
	link := c.getURL(route, query)

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetFolderPath(folderId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId, "path"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

func (c *Client) ListFolderChildren(folderId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId, "children"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	body, header := unpackageResponse(res)
	return body, header, err

}

func (c *Client) CreateFolder(parentId, name string) (*[]byte, *http.Header, error) {
	route := "folders"
	link := c.getURL(route, "")

	newFolder := map[string]string{
		"parent": parentId,
		"name":   name,
	}
	data, err := json.Marshal(newFolder)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal JSON for new folder")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

// Move a folder given its existing unique ID, the ID of its new parent and its
// new folder Name
func (c *Client) MoveFolder(folderId, newParentId, name string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"folders", folderId}, "/")
	link := c.getURL(route, "")

	content := map[string]string{"parent": newParentId, "name": name}
	data, err := json.Marshal(content)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal JSON for moving folder")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err

}

func (c *Client) DeleteFolder(folderId string, etags []string) error {
	route := strings.Join([]string{"folders", folderId}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DELETE", link, &newHeader)
	res.Body.Close()
	return err
}

func (c *Client) ShareFolder(folderId string) error {
	route := strings.Join([]string{"folders", folderId, "is_shared"}, "/")
	link := c.getURL(route, "")

	res, err := c.put(link, nil)
	res.Body.Close()
	return err
}
