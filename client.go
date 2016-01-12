package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	API = "api/v1.3"
)

type Client struct {
	// The hostname/IP of the AeroFS Appliance
	// Used when constructing the default API Prefix for all subsequent API calls
	// Ie. share.syncfs.com
	Host string

	// The OAuth token
	Token string

	// Contains the authorization token
	// For conditional file, and folder requests, the header is populated
	// with an ETag
	Header http.Header
}

// SDK-Client Constructor
// Constructs the HTTP header used for subsequent requests
// OAuth token stored in HTTP header
func NewClient(token, host string) (*Client, error) {
	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	header.Set("Content-Type", "application/json")
	header.Set("Endpoint-Consistency", "strict")

	c := Client{Host: host,
		Header: header,
		Token:  token}

	return &c, nil
}

// Construct a URL given well-defined parameters
// The Scheme should be constant but the client is able to reset the Host and
// Token
func (c *Client) getURL(path, query string) string {
	link := url.URL{Scheme: "https",
		Path: path,
		Host: c.Host,
	}

	if query != "" {
		link.RawQuery = query
	}

	return link.String()
}

// Resets the token for a given client
// Allows the third-party developer to construct 1 SDK-Client used to retrieve
// the values for multiple users
func (c *Client) SetToken(token string) {
	c.Header.Set("Authorization", "Bearer "+token)
}

// Wrappers for basic HTTP functions
// Use HTTPClient since the stdlib does not provide function prototypes for all
// request types
func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP GET Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

func (c *Client) post(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP POST request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

func (c *Client) put(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP PUT request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}
func (c *Client) del(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP DELETE Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}
