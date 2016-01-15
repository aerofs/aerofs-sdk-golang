package aerofsapi

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
//  - Refactor into a rest API and then SDK
//    - for the API, simple return a buffer of the body, the Header map and an
//      error
import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	API       = "api/v1.3"
	CHUNKSIZE = 5000
)

// A Client is used to communicate with an AeroFS Appliance and return the
// resultant responses
type Client struct {
	// The hostname/IP of the AeroFS Appliance
	// Used when constructing the default API Prefix for all subsequent API calls
	// Ie. share.syncfs.com
	Host string

	// The OAuth token
	Token string

	// Default header containing Token, Content-type and Endpoint-Consistency
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

// For a given HTTP-Response, this returns the associated body,header
func unpackageResponse(res *http.Response) (*[]byte, *http.Header, error) {
	body, _ := ioutil.ReadAll(res.Body)
	header := res.Header

	// For each API call, unpackage the HTTP response and return an error if a non
	// 2XX status code is retrieved
	if res.StatusCode >= 300 {
		err := errors.New(res.Status)
		return &body, &header, err
	}
	return &body, &header, nil
}

// Construct a URL given a route and query parameters
func (c *Client) getURL(route, query string) string {
	link := url.URL{Scheme: "https",
		Path: strings.Join([]string{API, route}, "/"),
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

// HTTP-GET
func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP GET Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

// HTTP-POST
func (c *Client) post(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP POST request")
	}

	request.Header = c.Header
	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	hClient := &http.Client{}
	return hClient.Do(request)
}

// HTTP-PUT
func (c *Client) put(url string, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP PUT request")
	}

	request.Header = c.Header

	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	hClient := &http.Client{}
	return hClient.Do(request)
}

// HTTP-DELETE
func (c *Client) del(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, errors.New("Unable to create HTTP DELETE Request")
	}

	request.Header = c.Header
	hClient := &http.Client{}
	return hClient.Do(request)
}

// Generic Handler for HTTP request
// Allows the passing of additional HTTP request header K/V pairs
func (c *Client) request(req, url string, options *http.Header, buffer io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(req, url, buffer)
	if err != nil {
		return nil, errors.New("Unable to create HTTP " + req + " Request")
	}

	// If header map passed in , add additional KV pairs
	request.Header = c.Header
	if options != nil && len(*options) > 0 {
		for k, v := range *options {
			for _, el := range v {
				request.Header.Add(k, el)
			}
		}
	}

	// If we are not sending data, delete the default content-Type
	if buffer == nil {
		request.Header.Del("Content-Type")
	}

	// TODO : Add extra field to signal serializing
	// Note : Determine if this has actual effect
	contentType := options.Get("Content-Type")
	if options.Get("Content-Type") != "" {
		request.Header.Set("Content-Type", contentType)
	}

	hClient := &http.Client{}
	return hClient.Do(request)
}

// Unmarshalls data from an HTTP Response into a given entity
func GetEntity(res *http.Response, entity interface{}) error {
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &entity)
}
