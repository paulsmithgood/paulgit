package version6

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

func Newrouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

type node struct {
	path string

	starchildren *node

	paramchildren *node

	children map[string]*node

	middlware []Middleware

	funtion_ Handlerfunc

	route string //完整的路由
}

type matchInfo struct {
	n *node

	param_map map[string]string

	mdls []Middleware
}

func Newnode(path string) *node {
	return &node{
		path: path,
	}
}
func (r *router) addrouter(methods, path string, handler Handlerfunc, middleware []Middleware) {
	if path == "" {
		panic("路由不可以为空!")
	}
	//	首先要先判断 这个trees 是否存在
	tree, ok := r.trees[methods]
	if !ok {
		//	说明tree 不存在
		root := Newnode("/")
		r.trees[methods] = root
		tree = root
	}
	if path == "/" {
		if tree.funtion_ != nil {
			panic("根节点不允许重复注册")
		}
		//	path 如果是 /的要进行特殊处理
		tree.funtion_ = handler
		tree.middlware = middleware
		return
	}
	if path[0] != '/' {
		panic("路由必须是以/为前导")
	}
	if path[len(path)-1] == '/' {
		panic("路由不能以/结束")
	}
	cur := tree
	path_list := strings.Split(strings.Trim(path, "/"), "/")
	for _, value := range path_list {
		if value == "" {
			panic("path 带多//,不允许")
		}
		cur = cur.findrouter_createrouter(value)
	}
	if cur.funtion_ != nil {
		panic("该节点不允许多次注册！")
	}
	cur.funtion_ = handler
	cur.middlware = middleware
	cur.route = path
}
func (n *node) findrouter_createrouter(path string) *node {
	if path == "*" {
		//	判断路径是否是*路径，通配符
		if n.paramchildren != nil {
			panic("不允许 * 和 : 并存")
		}
		if n.starchildren == nil {
			n.starchildren = Newnode("*")
		}
		return n.starchildren
	}
	if path[0] == ':' {
		if n.starchildren != nil {
			panic("不允许 有 * 和 : 并存")
		}
		if n.paramchildren != nil {
			if n.paramchildren.path != path {
				panic("不允许注册 有多个 :")
			}
		} else {
			n.paramchildren = Newnode(path)
		}
		return n.paramchildren
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = Newnode(path)
		n.children[path] = child
	}
	return child
}

func (r *router) findrouter(methods, path string) (*matchInfo, bool) {
	root, ok := r.trees[methods]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{n: root, mdls: root.middlware}, true
	}
	//当前的节点是root
	mi := &matchInfo{}
	cur := root
	path_string := strings.Split(strings.Trim(path, "/"), "/")
	for _, values := range path_string {
		var mathbool bool
		cur, mathbool, ok = cur.findchild(values)
		if !ok {
			//没找到
			return nil, false
		}
		if mathbool {
			//判断一下如果param_map 是 空的，就初始化一个
			mi.addValue(cur.path[1:], values)
		}
	}
	mi.n = cur
	//	到这里节点已经找到了====
	//	接下来要找这个符合这个节点的middleware---
	var middleware_list []Middleware
	if root.middlware != nil {
		//	根节点下面是否有middleware--有的话就先把他放进来
		middleware_list = append(middleware_list, root.middlware...)
	}
	//在root加点下 做匹配了
	middleware_list = append(middleware_list, root.findMiddleware(path_string)...)
	mi.mdls = middleware_list
	return mi, true
}

func (n *node) findchild(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramchildren != nil {
			return n.paramchildren, true, true
		}
		return n.starchildren, false, n.starchildren != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.paramchildren != nil {
			return n.paramchildren, true, true
		}
		return n.starchildren, false, n.starchildren != nil
	}
	return res, false, ok
}
func (n *node) findMiddleware(path_string []string) []Middleware {
	//	把已经分切好的path-string 传进来，然后再 挨个判断
	var middleware_list []Middleware
	var root_list []*node
	var root_list_temp []*node
	root_list = append(root_list, n)
	for _, values := range path_string {
		root_list_temp = []*node{}
		for _, root2 := range root_list {
			root_temp, middleware_list_root := root2.findmiddlewares(values)
			root_list_temp = append(root_list_temp, root_temp...)
			middleware_list = append(middleware_list, middleware_list_root...)
		}
		root_list = root_list_temp
	}

	return middleware_list
	//panic("implement me")
}

func (n *node) findmiddlewares(path string) ([]*node, []Middleware) {
	var middleware_list []Middleware
	var node_list []*node
	//层级关系，*<通配符<静态
	if n.starchildren != nil {
		node_list = append(node_list, n.starchildren)
		middleware_list = append(middleware_list, n.starchildren.middlware...)
	}
	if n.paramchildren != nil {
		node_list = append(node_list, n.paramchildren)
		middleware_list = append(middleware_list, n.paramchildren.middlware...)
	}
	Static, ok := n.children[path]
	if ok {
		//	找到了
		node_list = append(node_list, Static)
		middleware_list = append(middleware_list, Static.middlware...)
	}
	return node_list, middleware_list
}
func (m *matchInfo) addValue(key string, value string) {
	if m.param_map == nil {
		// 大多数情况，参数路径只会有一段
		m.param_map = map[string]string{key: value}
	}
	m.param_map[key] = value
}
