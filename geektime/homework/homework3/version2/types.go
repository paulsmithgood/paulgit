package version2

//创建一个返回 构造结果的结构体

type Query struct {
	Sql  string
	args []any
}

type QueryBuild interface {
	Build() (*Query, error) //用于结构体的级联调用，并且 实现 返回 query 结构体
}
