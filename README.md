# Tiger-Eye 🐅👁️

**Tiger-Eye** is a fast, lightweight, cross-platform network reconnaissance and service fingerprinting CLI tool written in Go (Golang). Designed for authorized security auditing and network discovery, it scans target hosts for open TCP ports, grabs service banners, inspects HTTP server headers and page titles, and probes RTSP media streams.

## Features
* **High-Speed Concurrent Scanning:** Uses Go routines and `sync.WaitGroup` for parallel port probing with configurable connection timeouts.
* **Service & Device Fingerprinting:** Automatically extracts HTTP `Server` headers, HTML `<title>` tags, and RTSP (IP camera/NVR) protocol responses.
* **Flexible Port Targeting:** Scan built-in common/IoT ports by default or pass custom comma-separated ports and ranges (`-ports 80,443,8000-8010`).
* **Cross-Platform Support:** Compiles cleanly into standalone executables for both **Linux** and **Windows** with zero external dependencies.

## Quick Start

### Run from Source
Requires Go 1.20 or higher:
```bash
git clone https://github.com
cd tiger-eye
go run main.go -target scanme.nmap.org
```

### Custom Ports & Ranges
```bash
go run main.go -target 192.168.1.1 -ports 22,80,443,554,8000-8010
```

## Build Standalone Executables
Generate single-file binaries that run without Go installed:

* **Linux:**
  ```bash
  GOOS=linux GOARCH=amd64 go build -o tiger-eye-linux main.go
  ./tiger-eye-linux -target 127.0.0.1
  ```
* **Windows:**
  ```bash
  GOOS=windows GOARCH=amd64 go build -o tiger-eye-windows.exe main.go
  tiger-eye-windows.exe -target 127.0.0.1
  ```

## Sample Output
```text
Tiger-Eye v0.2 - Fingerprinting scanme.nmap.org (14 ports)...
[OPEN] scanme.nmap.org:22    -> SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
[OPEN] scanme.nmap.org:80    -> Server: Apache/2.4.7 (Ubuntu) | Title: "Go ahead and ScanMe!"
Scan complete.
```

## Legal Disclaimer
This tool is intended strictly for educational purposes, authorized security assessments, and auditing networks you own or have explicit written permission to test. The author assumes no liability for misuse or unauthorized scanning.
