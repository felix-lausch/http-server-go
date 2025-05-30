package main

import "strings"

type TrieRouter struct {
	Root *RouterNode
	//TODO: what happens if i add POST for root path? doesnt really work

}

type RouterNode struct {
	Handler    func(req Request) *Response
	Method     HttpMethod
	Children   map[string]*RouterNode
	ParamChild *RouterNode
	Segment    string
	IsParam    bool
	ParamName  string
}

func NewTrieRouter() *TrieRouter {
	return &TrieRouter{
		Root: &RouterNode{
			Method:   GET,
			Children: map[string]*RouterNode{},
			Segment:  "/",
		},
	}
}

func (router *TrieRouter) AddHandler(pattern string, handler func(req Request) *Response) {
	current := router.Root
	segments := strings.SplitSeq(pattern, "/")

	for seg := range segments {
		if strings.HasPrefix(seg, ":") {
			if current.ParamChild == nil {
				paramChild := &RouterNode{
					Children:  map[string]*RouterNode{},
					Segment:   seg,
					IsParam:   true,
					ParamName: seg[1:],
				}
				current.ParamChild = paramChild
			}
			current = current.ParamChild
		} else {
			if child, ok := current.Children[seg]; ok {
				//child exists, move along
				current = child
			} else {
				//node with seg doesnt exist yet, create new node
				newChild := &RouterNode{
					Children: map[string]*RouterNode{},
					Segment:  seg,
				}
				current.Children[seg] = newChild
				current = newChild
			}
		}
	}

	current.Handler = handler
	current.Method = GET //TODO: pass in method
}

func (router *TrieRouter) HandleRequest(req Request) *Response {
	current := router.Root
	pathSegments := strings.Split(req.Path, "/")
	params := make(map[string]string)

	for _, seg := range pathSegments {
		if node, ok := current.Children[seg]; ok {
			current = node
		} else if current.ParamChild != nil {
			params[current.ParamChild.ParamName] = seg
			current = current.ParamChild
		} else {
			return NewResponse(NOT_FOUND, nil, "")
		}
	}

	if current.Handler == nil {
		return NewResponse(NOT_FOUND, nil, "")
	}

	if req.Method != current.Method {
		return NewResponse(METHOD_NOT_ALLOWED, nil, "")
	}

	req.PathParams = params
	return current.Handler(req)
}
