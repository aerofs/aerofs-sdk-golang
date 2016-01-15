package main

import (
	"fmt"
	// context.ClearHandler supposedly needed to prevent memory leak with a
	// non-Gorilla Mux
	"github.com/aerofs/aerofs-sdk-golang/aerofsapi"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("UNIQUEID"))

// Note that an HTTP Handler is a (HTTP.ResponseWriter, *http.Request)
func main() {
	fmt.Println("Main")
	r := httprouter.New()

	r.GET("/test", test_1)
	r.GET("/", arrive)
	r.GET("/redirect", redirect)
	r.GET("/tokenization", tokenization)
	/*	r.GET("/files", func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		http.FileServer(http.Dir("/")).ServeHTTP(rw, req)
	})*/

	http.ListenAndServe("localhost:13337", context.ClearHandler(r))
}

// Handler called when a user arrives
func arrive(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
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
	rw.Write([]byte(fmt.Sprintf(`A session has been generated for you with ID :
%s`, session.Name())))

}

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
	rw.Write([]byte(str))
}

func test_1(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	rw.Write([]byte(`This is a test response`))
}
