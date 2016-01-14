package aerofs

import (
	"encoding/json"
	"errors"
)

// Device, client wrapper
type DeviceClient struct {
	APIClient *Client
	Desc      Device
}

// Device descriptor
type Device struct {
	Id          string `json:"id"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	OSFamily    string `json:"os_family"`
	InstallDate string `json:"install_data"`
}

type DeviceStatus struct {
	Online   bool   `json:"online"`
	LastSeen string `json:"last_seen"`
}

// Return an existing device client given a deviceId
func NewDeviceClient(c *Client, deviceId string) (*DeviceClient, error) {
	body, _, err := c.GetDeviceMetadata(deviceId)
	if err != nil {
		return nil, err
	}
	device := Device{}
	err = json.Unmarshal(*body, &device)
	return &DeviceClient{c, device}, err
}

// Update the name of the device
func (c *DeviceClient) Update(name string) error {
	body, _, err := c.APIClient.UpdateDevice(name)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body, &c.Desc)
	if err != nil {
		return errors.New("Unable to demarshal updated device metadata")
	}

	return nil
}

// Retrieve the status of the current device
func (c *DeviceClient) Status() (*DeviceStatus, error) {
	body, _, err := c.APIClient.GetDeviceStatus(c.Desc.Id)
	if err != nil {
		return nil, err
	}

	deviceStatus := new(DeviceStatus)
	err = json.Unmarshal(*body, deviceStatus)
	if err != nil {
		return nil, errors.New("Unable to demarshal current device status")
	}

	return deviceStatus, err
}
