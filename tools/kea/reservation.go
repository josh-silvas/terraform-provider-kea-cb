package kea

import (
	"net"
	"net/http"
)

type (
	// Reservation : Represents a single reservation entry in Kea.
	Reservation struct {
		BootFileName string         `json:"boot-file-name,omitempty"`
		ClientID     string         `json:"client-id,omitempty"`
		CircuitID    string         `json:"circuit-id,omitempty"`
		DuID         string         `json:"duid,omitempty"`
		FlexID       string         `json:"flex-id,omitempty"`
		IPAddress    string         `json:"ip-address"`
		HwAddress    string         `json:"hw-address"`
		Hostname     string         `json:"hostname"`
		NextServer   string         `json:"next-server,omitempty"`
		OptionData   []OptionData   `json:"option-data,omitempty"`
		SubnetID     int            `json:"subnet-id"`
		UserContext  map[string]any `json:"user-context,omitempty"`
	}
)

// ReservationGetAll : Gets the remote option for the subnet4 list.
func (c *Client) ReservationGetAll(hostname string, subnetID int) ([]Reservation, error) {
	payload := Request{
		Command:   "reservation-get-all",
		Service:   []string{"dhcp4"},
		Arguments: map[string]any{"subnet-id": subnetID},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return nil, err
	}

	var ret struct {
		Hosts []Reservation `json:"hosts"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return nil, err
	}
	return ret.Hosts, nil
}

// ReservationAdd : Adds a reservation to the subnet4 list.
func (c *Client) ReservationAdd(hostname string, res Reservation) error {
	if net.ParseIP(res.IPAddress) == nil {
		return ErrInvalidIP
	}
	if val, err := net.ParseMAC(res.HwAddress); err == nil {
		res.HwAddress = val.String()
	} else {
		return ErrInvalidMAC
	}

	payload := Request{
		Command:   "reservation-add",
		Service:   []string{"dhcp4"},
		Arguments: map[string]any{"reservation": res},
	}

	req, err := c.make(http.MethodPost, hostname, payload, nil)
	if err != nil {
		return err
	}

	var ret struct {
		Options []OptionReq `json:"options"`
	}
	if _, err := c.do(req, &ret); err != nil {
		return err
	}
	return nil
}

// ReservationDel : Deletes a reservation to the subnet4 list.
func (c *Client) ReservationDel(hostname, ipAddress string, subnetID int) error {
	if net.ParseIP(ipAddress) == nil {
		return ErrInvalidIP
	}

	payload := Request{
		Command:   "reservation-del",
		Service:   []string{"dhcp4"},
		Arguments: map[string]any{"subnet-id": subnetID, "ip-address": ipAddress},
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
