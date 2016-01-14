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
	HasMore bool       `json:"has_more"`
	Users   []UserDesc `json:"data"`
}

// Wraps the data members of a User and the API Client required to mutate
type User struct {
	APIClient *Client `json:"-"`
	Desc      UserDesc
}

// The fields of a User
type UserDesc struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

// Return an existing user
func GetUser(client *Client, email string) (*User, error) {
	body, _, err := client.GetUser(email)
	if err != nil {
		return nil, err
	}

	u := User{APIClient: client}
	err = json.Unmarshal(*body, &u.Desc)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Create a new user and return
func CreateUser(client *Client, email, firstName, lastName string) (*User, error) {
	// CreateUser returns a Location of the new resource
	body, _, err := client.CreateUser(email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	u := User{APIClient: client}
	err = json.Unmarshal(*body, &u.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Update a users first, last Name
func (u *User) Update(newFirstName, newLastName string) error {
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
func (u *User) changePassword(password string) error {
	err := u.APIClient.ChangePassword(u.Desc.Email, password)
	return err
}

// Disable two-factor authentication
func (u *User) DisableTwoFactorAuth() error {
	err := u.APIClient.DisableTwoFactorAuth(u.Desc.Email)
	return err
}
