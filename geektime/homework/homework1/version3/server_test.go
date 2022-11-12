package version3

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := Newsdkhttpserver()

	s.Get("/", func(ctx *context) {
		ctx.resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *context) {
		ctx.resp.Write([]byte("hello, user"))
	})

	s.Start_(":8081")
}
