package kea

import (
	"net/http"
)

type (
	// Lease4 : Represents a single lease entry in Kea.
	Lease4 struct {
		ClientID  string `json:"client-id,omitempty"`
		Cltt      int    `json:"cltt"`
		FqdnFwd   bool   `json:"fqdn-fwd"`
		FqdnRev   bool   `json:"fqdn-rev"`
		Hostname  string `json:"hostname"`
		HwAddress string `json:"hw-address"`
		IPAddress string `json:"ip-address"`
		State     int    `json:"state"`
		SubnetID  int    `json:"subnet-id"`
		ValidLft  int    `json:"valid-lft"`
	}
)

// GetLease4All : Gets a list of leases from the Kea API.
//
// POST / {"command": "lease4-get-all","arguments":{"subnets":[2]},"service":["dhcp4"]}'
func (c *Client) GetLease4All(hostname string, subnetIDs []int) ([]Lease4, error) {
	payload := Request{Command: "lease4-get-all", Service: []string{"dhcp4"}}
	if subnetIDs != nil {
		payload.Arguments = map[string]any{"subnets": subnetIDs}
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Leases []Lease4 `json:"leases"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Leases, nil
}

// GetLease4ByIP : Gets a lease by IP Address from the Kea API.
//
// POST / {"command": "lease4-get","arguments":{"ip-address": "192.0.2.1"},"service":["dhcp4"]}'
func (c *Client) GetLease4ByIP(hostname string, ip string) (Lease4, error) {
	var ret Lease4
	payload := Request{
		Command:   "lease4-get",
		Arguments: map[string]any{"ip-address": ip},
		Service:   []string{"dhcp4"},
	}
	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return ret, err
	}
	if _, err := c.do(req, &ret); err != nil {
		return ret, err
	}
	return ret, nil
}

// GetLease4ByHost : Gets a lease by IP Address from the Kea API.
//
// POST / {"command": "lease4-get-by-hostnam","arguments":{"hostname": "192.0.2.1"},"service":["dhcp4"]}'
func (c *Client) GetLease4ByHost(hostname string, host string) ([]Lease4, error) {
	payload := Request{
		Command:   "lease4-get-by-hostname",
		Arguments: map[string]any{"hostname": host},
		Service:   []string{"dhcp4"},
	}
	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}
	var ret struct {
		Leases []Lease4 `json:"leases"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Leases, nil
}

// GetLease4ByMac : Gets a lease by IP Address from the Kea API.
//
// POST / {"command": "lease4-get-by-hw-address","arguments":{"hw-address": "192.0.2.1"},"service":["dhcp4"]}'
func (c *Client) GetLease4ByMac(hostname string, mac string) ([]Lease4, error) {
	payload := Request{
		Command:   "lease4-get-by-hw-address",
		Arguments: map[string]any{"hw-address": mac},
		Service:   []string{"dhcp4"},
	}
	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}
	var ret struct {
		Leases []Lease4 `json:"leases"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Leases, nil
}

// DelLease4 : Deletes a lease by IP Address from the Kea API.
//
// POST / {"command": "lease4-del","arguments":{"ip-address": "192.0.2.1"},"service":["dhcp4"]}'
func (c *Client) DelLease4(hostname string, ip string) (string, error) {
	payload := Request{
		Command:   "lease4-del",
		Arguments: map[string]any{"ip-address": ip},
		Service:   []string{"dhcp4"},
	}
	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return "", err
	}
	base, err := c.do(req, nil)
	if err != nil {
		return "", err
	}
	return base.Text, nil
}
