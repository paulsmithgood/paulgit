package version1

//自定义模式

type RawExpression struct {
	exprs string
	args  []any
	alias string
}

func (r RawExpression) expr() {

}

func (r RawExpression) selectable() {

}

func NewRaw(expr string, args ...any) RawExpression {
	return RawExpression{
		exprs: expr,
		args:  args,
	}
}

func (r RawExpression) As(alias string) RawExpression {
	r.alias = alias
	return r
}

func (r RawExpression) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}
