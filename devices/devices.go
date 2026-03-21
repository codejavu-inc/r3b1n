package devices

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var seenDevices = make(map[string]bool)

func getDevices() []string {
	cmd := exec.Command("ip", "neigh", "show", "dev", "wlo1")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	lines := strings.Split(string(output), "\n")
	var devices []string
	for _, line := range lines {
		if strings.Contains(line, "REACHABLE") || strings.Contains(line, "STALE") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				ip := parts[0]
				devices = append(devices, ip)
			}
		}
	}
	return devices
}

func getHostname(ip string) string {
	cmd := exec.Command("getent", "hosts", ip)
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

func askUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func MonitorDevices() string {
	fmt.Println("[*] Waiting for devices to connect...")
	for {
		devices := getDevices()
		for _, ip := range devices {
			if !seenDevices[ip] {
				seenDevices[ip] = true
				hostname := getHostname(ip)
				fmt.Println("\n[+] New device connected:")
				fmt.Println("    IP:", ip)
				fmt.Println("    Hostname:", hostname)
				if askUser("[?] Do you want to track this device?") {
					fmt.Println("[+] Tracking device:", ip)
					return ip
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}
