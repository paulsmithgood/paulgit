package version2

import (
	"reflect"
	"strings"
)

//实现一个 delete结构体

var _ QueryBuild = &Deleter[interface{}]{}

type Deleter[T any] struct {
	table string
	where []predicate
	wb    *strings.Builder
	args  []any
}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}
func (d *Deleter[T]) Where(predicates ...predicate) *Deleter[T] {
	d.where = predicates
	return d
}

func (d *Deleter[T]) Build() (*Query, error) {
	d.wb = &strings.Builder{}
	d.wb.WriteString("DELETE FROM ")
	if d.table == "" {
		//	reflect T
		var t T
		typ := reflect.TypeOf(t)
		d.wb.WriteByte('`')
		d.wb.WriteString(typ.Name())
		d.wb.WriteByte('`')
	} else {
		d.wb.WriteString(d.table)
	}

	if len(d.where) > 0 {
		//	说明有where
		d.wb.WriteString(" WHERE ")
		left := d.where[0]
		for i := 1; i < len(d.where); i++ {
			left = left.And(d.where[i])
		}
		Expressions := left

		err := d.ReadExpressions(Expressions)
		if err != nil {
			return nil, err
		}
	}

	d.wb.WriteByte(';')
	return &Query{Sql: d.wb.String(), args: d.args}, nil
}

func (d *Deleter[T]) ReadExpressions(expression Expression) error {
	switch typ := expression.(type) {
	case predicate:
		_, ok := typ.left.(predicate)
		if ok {
			d.wb.WriteByte('(')
		}
		err := d.ReadExpressions(typ.left)
		if err != nil {
			return err
		}
		d.wb.WriteByte(byte(32))
		d.wb.WriteString(typ.op)
		d.wb.WriteByte(byte(32))
		_, ok = typ.right.(predicate)
		if ok {
			d.wb.WriteByte(')')
		}
		err = d.ReadExpressions(typ.right)
		if err != nil {
			return err
		}
	case Column:
		d.wb.WriteByte('`')
		d.wb.WriteString(typ.column)
		d.wb.WriteByte('`')
	case Values:
		d.wb.WriteByte('?')
		d.addarg(typ.val)
	}
	return nil
}
func (d *Deleter[T]) addarg(val any) {
	if d.args == nil {
		d.args = make([]any, 0, 4)
	}
	d.args = append(d.args, val)
}
