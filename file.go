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

	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetFileContent(fileId, etag_IfRange string, fileRanges, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")

	newHeader := http.Header{}
	if len(fileRanges) > 0 {
		for _, v := range fileRanges {
			newHeader.Add("Range", v)
		}
	}

	if etag_IfRange != "" {
		newHeader.Set("If-Range", etag_IfRange)
	}

	if len(etags) > 0 {
		for _, v := range fileRanges {
			newHeader.Add("If-None-Match", v)
		}
	}

	res, err := c.request("GET", link, &newHeader, nil)
	body, header, err := unpackageResponse(res)
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

func (c *Client) MoveFile(fileId, parentId, name string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-Match": etags}
	newFile := map[string]interface{}{
		"parent": parentId,
		"name":   name,
	}

	data, err := json.Marshal(newFile)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) DeleteFile(fileid string, etags []string) error {
	route := strings.Join([]string{"files", fileid}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DEL", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	_, _, err = unpackageResponse(res)
	return err
}
