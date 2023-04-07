package ginn

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

/**
TODO 按照现在的设计，Engine与RouterFroup是一一对应的，这个是有问题的，一个Engine可以包括多个RouterGroup
*/
type Engine struct {
	// router map[string]HandlerFunc // path与handler的映射
	// router *router
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 分组
type RouterGroup struct {
	prefix      string // 一直累加
	middlewares []HandlerFunc
	// router      *router
	engine *Engine // 所有的group共享一个Engine实例
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)
	return engine
}

func Deafult() *Engine {
	engine := New()
	engine.Use(Recovery())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	// 获取到当前group的engine实例
	engine := group.engine

	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// key：method + path的结合
func (group *RouterGroup) addRouter(method string, path string, handler HandlerFunc) {
	// key := method + "_" + path
	// engine.router[key] = handler

	// 需要加上路由组的前缀
	newPath := group.prefix + path
	log.Printf("addRouter | %s - %s", method, newPath)
	// group.router.addRouter(method, ablosutePath, handler)
	group.engine.router.addRouter(method, newPath, handler)
}

func (gourp *RouterGroup) GET(path string, handler HandlerFunc) {
	gourp.addRouter("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRouter("POST", path, handler)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
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

	// 加入中间件
	var middlewares []HandlerFunc
	// fmt.Printf("ServeHTTP | path: %v, engine.prefix: %v, middlewares: %v\n", req.URL.Path, engine.prefix, engine.middlewares)
	// if strings.HasPrefix(req.URL.Path, engine.prefix) {
	// 	// fmt.Printf("ServeHTTP | engine.prefix: %v, middlewares: %v\n", engine.prefix, engine.middlewares)
	// 	fmt.Println("ServeHTTP | hasPrefix")
	// 	middlewares = append(middlewares, engine.middlewares...)
	// }

	for _, group := range engine.groups {
		fmt.Printf("ServeHTTP | path: %v, group.prefix: %v, middlewares: %v\n", req.URL.Path, group.prefix, group.middlewares)

		if strings.HasPrefix(strings.TrimSpace(req.URL.Path), strings.TrimSpace(group.prefix)) {
			fmt.Println("ServeHTTP | hasPrefix")
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	ctx := newContext(w, req)
	ctx.handlers = middlewares
	engine.router.handle(ctx)

}
