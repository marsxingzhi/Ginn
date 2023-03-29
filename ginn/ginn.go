package ginn

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc // path与handler的映射
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

// key：method + path的结合
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "_" + pattern
	engine.router[key] = handler
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

func (engine *Engine) Run(pattern string) error {
	return http.ListenAndServe(pattern, engine)
}

// 实现http.handler接口
// 只要有请求进来，就会执行这个方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "[ginn | Engine] URL | host: %v, schema: %v, path: %v\n", req.URL.Host, req.URL.Scheme, req.URL.Path)

	// for k, v := range req.Header {
	// 	fmt.Fprintf(w, "[Engine] | Header[%v] = %v\n", k, v)

	// }

	key := req.Method + "_" + req.URL.Path

	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND, Please check method: \"%v\", path: \"%v\" is correct.\n", req.Method, req.URL.Path)
	}

}
