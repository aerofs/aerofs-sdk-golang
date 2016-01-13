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
	"net/http"
	"strings"
)

// Retrieve information about a person invited to an AeroFS instance
func (c *Client) GetInvitee(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"invitees", email}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Create an invitation
func (c *Client) CreateInvitee(email_to, email_from string) (*[]byte,
	*http.Header, error) {
	route := "invitees"
	link := c.getURL(route, "")
	invitee := map[string]string{
		"email_to":   email_to,
		"email_from": email_from,
	}

	data, err := json.Marshal(invitee)
	if err != nil {
		return nil, nil, errors.New("Unable to serialize invitation request")
	}

	res, err := c.post(link, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header, err := unpackageResponse(res)
	return body, header, err
}

// Delete an unsatisfied invitation
func (c *Client) DeleteInvitee(email string) error {
	route := strings.Join([]string{"invitees", email}, "/")
	link := c.getURL(route, "")
	res, err := c.del(link)
	defer res.Body.Close()
	return err
}
