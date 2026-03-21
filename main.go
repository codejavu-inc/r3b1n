package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"r3b1n/devices"
	"r3b1n/hotspot"
	"r3b1n/tracker"
	"syscall"
)

func main() {
	fmt.Println("=== r3b1n ===")

	err := hotspot.HandleHotspot()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("[+] Hotspot ready ✔")

	// Monitor devices and select one
	targetIP := devices.MonitorDevices()
	fmt.Println("[+] Selected target:", targetIP)

	// CLI flags
	iface := flag.String("i", "wlo1", "Interface to track")
	proto := flag.String("p", "all", "Protocol to track (dns, tcp, udp, all)")
	out := flag.String("o", "output.txt", "Output file")
	flag.Parse()

	// Ctrl+C cleanup
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\n[!] Interrupt received")
		hotspot.DeleteHotspot()
		fmt.Println("[+] Exiting r3b1n")
		os.Exit(0)
	}()

	// Start tracker
	err = tracker.TrackDevice(*iface, targetIP, *proto, *out)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
