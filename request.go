package main

type Request struct {
	Method      HttpMethod
	Path        string
	HttpVersion string
	Headers     map[string]string
	QueryParams map[string][]string
	PathParams  map[string]string
	Body        string
}

func (req *Request) GetRoute() Route {
	return Route{req.Path, req.Method}
}
