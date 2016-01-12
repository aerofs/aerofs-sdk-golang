package aerofs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) ListPendingMembers(sid string, etags []string) (*[]SFPendingMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{"shares", sid, "pending"}, "/"),
	}

	if len(etags) > 0 {
		vals := url.Values{"If-None-Match": etags}
		link.RawQuery = vals.Encode()
	}

	res, err := c.get(link.String())
	if err != nil {
		return nil, err
	}

	sfpmList := []SFPendingMember{}
	err = GetEntity(res, &sfpmList)
	return &sfpmList, err
}

func (c *Client) GetPendingMember(id, email string) (*SFPendingMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{"shares", id, "pending", email}, "/"),
	}

	res, err := c.get(link.String())
	if err != nil {
		return nil, err
	}

	sfpm := SFPendingMember{}
	err = GetEntity(res, &sfpm)
	return &sfpm, err
}

func (c *Client) InviteToSharedFolder(sid, email string, permissions []string,
	note string) (*SFPendingMember, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{"shares", sid, "pending"}, "/"),
	}
	sfpm := SFPendingMember{Email: email, Permissions: permissions, Note: note}

	data, err := json.Marshal(sfpm)
	if err != nil {
		return nil, errors.New(fmt.Sprint(`Unable to parse given parameters : %s %s
%v`, email, note, permissions))
	}

	res, err := c.post(link.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	sfpm = SFPendingMember{}
	err = GetEntity(res, &sfpm)
	return &sfpm, err
}

func (c *Client) RemovePendingMember(sid, email string) error {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{"shares", sid, "pending", email}, "/"),
	}

	_, err := c.del(link.String())
	return err
}
