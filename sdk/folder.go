package aerofs

import (
	"encoding/json"
	"errors"
)

type FolderClient struct {
	APIClient *Client
	Desc      Folder

	// OnDemand fields must be explicitly stated in requests to retrieve items
	// For folders, this is specifically ParentPath, Children objects
	OnDemand []string
}

type Folder struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Parent    string     `json:"parent"`
	IsShared  bool       `json:"is_shared"`
	Sid       string     `json:"sid"`
	Path      ParentPath `json:"path"`
	ChildList Children   `json:"children}`
	Etag      string
}

// Return an existing FolderClient given a folderId and on-demand fields
func GetFolderClient(c *Client, folderId string, fields []string) (*FolderClient, error) {
	body, header, err := c.GetFolderMetadata(folderId, fields)
	if err != nil {
		return nil, err
	}

	f := FolderClient{APIClient: c, OnDemand: fields}
	err = json.Unmarshal(*body, &f.Desc)

	if err != nil {
		return nil, errors.New("Unable to unmarshal existing Folder")
	}
	f.Desc.Etag = header.Get("ETag")

	return &f, nil
}

// Load the most up to date path from the server
func (f *FolderClient) LoadPath() error {
	body, _, err := f.APIClient.GetFolderPath(f.Desc.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &f.Desc.Path)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved folder ParentPath")
	}

	return nil
}

// Load new Folder children from the server
func (f *FolderClient) LoadChildren() error {
	body, _, err := f.APIClient.GetFolderChildren(f.Desc.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &f.Desc.ChildList)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved folder ParentPath")
	}

	return nil
}

// Load new Folder metadata from the server
func (f *FolderClient) LoadMetadata() error {
	body, _, err := f.APIClient.GetFolderMetadata(f.Desc.Id, f.OnDemand)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &f.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved folder ParentPath")
	}

	return nil
}

// Update all Folder descriptor fields
func (f *FolderClient) Load() error {
	// Perform in a single call by retrieving all fields by setting the On-Demand
	// fields. This only performs one request vs. 3 by calling
	// load{Metadata,Path,Children}
	oldFields := f.OnDemand
	f.OnDemand = []string{"path", "children"}
	err := f.LoadMetadata()

	f.OnDemand = oldFields
	return err
}

// Delete the Folder
func (f *FolderClient) Delete() error {
	return f.APIClient.DeleteFolder(f.Desc.Id, []string{f.Desc.Etag})
}

// Move the existing folder to a new location
func (f *FolderClient) Move(newName, parentId string) error {
	body, header, err := f.APIClient.MoveFolder(f.Desc.Id, parentId, newName, []string{f.Desc.Etag})
	if err != nil {
		return err
	}
	err = json.Unmarshal(*body, &f.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal the new Folder location")
	}
	f.Desc.Etag = header.Get("ETag")
	return nil
}

// Share a folder
// Share and then update the folder to retrieve new SID, is_shared value
func (f *FolderClient) Share() error {
	err := f.APIClient.ShareFolder(f.Desc.Id)
	if err != nil {
		return err
	}
	return f.LoadMetadata()
}
