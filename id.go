package lights

import (
	"encoding/hex"
	"errors"
	"log"
	"net"
	"os"
	"strings"
)

// DeviceID represents a device's ID. The class provides some extra
// helpers for determining equality of IDs across systems (many IDs may)
// be shortened to the shortest unique ID (similar to git commit IDs).
type DeviceID struct {
	ID string
}

// PanicOrID either produces an ID or panics if one can't be generated.
func PanicOrID() string {
	id, err := NewID()
	if err != nil {
		log.Println("Could not generate a device ID", err)
		os.Exit(0)
	}
	return id.ID
}

// NewID produces the device ID for an agent. The ID will be taken from
// the following in order of priority:
//
// 1. INC_DEVICE_ID environmental variable
// 2. Network MAC address
func NewID() (*DeviceID, error) {
	id := os.Getenv("INC_DEVICE_ID")
	if id != "" {
		return &DeviceID{id}, nil
	}
	// Grab hardware ID for the "best" network interface
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	second := "" // Second choice ID
	for _, intf := range interfaces {
		// The order of preference is wlan0, then en0
		switch intf.Name {
		case "wlan0":
			return &DeviceID{AddrToID(intf.HardwareAddr)}, nil
		case "en0":
			second = AddrToID(intf.HardwareAddr)
		}
	}
	if second != "" {
		return &DeviceID{second}, nil
	}
	return nil, errors.New("Missing WiFi or Ethernet interface wlan0/en0")
}

// Equals returns true if the given ID exactly equals this ID.
func (d *DeviceID) Equals(id string) bool {
	return id == d.ID
}

// Matches returns true if the given ID matchs this ID. The provided ID
// can be a shortened version of this ID.
func (d *DeviceID) Matches(id string) bool {
	return strings.HasPrefix(d.ID, id) || strings.HasPrefix(id, d.ID)
}

// ContainedIn returns true if this ID is matched to any of IDs in the provided
// string slice.
func (d *DeviceID) ContainedIn(ids []string) bool {
	for _, id := range ids {
		if d.Matches(id) {
			return true
		}
	}
	return false
}

// Returns this device ID as a string.
func (d *DeviceID) String() string {
	return d.ID
}

// AddrToID converts a network hardware address to a string device ID.
func AddrToID(addr net.HardwareAddr) string {
	return hex.EncodeToString(addr)
}
