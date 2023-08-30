package kea

import (
	"net/http"
)

type (
	// OptionReq : Represents a single option-data entry in Kea.
	OptionReq struct {
		Code int    `json:"code"`
		Data string `json:"data,omitempty"`
	}
)

// RemoteOption4Set : Sets the remote option for the subnet4 list.
func (c *Client) RemoteOption4Set(hostname string, subnetID int, opts []OptionReq) ([]OptionReq, error) {
	payload := Request{
		Command: "remote-option4-subnet-set",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]int{{"id": subnetID}},
			"options": opts,
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Options []OptionReq `json:"options"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Options, nil
}

// RemoteOption4Del : Deletes the remote option for the subnet4 list.
func (c *Client) RemoteOption4Del(hostname string, subnetID int, opts []OptionReq) ([]OptionReq, error) {
	payload := Request{
		Command: "remote-option4-subnet-del",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]int{{"id": subnetID}},
			"options": opts,
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Options []OptionReq `json:"options"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Options, nil
}
