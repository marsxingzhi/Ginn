package ginn

import (
	"fmt"
	"log"
	"net/http"
)

/**
RouterGroup
1. 继承router路由功能
2. 实现分组功能

Engine
1. 继承RouterGroup的分组功能
2. 实现Run、ServeHTTP接口
*/

// type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

type H map[string]interface{}

type Engine struct {
	// router map[string]HandlerFunc // path与handler的映射
	// router *router
	*RouterGroup
}

// 分组
type RouterGroup struct {
	prefix      string // 一直累加
	middlewares []HandlerFunc
	router      *router
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix:      group.prefix + prefix,
		middlewares: group.middlewares,
		router:      group.router,
	}
	return newGroup
}

func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{router: newRouter()}
	return engine
}

// key：method + path的结合
func (group *RouterGroup) addRouter(method string, path string, handler HandlerFunc) {
	// key := method + "_" + path
	// engine.router[key] = handler

	// 需要加上路由组的前缀
	ablosutePath := group.prefix + path
	log.Printf("addRouter | %s - %s", method, ablosutePath)
	group.router.addRouter(method, ablosutePath, handler)
}

func (gourp *RouterGroup) GET(path string, handler HandlerFunc) {
	gourp.addRouter("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRouter("POST", path, handler)
}

func (engine *Engine) Run(path string) error {
	return http.ListenAndServe(path, engine)
}

// 实现http.handler接口
// 只要有请求进来，就会执行这个方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintf(w, "[ginn | Engine] URL | host: %v, schema: %v, path: %v\n", req.URL.Host, req.URL.Scheme, req.URL.Path)
	fmt.Printf("[ginn | Engine] URL | host: %v, schema: %v, path: %v\n", req.URL.Host, req.URL.Scheme, req.URL.Path)

	// for k, v := range req.Header {
	// 	fmt.Fprintf(w, "[Engine] | Header[%v] = %v\n", k, v)

	// }

	// key := req.Method + "_" + req.URL.Path

	// if handler, ok := engine.router[key]; ok {
	// 	handler(w, req)
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	fmt.Fprintf(w, "404 NOT FOUND, Please check method: \"%v\", path: \"%v\" is correct.\n", req.Method, req.URL.Path)
	// }

	ctx := newContext(w, req)
	engine.router.handle(ctx)

}
