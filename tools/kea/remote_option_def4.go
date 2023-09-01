package kea

import (
	"net/http"
	"strings"
)

type (
	// RemoteOptionDef4 : Represents a single remote option definition entry in Kea.
	RemoteOptionDef4 struct {
		Name        string `json:"name,omitempty"`
		Code        int    `json:"code"`
		Type        string `json:"type,omitempty"`
		Array       bool   `json:"array,omitempty"`
		RecordTypes string `json:"record-types,omitempty"`
		Space       string `json:"space,omitempty"`
		Encapsulate string `json:"encapsulate,omitempty"`
	}
)

// RemoteOptionDef4Set : Sets the remote option definition for the dhcp4 configuration.
func (c *Client) RemoteOptionDef4Set(hostname string, def RemoteOptionDef4) error {
	payload := Request{
		Command: "remote-option-def4-set",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":      map[string]string{"type": c.remote},
			"server-tags": []string{"all"},
			"option-defs": []RemoteOptionDef4{def},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return err
	}

	var ret interface{}
	if _, err := c.do(req, &ret); err != nil {
		return err
	}
	return nil
}

// RemoteOptionDef4Get : Gets the remote option definition from the dhcp4 configuration.
func (c *Client) RemoteOptionDef4Get(hostname, space string, code int) (*RemoteOptionDef4, error) {
	payload := Request{
		Command: "remote-option-def4-get",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":      map[string]string{"type": c.remote},
			"server-tags": []string{"all"},
			"option-defs": []RemoteOptionDef4{{Space: space, Code: code}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}
	type optDefResp struct {
		Count      int                `json:"count"`
		OptionDefs []RemoteOptionDef4 `json:"option-defs"`
	}

	ret := new(optDefResp)
	if _, err := c.do(req, ret); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	if len(ret.OptionDefs) == 0 {
		return nil, nil
	}
	return &ret.OptionDefs[0], nil
}

// RemoteOptionDef4Del : Deletes the remote option definition from the dhcp4 configuration.
func (c *Client) RemoteOptionDef4Del(hostname, space string, code int) error {
	payload := Request{
		Command: "remote-option-def4-del",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":      map[string]string{"type": c.remote},
			"server-tags": []string{"all"},
			"option-defs": []RemoteOptionDef4{{Space: space, Code: code}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return err
	}

	var ret interface{}
	if _, err := c.do(req, &ret); err != nil {
		return err
	}
	return nil
}
