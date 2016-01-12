package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
)

// Retrieve information about a person invited to an AeroFS instance
func (c *Client) GetInvitee(email string) (*Invitee, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees", email}, "/"),
	}

	res, err := c.get(url.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	invitee := new(Invitee)
	err = GetEntity(res, invitee)
	return invitee, err
}

// Create an invitation
func (c *Client) CreateInvitee(email_to, email_from string) (*Invitee, error) {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees"}, "/"),
	}

	newInvitee := new(Invitee)
	invitee := new(Invitee)
	invitee.EmailTo = email_to
	invitee.EmailFrom = email_from

	data, err := json.Marshal(invitee)
	if err != nil {
		return newInvitee, errors.New("Unable to serialize invitation request")
	}

	res, err := c.post(url.String(), bytes.NewBuffer(data))
	if err != nil {
		return newInvitee, err
	}

	err = GetEntity(res, newInvitee)
	return newInvitee, err
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	url := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "invitees", email}, "/"),
	}

	res, err := c.del(url.String())
	defer res.Body.Close()
	return err
}
