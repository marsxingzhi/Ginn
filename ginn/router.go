package ginn

import (
	"fmt"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "_" + pattern
	r.handlers[key] = handler
}

// func (r *router) GET(pattern string, handler HandlerFunc) {
// 	r.addRouter("GET", pattern, handler)
// }

// func (r *router) POST(pattern string, handler HandlerFunc) {
// 	r.addRouter("POST", pattern, handler)
// }

func (r *router) handle(ctx *Context) {

	key := ctx.Method + "_" + ctx.Path

	if handleFunc, ok := r.handlers[key]; ok {
		handleFunc(ctx)
	} else {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(ctx.Writer, "404 NOT FOUND, Please check method: \"%v\", path: \"%v\" is correct.\n", ctx.Method, ctx.Path)
	}
}
