package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var (
	Port           = flag.Int("port", 8090, "port for listening")
	BufferSize int = 1024 * 1024
	Host           = flag.String("host", "127.0.0.1:8090", "host:port")
	Mode           = flag.String("mode", "help", "select mode [client, server, help]")
)

func main() {
	flag.Parse()

	switch *Mode {
	case "client":
		client(*Host)
	case "server":
		server(*Port)
	default:
		help()
	}
}

func client(host string) {
	fmt.Println("Client mode", host)
	buffer := make([]byte, BufferSize)
	startTime := time.Now()
	conn, err := net.Dial("tcp", *Host)
	if err != nil {
		log.Printf("Dialing error: %v\n", err)
		return
	}
	defer conn.Close()
	var totalBytes int64
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Read error: %v\n", err)
			break
		}
		totalBytes += int64(n)
	}
	duration := time.Since(startTime)
	speed := float64(totalBytes) / duration.Seconds()
	fmt.Printf("Received %d bytes in %v. Speed: %.2f MB/s\n", totalBytes, duration, speed/(1024*1024))
}

func getBuffer() []byte {
	buffer := make([]byte, BufferSize)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = byte(i % 4)
	}
	return buffer
}

func handleConnect(conn net.Conn) {
	defer conn.Close()
	log.Printf("New clinet: %s\n", conn.RemoteAddr())
	buffer := getBuffer()
	startTime := time.Now()
	var totalBytes int64
	for i := 0; i < 100; i++ {
		n, err := conn.Write(buffer)
		if err != nil {
			log.Println("Handle connection error:", err)
			return
		}
		totalBytes += int64(n)
	}
	duration := time.Since(startTime)
	speed := float64(totalBytes) / duration.Seconds()
	fmt.Printf("Sent %d bytes in %v. Speed: %.2f MB/s\n", totalBytes, duration, speed/(1024*1024))

}

func server(port int) {
	fmt.Println("Server mode", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *Port))
	if err != nil {
		log.Println("Listen error:", err)
		return
	}
	defer lis.Close()
	fmt.Printf("Server listening on port %d...\n", *Port)
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnect(conn)
	}

}

func help() {
	usage := `Usage: net-speed -mode [mode] [args]
	- mode server: net-speed -mode server -port 8090 - run server and listen port 8090
	- mode client: net-speed -mode client -host 127.0.0.1:8090 - start diag`
	fmt.Println(usage)
}
