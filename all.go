package main

import (
	"fmt"
	api "github.com/aerofs/aerofs-sdk-golang/aerofsapi"
)

func main() {
	fmt.Println("TEST")
	b := api.AuthClient{}
	fmt.Println(b)
}
