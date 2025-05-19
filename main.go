package main

import (
	"fmt"
	"log"
	"net"
)

const PORT = 8080

func main() {
	// http.HandleFunc("/", handler)
	// http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer listener.Close()
	log.Println("Listening on port:", PORT)

	requestCount := 0
	// Accept connections indefinitely
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a new goroutine
		go handleConnection(conn, requestCount)
		requestCount++
		// log.Println("Number of handled requests:", requestCount)
	}
}

func handleConnection(conn net.Conn, count int) {
	defer conn.Close()

	log.Println("handling Connection:", count)

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading: %v", err)
	}

	log.Printf("Received:\n%s", string(buffer[:n]))

	// Send HTTP response
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: close\r\n" +
		"\r\n" +
		"Hello from your TCP server now speaking HTTP!\n" +
		"Requests handled: " + fmt.Sprint(count) + "\r\n"

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
