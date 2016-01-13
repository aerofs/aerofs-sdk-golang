package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Device specific API Calls
func (c *Client) ListDevices(email string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "devices"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

func (c *Client) GetDeviceMetadata(deviceId string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceId}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}

func (c *Client) UpdateDeviceMetadata(deviceName string) (*[]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceName}, "/")
	link := c.getURL(route, "")
	newDevice := map[string]string{
		"name": deviceName,
	}

	data, err := json.Marshal(newDevice)
	if err != nil {
		return nil, nil, errors.New("Unable to marshal new device")
	}

	res, err := c.put(link, bytes.NewBuffer(data))
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	body, header := unpackageResponse(res)
	return body, header, err
}
