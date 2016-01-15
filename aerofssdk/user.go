package aerofs

import (
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
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

// User, Client wrapper
type UserClient struct {
	APIClient *api.Client `json:"-"`
	Desc      User
}

// User descriptor
type User struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

// Custom User print
func (u User) String() string {
	return fmt.Sprintf("\nEmail : %s\n FN : %s\n LN : %s\n", u.Email, u.FirstName, u.LastName)
}

// Return an existing userClient given an email and client reference
func GetUserClient(client *api.Client, email string) (*UserClient, error) {
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
// Note that these users are not tied to the given client
func ListUsers(client *api.Client, limit int) (*[]User, error) {
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

// Create a new user Client and return
func CreateUserClient(client *api.Client, email, firstName, lastName string) (*UserClient, error) {
	body, _, err := client.CreateUser(email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	u := User{}
	err = json.Unmarshal(*body, &u)
	if err != nil {
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &UserClient{client, u}, nil
}

// Update a users first, last Name
func (u *UserClient) Update(newFirstName, newLastName string) error {
	body, _, err := u.APIClient.UpdateUser(u.Desc.Email, newFirstName, newLastName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &u.Desc)
	if err != nil {
		return errors.New("Unable to update user")
	}

	return nil
}

// Change the user's password
func (u *UserClient) changePassword(password string) error {
	return u.APIClient.ChangePassword(u.Desc.Email, password)
}

// Disable two-factor authentication
func (u *UserClient) DisableTwoFactorAuth() error {
	return u.APIClient.DisableTwoFactorAuth(u.Desc.Email)
}

// Delete the current user from the backend
func (u *UserClient) Delete() error {
	return u.APIClient.DeleteUser(u.Desc.Email)
}
