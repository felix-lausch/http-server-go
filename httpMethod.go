package main

import (
	"fmt"
	"strings"
)

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
