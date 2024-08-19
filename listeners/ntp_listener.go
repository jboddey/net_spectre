// spectre-go/listeners/ntp_listener.go

package listeners

import (
	"net"
	"spectre-go/models"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
)

// NTPListener processes NTP packets to update NTP-related information
type NTPListener struct {
	Log       *logrus.Logger
	DeviceSet map[string]*models.Device
}

// OnPacket processes each NTP packet to extract NTP version and server IP
func (l *NTPListener) OnPacket(packet gopacket.Packet) {
	// Check for UDP layer
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udpPacket, _ := udpLayer.(*layers.UDP)

		// Check for NTP (port 123)
		if udpPacket.DstPort == 123 || udpPacket.SrcPort == 123 {
			// Extract NTP data
			ntpLayer := packet.Layer(layers.LayerTypeNTP)
			if ntpLayer != nil {
				ntpPacket, _ := ntpLayer.(*layers.NTP)

				// Check that this is NTP client type
				if ntpPacket.Mode != 3 {
					return
				}

				// Get the MAC address from the Ethernet layer
				ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
				if ethernetLayer == nil {
					return
				}
				ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
				macAddr := ethernetPacket.SrcMAC.String()

				// Update or create the device entry
				device, found := l.DeviceSet[macAddr]
				if found {
					// Update last seen time
					device.LastSeen = time.Now()
				} else {
					// Ignore newly detected devices
					return
				}

				// Initialize NTPData if nil
				if device.NTPData == nil {
					device.NTPData = &models.NTPData{}
				}

				// Get server IP addresses
				// Check for IP layer (IPv4 or IPv6)
				ipLayer := packet.Layer(layers.LayerTypeIPv4)
				if ipLayer == nil {
					ipLayer = packet.Layer(layers.LayerTypeIPv6)
				}

				var dstIP net.IP

				if ipLayer != nil {

					switch ip := ipLayer.(type) {
					case *layers.IPv4:
						dstIP = ip.DstIP
					case *layers.IPv6:
						dstIP = ip.DstIP
					}

				} else {
					return
				}

				// Append version and server IP if not already present
				if !containsUint8(device.NTPData.Versions, uint8(ntpPacket.Version)) {
					device.NTPData.Versions = append(device.NTPData.Versions, uint8(ntpPacket.Version))
				}
				if !containsString(device.NTPData.ServerIPs, dstIP.String()) {
					device.NTPData.ServerIPs = append(device.NTPData.ServerIPs, dstIP.String())
				}
				device.NTPData.LastQuery = time.Now()

				// Log the updated NTP information
				l.Log.Infof("NTP Data for %s: Versions %v, Server IPs %v", macAddr, device.NTPData.Versions, device.NTPData.ServerIPs)
			}
		}
	}
}

// Helper function to check if a slice contains a specific uint8 value
func containsUint8(slice []uint8, value uint8) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Helper function to check if a slice contains a specific string value
func containsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
