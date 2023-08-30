package kea

import (
	"net/http"
)

type (
	// RemoteSubnet4 : Represents a single subnet4 entry in Kea.
	RemoteSubnet4 struct {
		FourO6Interface   string                 `json:"4o6-interface"`
		FourO6InterfaceID string                 `json:"4o6-interface-id"`
		FourO6Subnet      string                 `json:"4o6-subnet"`
		ID                int                    `json:"id"`
		Metadata          Metadata               `json:"metadata"`
		OptionData        []OptionData           `json:"option-data"`
		Pools             []Pool                 `json:"pools"`
		Relay             Relay                  `json:"relay"`
		SharedNetworkName interface{}            `json:"shared-network-name"`
		Subnet            string                 `json:"subnet"`
		UserContext       map[string]interface{} `json:"user-context"`
	}

	// RemoteSubnet4List : Represents a single subnet4 entry in Kea.
	RemoteSubnet4List struct {
		ID                int      `json:"id"`
		Metadata          Metadata `json:"metadata"`
		SharedNetworkName string   `json:"shared-network-name"`
		Subnet            string   `json:"subnet"`
	}

	// NewRemoteSubnet4 : Represents a single subnet4 entry in Kea.
	NewRemoteSubnet4 struct {
		ID                int               `json:"id"`
		Subnet            string            `json:"subnet"`
		SharedNetworkName *string           `json:"shared-network-name"`
		Pools             []Pool            `json:"pools"`
		OptionData        []OptionData      `json:"option-data"`
		Relay             Relay             `json:"relay,omitempty"`
		UserContext       map[string]string `json:"user-context,omitempty"`
	}

	// Relay : Represents a single relay entry in Kea.
	Relay struct {
		IPAddresses []string `json:"ip-addresses,omitempty"`
	}

	// Pool : Represents a single pool entry in Kea.
	Pool struct {
		Pool string `json:"pool"`
	}

	// OptionData : Represents a single option-data entry in Kea.
	OptionData struct {
		Code       *int    `json:"code,omitempty"`
		Data       string  `json:"data"`
		Name       string  `json:"name"`
		Space      *string `json:"space,omitempty"`
		AlwaysSend bool    `json:"always-send"`
	}
)

// RemoteSubnet4List : Gets a list of subnets from the Kea configuration-backend commands API.
//
// POST / {"command":"remote-subnet4-list","service":["dhcp4"],"arguments":{"remote":{"type":"postgresql"},"server-tags": ["all"]}}'
func (c *Client) RemoteSubnet4List(hostname string) ([]RemoteSubnet4List, error) {
	payload := Request{
		Command: "remote-subnet4-list",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":      map[string]string{"type": c.remote},
			"server-tags": []string{"all"},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Subnets []RemoteSubnet4List `json:"subnets"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Subnets, nil
}

// RemoteSubnet4GetByPrefix : Gets a list of subnets from the Kea configuration-backend commands API.
//
// POST / {"command":"remote-subnet4-list","service":["dhcp4"],"arguments":{"remote":{"type":"postgresql"},"server-tags": ["all"]}}'
func (c *Client) RemoteSubnet4GetByPrefix(hostname, prefix string) (RemoteSubnet4, error) {
	payload := Request{
		Command: "remote-subnet4-get-by-prefix",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]string{{"subnet": prefix}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return RemoteSubnet4{}, err
	}

	var ret struct {
		Subnets []RemoteSubnet4 `json:"subnets"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return RemoteSubnet4{}, err
	}
	return ret.Subnets[0], nil
}

// RemoteSubnet4GetByID : Gets a list of subnets from the Kea configuration-backend commands API.
//
// POST / {"command":"remote-subnet4-list","service":["dhcp4"],"arguments":{"remote":{"type":"postgresql"},"server-tags": ["all"]}}'
func (c *Client) RemoteSubnet4GetByID(hostname string, id int) (RemoteSubnet4, error) {
	payload := Request{
		Command: "remote-subnet4-get-by-id",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]int{{"id": id}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return RemoteSubnet4{}, err
	}

	var ret struct {
		Subnets []RemoteSubnet4 `json:"subnets"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return RemoteSubnet4{}, err
	}
	return ret.Subnets[0], nil
}

// RemoteSubnet4DelByPrefix : Deletes a subnet from the Kea configuration-backend commands API.
func (c *Client) RemoteSubnet4DelByPrefix(hostname, prefix string) (int, error) {
	payload := Request{
		Command: "remote-subnet4-del-by-prefix",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]string{{"subnet": prefix}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return 0, err
	}

	var ret struct {
		Count int `json:"count"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return 0, err
	}
	return ret.Count, nil
}

// RemoteSubnet4DelByID : Deletes a subnet from the Kea configuration-backend commands API.
func (c *Client) RemoteSubnet4DelByID(hostname string, id int) (int, error) {
	payload := Request{
		Command: "remote-subnet4-del-by-id",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":  map[string]string{"type": c.remote},
			"subnets": []map[string]int{{"id": id}},
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return 0, err
	}

	var ret struct {
		Count int `json:"count"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return 0, err
	}
	return ret.Count, nil
}

// RemoteSubnet4Set : Creates a new subnet using the Kea configuration-backend commands API.
func (c *Client) RemoteSubnet4Set(hostname string, subnets []NewRemoteSubnet4) ([]RemoteSubnet4List, error) {
	payload := Request{
		Command: "remote-subnet4-set",
		Service: []string{"dhcp4"},
		Arguments: map[string]any{
			"remote":      map[string]string{"type": c.remote},
			"server-tags": []string{"all"},
			"subnets":     subnets,
		},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Subnets []RemoteSubnet4List `json:"subnets"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Subnets, nil
}
