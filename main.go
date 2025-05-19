package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

const PORT = 8080

func main() {
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
	}
}

func handleConnection(conn net.Conn, count int) {
	defer conn.Close()

	log.Println("handling Connection:", count)

	req, err := ParseRequest(conn)
	if err != nil {
		log.Printf("Error parsing request: %s", err)
		return
	}

	log.Println(req)

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

func ParseRequest(conn net.Conn) (Request, error) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return Request{}, err
		// log.Printf("Error reading: %v", err)
	}

	content := string(buffer[:n])
	// log.Printf("Received:\n%s", content)

	splitContent := strings.Split(content, "\r\n\r\n")
	if len(splitContent) != 2 {
		return Request{}, errors.New("Request isn't http formatted")
	}

	headers := splitContent[0]
	headerLines := strings.Split(headers, "\r\n")
	firstLineSplit := strings.Split(headerLines[0], " ")

	return Request{
		firstLineSplit[0],
		firstLineSplit[1],
		firstLineSplit[2],
	}, nil
}

type Request struct {
	Method      string
	Path        string
	HttpVersion string
}

func (req Request) String() string {
	return fmt.Sprintf("Method: %v, Path: %v, HttpVersion: %v", req.Method, req.Path, req.HttpVersion)
}
