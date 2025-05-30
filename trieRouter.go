package main

import "strings"

type TrieRouter struct {
	Root *RouterNode
	//TODO: what happens if i add POST for root path? doesnt really work

}

type RouterNode struct {
	Handler  func(req Request) *Response
	Method   HttpMethod
	Children map[string]*RouterNode
	Segment  string
	// isParam  bool
	// paramName string
}

func NewTrieRouter() *TrieRouter {
	return &TrieRouter{
		Root: &RouterNode{
			Handler:  nil,
			Method:   GET,
			Children: map[string]*RouterNode{},
			Segment:  "/",
		},
	}
}

func (router *TrieRouter) AddHandler(pattern string, handler func(req Request) *Response) {
	currentNode := router.Root
	segments := strings.SplitSeq(pattern, "/")

	for seg := range segments {
		//TODO: determine isParam

		if child, ok := currentNode.Children[seg]; ok {
			currentNode = child
		} else {
			//node with seg doesnt exist yet, create new node
			newChild := RouterNode{
				Children: map[string]*RouterNode{},
				Segment:  seg,
			}
			currentNode.Children[seg] = &newChild
			currentNode = &newChild
		}
	}

	currentNode.Handler = handler
	currentNode.Method = GET //TODO: pass in method
}

func (router *TrieRouter) HandleRequest(req Request) *Response {
	root := router.Root
	pathSegments := strings.Split(req.Path, "/")

	for _, seg := range pathSegments {
		//TODO: wwhat if children is nil?

		if node, ok := root.Children[seg]; ok {
			root = node
		} else {
			return NewResponse(NOT_FOUND, nil, "")
		}
	}

	if root.Handler == nil {
		return NewResponse(NOT_FOUND, nil, "")
	}

	return root.Handler(req)
}
