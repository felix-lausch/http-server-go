package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
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

	router := NewRouter()
	AddRequestHandlers(router)

	// Accept connections indefinitely
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		// Handle each connection in a new goroutine
		go handleConnection(conn, router)
	}
}

func handleConnection(conn net.Conn, router *Router) {
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

// TODO: add query params and seperate them from path
type Request struct {
	Method      HttpMethod
	Path        string
	HttpVersion string
	Headers     map[string]string
	Body        string
}

func (req *Request) GetRoute() Route {
	return Route{req.Path, req.Method}
}

type Response struct {
	StatusCode  StatusCode
	HttpVersion string
	Headers     map[string]string
	Body        string
}

func NewResponse(statusCode StatusCode, headers map[string]string, body string) *Response {
	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-Length"] = fmt.Sprint(len(body))
	headers["Server"] = "felixGoServer/0.1"
	headers["Connection"] = "close"

	return &Response{
		statusCode,
		"HTTP/1.1",
		headers,
		body,
	}
}

func (r Response) String() string {
	statusCodeString := strings.ReplaceAll(r.StatusCode.String(), "_", " ")
	res := fmt.Sprintf("%v %v %v\r\n", r.HttpVersion, int(r.StatusCode), statusCodeString)

	for k, v := range r.Headers {
		res += fmt.Sprintf("%v: %v\r\n", k, v)
	}

	res += "\r\n"
	res += r.Body

	return res
}

func FormatResponse(statusCode StatusCode, body string) []byte {
	res := NewResponse(statusCode, make(map[string]string), body)
	return []byte(res.String())
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

//go:generate stringer -type=StatusCode
type StatusCode int

const (
	OK                    StatusCode = 200
	CREATED               StatusCode = 201
	ACCEPTED              StatusCode = 202
	NO_CONTENT            StatusCode = 204
	FOUND                 StatusCode = 302
	BAD_REQUEST           StatusCode = 400
	UNAUTHORIZED          StatusCode = 401
	FORBIDDEN             StatusCode = 403
	NOT_FOUND             StatusCode = 404
	METHOD_NOT_ALLOWED    StatusCode = 405
	CONTENT_TOO_LARGE     StatusCode = 413
	URI_TOO_LONG          StatusCode = 414
	INTERNAL_SERVER_ERROR StatusCode = 500
)

// var (
// 	statusCodeMap = map[string]StatusCode{
// 		"OK":     OK,
// 		"CREATED":     CREATED,
// 		"ACCEPTED":     ACCEPTED,
// 		"NO_CONTENT":     NO_CONTENT,
// 		"FOUND":     FOUND,
// 		"BAD_REQUEST":     BAD_REQUEST,
// 		"UNAUTHORIZED":     UNAUTHORIZED,
// 		"FORBIDDEN":     FORBIDDEN,
// 		"NOT_FOUND":     NOT_FOUND,
// 		"METHOD_NOT_ALLOWED":     METHOD_NOT_ALLOWED,
// 		"CONTENT_TOO_LARGE":     CONTENT_TOO_LARGE,
// 		"NO_CONTENT":     NO_CONTENT,
// 		"NO_CONTENT":     NO_CONTENT,
// 		"NO_CONTENT":     NO_CONTENT,
// 	}
// )

func AddRequestHandlers(router *Router) {
	indexHandler := func(req Request) *Response {
		return ServeStaticHtml("static/index.html")
	}

	aboutHandler := func(req Request) *Response {
		return ServeStaticHtml("static/about.html")
	}

	router.RequestHandlers[Route{"/", GET}] = indexHandler
	router.RequestHandlers[Route{"/about", GET}] = aboutHandler
}

func ServeStaticHtml(file string) *Response {
	content, err := os.ReadFile(file)
	if err != nil {
		return NewResponse(INTERNAL_SERVER_ERROR, nil, "Internal Server Error")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "text/html"

	return NewResponse(200, headers, string(content))
}
