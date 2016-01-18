package main

import (
	"fmt"
	"github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	sdk "github.com/aerofs/aerofs-sdk-golang/aerofssdk"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

// Non-persistent datastore for session information
// For persistence, use an actual DB or FileSystemStore
var store = sessions.NewCookieStore([]byte("UNIQUEID"))

// A default handler at the root of the website
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

// A login handler is required so we can get the user's email for various future
// requests
func loginEntryHandler(w http.ResponseWriter, r *http.Request) {
	signIn := "templates/signin.html"
	t, _ := template.ParseFiles(signIn)
	t.Execute(w, nil)
}

// Handler for when a user submits their email
// The user is redirected to the AeroFS Appliance, where they must grant the App
// requested permissions
func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Get new session
	session, _ := store.Get(r, "session-name")

	// Assumes a valid email was given
	r.ParseForm()
	session.Values["email"] = r.Form.Get("email")

	// Redirect User to AeroFS Appliance to retrieve Authorization Code
	ac, err := aerofsapi.NewAuthClient("appconfig.json",
		"http://localhost:13337/tokenization",
		"uniqueState", []string{"files.read", "files.write", "user.read", "user.write", "user.password"})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	aeroUrl := ac.GetAuthorizationUrl()
	logger.Printf("Sending user %s to the AeroFS Appliance at %s", session.Values["email"], aeroUrl)
	session.Save(r, w)
	http.Redirect(w, r, aeroUrl, 301)
}

func MiscHandler(w http.ResponseWriter, r *http.Request) {
	logger.Print("In Misc Handler")
	w.Write([]byte(`You are on ` + r.URL.Path + ". Random Path"))
	logger.Print("Leaving Misc Handler")
}

// Enumerate the users devices
func yourDevicesHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	token := session.Values["token"].(string)
	ac, err := aerofsapi.NewAuthClient("appconfig.json", "", "", []string{})
	a, _ := aerofsapi.NewClient(token, ac.AeroUrl)
	devices, _ := sdk.ListDevices(a, session.Values["email"].(string))
	logger.Print(*devices)
	t, err := template.ParseFiles("templates/userDevices.tmpl")
	logger.Print("Attempting to parse user devices page")
	if err != nil {
		logger.Println("Unable to retrieve template file")
		http.Error(w, err.Error(), 500)
		return
	}
	t.Execute(w, devices)
	session.Save(r, w)
}

// Enumerate the user's files
func totalUsersHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	token := session.Values["token"].(string)
	ac, err := aerofsapi.NewAuthClient("appconfig.json", "", "", []string{})
	a, _ := aerofsapi.NewClient(token, ac.AeroUrl)
	devices, _ := sdk.ListDevices(a, session.Values["email"].(string))
	logger.Print(*devices)
	t, err := template.ParseFiles("templates/userDevices.tmpl")
	logger.Print("Attempting to parse user devices page")
	if err != nil {
		logger.Println("Unable to retrieve template file")
		http.Error(w, err.Error(), 500)
		return
	}
	t.Execute(w, devices)
	session.Save(r, w)
}

// Enumerate the total number of users on the system
func yourFilesHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	token := session.Values["token"].(string)

	ac, err := aerofsapi.NewAuthClient("appconfig.json", "", "", []string{})
	a, _ := aerofsapi.NewClient(token, ac.AeroUrl)

	logger.Print("Attempting to parse user files page")
	t, err := template.ParseFiles("templates/userFiles.tmpl")
	if err != nil {
		logger.Println("Unable to retrieve template file")
		http.Error(w, err.Error(), 500)
		return
	}

	// Retrieve children of root folder
	folder, err := sdk.GetFolderClient(a, "root", []string{})
	if err != nil {
		logger.Println("Unable to retrieve file client for file.")
		http.Error(w, err.Error(), 500)
		return
	}

	folder.LoadPath()
	folder.LoadChildren()
	logger.Print(folder.Desc.ChildList.Files)
	t.Execute(w, folder.Desc)
	session.Save(r, w)
}

// Receive a Token after user accepts permissions
func tokenization(rw http.ResponseWriter, req *http.Request) {

	// Retrieve session-id so we can store corresponding token with it
	session, err := store.Get(req, "session-name")
	ac, err := aerofsapi.NewAuthClient("appconfig.json",
		"http://localhost:13337/tokenization", "uniqueState", []string{})

	// disregard state
	code := req.URL.Query().Get("code")
	token, _, err := ac.GetAccessToken(code)
	logger.Printf("USER : Email : %s | Code : %s | Token : %s",
		session.Values["email"], code, token)
	if err != nil {
		logger.Println("Unable to get correct access token")
	}

	session.Values["token"] = token
	session.Save(req, rw)
	logger.Printf("%v", session.Values)

	// Make an API Client to use with the SDK
	a, _ := aerofsapi.NewClient(token, ac.AeroUrl)

	// Printer list of users, user's devices
	users, _ := sdk.ListUsers(a, 100)
	devices, _ := sdk.ListDevices(a, session.Values["email"].(string))
	us := fmt.Sprintf("%v %v", *users, *devices)
	rw.Write([]byte(us))

	// Retrieve appconfig.json file and print it
	file, err := sdk.GetFileClient(a,
		"568e2b4ca47d340d5cb9fcb85c07f2a04e86ed3b4c0d4d43ac3a04a076025f16",
		[]string{})
	if err != nil {
		logger.Println("Unable to retrieve file client for file.")
		http.Error(rw, err.Error(), 500)
		return
	}

	content, err := file.GetContent()
	if err != nil {
		logger.Println("Unable to retrieve file content")
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.Write(*content)

	// Retrieve list of files/folders in root
	folder, err := sdk.GetFolderClient(a, "root", []string{})
	if err != nil {
		logger.Println("Unable to retrieve file client for file.")
		http.Error(rw, err.Error(), 500)
		return
	}

	folder.LoadPath()
	folder.LoadChildren()
	fol := fmt.Sprintf("%v", folder.Desc)
	f := fmt.Sprintf("%v", file.Desc)
	rw.Write([]byte(f + "\n\n\n" + fol))
}
