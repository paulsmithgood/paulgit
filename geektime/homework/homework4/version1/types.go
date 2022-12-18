package version1

// 这个query 我们用于定义 build
type Query struct {
	SQL  string
	Args []any
}

type QueryBuild interface {
	// 用于结构体的级联调用，并且 实现 返回 query 结构体
	Build() (*Query, error)
}
