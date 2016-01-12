package aerofs

import (
	"net/url"
	"strings"
)

func (c *Client) GetDeviceStatus(deviceID string) (*DeviceStatus, error) {
	link := url.URL{Scheme: "https",
		Host: c.Host,
		Path: strings.Join([]string{API, "devices", deviceID, "status"}, "/"),
	}

	res, err := c.get(link.String())
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	deviceStatus := DeviceStatus{}
	err = GetEntity(res, &deviceStatus)
	return &deviceStatus, err

}
