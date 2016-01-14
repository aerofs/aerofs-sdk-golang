package aerofs

import (
	"encoding/json"
	"errors"
	"fmt"
)

// The User object is used to easily modify backend users assuming
// the object has a reference to a given APIClient. Methods are able to modify
// internal user state as well as backend state such as user password
// Each object has a corresponding Descriptor struct containing its members

// The response structure returned from a ListUser(..) call
type userListResponse struct {
	HasMore bool   `json:"has_more"`
	Users   []User `json:"data"`
}

// Wraps the data members of a User and the API Client required to mutate
type UserClient struct {
	APIClient *Client `json:"-"`
	Desc      User
}

type User struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

func (u User) String() string {
	return fmt.Sprintf("\nEmail : %s\n FN : %s\n LN : %s\n", u.Email, u.FirstName, u.LastName)
}

// Return an existing user
func GetUser(client *Client, email string) (*UserClient, error) {
	body, _, err := client.GetUser(email)
	if err != nil {
		return nil, err
	}

	u := UserClient{APIClient: client}
	err = json.Unmarshal(*body, &u.Desc)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Get a list of Users
func ListUsers(client *Client, limit int) (*[]User, error) {
	body, _, err := client.ListUsers(limit, nil, nil)
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = json.Unmarshal(*body, &users)
	if err != nil {
		return nil, errors.New("Unable to unmarshal a retrieved list of users")
	}
	return &users, nil
}

// Create a new user and return
func CreateUser(client *Client, email, firstName, lastName string) (*User, error) {
	// CreateUser returns a Location of the new resource
	body, _, err := client.CreateUser(email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	u := User{}
	err = json.Unmarshal(*body, &u)
	if err != nil {
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Update a users first, last Name
func (u *UserClient) Update(newFirstName, newLastName string) error {
	body, _, err := u.APIClient.UpdateUser(u.Desc.Email, newFirstName, newLastName)
	if err != nil {
		return err
	}

	// TODO : When unmarshalling a failure occurs, is it possible for the body we
	// unmarshal to be mutated?
	err = json.Unmarshal(*body, &u.Desc)
	if err != nil {
		return errors.New("Unable to update user")
	}

	return nil
}

// Change the user's password
func (u *UserClient) changePassword(password string) error {
	err := u.APIClient.ChangePassword(u.Desc.Email, password)
	return err
}

// Disable two-factor authentication
func (u *UserClient) DisableTwoFactorAuth() error {
	err := u.APIClient.DisableTwoFactorAuth(u.Desc.Email)
	return err
}
