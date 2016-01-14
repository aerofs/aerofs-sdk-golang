package aerofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// User Related Calls

func (c *Client) ListUsers(limit int, after, before *int) (*[]byte, *http.Header, error) {
	route := "users"
	query := url.Values{}
	query.Set("limit", strconv.Itoa(limit))
	if before != nil {
		query.Set("before", strconv.Itoa(*before))
	}
	if after != nil {
		query.Set("after", strconv.Itoa(*after))
	}

	link := c.getURL(route, query.Encode())
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetUser(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) CreateUser(email, firstName, lastName string) (*[]byte,
	*http.Header, error) {
	route := "users"
	link := c.getURL(route, "")

	user := map[string]string{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
	}
	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	//	body, header, err := unpackageResponse(res)
	return unpackageResponse(res)
	//body, header, err
}

func (c *Client) UpdateUser(email, firstName, lastName string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	user := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal User data")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) DeleteUser(email string) error {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

func (c *Client) ChangePassword(email, password string) error {
	route := strings.Join([]string{"users", email, "password"}, "/")
	link := c.getURL(route, "")
	data := []byte(`"` + password + `"`)

	_, err := c.put(link, bytes.NewBuffer(data))
	return err
}

func (c *Client) DisablePassword(email string) error {
	route := strings.Join([]string{"users", email, "password"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

func (c *Client) CheckTwoFactorAuth(email string) error {
	route := strings.Join([]string{"users", email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.get(link)
	return err
}

func (c *Client) DisableTwoFactorAuth(email string) error {
	route := strings.Join([]string{"users", email, "two_factor"}, "/")
	link := c.getURL(route, "")

	_, err := c.del(link)
	return err
}

// Shared Folder Invitation Calls

func (c *Client) ListSFInvitations(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) ViewPendingSFInvitation(email, sid string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) AcceptSFInvitation(email, sid string, external int) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	query := url.Values{}
	query.Set("external", strconv.Itoa(external))
	link := c.getURL(route, query.Encode())

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Ignore an existing invitation to a shared folder
func (c *Client) IgnoreSFInvitation(email, sid string) error {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}

// Inviteee Related Calls

func (c *Client) GetInvitee(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"invitees", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) CreateInvitee(email_to, email_from string) (*[]byte,
	*http.Header, error) {
	route := "invitees"
	link := c.getURL(route, "")
	invitee := map[string]string{
		"email_to":   email_to,
		"email_from": email_from,
	}

	data, err := json.Marshal(invitee)
	if err != nil {
		return nil, nil, errors.New("Unable to serialize invitation request")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	route := strings.Join([]string{"invitees", email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	defer res.Body.Close()
	return err
}

// Group related calls

func (c *Client) ListGroups(offset, results int) (*[]byte, *http.Header, error) {
	route := "groups"
	query := url.Values{}
	query.Set("offset", strconv.Itoa(offset))
	query.Set("results", strconv.Itoa(results))
	link := c.getURL(route, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) CreateGroup(groupName string) (*[]byte, *http.Header, error) {
	route := "groups"
	link := c.getURL(route, "")
	// TODO : Is this preferred to constructing a map, then marshalling?
	// robust vs. bootstrap
	newGroup := []byte(fmt.Sprintf(`{"name" : %s}`, groupName))

	res, err := c.post(link, bytes.NewBuffer(newGroup))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetGroup(groupId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"request", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err

}

func (c *Client) DeleteGroup(groupId string) error {
	route := strings.Join([]string{API, "groups", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}

// GroupMember calls

func (c *Client) ListGroupMembers(groupId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) AddGroupMember(groupId, name string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members"}, "/")
	link := c.getURL(route, "")
	newMember := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal provided group member")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetGroupMember(groupId, email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemoveMember(groupId, email string) error {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}

// File related calls

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

// Folder calls

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

	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) ListFolderChildren(folderId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId, "children"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
	return body, header, err

}

func (c *Client) DeleteFolder(folderId string, etags []string) error {
	route := strings.Join([]string{"folders", folderId}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DELETE", link, &newHeader, nil)
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

// SharedFolder calls

func (c *Client) ListSharedFolders(email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "shares"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) ListSharedFolderMetadata(sid string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", sid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) CreateSharedFolder(name string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares"}, "/")
	link := c.getURL(route, "")
	data := []byte(fmt.Sprintf(`{"name" : %s"}`, name))

	res, err := c.post(link, bytes.NewBuffer(data))

	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// SharedFolder Member calls

func (c *Client) ListSFMember(id string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members"}, "/")
	newHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) AddSFMember(id, email string, permissions []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members"}, "/")
	link := c.getURL(route, "")

	newMember := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
	}
	data, err := json.Marshal(newMember)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) SetSFMemberPermissions(id, email string, permissions, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	newPerms := map[string]interface{}{
		"permissions": permissions,
	}
	data, err := json.Marshal(newPerms)
	if err != nil {
		return nil, nil, err
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemoveSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DEL", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// SharedFolder pending member calls

func (c *Client) ListPendingMembers(sid string, etags []string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"shares", sid, "pending"}, "/")
	reqHeader := http.Header{"If-None-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &reqHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetPendingMember(id, email string) (*[]byte, *http.Header,
	error) {
	route := strings.Join([]string{"shares", id, "pending", email}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// SharedFolder Invitation calls

func (c *Client) InviteToSharedFolder(sid, email string, permissions []string,
	note string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", sid, "pending"}, "/")
	link := c.getURL(route, "")

	newInvitee := map[string]interface{}{
		"email":       email,
		"permissions": permissions,
		"note":        note,
	}
	data, err := json.Marshal(newInvitee)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal new invitee")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) RemovePendingMember(sid, email string) error {
	route := strings.Join([]string{"shares", sid, "pending", email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	_, _, err = unpackageResponse(res)
	return err
}

// SharedFolder Group Member calls

// List all associated groups for a shared folder with a given identifier
func (c *Client) ListSFGroups(sid string) (*[]byte, *http.Header, error) {
	path := strings.Join([]string{"shares", sid, "groups"}, "/")
	link := c.getURL(path, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
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

	body, header, err := unpackageResponse(res)
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

	res, err := c.del(link)
	_, _, err = unpackageResponse(res)
	return err
}

// Device specific API Calls

func (c *Client) ListDevices(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "devices"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetDeviceMetadata(deviceId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceId}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) UpdateDeviceMetadata(deviceName string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceName}, "/")
	link := c.getURL(route, "")
	newDevice := map[string]string{
		"name": deviceName,
	}

	data, err := json.Marshal(newDevice)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal new device")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetDeviceStatus(deviceId string) (*[]byte, *http.Header,
	error) {
	route := strings.Join([]string{"devices", deviceId, "status"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	body, header, err := unpackageResponse(res)
	return body, header, err
}
