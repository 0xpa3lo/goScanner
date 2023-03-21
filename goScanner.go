package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

func generatePortsRange(n int) []int {
	ports := make([]int, n)
	for i := 1; i <= n; i++ {
		ports[i-1] = i
	}
	return ports
}

func main() {
	var host string

	fmt.Print("Enter the IP address to scan: ")
	fmt.Scanln(&host)

	fmt.Printf("\nScanning %s\n", host)

	portsToScan := generatePortsRange(1000)
	results := scanPorts(host, portsToScan)

	if len(results) == 0 {
		fmt.Println("No open ports found.")
	} else {
		fmt.Println("Open ports found:", results)
	}
}

func scanPorts(host string, ports []int) []int {
	openPorts := make([]int, 0)
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(port int) {
			defer wg.Done()

			serverInfo, isOpen := isPortOpenAndGrabServerInfo(host, port)
			if isOpen {
				openPorts = append(openPorts, port)
				if serverInfo != "" {
					fmt.Printf("Port is open %d: %s\n", port, serverInfo)
				}
			}
		}(port)
	}

	wg.Wait()

	sort.Ints(openPorts)
	return openPorts
}

func isPortOpenAndGrabServerInfo(host string, port int) (string, bool) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return "", false
	}

	defer conn.Close()

	serverInfo := ""

	if port == 80 {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\nHost: %s\r\n\r\n", host)

		reader := bufio.NewReader(conn)
		resp, err := http.ReadResponse(reader, nil)

		if err == nil {
			serverInfo = resp.Header.Get("Server")
		}

	} else {
		// For other ports, just read the banner
		reader := bufio.NewReader(conn)
		banner, err := reader.ReadString('\n')
		if err == nil {
			serverInfo = strings.TrimSpace(banner)
		}
	}

	return serverInfo, true
}
