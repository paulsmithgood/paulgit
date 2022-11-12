package version3

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	trees map[string]*node
}

func NewRouter() *router {
	return &router{trees: map[string]*node{}}
}

type node struct {
	path string

	children map[string]*node

	regchild *node
	regexp   *regexp.Regexp

	paramchild *node

	starchild *node

	function handlerfunc
}

func newnode(path string) *node {
	return &node{path: path}
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

//注册路由
func (r *router) AddRouter(methods, path string, function handlerfunc) {
	if path == "" {
		panic("路由不可以为空!")
	}
	//	用于注册路由
	//	首先先判定树 存不存在
	root, ok := r.trees[methods]
	if !ok {
		//	说明树不存在，需要添加一个新的树
		root = &node{path: "/"}
		r.trees[methods] = root
	}
	//直接给根节点附上 function,返回 不参与后面的逻辑
	if path == "/" {
		if root.function != nil {
			panic("根节点不允许重复注册")
		}
		root.function = function
		return
	}
	//	当前的节点是root
	// 异常问题处理：
	//	路由必须以/为前导
	if path[0] != '/' {
		panic("路由必须是以/为前导")
	}
	if path[len(path)-1] == '/' {
		panic("路由不能以/结束")
	}
	//	当前的节点是root
	path_list := strings.Split(strings.Trim(path, "/"), "/")
	fmt.Println(path_list, path)
	for _, values := range path_list {
		if values == "" {
			panic("路由不能存在连续两个//")
		}
		//child 是本节点返回的 下一个节点
		child := root.findrouter_createrouter(values)
		root = child

	}
	//把方法赋值给 最后的节点
	if root.function != nil {
		panic("路由不能重复注册")
	}
	root.function = function
}

//用于在本节点去寻找符合path 的 另外的节点，并且返回
func (n *node) findrouter_createrouter(path string) *node {
	//	因为路径参数，通配符路径，正则路径 三者是互斥的
	if path == "*" {
		//	因为path=*
		if n.paramchild != nil {
			panic("通配符路径 与 路径参数 不可以同时存在")
		}
		if n.regchild != nil {
			panic("通配符路径 与 正则路径不可以同时存在")
		}
		if n.starchild != nil {
			return n.starchild
		}
		newnode_star := newnode(path)
		n.starchild = newnode_star
		return newnode_star
	}
	if path[0] == ':' {
		if n.starchild != nil {
			panic("不允许1")
		}

		//	因为 正则表达式 参数路径都有这个: /user/:id /user/:id(^[0-9]+$)
		check_ := strings.Contains(path[1:], ")")
		if check_ {
			if n.regchild != nil {
				panic("正则 不允许多次注册")
			}
			if n.paramchild != nil {
				panic("不允许3")
			}
			index := strings.Index(path, "(")
			if index == -1 {
				panic("正则出现错误,没有找到(")
			}
			//	说明他是正则，因为他最后一个是 )
			pathname := path[1:index]
			reg := path[index+1 : len(path)-1]
			regexp_, err := regexp.Compile(reg)
			if err != nil {
				panic("正则出现错误,compile失败")
			}
			newnode_reg := newnode(pathname)
			newnode_reg.regexp = regexp_
			n.regchild = newnode_reg
			return newnode_reg
		} else {
			//if n.starchild != nil {
			//	panic("不允许*")
			//}
			if n.regchild != nil {
				panic("不允许")
			}
			//	说明他是参数路径，因为他没有)
			pathname := path[1:len(path)]
			if n.paramchild == nil {
				//因为 paramschild 是空的，所以要创建一个新的
				newnode_param := newnode(pathname)
				n.paramchild = newnode_param
				return newnode_param
			} else {
				//	不为空可能不会存在一个什么样的问题 同名或者不同名
				if n.paramchild.path == pathname {
					return n.paramchild
				} else {
					panic("不允许 在同一个节点注册多个 参数节点")
				}
			}

		}
	}
	if n.children == nil {
		//	说明没有下一层
		n.children = make(map[string]*node, 0)
	}
	child, ok := n.children[path]
	if !ok {
		//	说明不存在，那我们就创建一个
		child = newnode(path)
		n.children[path] = child
	}
	return child
}

//在我们的路由树里面去找到我们要的路由
func (r *router) Findrouter(methods, path string) (*matchInfo, bool) {
	root, ok := r.trees[methods]
	if !ok {
		//连methods *Node 都没有找到
		return nil, false
	}
	//	root 是本节点
	if path == "/" {
		//	因为是根节点，直接return
		return &matchInfo{n: root, pathParams: nil}, true
	}
	var star_node *node
	var paramsmap map[string]string
	path_strings := strings.Split(strings.Trim(path, "/"), "/")
	for _, values := range path_strings {
		child, paramsok, starsok, found := root.findchildrouter(values)
		if !found {
			if star_node != nil {
				break
			}
			return nil, false
		}
		if starsok {
			star_node = child
		}
		if paramsok {
			if paramsmap == nil {
				//对map 进行初始化
				paramsmap = make(map[string]string, 0)
			}
			paramsmap[child.path] = values
		}
		root = child
	}
	//	说明遍历结束

	//if paramsmap!=nil{
	//	return &matchInfo{n: root,pathParams: paramsmap},true
	//}else {
	//	return &matchInfo{n:root,pathParams: nil},true
	//}
	return &matchInfo{n: root, pathParams: paramsmap}, true
}

//child,paramsok,starsok,found
func (n *node) findchildrouter(path string) (*node, bool, bool, bool) {
	//	本节点是n 从本节点开始向下寻找符合要求的节点
	//	首先先判定是否有静态节点
	if n.children != nil {
		//	有静态children ---选
		child, ok := n.children[path]
		if !ok {
			//fmt.Println("没有找到对应map 静态*node 往下找")
			if n.regchild != nil {
				//fmt.Println("找到了正则节点，判断一下是否符合要求")
				if n.regchild.regexp.MatchString(path) {
					return n.regchild, true, false, true
				}
			}
			if n.paramchild != nil {
				//fmt.Println("找到了参数节点")
				return n.paramchild, true, false, true
			}
			if n.starchild != nil {
				//fmt.Println("找到了*节点")
				return n.starchild, false, true, true
			}
		} else {
			return child, false, false, true
		}
	} else {
		//fmt.Println("在 本节点下没有map 的static children")
		if n.regchild != nil {
			//fmt.Println("找到了正则节点，判断一下是否符合要求")
			if n.regchild.regexp.MatchString(path) {
				return n.regchild, true, false, true
			}
		}
		if n.paramchild != nil {
			//fmt.Println("找到了参数节点")
			return n.paramchild, true, false, true
		}
		if n.starchild != nil {
			//fmt.Println("找到了*节点")
			return n.starchild, false, true, true
		}
	}
	return nil, false, false, false
}
