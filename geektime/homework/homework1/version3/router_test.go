package version3

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

var func_ handlerfunc = func(ctx *context) {

}

func TestRouter_AddRouter(t *testing.T) {
	testRouter := []struct {
		methods string
		path    string
	}{
		// 静态节点测试
		{
			methods: http.MethodGet,
			path:    "/",
		},
		{
			methods: http.MethodGet,
			path:    "/user",
		},
		{
			methods: http.MethodGet,
			path:    "/user/home",
		},
		{
			methods: http.MethodGet,
			path:    "/order/detail",
		},
		{
			methods: http.MethodGet,
			path:    "/order/create",
		},
		{
			methods: http.MethodPost,
			path:    "/login",
		},
		//	通配符测试用例
		{
			methods: http.MethodGet,
			path:    "/order/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*/abc",
		},
		{
			methods: http.MethodGet,
			path:    "/*/abc/*",
		},
		// 参数路由
		{
			methods: http.MethodGet,
			path:    "/param/:id",
		},
		{
			methods: http.MethodGet,
			path:    "/param/:id/detail",
		},
		{
			methods: http.MethodGet,
			path:    "/param/:id/*",
		},
		{
			methods: http.MethodDelete,
			path:    "/reg/:id(.*)",
		},
		{
			methods: http.MethodDelete,
			path:    "/:name(^.+$)/abc",
		},
	}

	router_ := NewRouter()
	var function_ handlerfunc = func(ctx *context) {
		fmt.Println("uryyb")
	}
	for _, values := range testRouter {
		router_.AddRouter(values.methods, values.path, function_)
	}
	//	使用我们的addrouter 创建了一颗树
	regexp1, err := regexp.Compile(".*")
	if err != nil {
		panic("编译失败")
	}
	regexp2, err := regexp.Compile("^.+$")
	if err != nil {
		panic("编译失败")
	}
	Wanttree := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:     "/",
				function: function_,
				starchild: &node{
					path:     "*",
					function: function_,
					starchild: &node{
						path:     "*",
						function: function_,
					},
					children: map[string]*node{
						"abc": &node{
							path:     "abc",
							function: function_,
							starchild: &node{
								path:     "*",
								function: function_,
							},
						},
					},
				},
				children: map[string]*node{
					"user": &node{
						path:     "user",
						function: function_,
						children: map[string]*node{
							"home": &node{
								path:     "home",
								function: function_,
							},
						},
					},
					"order": &node{
						path: "order",
						starchild: &node{
							path:     "*",
							function: function_,
						},
						children: map[string]*node{
							"detail": &node{
								path:     "detail",
								function: function_,
							},
							"create": &node{
								path:     "create",
								function: function_,
							},
						},
					},
					"param": &node{
						path: "param",
						paramchild: &node{
							path:     "id",
							function: function_,
							children: map[string]*node{
								"detail": &node{
									path:     "detail",
									function: function_,
								},
							},
							starchild: &node{
								path:     "*",
								function: function_,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"login": &node{
						path:     "login",
						function: function_,
					},
				},
			},
			http.MethodDelete: &node{
				path: "/",
				regchild: &node{
					path: "name",
					//function: function_,
					regexp: regexp2,
					children: map[string]*node{
						"abc": &node{
							path:     "abc",
							function: function_,
						},
					},
				},
				children: map[string]*node{
					"reg": &node{
						path: "reg",
						regchild: &node{
							path:     "id",
							function: function_,
							regexp:   regexp1,
						},
					},
				},
			},
		},
	}
	fmt.Println(router_, Wanttree)
	msg, ok := Wanttree.equal(router_)
	assert.True(t, ok, msg)

	//	panic 测试
	newrouter_ := NewRouter()

	// 空字符串
	assert.PanicsWithValue(t, "路由不可以为空!", func() {
		newrouter_.AddRouter(http.MethodGet, "", function_)
	})
	// 前导没有 /
	assert.PanicsWithValue(t, "路由必须是以/为前导", func() {
		newrouter_.AddRouter(http.MethodGet, "a/b/c", function_)
	})
	// 后缀有 /
	assert.PanicsWithValue(t, "路由不能以/结束", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/c/", function_)
	})
	// 根节点重复注册
	newrouter_.AddRouter(http.MethodGet, "/", function_)
	assert.PanicsWithValue(t, "根节点不允许重复注册", func() {
		newrouter_.AddRouter(http.MethodGet, "/", function_)
	})
	// 普通节点重复注册
	newrouter_.AddRouter(http.MethodGet, "/a/b/c", function_)
	assert.PanicsWithValue(t, "路由不能重复注册", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/c", function_)
	})
	assert.PanicsWithValue(t, "路由不能存在连续两个//", func() {
		newrouter_.AddRouter(http.MethodGet, "/a//b", function_)
	})
	//assert.PanicsWithValue(t, "路由不能存在连续两个//", func() {
	//	newrouter_.AddRouter(http.MethodGet, "//a/b", function_)
	//})
	//同时注册通配符路由，参数路由，正则路由
	assert.PanicsWithValue(t, "不允许1", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/*", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/:id", function_)
	})
	//newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "不允许1", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/*", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "不允许1", func() {
		newrouter_.AddRouter(http.MethodGet, "/*", function_)
		newrouter_.AddRouter(http.MethodGet, "/:id", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "通配符路径 与 路径参数 不可以同时存在", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/*", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "通配符路径 与 路径参数 不可以同时存在", func() {
		newrouter_.AddRouter(http.MethodGet, "/:id", function_)
		newrouter_.AddRouter(http.MethodGet, "/*", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "不允许3", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "通配符路径 与 正则路径不可以同时存在", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/*", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "不允许", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id", function_)
	})
	// 参数冲突
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "不允许 在同一个节点注册多个 参数节点", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/c/:id", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/c/:name", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "正则出现错误,没有找到(", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id.*)", function_)
		//newrouter_.AddRouter(http.MethodGet, "/a/b/c/:name", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "正则 不允许多次注册", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.*)", function_)
		//newrouter_.AddRouter(http.MethodGet, "/a/b/c/:name", function_)
	})
	newrouter_ = NewRouter()
	assert.PanicsWithValue(t, "正则出现错误,compile失败", func() {
		newrouter_.AddRouter(http.MethodGet, "/a/b/:id(.sdfkkefi****)", function_)
	})

}
func (router1 *router) equal(router2 *router) (string, bool) {
	for index, values := range router1.trees {
		tree, ok := router2.trees[index]
		if !ok {
			return fmt.Sprintf("说明连 methods tree 都没有找到"), false
		}
		//	把找到的tree 传到下一层进行比较
		msg, ok := values.equal(tree)
		if !ok {
			return msg, false
		}
	}
	return "", true
}
func (node1 *node) equal(node2 *node) (string, bool) {
	//比较两个node--------
	if node1.path != node2.path {
		return fmt.Sprintf("两个节点的 path 不相同"), false
	}
	if len(node1.children) != len(node2.children) {
		return fmt.Sprintf("两个节点的 children 个数不相同"), false
	}
	node1_function := reflect.ValueOf(node1.function)
	node2_function := reflect.ValueOf(node2.function)
	if node1_function != node2_function {
		return fmt.Sprintf("两个节点的 function 不相同"), false
	}
	for index, values := range node1.children {
		values2 := node2.children[index]
		msg, ok := values.equal(values2)
		if !ok {
			return msg, false
		}
	}
	if node1.regchild != nil {
		node1reg := node1.regchild.regexp.String()
		node2reg := node2.regchild.regexp.String()
		//fmt.Println(node1.regchild.regexp, node2.regchild.regexp, "zz")
		if node1reg != node2reg {
			return "正则表达式不同--", false
		}
		msg, ok := node1.regchild.equal(node2.regchild)
		if !ok {
			return msg, false
		}
	}
	if node1.paramchild != nil {
		msg, ok := node1.paramchild.equal(node2.paramchild)
		if !ok {
			return msg, false
		}
	}
	if node1.starchild != nil {
		msg, ok := node1.starchild.equal(node2.starchild)
		if !ok {
			return msg, false
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
		//{
		//	method: http.MethodGet,
		//	path:   "/param/delete",
		//},
		// 正则
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:id([0-9]+)/home",
		},
	}

	mockHandler := func(ctx *context) {}

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
					function: mockHandler,
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
					function: mockHandler,
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
					function: mockHandler,
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
					function: mockHandler,
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
					function: mockHandler,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     "*",
					function: mockHandler,
				},
			},
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
					path:     ":id",
					function: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
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
					function: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
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
					function: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /reg/:id(.*)
			name:   ":id(.*)",
			method: http.MethodDelete,
			path:   "/reg/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     ":id(.*)",
					function: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /:id([0-9]+)/home
			name:   ":id([0-9]+)",
			method: http.MethodDelete,
			path:   "/123/home",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path:     ":id(.*)",
					function: mockHandler,
				},
				pathParams: map[string]string{"id": "123"},
			},
		},
		{
			// 未命中 /:id([0-9]+)/home
			name:   "not :id([0-9]+)",
			method: http.MethodDelete,
			path:   "/abc/home",
		},
	}

	r := NewRouter()
	for _, tr := range testRoutes {
		r.AddRouter(tr.method, tr.path, mockHandler)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.Findrouter(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.pathParams, mi.pathParams)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.function)
			nVal := reflect.ValueOf(n.function)
			assert.Equal(t, wantVal, nVal)
		})
	}
}

func BenchmarkRouter_Findrouter(b *testing.B) {
	testRouter := []struct {
		methods string
		path    string
	}{
		// 静态节点测试
		{
			methods: http.MethodGet,
			path:    "/",
		},
		{
			methods: http.MethodGet,
			path:    "/user",
		},
		{
			methods: http.MethodGet,
			path:    "/user/home",
		},
		{
			methods: http.MethodGet,
			path:    "/order/detail",
		},
		{
			methods: http.MethodGet,
			path:    "/order/create",
		},
		{
			methods: http.MethodPost,
			path:    "/login",
		},
		//	通配符测试用例
		{
			methods: http.MethodGet,
			path:    "/order/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*/*",
		},
		{
			methods: http.MethodGet,
			path:    "/*/abc",
		},
		{
			methods: http.MethodGet,
			path:    "/*/abc/*",
		},
		// 参数路由
		{
			methods: http.MethodGet,
			path:    "/param/:id",
		},
		{
			methods: http.MethodGet,
			path:    "/param/:id/detail",
		},
		{
			methods: http.MethodGet,
			path:    "/param/:id/*",
		},
		{
			methods: http.MethodDelete,
			path:    "/reg/:id(.*)",
		},
		{
			methods: http.MethodDelete,
			path:    "/:name(^.+$)/abc",
		},
	}

	router_ := NewRouter()
	var function_ handlerfunc = func(ctx *context) {
		fmt.Println("uryyb")
	}
	for _, values := range testRouter {
		router_.AddRouter(values.methods, values.path, function_)
	}

	for i := 0; i < b.N; i++ {
		router_.Findrouter(http.MethodGet, "/param/ww")
	}
}

//API server listening at: 127.0.0.1:49282
//goos: windows
//goarch: amd64
//pkg: geektime_go_go/go_/Exercise/geektime/homework1/version2
//cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
//BenchmarkRouter_Findrouter
//[user] /user
//[user] /user
//[user] /user
//[user] /user
//[user] /user
//BenchmarkRouter_Findrouter-8     3261686               310.6 ns/op
//PASS
