package aerofs

import (
	"fmt"
	"testing"
)

const (
	userToken  = "2a09580d057348d9a1382b866389b1ae"
	adminToken = "3d2a1005a27a4115946fe308eb30785f"
)

func TestCreateClient(t *testing.T) {
	_, err := NewClient(adminToken, "share.syncfs.com")
	if err != nil {
		t.Fatal("Unable to create API client for testing")
	}
}

func TestListUsers(t *testing.T) {
	c, _ := NewClient(adminToken, "share.syncfs.com")
	b, _, e := c.ListUsers(100, nil, nil)
	if e != nil {
		fmt.Println("Error retrieving list of users")
		fmt.Println(e)
	}
	fmt.Println("GetUsers")
	fmt.Println(string(*b))

}
func TestB(t *testing.T) {
	// Test ListUsers
	c, err := NewClient(adminToken, "share.syncfs.com")
	if err != nil {
		fmt.Println("BAD")
	}

	b, _, e := c.ListUsers(100, nil, nil)
	if e != nil {
		fmt.Println("Error retrieving list of users")
		fmt.Println(e)
	}
	fmt.Println("GetUsers")
	fmt.Println(string(*b))

	// Retrieve single user
	b, _, e = c.GetUser("daniel.cardoza@aerofs.com")
	fmt.Println("GetUser")
	fmt.Println(string(*b))

	// Create a user
	b, _, e = c.CreateUser("test_email@rivend99.com", "First", "Last")
	fmt.Println("CreateUser")
	fmt.Println(e)
	fmt.Println(string(*b))

	b, _, e = c.UpdateUser("test_email@rivendell.com", "Suaron", "The Worst")
	fmt.Println("UpdateUser")
	fmt.Println(e)
	fmt.Println(string(*b))

	e = c.DeleteUser("test_email@rivend.com")
	fmt.Println("DeleteUser")
	fmt.Println(e)

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
	/*
		pp, err := c.GetFolderPath("root")
		fmt.Println("GetFolderPath")
		fmt.Println(err)
		fmt.Println(pp)
	*/
}
