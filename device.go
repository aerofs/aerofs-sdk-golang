package aerofs

// This is the entrypoint class for making connections with an AeroFS Appliance
// A received OAuth Token is required for authentication
// TODO :
//  - reformat the Path construction per each URL object to remove extraneous
//  code
import (
	"net/url"
	"strings"
)

// Device specific API Calls
func (c *Client) ListDevices(email string) (*[]Device, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "users", email, "devices"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	devices := []Device{}
	err = GetEntity(res, &devices)
	return &devices, err

}

func (c *Client) GetDeviceMetadata(deviceID string) (*Device, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "devices", deviceID}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	device := Device{}
	err = GetEntity(res, &device)
	return &device, err
}
