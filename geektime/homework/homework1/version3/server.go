package version3

import (
	"net"
	"net/http"
)

var _ server = &sdkhttpserver{}

//用于 struct 是否 符合 server 的所有接口
type handlerfunc func(ctx *context)

type server interface {
	http.Handler
	AddRouter(methods, path string, function handlerfunc) //用于注册路由
	Start_(address string)                                //启动server
	Get(path string, handlerfunc handlerfunc)
}

type sdkhttpserver struct {
	*router
}

func Newsdkhttpserver() server {
	return &sdkhttpserver{router: NewRouter()}
}

func (s *sdkhttpserver) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//路由请求的判定
	methods := request.Method
	path := request.URL.Path
	child, ok := s.Findrouter(methods, path)
	if !ok || child.n.function == nil {
		writer.Write([]byte("NOT FOUND"))
		return
	}
	child.n.function(newcontext_map(writer, request, child.pathParams))
}

func (s *sdkhttpserver) Start_(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic("服务器启动失败")
	}
	http.Serve(listen, s)
}

func (s *sdkhttpserver) Get(path string, handlerfunc handlerfunc) {
	s.AddRouter(http.MethodGet, path, handlerfunc)
}
