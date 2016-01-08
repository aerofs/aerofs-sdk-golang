package aerofs

import (
	"fmt"
	"testing"
)

func TestBasic(t *testing.T) {
	config := AppConfig{"00d63fe6-5baf-4908-807e-d072885828f4",
		"25cf762e-023b-47e6-af91-b76506da5725",
		"http://pulsar", []string{}}
	url := "share.syncfs.com"
	authClient := AuthClient{config, url}
	code := "00fa2d4b3ee94befbb74c3a6b5404d8b"
	token, err := authClient.GetAccessToken(code)
	fmt.Println(token, err)
}
