package tracker

import (
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// TrackDevice captures all traffic but filters output based on protocol
func TrackDevice(iface string, targetIP string, protocol string, outputFile string) error {
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	fmt.Printf("[*] Capturing all traffic for %s on %s\n", targetIP, iface)
	fmt.Println("[*] Output filter:", protocol)
	fmt.Println("[*] Saving output to", outputFile)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}
		ip, _ := ipLayer.(*layers.IPv4)

		// Only process packets for the selected device
		if ip.SrcIP.String() != targetIP && ip.DstIP.String() != targetIP {
			continue
		}

		line := ""

		// DNS
		if (protocol == "dns" || protocol == "all") && packet.Layer(layers.LayerTypeDNS) != nil {
			dnsLayer := packet.Layer(layers.LayerTypeDNS)
			dns, _ := dnsLayer.(*layers.DNS)
			if dns.QR == false { // DNS request
				for _, q := range dns.Questions {
					line = fmt.Sprintf("[DNS] %s -> %s (%s)\n", ip.SrcIP, string(q.Name), q.Type)
					file.WriteString(line)
					fmt.Print(line)
				}
			}
		}

		// TCP
		if (protocol == "tcp" || protocol == "all") && packet.Layer(layers.LayerTypeTCP) != nil {
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			tcp, _ := tcpLayer.(*layers.TCP)
			line = fmt.Sprintf("[TCP] %s:%d -> %s:%d\n", ip.SrcIP, tcp.SrcPort, ip.DstIP, tcp.DstPort)
			file.WriteString(line)
			fmt.Print(line)
		}

		// UDP
		if (protocol == "udp" || protocol == "all") && packet.Layer(layers.LayerTypeUDP) != nil {
			udpLayer := packet.Layer(layers.LayerTypeUDP)
			udp, _ := udpLayer.(*layers.UDP)
			line = fmt.Sprintf("[UDP] %s:%d -> %s:%d\n", ip.SrcIP, udp.SrcPort, ip.DstIP, udp.DstPort)
			file.WriteString(line)
			fmt.Print(line)
		}

		// Optionally: ICMP, ARP, etc. can be added here for -p all
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
