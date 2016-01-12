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

func (c *Client) ListGroups(offset, results int) ([]Group, error) {
	query := url.Values{}
	query.Set("offset", string(offset))
	query.Set("results", string(results))

	link := url.URL{Scheme: "https",
		Host:     c.Host,
		Path:     strings.Join([]string{API, "groups"}, "/"),
		RawQuery: query.Encode(),
	}
	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	Groups := []Group{}
	err = GetEntity(res, &Groups)
	return Groups, err
}

// Create a new user group given a new groupname
func (c *Client) CreateGroup(groupName string) (*Group, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	group := new(Group)
	err = GetEntity(res, group)
	return group, err
}

func (c *Client) GetGroup(groupID string) (*Group, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "request", groupID}, "/"),
	}

	res, err := c.post(link.String(), nil)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	group := new(Group)
	err = GetEntity(res, group)
	return group, err
}

func (c *Client) DeleteGroup(groupID string) error {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID}, "/"),
	}

	res, err := c.del(link.String())
	defer res.Body.Close()
	return err
}
