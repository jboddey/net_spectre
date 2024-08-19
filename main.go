// spectre-go/main.go

package main

import (
	"os"
	"os/signal"
	"syscall"

	"spectre-go/listeners"
	"spectre-go/models"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set up logrus logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set the network interface to listen on
	iface := "wlp0s20f3" // Replace with your interface name
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening device %s: %v", iface, err)
	}
	defer handle.Close()

	log.Infof("Listening on %s", iface)

	// Initialize the device set as map of pointers
	deviceSet := make(map[string]*models.Device)

	// Initialize listeners
	newDeviceListener := &listeners.NewDeviceListener{
		Log:       log,
		DeviceSet: deviceSet,
	}

	ipAddressListener := &listeners.IPAddressListener{
		Log:       log,
		DeviceSet: deviceSet,
	}

	dhcpListener := &listeners.DHCPListener{
		Log:       log,
		DeviceSet: deviceSet,
	}

	ntpListener := &listeners.NTPListener{
		Log:       log,
		DeviceSet: deviceSet,
	}

	// List of registered listeners
	listenerList := []listeners.PacketListener{
		newDeviceListener,
		ipAddressListener,
		dhcpListener,
		ntpListener,
	}

	// Set up channel for packet capture
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Handle SIGINT and SIGTERM to gracefully shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Info("Shutting down...")
		handle.Close()
		os.Exit(0)
	}()

	// Capture and distribute packets to listeners
	for packet := range packetSource.Packets() {
		for _, listener := range listenerList {
			listener.OnPacket(packet)
		}
	}
}
