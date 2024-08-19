// spectre-go/listeners/ip_address_listener.go

package listeners

import (
	"net"
	"spectre-go/models"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// IPAddressListener updates the IP address of devices in the device set
type IPAddressListener struct {
	Log       *logrus.Logger
	DeviceSet map[string]*models.Device
}

// OnPacket processes each packet to update the IP address of devices
func (l *IPAddressListener) OnPacket(packet gopacket.Packet) {

	// Check for ARP layer to get IP addresses
	arpLayer := packet.Layer(layers.LayerTypeARP)
	if arpLayer != nil {

		// Cast arpPacket to ARP layer type
		arpPacket, _ := arpLayer.(*layers.ARP)

		// We are only interested in ARP requests or replies
		if arpPacket.Operation == layers.ARPRequest || arpPacket.Operation == layers.ARPReply {

			// Convert MAC address and IP address to strings
			macAddr := net.HardwareAddr(arpPacket.SourceHwAddress).String()
			ipAddr := net.IP(arpPacket.SourceProtAddress).String()

			// Ignore 0.0.0.0
			if ipAddr == "0.0.0.0" {
				return
			}

			// Update or add the device's IP address
			device, found := l.DeviceSet[macAddr]

			if found {
				// Update IP address if it has changed
				if device.IPAddress != ipAddr {
					device.IPAddress = ipAddr

					// Update the last seen value
					device.LastSeen = time.Now()

					// Update the map with the modified device
					l.DeviceSet[macAddr] = device
					l.Log.Infof("Updated IP address for %s is %s", macAddr, ipAddr)
				}
			}
		}
	}
}
