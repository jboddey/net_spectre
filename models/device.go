// spectre-go/models/device.go

package models

import "time"

// NTPData stores NTP-related information for a device
type NTPData struct {
	Versions  []uint8
	ServerIPs []string
	LastQuery time.Time
}

// Device represents a network device with additional network information
type Device struct {
	MACAddress string
	IPAddress  string
	Gateway    string
	Netmask    string
	LastSeen   time.Time
	NTPData    *NTPData
}
