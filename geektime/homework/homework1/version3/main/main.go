package main

import (
	"geektime_go_go/go_/Exercise/geektime/homework1/version2"
	"net/http"
)

/*
@@2022-11-11 19.49--zh
*/

func main() {
	server_ := version2.Newsdkhttpserver()
	//server_.AddRouter(http.MethodGet, "/", version2.Func1)
	//server_.AddRouter(http.MethodDelete, "/a/*", version2.Func2)
	//server_.AddRouter(http.MethodPost, "/:username", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/reg/:email(^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$)", version2.Func3)

	//server_.AddRouter(http.MethodPut, "/*", version2.Func1)
	//server_.AddRouter(http.MethodPut, "/*/*", version2.Func1)

	//server_.AddRouter(http.MethodGet, "/:id", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/:ok", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/reg/:email(^\\w+([-+.]\\w+)*W\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$)", version2.Func3) //需要处理一下 同一节点多个的问题
	//server_.AddRouter(http.MethodPut, "/user/:id", version2.Func3)
	//server_.AddRouter(http.MethodPut, "/user/:username", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/a/*", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/a/:id", version2.Func3)

	server_.AddRouter(http.MethodPost, "/order/*", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/a/b/*", version2.Func3)
	//server_.AddRouter(http.MethodGet, "/a/b/:id(.*)", version2.Func3)
	server_.Start_(":8081")
}
