package version6

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		// 通配符测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}

	mockHandler := func(ctx *Context) {}
	r := Newrouter()
	for _, tr := range testRoutes {
		r.addrouter(tr.method, tr.path, mockHandler, nil)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {path: "user", children: map[string]*node{
						"home": {path: "home", funtion_: mockHandler},
					}, funtion_: mockHandler},
					"order": {path: "order", children: map[string]*node{
						"detail": {path: "detail", funtion_: mockHandler},
					}, starchildren: &node{path: "*", funtion_: mockHandler}},
					"param": {
						path: "param",
						paramchildren: &node{
							path: ":id",
							starchildren: &node{
								path:     "*",
								funtion_: mockHandler,
							},
							children: map[string]*node{"detail": {path: "detail", funtion_: mockHandler}},
							funtion_: mockHandler,
						},
					},
				},
				starchildren: &node{
					path: "*",
					children: map[string]*node{
						"abc": {
							path:         "abc",
							starchildren: &node{path: "*", funtion_: mockHandler},
							funtion_:     mockHandler},
					},
					starchildren: &node{path: "*", funtion_: mockHandler},
					funtion_:     mockHandler},
				funtion_: mockHandler},
			http.MethodPost: {path: "/", children: map[string]*node{
				"order": {path: "order", children: map[string]*node{
					"create": {path: "create", funtion_: mockHandler},
				}},
				"login": {path: "login", funtion_: mockHandler},
			}},
		},
	}
	msg, ok := wantRouter.equal(*r)
	assert.True(t, ok, msg)

	// 非法用例
	r = Newrouter()

	// 空字符串
	assert.PanicsWithValue(t, "路由不可以为空!", func() {
		r.addrouter(http.MethodGet, "", mockHandler, nil)
	})

	// 前导没有 /
	assert.PanicsWithValue(t, "路由必须是以/为前导", func() {
		r.addrouter(http.MethodGet, "a/b/c", mockHandler, nil)
	})

	// 后缀有 /
	assert.PanicsWithValue(t, "路由不能以/结束", func() {
		r.addrouter(http.MethodGet, "/a/b/c/", mockHandler, nil)
	})

	// 根节点重复注册
	r.addrouter(http.MethodGet, "/", mockHandler, nil)
	assert.PanicsWithValue(t, "根节点不允许重复注册", func() {
		r.addrouter(http.MethodGet, "/", mockHandler, nil)
	})
	// 普通节点重复注册
	r.addrouter(http.MethodGet, "/a/b/c", mockHandler, nil)
	assert.PanicsWithValue(t, "该节点不允许多次注册！", func() {
		r.addrouter(http.MethodGet, "/a/b/c", mockHandler, nil)
	})

	// 多个 /
	assert.PanicsWithValue(t, "path 带多//,不允许", func() {
		r.addrouter(http.MethodGet, "/a//b", mockHandler, nil)
	})
	//assert.PanicsWithValue(t, "path 带多//,不允许", func() {
	//	r.addrouter(http.MethodGet, "//a/b", mockHandler, nil)
	//})

	// 同时注册通配符路由和参数路由
	assert.PanicsWithValue(t, "不允许 有 * 和 : 并存", func() {
		r.addrouter(http.MethodGet, "/a/*", mockHandler, nil)
		r.addrouter(http.MethodGet, "/a/:id", mockHandler, nil)
	})
	assert.PanicsWithValue(t, "不允许 * 和 : 并存", func() {
		r.addrouter(http.MethodGet, "/a/b/:id", mockHandler, nil)
		r.addrouter(http.MethodGet, "/a/b/*", mockHandler, nil)
	})
	r = Newrouter()
	assert.PanicsWithValue(t, "不允许 有 * 和 : 并存", func() {
		r.addrouter(http.MethodGet, "/*", mockHandler, nil)
		r.addrouter(http.MethodGet, "/:id", mockHandler, nil)
	})
	r = Newrouter()
	assert.PanicsWithValue(t, "不允许 * 和 : 并存", func() {
		r.addrouter(http.MethodGet, "/:id", mockHandler, nil)
		r.addrouter(http.MethodGet, "/*", mockHandler, nil)
	})

	// 参数冲突
	assert.PanicsWithValue(t, "不允许注册 有多个 :", func() {
		r.addrouter(http.MethodGet, "/a/b/c/:id", mockHandler, nil)
		r.addrouter(http.MethodGet, "/a/b/c/:name", mockHandler, nil)
	})
}

