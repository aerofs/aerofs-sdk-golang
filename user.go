package aerofs

import (
	"encoding/json"
	"errors"
)

type User struct {
	APIClient   *Client        `json:"-"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

// Return an existing user
func newUser(client *Client, email string) (*User, error) {
	body, _, err := client.GetUser(email)
	if err != nil {
		return nil, err
	}

	u := User{APIClient: client}
	err = json.Unmarshal(*body, &u)
	if err != nil {
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Create a new user and return
func createUser(client *Client, email, firstName, lastName string) (*User, error) {
	// CreateUser returns a Location of the new resource
	body, _, err := client.CreateUser(email, firstName, lastName)
	if err != nil {
		return nil, err
	}

	u := User{APIClient: client}
	err = json.Unmarshal(*body, &u)
	if err != nil {
		return nil, errors.New("Unable to unmarshal new User")
	}

	return &u, nil
}

// Update a users first, last Name
func (u *User) update(newFirstName, newLastName string) error {
	body, _, err := u.APIClient.UpdateUser(u.Email, newFirstName, newLastName)
	if err != nil {
		return err
	}

	// TODO : When unmarshalling a failure occurs, is it possible for the body we
	// unmarshal to has changes inside?
	err = json.Unmarshal(*body, u)
	if err != nil {
		return errors.New("Unable to update user")
	}

	return nil
}

func (u *User) changePassword(password string) error {
	err := u.APIClient.ChangePassword(u.Email, password)
	if err != nil {
		return err
	}
	return nil
}
