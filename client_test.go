package aerofs

import (
	"fmt"
	"testing"
)

// Non-admin token : 2a09580d057348d9a1382b866389b1ae
func TestB(t *testing.T) {
	// Test ListUsers
	//userToken := "2a09580d057348d9a1382b866389b1ae"
	adminToken := "3d2a1005a27a4115946fe308eb30785f"
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
	/*
		// Create an invitaiton
		invite, e := c.CreateInvitee("danielpcardoza@gmail.com",
			"daniel.cardoza@aerofs.com")
		fmt.Println(e)
		fmt.Println(invite)

		// Get invitation list
		in, e := c.GetInvitee("daniel.cardoza@aerofs.com")
		fmt.Println("GetInvitee")
		fmt.Println(e)
		fmt.Println(in)

		c.SetToken(userToken)
		// Get root folder data
		f, err := c.GetFolderMetadata("root", []string{"children"})
		fmt.Println("GetRootFolder")
		fmt.Println(err)
		fmt.Println(f)

		pp, err := c.GetFolderPath("root")
		fmt.Println("GetFolderPath")
		fmt.Println(err)
		fmt.Println(pp)
	*/
}
