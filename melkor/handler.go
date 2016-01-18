package main

import (
	//"github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	//sdk "github.com/aerofs/aerofs-sdk-golang/aerofssdk"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

// Non-persistent datastore for session information
// For persistence, use an actual DB or FileSystemStore
var store = sessions.NewCookieStore([]byte("UNIQUEID"))

// Redirect the user to either signin or the homepage depending on if
// a session exists for the user
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	redirect := *r.URL
	redirect.Path = "login"

	// If the session cookie is present, go to home
	for _, r := range r.Cookies() {
		if r.Name == "session-name" {
			redirect.Path = "home"
			break
		}
	}
	http.Redirect(w, r, redirect.String(), 301)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

}

func MiscHandler(w http.ResponseWriter, r *http.Request) {
	logger.Print("In Misc Handler")
	w.Write([]byte(`You are on ` + r.URL.Path + ". Random Path"))
	logger.Print("Leaving Misc Handler")
}

// Handler called when a user arrives
func arrive(rw http.ResponseWriter, req *http.Request) {

	fmt.Println("In arrive")
	session, err := store.Get(req, "ASODASDLASDL")
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	fmt.Println(session)
	// Set some session values.
	session.Values["foo"] = "bar"
	session.Values[42] = 43
	fmt.Println(session)
	// Save it before we write to the response/return from the handler.
	session.Save(req, rw)
	rw.Write([]byte(fmt.Sprintf(`A session has been generated for you with ID : %s`, session.Name())))

}

/*
// Have a user redirected to AeroFS site to give permissions and so we can get a
// token back
func redirect(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	fmt.Println("In redirect writer")
	ac, err := aerofsapi.NewAuthClient("appconfig.json",
		"http://localhost:13337/tokenization", "uniqueState", []string{"files.read",
			"files.write", "user.read", "user.write", "user.password"})

	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}

	aeroUrl := ac.GetAuthorizationUrl()
	fmt.Println("URL is", aeroUrl)
	http.Redirect(rw, req, aeroUrl, 301)
}

// Receive a Token after used accepts permissions
func tokenization(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	for a, e := range req.URL.Query() {
		fmt.Println("%s : %s", a, e)
	}
	str := fmt.Sprintf("%v", req.URL.Query())
	fmt.Println(str)
	ac, err := aerofsapi.NewAuthClient("appconfig.json",
		"http://localhost:13337/tokenization", "uniqueState", []string{"files.read",
			"files.write", "user.read", "user.write", "user.password"})
	code := req.URL.Query()["code"][0]
	fmt.Println(code)
	token, _, err := ac.GetAccessToken(code)
	if err != nil {
		fmt.Println("Unable to get correct access token")
	}
	fmt.Println("Token is", token)
	a, _ := aerofsapi.NewClient(token, ac.AeroUrl)
	users, _ := sdk.ListUsers(a, 100)
	devices, _ := sdk.ListDevices(a, "daniel.cardoza@aerofs.com")
	us := fmt.Sprintf("%v %v", *users, *devices)
	rw.Write([]byte(us))
}
*/
