package aerofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// List all associated groups for a shared folder with a given identifier
func (c *Client) ListSFGroups(sid string) (*[]byte, *http.Header, error) {
	path := strings.Join([]string{"shares", sid, "groups"}, "/")
	link := c.getURL(path, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

// Retrieve information for a group associated with a shared folder
// As of now, this only returns the new permissions associated with each group
// and the two argument
func (c *Client) GetSFGroups(sid, gid string) (*[]byte, *http.Header, error) {
	path := strings.Join([]string{"shares", sid, "members", gid}, "/")
	link := c.getURL(path, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

// Construct a new group for an existing Shared Folder
func (c *Client) AddGroupToSharedFolder(sid string, permissions []string) (*[]byte, *http.Header, error) {
	path := strings.Join([]string{"shares", sid, "groups"}, "/")
	link := c.getURL(path, "")
	reqBody := map[string]interface{}{
		"id":          sid,
		"permissions": permissions,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, errors.New(`Unable to marshal passed in SharedFolderGroupMember`)
	}
	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

// Modify the existing permissions of a group for an existing shared folder
func (c *Client) SetSFGroupPermissions(sid, gid string, permissions []string) (*SFGroupMember,
	error) {
	path := strings.Join([]string{"shares", sid, "groups", gid}, "/")
	link := c.getURL(path, "")

	permsList := PermissionList{Permissions: permissions}
	data, err := json.Marshal(permsList)
	if err != nil {
		return nil, errors.New("Unable to marshal given list of permissions")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newGroup := SFGroupMember{}
	err = GetEntity(res, &newGroup)
	return &newGroup, err
}

// Remove an existing group from its associated shared folder
func (c *Client) RemoveSFGroup(sid, gid string) error {
	path := strings.Join([]string{"shares", sid, "groups", gid}, "/")
	link := c.getURL(path, "")

	_, err := c.del(link)
	return err
}
