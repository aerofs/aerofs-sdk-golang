package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	"net/http"
)

type SFMemberClient struct {
	APIClient *api.Client
	Desc      SFMember
	Etag      string
}

// A Shared Folder Member is a member of an AeroFS Shared Folder
type SFMember struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Permissions []string `json:"permissions"`
	Sid         string
}

// Retrieve a list of SharedFolder member descriptors
// TOD : Should an Etag be return for each one?
func ListSFMember(c *api.Client, sid string, etags []string) (*[]SFMember, error) {
	body, _, err := c.ListSFMembers(sid, etags)
	if err != nil {
		return nil, err
	}
	sfmembers := []SFMember{}
	err = json.Unmarshal(*body, &sfmembers)
	if err != nil {
		return nil, errors.New("Unable to demarshal the list of retrieved SharedFolder members")
	}
	for _, v := range sfmembers {
		v.Sid = sid
	}
	return &sfmembers, nil
}

// Return an existing SFMemberClient given its shared folder and user email
func GetSFMemberClient(c *api.Client, sid, email string, etags []string) (*SFMemberClient, error) {
	body, header, err := c.GetSFMember(sid, email, etags)
	if err != nil {
		return nil, err
	}
	sfmClient := SFMemberClient{APIClient: c, Etag: header.Get("ETag"), Desc: SFMember{Sid: sid, Email: email}}
	err = json.Unmarshal(*body, &sfmClient.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal retrieved SFMember")
	}

	return &sfmClient, nil
}

func (sfm *SFMemberClient) reserialize(buffer *[]byte, header *http.Header) error {
	err := json.Unmarshal(*buffer, &sfm.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal retrieved SFMember")
	}
	sfm.Etag = header.Get("ETag")
	return nil
}

// Update a SFMember's permissions
// TODO : Does it make sense for a user to modify their own?
func (sfm *SFMemberClient) UpdatePermissions(newPermissions []string) error {
	body, header, err := sfm.APIClient.SetSFMemberPermissions(sfm.Desc.Sid, sfm.Desc.Email,
		newPermissions, []string{sfm.Etag})
	if err != nil {
		return err
	}

	return sfm.reserialize(body, header)
}

// Retrieve configuration from backend
func (sfm *SFMemberClient) Load() error {
	body, header, err := sfm.APIClient.GetSFMember(sfm.Desc.Sid, sfm.Desc.Email,
		[]string{sfm.Etag})

	if err != nil {
		return err
	}

	return sfm.reserialize(body, header)
}
