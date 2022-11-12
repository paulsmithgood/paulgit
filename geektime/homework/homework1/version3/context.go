package version3

import (
	"net/http"
)

type context struct {
	req      *http.Request
	resp     http.ResponseWriter
	Parammap map[string]string
}

func newcontext(req *http.Request, resp http.ResponseWriter) *context {
	return &context{
		req:  req,
		resp: resp,
	}
}

func newcontext_map(resp http.ResponseWriter, req *http.Request, map_ map[string]string) *context {
	return &context{
		req:      req,
		resp:     resp,
		Parammap: map_,
	}
}
