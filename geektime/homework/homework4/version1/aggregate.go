package version1

type fn string

const (
	fnsum   fn = "SUM"
	fncount fn = "COUNT"
	fnavg   fn = "AVG"
)

func (fn fn) String() string {
	return string(fn)
}

// select sum(age) as ye ,count(age)
// 在这里需要考虑 三种
type Aggregate struct {
	fn      fn     //方法名
	columns string //字段名
	alias   string //别名
}

func (a Aggregate) expr() {

}

func (a Aggregate) selectable() {}

func Sum(col string) Aggregate {
	return Aggregate{
		fn:      fnsum,
		columns: col,
	}
}
func Avg(col string) Aggregate {
	return Aggregate{
		fn:      fnavg,
		columns: col,
	}
}

func Count(col string) Aggregate {
	return Aggregate{
		fn:      fncount,
		columns: col,
	}
}

func (a Aggregate) As(alias string) Aggregate {
	a.alias = alias
	return a
}

func (left Aggregate) EQ(right any) Predicate {
	return Predicate{
		left:  left,
		op:    Eqop,
		right: exprOf(right),
	}
}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return Values{val: e}
	}
}
