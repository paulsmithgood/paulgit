package version2

type Column struct {
	column string
}

func NewColumn(name string) Column {
	return Column{column: name}
}

func (c Column) expr() {

}

func (left Column) Eq(values any) predicate {
	return predicate{
		left:  left,
		op:    "=",
		right: Values{val: values},
	}
}
