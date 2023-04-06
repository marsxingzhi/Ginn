package ginn

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*TrieNode
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*TrieNode),
	}
}

/**
method作为根节点，HTTP一共有9个方法，GET、POST、HEAD...

1. 将method与url作为key，handle作为value保存到map中
2. 将请求路由插入到前缀树
*/
func (r *router) addRouter(method string, path string, handler HandlerFunc) {
	key := method + "_" + path
	r.handlers[key] = handler

	root, ok := r.roots[method]
	if !ok {
		root = &TrieNode{}
		r.roots[method] = root
	}
	root.insert(path, parsePath(path), 0)

}

// func (r *router) getRoute(method string, path string) (*node, map[string]string)
func (r *router) getRoute(method, path string) *TrieNode {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	vals := parsePath(path)
	if node, exist := root.search(vals, 0); exist {
		return node
	}
	return nil
}

func parsePath(path string) []string {
	parts := strings.Split(path, "/")

	res := make([]string, 0)

	for _, item := range parts {
		if item != "" {
			res = append(res, item)
		}
	}
	return res
}

// func (r *router) GET(pattern string, handler HandlerFunc) {
// 	r.addRouter("GET", pattern, handler)
// }

// func (r *router) POST(pattern string, handler HandlerFunc) {
// 	r.addRouter("POST", pattern, handler)
// }

func (r *router) handle(ctx *Context) {

	// key := ctx.Method + "_" + ctx.Path

	node := r.getRoute(ctx.Method, ctx.Path)

	if node != nil {

		key := ctx.Method + "_" + node.path

		if handleFunc, ok := r.handlers[key]; ok {
			handleFunc(ctx)
		} else {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(ctx.Writer, "404 NOT FOUND, Please check method: \"%v\", path: \"%v\" is correct.\n", ctx.Method, ctx.Path)
		}

	} else {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(ctx.Writer, "404 NOT FOUND, Please check method: \"%v\", path: \"%v\" is correct.\n", ctx.Method, ctx.Path)
	}

}
