package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"net/url"
	"strings"
)

// File Related Operations
func (c *Client) GetFileMetadata(fileId string) (*File, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "files", fileId}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	file := File{}
	file.Etag = res.Header.Get("ETag")
	err = GetEntity(res, &file)
	return &file, err
}

func (c *Client) GetFilePath(fileId string) (*File, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "files", fileId, "path"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	file := File{}
	err = GetEntity(res, &file)
	return &file, err
}
