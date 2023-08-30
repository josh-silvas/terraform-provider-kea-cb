package kea

import (
	"net/http"
)

type (
	// Heartbeat : Represents the HA heartbeat information.
	Heartbeat struct {
		DateTime          string   `json:"date-time"`
		Scopes            []string `json:"scopes"`
		State             string   `json:"state"`
		UnsentUpdateCount int      `json:"unsent-update-count"`
	}
)

// HAHeartbeat : Gets HA status of the dhcp4 cluster..
//
// POST / {"command": "ha-heartbeat","service": ["dhcp4"]}'
func (c *Client) HAHeartbeat(hostname string) (Heartbeat, error) {
	var res Heartbeat
	payload := Request{Command: "ha-heartbeat", Service: []string{"dhcp4"}}
	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return res, err
	}

	if _, err := c.do(req, &res); err != nil {
		return res, err
	}
	return res, nil
}
