package version3

import "fmt"

func Func1(ctx *context) {
	ctx.resp.Write([]byte("uryyb"))
}
func Func2(ctx *context) {
	ctx.resp.Write([]byte("欢迎访问 func2\n"))
	ctx.resp.Write([]byte(fmt.Sprintf("path：%s", ctx.req.URL.Path)))
}
func Func3(ctx *context) {
	ctx.resp.Write([]byte(fmt.Sprintf("path：%s\n", ctx.req.URL.Path)))
	ctx.resp.Write([]byte(fmt.Sprintf("MAP：%s\n", ctx.Parammap)))
}
func Basefunc(ctx *context) {
	ctx.resp.Write([]byte("访问到了这个路由"))
}
