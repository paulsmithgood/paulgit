package version2

var _ Expression = &predicate{}
var _ Expression = &Column{}

type predicate struct {
	left  Expression
	op    string
	right Expression
}

// 需要让predicate 实现expression的接口
type Expression interface {
	expr()
}

func (p predicate) expr() {

}

func (left predicate) Or(right predicate) predicate {
	return predicate{
		left:  left,
		op:    "OR",
		right: right,
	}
}
func (left predicate) And(right predicate) predicate {
	return predicate{
		left:  left,
		op:    "AND",
		right: right,
	}
}
func Not(predicates predicate) predicate {
	return predicate{
		op:    "NOT",
		right: predicates,
	}
}

type Values struct {
	val any
}

func (v Values) expr() {

}
