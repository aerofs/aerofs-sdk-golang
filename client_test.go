package aerofs

import (
	"fmt"
	"testing"
)

// Non-admin token : 2a09580d057348d9a1382b866389b1ae
func TestB(t *testing.T) {
	// Test ListUsers
	userToken := "2a09580d057348d9a1382b866389b1ae"
	adminToken := "3d2a1005a27a4115946fe308eb30785f"
	c, err := NewClient(adminToken, "share.syncfs.com")

	if err != nil {
		fmt.Println("BAD")
	}
	//before := 100
	//after := 10
	//	a, e := c.ListUsers(1000, &after, &before)
	a, e := c.ListUsers(1000, nil, nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("GetUsers")
	fmt.Println(a)

	// Retrieve single user
	u, e := c.GetUser("daniel.cardoza@aerofs.com")
	fmt.Println("GetUser")
	fmt.Println(*u)

	// Create a user
	user := User{"test.email@yahoo.com", "Test_firstname", "Test_lastname",
		[]SharedFolder{}, []Invitation{}}
	h, e := c.CreateUser(user)
	fmt.Println(e)
	fmt.Println("CreateUser")
	fmt.Println(*h)

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
}
