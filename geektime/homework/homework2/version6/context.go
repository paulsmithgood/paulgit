package version6

import "net/http"

type Context struct {
	Resp      http.ResponseWriter
	Req       *http.Request
	Param_Map map[string]string

	RespData     []byte
	RespDatacode int

	MatchedRoute string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Resp: w,
		Req:  r,
	}
}
