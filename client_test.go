package aerofs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

// These unit tests test against a local AeroFS Test Appliance instance.
// To execute these tests, tokens for OAuth 2.0 authentication must be provided
// This can be done manually, by creating a 3rd-Party Application, and using the
// AuthClient to generate corresponding tokens
// TODO : Implement a teardown method
const (
	// [files.read,files.write,acl.read,acl.write,acl.invitations,user.read,user.write]
	userToken = "2a09580d057348d9a1382b866389b1ae"

	// [files.read,files.write,acl.read,acl.write,acl.invitations,user.read,user.write,organization.admin]
	adminToken = "3d2a1005a27a4115946fe308eb30785f"

	// The default Hostname for the local test-appliance
	appHost = "share.syncfs.com"
)

// Perform test teardown and setup
func TestMain(m *testing.M) {
	rand.Seed(int64(os.Getpid()))
	os.Exit(m.Run())
}

// Create a new APIClient
func TestAPICreateClient(t *testing.T) {
	_, err := NewClient(adminToken, appHost)
	if err != nil {
		t.Fatal("Unable to create API client for testing")
	}
}

// User Unittests

// Create a new User
func TestAPI_CreateUser(t *testing.T) {
	c, _ := NewClient(adminToken, appHost)
	email := fmt.Sprintf("test_email%d@moria.com", rand.Intn(10000))
	firstName := "Gimli"
	lastName := "Son of Gloin"

	b, _, e := c.CreateUser(email, firstName, lastName)
	if e != nil {
		t.Log("Error when attempting to create a user")
		t.Fatal(e)
	}

	t.Log("Successfully created the following new user")
	desc := User{}
	json.Unmarshal(*b, &desc)
	t.Log(desc)
}

// List a set of Users
func TestAPI_ListUsers(t *testing.T) {
	c, _ := NewClient(adminToken, appHost)
	b, _, e := c.ListUsers(100, nil, nil)
	if e != nil {
		t.Log("Error when attempting to list users")
		t.Fatal(e)
	}

	t.Log("Successfully listed a set of users")
	desc := userListResponse{}
	json.Unmarshal(*b, &desc)
	t.Log(desc.Users)
}

// Update an existing user
// Create a user, update their credentials and ensure they match
func TestAPI_UpdateUser(t *testing.T) {
	c, _ := NewClient(adminToken, appHost)

	email := fmt.Sprintf("test_email%d@moria.com", rand.Intn(10000))
	origUser := User{email, "Gimli", "Son of Gloin", []SharedFolder{}, []Invitation{}}
	new_firstName := "Sarumon"
	new_lastName := "Of Isengard"

	_, _, e := c.CreateUser(email, origUser.FirstName, origUser.LastName)
	if e != nil {
		t.Log("Error when attempting to create a user")
		t.Fatal(e)
	}

	b, _, e := c.UpdateUser(email, new_firstName, new_lastName)
	if e != nil {
		t.Log("Error when attempting to update a user")
		t.Fatal(e)
	}

	newUser := User{}
	e = json.Unmarshal(*b, &newUser)
	if e != nil {
		t.Log("Error when attempting to unmarshal UserDescriptor")
		t.Fatal(e)
	}

	if reflect.DeepEqual(origUser, newUser) {
		t.Fatalf("New user %v is same from %v", newUser, origUser)
	}
	t.Log("New user %v is different from %v", newUser, origUser)
}

// Retrieve an uploadId for an existing File
func TestAPI_GetUploadId(t *testing.T) {
	c, _ := NewClient(userToken, appHost)
	data, _, err := c.ListFolderChildren("root")
	if err != nil {
		t.Fatal("Error retrieving list of root Children")
	}

	children := Children{}
	err = json.Unmarshal(*data, &children)
	if err != nil {
		t.Fatal("Error unmarshalling the children of root folder")
	}

	var fileId string
	var etag string
	for _, f := range children.Files {
		if f.Name == "appconfig.json" {
			fileId = f.Id
			etag = f.Etag
		}
	}

	t.Logf("FileId for appconfig.json is %s:%s", fileId, etag)
	uploadId, err := c.GetFileUploadId(fileId, []string{etag})
	if err != nil {
		t.Fatalf("Unable to get file upload id because of %s", err)
	}
	t.Logf("UploadId for appconfig.json is %s", uploadId)
}

// Verify a user is created and then Deleted
func TestAPI_CreateDeleteUser(t *testing.T) {

}

// Vertify an invitee is invited and then the invitation is deleted
func TestAPI_CreateDeleteInvitee(t *testing.T) {

}

// Verify an invitee is created and can be retrieved
func TestAPI_CreateGetInvitees(t *testing.T) {

}

/*
func TestB(t *testing.T) {
		// Retrieve single user
	b, _, e = c.GetUser("daniel.cardoza@aerofs.com")
	fmt.Println("GetUser")
	fmt.Println(string(*b))

	// Get invitation list
	b, _, e = c.GetInvitee("daniel.cardoza@aerofs.com")
	fmt.Println("GetInvitee")
	fmt.Println(e)
	fmt.Println(string(*b))

	// Create an invitaiton
	b, h, e := c.CreateInvitee("danielpcardoza@gmail.com", "daniel.cardoza@aerofs.com")
	fmt.Println("CreateInvitee")
	fmt.Println(e)
	fmt.Println(string(*b))
	fmt.Println(h)

	c.SetToken(userToken)
	e = c.DeleteInvitee("danielcardoza@gmail.com")
	fmt.Println("DeleteInvitee")
	fmt.Println(e)

	c.SetToken(adminToken)
	c.SetToken(userToken)

	// Get root folder data
	b, h, err = c.GetFolderMetadata("root", []string{"children"})
	fmt.Println("GetRootFolder")
	fmt.Println(e)
	fmt.Println(h)
	fmt.Println(string(*b))

	// Get root folder data
	b, h, err = c.ListFolderChildren("root")
	fmt.Println("GetRootFolderChildren")
	fmt.Println(err)
	fmt.Println(h)
	fmt.Println(string(*b))

	// Create a folder with root as parent
	b, h, err = c.CreateFolder("appdata", "Moria")
	fmt.Println("CreateFolder")
	fmt.Println(err)
	fmt.Println(h)
	fmt.Println(string(*b))

	// Create a sharedFolder
	b, h, err = c.CreateSharedFolder("TheShire_Shared")
	fmt.Println("CreateSharedFolder")
	fmt.Println(err)
	fmt.Println(h)
	fmt.Println(string(*b))
		pp, err := c.GetFolderPath("root")
		fmt.Println("GetFolderPath")
		fmt.Println(err)
		fmt.Println(pp)
}*/
