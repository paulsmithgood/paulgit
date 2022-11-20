package version6

import (
	"log"
	"net"
	"net/http"
)

type Handlerfunc func(ctx *Context)

type server interface {
	http.Handler
	addrouter(methods, path string, handler Handlerfunc, middleware ...Middleware)
	Start_(address string)
}

type Http_Server struct {
	router *router
	Mdls   []Middleware //这里的middleware 是 适用于全局的middelware
}

func (h *Http_Server) Use_middleware(middleware ...Middleware) {
	//	这里其实是添加全局middleware
	if h.Mdls == nil {
		//	初始化
		h.Mdls = make([]Middleware, 0)
	}
	h.Mdls = append(h.Mdls, middleware...)
}

func NewHttp_Server() *Http_Server {
	return &Http_Server{router: Newrouter()}
}

func (h *Http_Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//这里说明已经请求了。
	ctx := NewContext(writer, request)
	root := h.server
	//	在这里做拼接
	for i := len(h.Mdls) - 1; i >= 0; i-- {
		root = h.Mdls[i](root)
	}
	var flashrespfunction Middleware = func(next Handlerfunc) Handlerfunc {
		return func(ctx *Context) {
			next(ctx)
			if ctx.RespDatacode > 0 {
				ctx.Resp.WriteHeader(ctx.RespDatacode)
			}
			_, err := ctx.Resp.Write(ctx.RespData)
			if err != nil {
				log.Fatalf("写入失败,原因是：%s\n", err)
			}
		}
	}
	root = flashrespfunction(root)
	root(ctx)
}

func (h *Http_Server) server(ctx *Context) {
	mi, ok := h.router.findrouter(ctx.Req.Method, ctx.Req.URL.Path)
	//fmt.Println(mi, ok)
	if !ok {
		//	说明没有找到--
		ctx.RespDatacode = http.StatusNotFound
		ctx.RespData = []byte("Not found")
		return
	}
	ctx.Param_Map = mi.param_map
	ctx.MatchedRoute = mi.n.route

	//	把mi 里面的 middelware 封装 在 mi.handler
	root := mi.n.funtion_
	for i := len(mi.mdls) - 1; i >= 0; i-- {
		root = mi.mdls[i](root)
	}

	//	拼接完成
	root(ctx)

}

func (h *Http_Server) addrouter(methods, path string, handler Handlerfunc, middleware ...Middleware) {
	h.router.addrouter(methods, path, handler, middleware)
}
func (h *Http_Server) Get(path string, handler Handlerfunc, middleware ...Middleware) {
	h.router.addrouter(http.MethodGet, path, handler, middleware)
}
func (h *Http_Server) Start_(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	http.Serve(l, h)
}
