package main

type Router struct {
	RequestHandlers map[Route]func(req Request) *Response
}

type Route struct {
	Path   string
	Method HttpMethod
}

func NewRouter() *Router {
	return &Router{make(map[Route]func(req Request) *Response)}
}

func (router *Router) AddHandler(route Route, handler func(req Request) *Response) {
	router.RequestHandlers[route] = handler
}

func (router *Router) HandleRequest(req Request) *Response {
	if handler, ok := router.RequestHandlers[req.GetRoute()]; ok {
		//execute request handler and return it's result
		return handler(req)
	}

	//TODO: possibly return 405 if path exists only for other method?

	//if route hasn't been registered return 405
	return NewResponse(NOT_FOUND, nil, "")
}
