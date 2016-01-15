package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

type FileClient struct {
	APIClient *api.Client
	Desc      File
	OnDemand  []string
}

type File struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Parent       string     `json:"parent"`
	LastModified string     `json:"last_modified"`
	Size         int        `json:"size"`
	Mime         string     `json:"mime_type"`
	Etag         string     `json:"etag"`
	Path         ParentPath `json:"path"`
	ContentState string     `json:"content_state"`
}

// Construct a FileClient given a file identifier and APIClient
func GetFileClient(c *api.Client, fileId string, fields []string) (*FileClient, error) {
	body, header, err := c.GetFileMetadata(fileId, fields)
	if err != nil {
		return nil, err
	}

	f := FileClient{APIClient: c, OnDemand: fields}
	err = json.Unmarshal(*body, &f.Desc)

	if err != nil {
		return nil, errors.New("Unable to unmarshal existing File")
	}
	f.Desc.Etag = header.Get("ETag")
	return &f, nil
}

//
func (f *FileClient) LoadPath() error {
	body, header, err := f.APIClient.GetFilePath(f.Desc.Id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &f.Desc.Path)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved file ParentPath")
	}

	f.Desc.Etag = header.Get("ETag")
	return nil
}

func (f *FileClient) Move(newName, parentId string) error {
	body, header, err := f.APIClient.MoveFile(f.Desc.Id, parentId, newName,
		[]string{f.Desc.Etag})
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &f.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal the new File location")
	}
	f.Desc.Etag = header.Get("Etag")
	return nil
}
