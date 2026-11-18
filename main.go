package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func scanPort(ip string, port int, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	fmt.Fprintf(conn, "HEAD / HTTP/1.0\r\n\r\n")

	reader := bufio.NewReader(conn)
	banner, _ := reader.ReadString('\n')
	banner = strings.TrimSpace(banner)

	if banner != "" {
		fmt.Printf("[OPEN] %s | Banner: %s\n", address, banner)
	} else {
		fmt.Printf("[OPEN] %s | (No text banner returned)\n", address)
	}
}

func main() {
	target := flag.String("target", "127.0.0.1", "Target IP address to scan")
	flag.Parse()

	ports := []int{21, 22, 23, 25, 80, 443, 445, 3389, 8080, 8443}

	fmt.Printf("Tiger-Eye v0.1 - Scanning %s...\n", *target)
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go scanPort(*target, port, 2*time.Second, &wg)
	}

	wg.Wait()
	fmt.Println("Scan complete.")
}
