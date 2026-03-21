package hotspot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var CreatedByUs bool = false

func IsHotspotActive() (bool, string) {
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,TYPE,DEVICE", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "wifi") {
			return true, line
		}
	}
	return false, ""
}

func AskUser(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question + " (y/n): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func CreateHotspot(ssid string, password string) error {
	fmt.Println("[+] Creating hotspot...")
	cmd := exec.Command(
		"nmcli", "device", "wifi", "hotspot",
		"ifname", "wlo1",
		"ssid", ssid,
		"password", password,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		CreatedByUs = true
	}
	return err
}

func DeleteHotspot() {
	if !CreatedByUs {
		fmt.Println("[*] Hotspot was not created by r3b1n, skipping delete")
		return
	}
	fmt.Println("\n[!] Cleaning up hotspot...")
	cmd := exec.Command("nmcli", "connection", "delete", "Hotspot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("[+] Hotspot removed")
}

func HandleHotspot() error {
	active, info := IsHotspotActive()
	if active {
		fmt.Println("[+] Hotspot already active:")
		fmt.Println("    ", info)
		use := AskUser("[?] Do you want to use the current hotspot?")
		if use {
			fmt.Println("[+] Using existing hotspot")
			return nil
		} else {
			fmt.Println("[!] Please disable hotspot manually and rerun")
			os.Exit(0)
		}
	}

	fmt.Println("[!] No active hotspot found")
	create := AskUser("[?] Do you want to create a hotspot?")
	if !create {
		fmt.Println("[!] Exiting...")
		os.Exit(0)
	}

	ssid := "r3b1n"
	password := "11111111"
	err := CreateHotspot(ssid, password)
	if err != nil {
		return err
	}
	fmt.Println("[+] Hotspot created successfully")
	return nil
}
