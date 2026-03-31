# r3b1n

**r3b1n** is a lightweight Go-based hotspot tracker and network sniffer for Linux.  
It allows you to create a Wi-Fi hotspot, monitor connected devices, and capture all network traffic, including DNS, TCP, and UDP, per device.  

---

## ⚡ Features

- Automatically **creates a hotspot** if none exists.
- Detects and **lists devices** connected to the hotspot in real-time.
- Allows selecting a **device to track**.
- Captures **all network traffic** from the selected device:
  - DNS requests
  - TCP connections
  - UDP connections
- `-p` flag filters **output** by protocol (`dns`, `tcp`, `udp`, or `all`).
- Traffic is **saved to a user-specified file**.
- Cleanly deletes the hotspot when the script is stopped (Ctrl+C).
- Modular structure:
  - `hotspot/` → hotspot management
  - `devices/` → device detection
  - `tracker/` → packet capture

---

## ⚙️ Requirements

- Debian-based Linux
- Go >= 1.20
- Root privileges (required to capture packets and manage hotspot)
- `nmcli` for hotspot management
- `libpcap` installed (`sudo apt install libpcap-dev`)
- `getent` for hostname resolution (standard on most Linux)

---

## 💻 Installation

Clone the repository:

```bash
git clone https://github.com/codejavu-inc/r3b1n.git
cd r3b1n
```
## Usage
```
sudo ./r3b1n -i wlo1 -p dns -o alltraffic.txt
```
The hotspot password id `11111111`
