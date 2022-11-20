package main

/*
version6
@@2022-11-18 07.31--zh

*/
import (
	"fmt"
	"go_/exercise/geektime/homework2/version6"
)

func main() {
	server_ := version6.NewHttp_Server()
	server_.Get("/order/*", func(ctx *version6.Context) {
		fmt.Println("hi")
	})
	server_.Get("/user", func(ctx *version6.Context) {
		fmt.Println("hi")
	})
	server_.Start_(":8083")
}
