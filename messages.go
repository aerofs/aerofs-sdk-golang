package aerofs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Structures used when communicating with an AeroFS Appliance

type AppConfig struct {
	Id       string
	Secret   string
	Redirect string
	Scopes   []string
}

type NewName struct {
	Grant       string
	Authcode    string
	ID          string
	Secret      string
	RedirectURL string
}

type Access struct {
	Token      string `json:"access_token"`
	TokenType  string `json:"token_type"`
	ExpireTime int    `json:"expires_in"`
	Scope      string `json:"scope"`
}

type SharedFolder struct {
	Id         string            `json:"id"`
	Name       string            `json:"email"`
	External   bool              `json:"is_external"`
	Members    []SFMember        `json:"members"`
	Groups     []SFGroupMember   `json:"groups"`
	Pending    []SFPendingMember `json:"pending"`
	Permission string            `json:"caller_effective_permissions"`
}

type SFMember struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Permissions string `json:"permissions"`
}

type SFGroupMember struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type SFPendingMember struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Inviter     string `json:"invited_by"`
	Permissions string `json:"permissions"`
	Note        string `json:"note"`
}

type Group struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Members []GroupMember `json:"members"`
}

type GroupMember struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Shares      []SharedFolder `json:"shares"`
	Invitations []Invitation   `json:"invitations"`
}

type Invitee struct {
	EmailTo    string `json:"email_to"`
	EmailFrom  string `json:"email_from"`
	SignupCode string `json:"signup_code,omitempty"`
}

type Invitation struct {
	Id          string   `json:"shared_id"`
	Name        string   `json:"shared_name"`
	Inviter     string   `json:"invited_by"`
	Permissions []string `json:"permissions"`
}

// Response specific structures

type ListUserResponse struct {
	HasMore bool   `json:"has_more"`
	Users   []User `json:"data"`
}

// Unmarshalls data from an HTTP response given a response struct
func GetEntity(res *http.Response, entity interface{}) error {
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Unable to parse HTTP Response")
	}

	return json.Unmarshal(data, entity)
}
