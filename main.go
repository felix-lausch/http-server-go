package main

import (
	"fmt"
	"os"
)

const PORT = 8080

func main() {
	server := NewHttpServer(PORT)
	router := server.Router //TODO: like this or pass router in?

	dynamicHandler := func(req *Request) *Response {
		body := "Recieved:\r\n"

		if req.PathParams != nil {
			for k, v := range req.PathParams {
				body += fmt.Sprintf("Key: %v Value: %v \r\n", k, v)
			}
			body += "\r\n"
		}

		return NewResponse(OK, nil, body)
	}

	indexHandler := func(req *Request) *Response {
		return ServeStaticHtml("static/index.html")
	}

	aboutHandler := func(req *Request) *Response {
		return ServeStaticHtml("static/about.html")
	}

	router.AddHandler("/", indexHandler)
	router.AddHandler("/about", aboutHandler)
	router.AddHandler("/articles/:id", dynamicHandler)
	router.AddHandler("/posts/:postId/comments/:commentId", dynamicHandler)

	server.Start()
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
