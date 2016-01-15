package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// This file maps all routes exposed on the AeroFS API

// User Related Calls

func (c *Client) ListUsers(limit int, after, before *string) (*[]byte, *http.Header, error) {
	route := "users"
	query := url.Values{}
	query.Set("limit", strconv.Itoa(limit))
	if before != nil {
		query.Set("before", *before)
	}
	if after != nil {
		query.Set("after", *after)
	}

	link := c.getURL(route, query.Encode())
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetUser(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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
	return unpackageResponse(res)
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
	defer res.Body.Close()

	return unpackageResponse(res)
}

func (c *Client) DeleteUser(email string) error {
	route := strings.Join([]string{"users", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	res.Body.Close()

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

	return unpackageResponse(res)
}

func (c *Client) ViewPendingSFInvitation(email, sid string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	route := strings.Join([]string{"invitees", email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	defer res.Body.Close()
	_, _, err = unpackageResponse(res)
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

func (c *Client) GetGroup(groupId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"request", groupId}, "/")
	link := c.getURL(route, "")

	res, err := c.post(link, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

func (c *Client) GetGroupMember(groupId, email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) RemoveMember(groupId, email string) error {
	route := strings.Join([]string{"groups", groupId, "members", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	_, _, err = unpackageResponse(res)
	return err
}

// File related calls

func (c *Client) GetFileMetadata(fileId string, fields []string) (*[]byte,
	*http.Header, error) {
	route := strings.Join([]string{"files", fileId}, "/")
	query := url.Values{"fields": fields}
	link := c.getURL(route, query.Encode())

	newHeader := http.Header{}
	newHeader.Set("Content-Type", "application/octet-stream")

	res, err := c.request("GET", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetFilePath(fileId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId, "path"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

//
func (c *Client) GetFileContent(fileId, rangeEtag string, fileRanges, matchEtags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")

	// Construct header
	newHeader := http.Header{}
	if len(fileRanges) > 0 {
		for _, v := range fileRanges {
			newHeader.Add("Range", v)
		}
	}

	if rangeEtag != "" {
		newHeader.Set("If-Range", rangeEtag)
	}

	if len(matchEtags) > 0 {
		for _, v := range matchEtags {
			newHeader.Add("If-None-Match", v)
		}
	}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Instantiate a newfile
func (c *Client) CreateFile(parentId, fileName string) (*[]byte, *http.Header,
	error) {
	route := "files"
	link := c.getURL(route, "")

	newFile := map[string]string{
		"parent": parentId,
		"name":   fileName,
	}
	data, err := json.Marshal(newFile)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Functions that upload content to a file

// Retrieve a files UploadID to be used for future content uploads
// Upload Identifiers are only valid for ~24 hours
func (c *Client) GetFileUploadId(fileId string, etags []string) (string, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	newHeader.Set("Content-Range", "bytes */*")
	newHeader.Set("Content-Length", "0")
	if len(etags) > 0 {
		for _, v := range etags {
			newHeader.Add("If-Match", v)
		}
	}

	res, err := c.request("PUT", link, &newHeader, nil)
	defer res.Body.Close()

	_, h, err := unpackageResponse(res)
	if err != nil {
		return "", err
	}

	fmt.Println(h)
	return h.Get("Upload-ID"), nil
}

// Retrieve the list of bytes already transferred by an unfinished upload
func (c *Client) GetUploadBytesSize(fileId, uploadId string, etags []string) (int, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	newHeader.Set("Content-Range", "bytes /*/")
	newHeader.Set("Upload-ID", uploadId)
	newHeader.Set("Content-Length", "0")
	if len(etags) > 0 {
		for _, v := range etags {
			newHeader.Add("If-Match", v)
		}
	}

	res, err := c.request("PUT", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return 0, nil
	}

	bytesUploaded, err := strconv.Atoi(res.Header.Get("Range"))
	if err != nil {
		return 0, errors.New("Unable to parse value of bytes transferred from HTTP-Header")
	}
	return bytesUploaded, err
}

// Upload a single file chunk
func (c *Client) UploadFileChunk(fileId, uploadId string, chunks *[]byte, startIndex, lastIndex int) (*http.Header, error) {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")
	byteRange := fmt.Sprintf("bytes %d-%d/*", startIndex, lastIndex)
	newHeader := http.Header{}
	newHeader.Set("Content-Range", byteRange)
	newHeader.Set("Upload-ID", uploadId)
	newHeader.Set("Content-Length", "0")

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(*chunks))
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	_, h, e := unpackageResponse(res)
	return h, e
}

// Upload a file
func (c *Client) UploadFile(fileId, uploadId string, file io.Reader, etags []string) error {
	route := strings.Join([]string{"files", fileId, "content"}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-Match": etags}
	newHeader.Set("Upload-ID", uploadId)

	startIndex := 0
	endIndex := 0
	chunk := make([]byte, CHUNKSIZE)

	// Iterate over the file in CHUNKSIZE pieces until EOF occurs
	// Do not set "Instance-Length" of "Content-Range" until we hit EOF and use
	// the format "bytes <startIndex>-<endIndex>/* for intermediary uploads
ChunkLoop:
	for {
		size, fileErr := file.Read(chunk)
		endIndex += size - 1

		switch {
		// If we have read all chunks, set Content-Range instance-Length
		case fileErr == io.EOF:
			byteRange := fmt.Sprintf("bytes %d-%d/%d", startIndex, endIndex, endIndex)
			newHeader.Set("Content-Range", byteRange)
		case fileErr != nil:
			return fileErr
		case fileErr == nil:
			byteRange := fmt.Sprintf("bytes %d-%d/*", startIndex, endIndex)
			newHeader.Set("Content-Range", byteRange)
		}

		res, httpErr := c.request("PUT", link, &newHeader, bytes.NewBuffer(chunk))
		if httpErr != nil {
			return httpErr
		}
		defer res.Body.Close()
		if fileErr == io.EOF {
			break ChunkLoop
		}

		startIndex += size
	}
	return nil
}

func (c *Client) MoveFile(fileId, parentId, name string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"files", fileId}, "/")
	link := c.getURL(route, "")

	newHeader := http.Header{}
	if len(etags) > 0 {
		for _, v := range etags {
			newHeader.Add("If-Match", v)
		}
	}

	newFile := map[string]string{
		"parent": parentId,
		"name":   name,
	}

	data, err := json.Marshal(newFile)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal the given file")
	}

	res, err := c.request("PUT", link, &newHeader, bytes.NewBuffer(data))
	return unpackageResponse(res)
}

// There must be at least one etag present
func (c *Client) DeleteFile(fileid string, etags []string) error {
	route := strings.Join([]string{"files", fileid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-Match": etags}

	res, err := c.request("DEL", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	_, _, err = unpackageResponse(res)
	return err
}

// Folder calls

// Note
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

	return unpackageResponse(res)
}

func (c *Client) GetFolderPath(folderId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId, "path"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetFolderChildren(folderId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId, "children"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

// Move a folder given its existing unique ID, the ID of its new parent and its
// new folder Name
func (c *Client) MoveFolder(folderId, newParentId, newFolderName string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"folders", folderId}, "/")
	link := c.getURL(route, "")

	content := map[string]string{"parent": newParentId, "name": newFolderName}
	data, err := json.Marshal(content)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal JSON for moving folder")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) DeleteFolder(folderId string, etags []string) error {
	route := strings.Join([]string{"folders", folderId}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DELETE", link, &newHeader, nil)
	defer res.Body.Close()
	if err != nil {
		_, _, err = unpackageResponse(res)
	}

	return err
}

func (c *Client) ShareFolder(folderId string) error {
	route := strings.Join([]string{"folders", folderId, "is_shared"}, "/")
	link := c.getURL(route, "")

	res, err := c.put(link, nil)
	defer res.Body.Close()
	if err != nil {
		_, _, err = unpackageResponse(res)
	}

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

	return unpackageResponse(res)
}

func (c *Client) ListSharedFolderMetadata(sid string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", sid}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{"If-None-Match": etags}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}
	return unpackageResponse(res)
}

func (c *Client) CreateSharedFolder(name string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares"}, "/")
	link := c.getURL(route, "")
	data := []byte(fmt.Sprintf(`{"name" : %s"}`, name))

	res, err := c.post(link, bytes.NewBuffer(data))

	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// SharedFolder Member calls

func (c *Client) ListSFMembers(id string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members"}, "/")
	newHeader := http.Header{}
	if len(etags) > 0 {
		newHeader = http.Header{"If-None-Match": etags}
	}
	link := c.getURL(route, "")

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	link := c.getURL(route, "")
	newHeader := http.Header{}
	if len(etags) > 0 {
		newHeader = http.Header{"If-None-Match": etags}
	}

	res, err := c.request("GET", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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
		return nil, nil, errors.New("Unable to marshal new ShareFolder member")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

func (c *Client) RemoveSFMember(id, email string, etags []string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"shares", id, "members", email}, "/")
	newHeader := http.Header{"If-Match": etags}
	link := c.getURL(route, "")

	res, err := c.request("DEL", link, &newHeader, nil)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

func (c *Client) GetPendingMember(id, email string) (*[]byte, *http.Header,
	error) {
	route := strings.Join([]string{"shares", id, "pending", email}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

func (c *Client) RemovePendingMember(sid, email string) error {
	route := strings.Join([]string{"shares", sid, "pending", email}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()

	if err == nil {
		_, _, err = unpackageResponse(res)
	}
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}

// Modify the existing permissions of a group for an existing shared folder
func (c *Client) SetSFGroupPermissions(sid, gid string, permissions []string) (*[]byte, *http.Header, error) {
	path := strings.Join([]string{"shares", sid, "groups", gid}, "/")
	link := c.getURL(path, "")

	permsList := map[string][]string{
		"permissions": permissions,
	}
	data, err := json.Marshal(permsList)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal given list of permissions")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

// Remove an existing group from its associated shared folder
func (c *Client) RemoveSFGroup(sid, gid string) error {
	path := strings.Join([]string{"shares", sid, "groups", gid}, "/")
	link := c.getURL(path, "")

	res, err := c.del(link)
	if err == nil {
		_, _, err = unpackageResponse(res)
	}
	return err
}

// Device specific API Calls

func (c *Client) _ListDevices(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "devices"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetDeviceMetadata(deviceId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceId}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) UpdateDevice(deviceName string) (*[]byte, *http.Header, error) {
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

	return unpackageResponse(res)
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

	return unpackageResponse(res)
}
