package version1

var _ Expression = &Predicate{}
var _ Expression = &Values{}

type op string

const (
	Eqop  op = "="
	Andop op = "AND"
	GTop  op = ">"
	LTop  op = "<"
	NOTop op = "NOT"
	ORTop op = "OR"
)

func (o op) String() string {
	return string(o)
}

type Expression interface {
	expr()
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {

}

func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    Andop,
		right: right,
	}
}

func Not(right Predicate) Predicate {
	return Predicate{
		op:    NOTop,
		right: right,
	}
}
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    ORTop,
		right: right,
	}
}

type Values struct {
	val any
}

func (v Values) expr() {

}
