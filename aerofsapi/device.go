package aerofsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// This file maps all routes exposed on the AeroFS API

// Device specific API Calls

func (c *Client) ListDevices(email string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"users", email, "devices"}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) GetDeviceMetadata(deviceId string) ([]byte, *http.Header, error) {
	route := strings.Join([]string{"devices", deviceId}, "/")
	link := c.getURL(route, "")

	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}

func (c *Client) UpdateDevice(deviceName string) ([]byte, *http.Header, error) {
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

	return unpackageResponse(res)
}

func (c *Client) GetDeviceStatus(deviceId string) ([]byte, *http.Header,
	error) {
	route := strings.Join([]string{"devices", deviceId, "status"}, "/")
	link := c.getURL(route, "")
	res, err := c.get(link)
	defer res.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	return unpackageResponse(res)
}