func (r router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为 nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}

	nhv := reflect.ValueOf(n.funtion_)
	yhv := reflect.ValueOf(y.funtion_)
	if nhv != yhv {
		return fmt.Sprintf("%s 节点 handler 不相等 x %s, y %s", n.path, nhv.Type().String(), yhv.Type().String()), false
	}

	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}

	if n.starchildren != nil {
		str, ok := n.starchildren.equal(y.starchildren)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}
	}

	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
	}

	mockHandler := func(ctx *Context) {}

	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchInfo
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "/",
					funtion_: mockHandler,
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "user",
					funtion_: mockHandler,
				},
			},
		},
		{
			name:   "no handler",
			method: http.MethodPost,
			path:   "/order",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "create",
					funtion_: mockHandler,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:   "star match",
			method: http.MethodPost,
			path:   "/order/delete",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "*",
					funtion_: mockHandler,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:   "star in middle",
			method: http.MethodGet,
			path:   "/user/Tom/home",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "home",
					funtion_: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/param/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "id",
					funtion_: mockHandler,
				},
				param_map: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/*
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/abc",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "*",
					funtion_: mockHandler,
				},
				param_map: map[string]string{"id": "123"},
			},
		},

		{
			// 命中 /param/:id/detail
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/detail",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "detail",
					funtion_: mockHandler,
				},
				param_map: map[string]string{"id": "123"},
			},
		},
	}

	r := Newrouter()
	for _, tr := range testRoutes {
		r.addrouter(tr.method, tr.path, mockHandler, nil)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findrouter(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.param_map, mi.param_map)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.funtion_)
			nVal := reflect.ValueOf(n.funtion_)
			assert.Equal(t, wantVal, nVal)
		})
	}
}

func Test_findRoute_Middleware(t *testing.T) {
	var mdlBuilder = func(i byte) Middleware {
		return func(next Handlerfunc) Handlerfunc {
			return func(ctx *Context) {
				ctx.RespData = append(ctx.RespData, i)
				next(ctx)
			}
		}
	}
	mdlsRoute := []struct {
		method string
		path   string
		mdls   []Middleware
	}{
		{
			method: http.MethodGet,
			path:   "/a/b",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b')},
		},
		{
			method: http.MethodGet,
			path:   "/a/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('*')},
		},
		{
			method: http.MethodGet,
			path:   "/a/b/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('*')},
		},
		{
			method: http.MethodPost,
			path:   "/a/b/*",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('*')},
		},
		{
			method: http.MethodPost,
			path:   "/a/*/c",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('*'), mdlBuilder('c')},
		},
		{
			method: http.MethodPost,
			path:   "/a/b/c",
			mdls:   []Middleware{mdlBuilder('a'), mdlBuilder('b'), mdlBuilder('c')},
		},
		{
			method: http.MethodDelete,
			path:   "/*",
			mdls:   []Middleware{mdlBuilder('*')},
		},
		{
			method: http.MethodDelete,
			path:   "/",
			mdls:   []Middleware{mdlBuilder('/')},
		},
	}
	r := Newrouter()
	for _, mdlRoute := range mdlsRoute {
		r.addrouter(mdlRoute.method, mdlRoute.path, nil, mdlRoute.mdls)
	}
	testCases := []struct {
		name   string
		method string
		path   string
		// 我们借助 ctx 里面的 RespData 字段来判断 middleware 有没有按照预期执行
		wantResp string
	}{
		{
			name:   "static, not match",
			method: http.MethodGet,
			path:   "/a",
		},
		{
			name:     "static, match",
			method:   http.MethodGet,
			path:     "/a/c",
			wantResp: "a*",
		},
		{
			name:     "static and star",
			method:   http.MethodGet,
			path:     "/a/b",
			wantResp: "a*ab",
		},
		{
			name:     "static and star",
			method:   http.MethodGet,
			path:     "/a/b/c",
			wantResp: "a*abab*",
		},
		{
			name:     "abc",
			method:   http.MethodPost,
			path:     "/a/b/c",
			wantResp: "a*cab*abc",
		},
		{
			name:     "root",
			method:   http.MethodDelete,
			path:     "/",
			wantResp: "/",
		},
		{
			name:     "root star",
			method:   http.MethodDelete,
			path:     "/a",
			wantResp: "/*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, _ := r.findrouter(tc.method, tc.path)
			mdls := mi.mdls
			var root Handlerfunc = func(ctx *Context) {
				// 使用 string 可读性比较高
				assert.Equal(t, tc.wantResp, string(ctx.RespData))
			}
			for i := len(mdls) - 1; i >= 0; i-- {
				root = mdls[i](root)
			}
			// 开始调度
			root(&Context{
				RespData: make([]byte, 0, len(tc.wantResp)),
			})
		})
	}

}
