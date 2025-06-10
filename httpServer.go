package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
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
	log.Println("handling Connection:", conn.LocalAddr())

	// Set a reasonable deadline for the entire connection
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	for {
		req, err := ParseRequest(conn)
		if err != nil {
			if err == io.EOF {
				log.Println("Client closed connection")
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Connection timeout")
			} else {
				log.Printf("Error parsing request: %v", err)
				// Send 400 Bad Request if parsing failed
				res := NewResponse(400, nil, "Bad Request")
				conn.Write([]byte(res.String()))
			}
			return
		}

		log.Println(req)
		res := router.HandleRequest(req)

		// Set Connection header based on client request
		if req.Headers["Connection"] == "keep-alive" {
			res.Headers["Connection"] = "keep-alive"
			res.Headers["Keep-Alive"] = "timeout=30"
		}

		// Send HTTP response
		_, err = conn.Write([]byte(res.String()))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}

		// Break loop if client requested connection close
		if req.Headers["Connection"] == "close" {
			return
		}
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

	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading header line: %v", err)
		}

		if headerLine == "\r\n" {
			break
		}

		headerLineSplit := strings.Split(headerLine, ": ")
		req.Headers[headerLineSplit[0]] = headerLineSplit[1]
	}

	// Handle body if present
	// if contentLenStr := req.Headers["Content-Length"]; contentLenStr != "" {
	// contentLength, err := strconv.Atoi(contentLenStr)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid Content-Length: %v", err)
	// }

	// buffered, err := io.ReadAll(io.LimitReader(reader, 2048))
	// req.Body = string(buffered)

	// buffer := make([]byte, 2048)

	// n, err := reader.Read(buffer)
	// if err != nil {
	// 	return nil, err
	// }

	// log.Println("N:", n)

	// req.Body = string(buffer[:n])

	// _, err = io.ReadFull(reader, body)
	// if err != nil {
	// 	return nil, err
	// }

	// req.Body = string(body)
	// }

	return req, nil
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
