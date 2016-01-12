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

// Retrieve the metadata of a specified folder
// Path and children are on demand fields
func (c *Client) GetFolderMetadata(folderId string, fields []string) (Folder, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "folders", folderId}, "/"),
	}

	if len(fields) > 0 {
		v := url.Values{"fields": fields}
		link.RawQuery = v.Encode()
	}

	folder := Folder{}
	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return folder, err
	}

	// TODO : The unmarshalling interface could be redefined for a folder to
	// include the ETag. However, there are so few instances of this it might not
	// be worth it
	folder.Etag = res.Header.Get("ETag")
	err = GetEntity(res, &folder)
	return folder, err
}

func (c *Client) GetFolderPath(folderId string) (*ParentPath, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "folders", folderId, "path"}, "/"),
	}

	pp := ParentPath{}
	res, err := c.get(link.String())
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	err = GetEntity(res, &pp)
	return &pp, err
}
