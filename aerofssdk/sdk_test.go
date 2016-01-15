package aerofssdk

import (
	"fmt"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	"math/rand"
	"os"
	"strings"
	"testing"
)

// Constants used for establishing appliance connections
const (
	UserToken  = "2a09580d057348d9a1382b866389b1ae"
	AdminToken = "3d2a1005a27a4115946fe308eb30785f"
	AppHost    = "share.syncfs.com"
)

// Perform teardown
func TestMain(m *testing.M) {
	err := rmUsers()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

// Remove all of test-generated users
// Note that users with <email>@aerofs.com are persisted and not removed
func rmUsers() error {
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")
	users, e := ListUsers(c, 1000)
	if e != nil {
		return e
	}

	for _, user := range *users {
		uClient := UserClient{c, user}
		if !strings.Contains(uClient.Desc.Email, "aerofs.com") {
			err := uClient.Delete()
			if err != nil {
				fmt.Printf("Unable to remove users")
				return err
			}
		}
	}
	return nil
}

// Test the creation of an existing user given an Email
func TestGetUser(t *testing.T) {
	t.Logf("Retrieving Existing user")
	c, _ := api.NewClient(UserToken, "share.syncfs.com")
	email := "daniel.cardoza@aerofs.com"

	t.Logf("Retrieving an existing user %s", email)
	u, e := GetUserClient(c, email)

	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to create new user with email %s", email)
	} else {
		t.Logf("Successfully created a new user with email %s", email)
		t.Log(*u)
	}
}

// Create a new user
func TestCreateUser(t *testing.T) {
	t.Logf("Creating new user")
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")

	t.Logf("Creating a new user")
	rand.Seed(int64(os.Getpid()))
	email := fmt.Sprintf("elrond.elf%d@middleearth.org", rand.Intn(100))
	firstName := "Melkor"
	lastName := "Bauglir"
	u, e := CreateUserClient(c, email, firstName, lastName)

	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to create new user")
	} else if u.Desc.Email == email && u.Desc.FirstName == firstName && u.Desc.LastName == lastName {
		t.Logf("Successfully created a new user")
		t.Log(*u)
	} else {
		t.Fatal("User created with incorrect fields")
		if u != nil {
			t.Fatal(*u)
		}
	}
}

// Update an already existing user
// Retrieve the user's new credentials and ensure they are updated
func TestUpdateUser(t *testing.T) {
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")

	email := fmt.Sprintf("melkor.morgoth%d@gmail.com", rand.Intn(10000))
	firstName := "Melkor"
	lastName := "Bauglir"
	u, e := CreateUserClient(c, email, firstName, lastName)
	if e != nil {
		t.Fatalf("Unable to create new user : %s", e)
	}

	// Update user {first,last}name
	t.Log(*u)
	e = u.Update("Eru", "Iluvatar")
	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to update user")
	} else {
		t.Logf("Successfully updated user")
		t.Log(*u)
	}

}

// Retrieve a list of backend users
func TestListUsers(t *testing.T) {
	c, _ := api.NewClient(AdminToken, "share.syncfs.com")
	u, e := ListUsers(c, 1000)
	if e != nil {
		t.Fatalf("Unable to retrieve a list of users : %s", e)
	}
	if u != nil {
		t.Logf("There are %d users", len(*u))
		t.Log(*u)
	}

}

func TestGetFolder(t *testing.T) {
	c, _ := api.NewClient(UserToken, "share.syncfs.com")
	f, e := GetFolderClient(c, "root", []string{"path", "children"})
	if e != nil {
		t.Fatalf("Unable to retrieve a FolderClient : %s", e)
	}

	f.LoadChildren()
	f.LoadMetadata()
	t.Log(*f)
}

func TestGetFile(t *testing.T) {
	c, _ := api.NewClient(UserToken, "share.syncfs.com")
	fileId := "568e2b4ca47d340d5cb9fcb85c07f2a04e86ed3b4c0d4d43ac3a04a076025f16"
	f, e := GetFileClient(c, fileId, []string{"path", "children"})
	if e != nil {
		t.Fatalf("Unable to retrieve a FileCLient : %s", e)
	}
	t.Log(*f)
	f.LoadPath()
	t.Log(*f)
}

func TestMoveFile(t *testing.T) {

}
