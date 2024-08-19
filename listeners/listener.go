// spectre-go/listeners/listener.go

package listeners

import "github.com/google/gopacket"

// PacketListener is the interface that all listeners must implement
type PacketListener interface {
	OnPacket(packet gopacket.Packet)
}
