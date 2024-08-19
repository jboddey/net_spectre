// spectre-go/listeners/dhcp_listener.go

package listeners

import (
	"net"
	"spectre-go/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// DHCPListener processes DHCP packets to update network information
type DHCPListener struct {
	Log       *logrus.Logger
	DeviceSet map[string]*models.Device
}

// OnPacket processes each DHCP packet to extract network information
func (l *DHCPListener) OnPacket(packet gopacket.Packet) {

	// Check for DHCPv4 layer
	dhcpLayer := packet.Layer(layers.LayerTypeDHCPv4)
	if dhcpLayer != nil {

		// Cast dhcpLayer to the DHCPv4 layer type
		dhcpPacket, _ := dhcpLayer.(*layers.DHCPv4)

		// Process DHCP packets
		for _, option := range dhcpPacket.Options {

			// Process each message type differently
			switch option.Type {
			case layers.DHCPOptMessageType:

				// Message Type option to determine if it's a DHCP Offer or ACK
				if len(option.Data) > 0 && option.Data[0] == byte(layers.DHCPMsgTypeAck) {
					// DHCP ACK indicates that the packet is a response with network info
					macAddr := dhcpPacket.ClientHWAddr.String()
					ipAddr := dhcpPacket.YourClientIP.String()

					// Update or create the device entry
					device, found := l.DeviceSet[macAddr]
					if found {
						device.IPAddress = ipAddr
						l.Log.Infof("Updated IP address for %s is %s", macAddr, ipAddr)
					}
				}
			case layers.DHCPOptRouter:
				// Router option provides the gateway
				macAddr := dhcpPacket.ClientHWAddr.String()
				device, found := l.DeviceSet[macAddr]
				if found {
					if device.Gateway == "" { // Check if Gateway is not set
						device.Gateway = net.IP(option.Data).String()
						l.Log.Infof("Detected gateway for %s: %s", device.MACAddress, device.Gateway)
					}
				}
			case layers.DHCPOptSubnetMask:
				// Subnet Mask option
				macAddr := dhcpPacket.ClientHWAddr.String()
				device, found := l.DeviceSet[macAddr]
				if found {
					if device.Netmask == "" { // Check if Netmask is not set
						device.Netmask = net.IP(option.Data).String()
						l.Log.Infof("Detected netmask for %s: %s", device.MACAddress, device.Netmask)
					}
				}
			}
		}
	}
}
