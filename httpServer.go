package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
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
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal("Failed to load https certificate or private key:", err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Listening on port:", PORT)

	defer tcpListener.Close()

	tlsListener := tls.NewListener(tcpListener, tlsConfig)

	// Accept connections indefinitely
	for {
		conn, err := tlsListener.Accept()
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

	path, queryParams := ParseRequestTarget(startLineSplit[1])

	headers := make(map[string]string, len(headerLines[1:]))
	for _, headerLine := range headerLines[1:] {
		splitHeaderLine := strings.Split(headerLine, ": ")
		headers[splitHeaderLine[0]] = splitHeaderLine[1]
	}

	return Request{
		Method:      method,
		Path:        path,
		HttpVersion: startLineSplit[2],
		QueryParams: queryParams,
		Headers:     headers,
		Body:        splitContent[1],
	}, nil
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
