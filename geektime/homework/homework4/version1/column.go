package version1

var _ Expression = &Columns{}

type Columns struct {
	name  string
	alias string
}

func (c Columns) expr() {

}
func (c Columns) selectable() {

}

func NewColumns(name string) Columns {
	return Columns{
		name: name,
	}
}

func (c Columns) As(name string) Columns {
	//	对外的别名接口
	c.alias = name
	return c
}
func (left Columns) Eq(right any) Predicate {
	return Predicate{
		left:  left,
		op:    Eqop,
		right: Values{val: right},
	}
}
func (left Columns) GT(right any) Predicate {
	return Predicate{
		left:  left,
		op:    GTop,
		right: Values{val: right},
	}
}
func (left Columns) LT(right any) Predicate {
	return Predicate{
		left:  left,
		op:    LTop,
		right: Values{val: right},
	}
}
