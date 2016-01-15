package aerofsapi

import (
	"reflect"
	"testing"
)

// Test that an autogenerated application configure file unmarshals into a
// correct Authorization Client
func TestAuth_ParseAppConfig(t *testing.T) {
	fileName := "test_appconfig.json"
	secret := "b706384f-e65f-4058-9c13-e052044c408b"
	id := "f3c121ed-1beb-4b52-9145-70ca55d82af2"
	url := "share.syncfs.com"

	correctClient := AuthClient{Config: AuthConfig{Id: id, Secret: secret, AeroUrl: url}}
	authClient, err := NewAuthClient(fileName, "", "", []string{})
	if err != nil {
		t.Fatal("Unable to create new AuthClient from file : %s", err)
	}

	if !reflect.DeepEqual(correctClient, *authClient) {
		t.Fatalf("The two clients are different %v : %v", correctClient, *authClient)
	}

	t.Log("The clients are the same")
	t.Log(*authClient)
}

/*
func TestURL(t *testing.T) {
	config := AppConfig{"764d972d-5717-4a98-b9de-aa41d13da7a2",
		"b03927cb-8d07-422b-98a3-a9b6483185e8",
		"http://blackhole",
		[]string{"files.read", "files.write", "acl.read", "acl.write",
			"acl.invitations", "user.read", "user.write"}}
	authClient := AuthClient{config, "share.syncfs.com"}
	url := authClient.GetAuthCode()
	fmt.Println(url)
	exe := exec.Command("open", url)
	exe.Run()
}

func TestBasic(t *testing.T) {
	config := AppConfig{"764d972d-5717-4a98-b9de-aa41d13da7a2",
		"b03927cb-8d07-422b-98a3-a9b6483185e8",
		"http://blackhole",
		[]string{"files.read", "files.write"}}
	authClient := AuthClient{config, "share.syncfs.com"}
	code := "e731374c7999450ca236859a0968b310"
	token, scopes, err := authClient.GetAccessToken(code)
	fmt.Println(token, scopes, err)
}*/
