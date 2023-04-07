package ginn

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	// http.ResponseWriter // 如果这么写的话，无法将外部的writer赋值给Context
	// *http.Request

	Writer http.ResponseWriter
	Req    *http.Request

	Method string
	Path   string

	handlers []HandlerFunc // 挂载到Context实例上的handler，包括中间件和路由映射的handler
	index    int           // 执行到第几个handlers
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
		index:  -1,
	}
}

func (ctx *Context) SetHeader(key, val string) {
	ctx.Writer.Header().Set(key, val)
}

func (ctx *Context) SetStatus(code int) {
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) JSON(code int, obj interface{}) {
	// ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.SetHeader("Content-Type", "application/json;charset=utf-8")
	ctx.SetStatus(code)

	// encode返回
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}

}

func (ctx *Context) Next() {
	ctx.index++
	for n := len(ctx.handlers); ctx.index < n; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}
