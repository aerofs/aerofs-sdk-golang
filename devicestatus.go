package aerofs

import (
	"net/http"
	"strings"
)

func (c *Client) GetDeviceStatus(deviceId string) (*[]byte, *http.Header,
	error) {
	route := strings.Join([]string{"devices", deviceId, "status"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	body, header, err := unpackageResponse(res)
	return body, header, err
}
