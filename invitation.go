package aerofs

import (
	"net/url"
	"strings"
)

// Return a list of Invitations to shared folders for a given user
func (c *Client) ListSFInvitations(email string) (*[]Invitation, error) {
	route := strings.Join([]string{"users", email, "invitations"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	if err != nil {
		return nil, err
	}

	invitations := []Invitation{}
	err = GetEntity(res, &invitations)
	return &invitations, err
}

// View invitation metadata for a pending invitation
func (c *Client) ViewPendingSFInvitation(email, sid string) (*Invitation, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	invitation := Invitation{}
	err = GetEntity(res, &invitation)
	return &invitation, err
}

// Accept a pending Shared Folder invitation
func (c *Client) AcceptSFInvitation(email, sid string) (*SharedFolder, error) {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	query := url.Values{"external": []string{"0"}}
	link := c.getURL(route, query.Encode())

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	sf := SharedFolder{}
	err = GetEntity(res, &sf)
	return &sf, err
}

// Ignore an existing invitation to a shared folder
func (c *Client) IgnoreSFInvitation(email, sid string) error {
	route := strings.Join([]string{"users", email, "invitations", sid}, "/")
	link := c.getURL(route, "")

	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
