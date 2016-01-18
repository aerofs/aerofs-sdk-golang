package aerofssdk

import (
	"encoding/json"
	"errors"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

type GroupMemberClient struct {
	APIClient *api.Client
	Desc      GroupMember
}

type GroupMember struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	GroupId   string
}

func ListGroupMembers(c *api.Client, groupId string) (*[]GroupMember, error) {
	body, _, err := c.ListGroupMembers(groupId)
	if err != nil {
		return nil, err
	}

	groupMembers := []GroupMember{}
	err = json.Unmarshal(*body, &groupMembers)
	if err != nil {
		return nil, errors.New("Unable to unmarshal list of group members")
	}

	return &groupMembers, nil
}

func GetGroupMember(c *api.Client, groupId, memberEmail string) (*GroupMemberClient, error) {
	body, _, err := c.GetGroupMember(groupId, memberEmail)
	if err != nil {
		return nil, err
	}

	g := GroupMemberClient{APIClient: c, Desc: GroupMember{GroupId: groupId}}
	err = json.Unmarshal(*body, &g.Desc)
	if err != nil {
		return nil, errors.New("Unable to unmarshal group member")
	}

	return &g, nil
}

// Update the groupMember information
func (g *GroupMemberClient) Load() error {
	body, _, err := g.APIClient.GetGroupMember(g.Desc.GroupId, g.Desc.Email)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &g.Desc)
	if err != nil {
		return errors.New("Unable to unmarshal group member")
	}

	return nil
}
