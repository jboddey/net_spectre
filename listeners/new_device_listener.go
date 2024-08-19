// spectre-go/listeners/new_device_listener.go

package listeners

import (
	"spectre-go/models"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// NewDeviceListener logs MAC addresses of new devices found on the network
type NewDeviceListener struct {
	Log       *logrus.Logger
	DeviceSet map[string]*models.Device
}

// OnPacket processes each packet to check for new devices
func (listener *NewDeviceListener) OnPacket(packet gopacket.Packet) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		macAddr := ethernetPacket.SrcMAC.String()

		device := listener.DeviceSet[macAddr]
		seenBefore := (device != nil)

		// Check if this MAC address has been seen before
		if !seenBefore {
			// New device found, create a Device object and log the MAC address
			device = &models.Device{
				MACAddress: macAddr,
			}
			listener.Log.Infof("New device found: %s", macAddr)
		}
		if device != nil {

			// Update the last seen value
			device.LastSeen = time.Now()

			// Update the device in the device list
			listener.DeviceSet[macAddr] = device
		}
	}
}
