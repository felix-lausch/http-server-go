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

	var res []byte
	req, err := ParseRequest(conn)
	if err != nil {
		log.Printf("Error parsing request: %s", err)
		res = FormatResponse(400)
	} else {
		log.Println(req)
		//TODO: handle req & pass result into FormatResponse

		// Send HTTP response
		res = FormatResponse(200)
	}

	_, err = conn.Write(res)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func ParseRequest(conn net.Conn) (Request, error) {
	buffer := make([]byte, 4096) //TODO: make reading smarter? what if body is larger than this?
	n, err := conn.Read(buffer)
	if err != nil {
		return Request{}, err
	}

	content := string(buffer[:n])
	log.Printf("Received:\n%s", content)

	splitContent := strings.Split(content, "\r\n\r\n")
	if len(splitContent) != 2 {
		return Request{}, errors.New("request isn't http formatted")
	}

	requestInfo := splitContent[0]
	headerLines := strings.Split(requestInfo, "\r\n")

	startLineSplit := strings.Split(headerLines[0], " ")
	if len(startLineSplit) != 3 {
		return Request{}, fmt.Errorf("http start line is not correctly formatted: %v", startLineSplit)
	}

	method, err := ParseHttpMethod(startLineSplit[0])
	if err != nil {
		return Request{}, err
	}

	headers := make(map[string]string, len(headerLines[1:]))
	for _, headerLine := range headerLines[1:] {
		splitHeaderLine := strings.Split(headerLine, ": ")
		headers[splitHeaderLine[0]] = splitHeaderLine[1]
	}

	return Request{
		Method:      method,
		Path:        startLineSplit[1],
		HttpVersion: startLineSplit[2],
		Headers:     headers,
		Body:        splitContent[1],
	}, nil
}

type Request struct {
	Method      HttpMethod
	Path        string
	HttpVersion string
	Headers     map[string]string
	Body        string
}

func FormatResponse(statusCode int) []byte {
	var responseString string
	if statusCode == 400 {
		responseString = "HTTP/1.1 400 BAD REQUEST\r\n" +
			"Connection: close\r\n"
	} else {
		responseString = "HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Connection: close\r\n" +
			"Server: felixGoServer/0.1\r\n" +
			"\r\n" +
			"Hello from my custom GOlang server!"
	}

	return []byte(responseString)
}

//go:generate stringer -type=HttpMethod
type HttpMethod int

const (
	GET HttpMethod = iota
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
)

var (
	httpMethodsMap = map[string]HttpMethod{
		"GET":     GET,
		"HEAD":    HEAD,
		"POST":    POST,
		"PUT":     PUT,
		"DELETE":  DELETE,
		"CONNECT": CONNECT,
		"OPTIONS": OPTIONS,
		"TRACE":   TRACE,
		"PATCH":   PATCH,
	}
)

func ParseHttpMethod(str string) (HttpMethod, error) {
	method, ok := httpMethodsMap[strings.ToUpper(str)]
	if !ok {
		return -1, fmt.Errorf("invalid HTTP method: %q", str)
	}

	return method, nil
}
