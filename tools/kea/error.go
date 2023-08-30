package kea

import "errors"

var (
	// ErrInvalidIP : Invalid IP address
	ErrInvalidIP = errors.New("invalid IP address")
	// ErrInvalidMAC : Invalid MAC address
	ErrInvalidMAC = errors.New("invalid MAC address")
)
