package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func extractTitle(html string) string {
	re := regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		title = strings.ReplaceAll(title, "\n", " ")
		title = strings.ReplaceAll(title, "\r", "")
		return title
	}
	return ""
}

func probeService(conn net.Conn, port int) string {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Check RTSP (common for IP cameras/NVRs on port 554 or 8554)
	if port == 554 || port == 8554 {
		fmt.Fprintf(conn, "OPTIONS * RTSP/1.0\r\nCSeq: 1\r\n\r\n")
		reader := bufio.NewReader(conn)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			return fmt.Sprintf("RTSP Service (%s)", line)
		}
		return "RTSP Port Open (No banner)"
	}

	// Send HTTP GET request to inspect Server headers and HTML <title>
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: localhost\r\nUser-Agent: Tiger-Eye/0.2\r\n\r\n")
	buf := make([]byte, 4096)
	n, _ := io.ReadFull(conn, buf)
	if n == 0 {
		// Fallback: read whatever initial banner arrived (e.g., SSH, FTP, SMTP)
		n, _ = conn.Read(buf)
	}
	if n == 0 {
		return "(No text banner returned)"
	}

	response := string(buf[:n])
	lines := strings.Split(response, "\n")
	firstLine := strings.TrimSpace(lines[0])

	var serverHeader string
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "server:") {
			serverHeader = strings.TrimSpace(line[7:])
			break
		}
	}

	title := extractTitle(response)

	var details []string
	if serverHeader != "" {
		details = append(details, fmt.Sprintf("Server: %s", serverHeader))
	}
	if title != "" {
		details = append(details, fmt.Sprintf("Title: \"%s\"", title))
	}
	if len(details) > 0 {
		return strings.Join(details, " | ")
	}
	return firstLine
}

func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()

	fingerprint := probeService(conn, port)
	fmt.Printf("[OPEN] %-21s -> %s\n", address, fingerprint)
}

func parsePorts(portStr string) []int {
	var ports []int
	parts := strings.Split(portStr, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if strings.Contains(p, "-") {
			bounds := strings.Split(p, "-")
			if len(bounds) == 2 {
				start, err1 := strconv.Atoi(bounds[0])
				end, err2 := strconv.Atoi(bounds[1])
				if err1 == nil && err2 == nil && start <= end {
					for i := start; i <= end; i++ {
						ports = append(ports, i)
					}
				}
			}
		} else {
			if port, err := strconv.Atoi(p); err == nil {
				ports = append(ports, port)
			}
		}
	}
	return ports
}

func main() {
	target := flag.String("target", "127.0.0.1", "Target IP address or hostname to scan")
	portList := flag.String("ports", "21,22,23,25,80,443,445,554,3000,3001,3389,8000,8080,8443,8554,37777", "Comma-separated ports or ranges (e.g., 80,443,8000-8010)")
	flag.Parse()

	ports := parsePorts(*portList)
	fmt.Printf("Tiger-Eye v0.2 - Fingerprinting %s (%d ports)...\n", *target, len(ports))

	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go scanPort(*target, port, 2*time.Second, &wg)
	}

	wg.Wait()
	fmt.Println("Scan complete.")
}
