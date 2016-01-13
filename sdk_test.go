package aerofs

import (
	"fmt"
	"math/rand"
	"testing"
)

// Test the creation of an existing user given an Email
func TestGetUser(t *testing.T) {
	t.Logf("Retrieving Existing user")
	c, _ := NewClient(userToken, "share.syncfs.com")
	email := "daniel.cardoza@aerofs.com"

	t.Logf("Creating a new user %s", email)
	u, e := GetUser(c, email)

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
	c, _ := NewClient(adminToken, "share.syncfs.com")

	t.Logf("Creating a new user")
	email := fmt.Sprintf("elrond.rivendell%d@middleearth.org", rand.Intn(10000))
	u, e := CreateUser(c, email, "Elrond", "Rivendell")

	if e != nil {
		t.Log(e)
		t.Fatalf("Unable to create new user")
	} else {
		t.Logf("Successfully created a new user")
		t.Log(*u)
	}
}
