package main

import (
	"fmt"
	"strings"
)

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
