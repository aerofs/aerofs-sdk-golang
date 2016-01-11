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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	API = "api/v1.3"
)

type Client struct {
	// The hostname/IP of the AeroFS Appliance
	// Used when constructing the default API Prefix for all subsequent API calls
	// Ie. share.syncfs.com
	Host string

	// The OAuth token
	Token string

	// Contains the authorization token
	// For conditional file, and folder requests, the header is populated
	// with an ETag
	Header http.Header
}

// SDK-Client Constructor
// Constructs the HTTP header used for subsequent requests
// OAuth token stored in HTTP header
func NewClient(token, host string) (*Client, error) {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("Endpoint-Consistency", "strict")

	c := Client{Host: host,
		Header: header,
		Token:  token}

	return &c, nil
}

// Resets the token for a given client
// Allows the third-party developer to construct 1 SDK-Client used to retrieve
// the values for multiple users
func (c *Client) SetToken(token string) {
	c.Header.Set("Authorization", "Bearer "+token)
}

// Wrappers for basic HTTP functions
// DELETE, PUT can only be performed by an httpClient
func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP GET Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

func (c *Client) post(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP POST request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

func (c *Client) del(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP DELETE Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

// USER API Calls
// The format is to construct a URL, marshal the request
// and unmarshal the parsed HTTP response

// Retrieve array of Appliance users
func (c *Client) ListUsers(limit int) (*[]User, error) {
	query := url.Values{"limit": []string{string(limit)}}
	link := url.URL{Scheme: "https",
		Host:     c.Host,
		Path:     strings.Join([]string{API, "users"}, "/"),
		RawQuery: query.Encode(),
	}
	fmt.Println(link.String())
	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	listResponse := ListUserResponse{}
	err = GetEntity(res, &listResponse)
	return &listResponse.Users, err
}

// Retrieve a single user
func (c *Client) GetUser(email string) (*User, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users", email}, "/"),
	}
	fmt.Println(url)

	res, err := c.get(url.String())
	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	user := User{}
	err = GetEntity(res, &user)
	return &user, err

}

// TODO : Should Shares, Invitations be ignored?
func (c *Client) CreateUser(user User) (*User, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users"}, "/"),
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	res, err := c.post(url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newUser := User{}
	err = GetEntity(res, &newUser)
	return &newUser, err
}

// Retrieve information about a person invited to an AeroFS instance
func (c *Client) GetInvitee(email string) (*Invitee, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees", email}, "/"),
	}

	res, err := c.get(url.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	invitee := new(Invitee)
	err = GetEntity(res, invitee)
	return invitee, err
}

// Create an invitation
func (c *Client) CreateInvitee(email_to, email_from string) (*Invitee, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees"}, "/"),
	}

	newInvitee := new(Invitee)
	invitee := new(Invitee)
	invitee.EmailTo = email_to
	invitee.EmailFrom = email_from

	data, err := json.Marshal(invitee)
	if err != nil {
		return newInvitee, errors.New("Unable to serialize invitation request")
	}

	res, err := c.post(url.String(), bytes.NewBuffer(data))
	if err != nil {
		return newInvitee, err
	}

	err = GetEntity(res, newInvitee)
	return newInvitee, err
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees", email}, "/"),
	}

	res, err := c.del(url.String())
	defer res.Body.Close()
	return err
}

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

// Group related Operations

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

// GroupMember Functions
func (c *Client) ListGroupMembers(groupID string) (*[]GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	groupMembers := []GroupMember{}
	err = GetEntity(res, &groupMembers)
	return &groupMembers, err
}

func (c *Client) AddGroupMember(groupID string, groupMember GroupMember) (*GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members"}, "/"),
	}

	data, err := json.Marshal(groupMember)
	if err != nil {
		return nil, errors.New("Unable to marshal provided group member")
	}

	res, err := c.post(link.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	newMember := GroupMember{}
	err = GetEntity(res, &newMember)
	return &newMember, err
}

func (c *Client) GetGroupMember(groupID, email string) (*GroupMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members", email}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	groupMember := GroupMember{}
	err = GetEntity(res, &groupMember)
	return &groupMember, err
}

func (c *Client) RemoveMember(groupID, email string) error {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "groups", groupID, "members", email}, "/"),
	}

	res, err := c.del(link.String())
	defer res.Body.Close()
	return err
}
