package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	DEFAULT_BUFFER_SIZE = 2048
)

type HttpServer struct {
	Port   int
	Router *TrieRouter
}

func NewHttpServer(port int) *HttpServer {
	return &HttpServer{
		port,
		NewTrieRouter(),
	}
}

func (s *HttpServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Listening on port:", PORT)

	defer listener.Close()

	// Accept connections indefinitely
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a new goroutine
		go handleConnection(conn, s.Router)
	}
}

func handleConnection(conn net.Conn, router *TrieRouter) {
	defer conn.Close()
	var res *Response

	log.Println("handling Connection:", conn.LocalAddr())

	req, err := ParseRequest(conn)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing request: %s", err)

		log.Print(errMsg)
		res = NewResponse(400, nil, errMsg)
	} else {
		log.Println(req)
		res = router.HandleRequest(req)
	}

	// Send HTTP response
	_, err = conn.Write([]byte(res.String()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func ParseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	startLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	startLineSplit := strings.Split(strings.TrimSpace(startLine), " ")
	if len(startLineSplit) < 3 {
		return nil, fmt.Errorf("http start line is not correctly formatted: %v", startLine)
	}

	method, err := ParseHttpMethod(startLineSplit[0])
	if err != nil {
		return nil, err
	}

	path, queryParams := ParseRequestTarget(startLineSplit[1])

	req := &Request{
		Method:      method,
		Path:        path,
		Headers:     make(map[string]string),
		QueryParams: queryParams,
		HttpVersion: startLineSplit[2],
	}

	//Read headers line by line
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header line: %v", err)
		}

		if headerLine == "\r\n" {
			break
		}

		headerLine = strings.TrimSuffix(headerLine, "\r\n")
		headerLineSplit := strings.Split(headerLine, ": ")

		req.Headers[headerLineSplit[0]] = headerLineSplit[1]
	}

	body, err := ParseBody(req.Headers, reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing body: %v", err)
	}

	req.Body = body
	return req, nil
}

func ParseBody(requestHeaders map[string]string, reader *bufio.Reader) (string, error) {
	bufferSize := DEFAULT_BUFFER_SIZE

	if len, ok := requestHeaders["Content-Length"]; ok {
		contentLength, err := strconv.Atoi(len)
		if err != nil {
			return "", fmt.Errorf("error parsing content length header: %v", err)
		}

		bufferSize = contentLength
	}

	buffer := make([]byte, bufferSize)

	n, err := reader.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("error reading request body: %v", err)
	}

	//TODO: handle reading bodies larger than default buffer size

	return string(buffer[:n]), nil
}

func ParseRequestTarget(requestTarget string) (path string, queryArgs map[string][]string) {
	splitRequestTarget := strings.Split(requestTarget, "?")

	if len(splitRequestTarget) == 1 {
		return splitRequestTarget[0], make(map[string][]string)
	}

	splitArgs := strings.Split(splitRequestTarget[1], "&")

	args := make(map[string][]string, len(splitArgs))
	for _, arg := range splitArgs {
		keyValue := strings.SplitN(arg, "=", 2)

		if val, ok := args[keyValue[0]]; ok {
			args[keyValue[0]] = append(val, keyValue[1])
		} else {
			args[keyValue[0]] = []string{keyValue[1]}
		}
	}

	return splitRequestTarget[0], args
}
