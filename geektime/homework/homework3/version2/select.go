package version2

import (
	"fmt"
	"reflect"
	"strings"
)

var _ QueryBuild = &Selector[interface{}]{}

type Selector[T any] struct {
	wb    *strings.Builder
	table string
	where []predicate
	args  []any
}

func (s *Selector[T]) From(table string) *Selector[T] {
	//	只要调用这个接口，就会把 table 赋予给这个值
	s.table = table
	return s
}

func (s *Selector[T]) Where(predicates ...predicate) *Selector[T] {
	s.where = predicates
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	fmt.Println(s)
	s.wb = &strings.Builder{}
	s.wb.WriteString("SELECT * FROM ")
	if s.table == "" {
		var t T //通过反射拿到T 结构体的名字
		s.wb.WriteByte('`')
		s.wb.WriteString(reflect.TypeOf(t).Name())
		s.wb.WriteByte('`')
	} else {
		s.wb.WriteString(s.table)
	}
	//到这里，from 已经构造好了，开始构造 where

	if len(s.where) > 0 {
		s.wb.WriteString(" WHERE ") //说明有where
		//可能存在不定参数，通过循环先把他们都连起来
		left := s.where[0]
		for i := 1; i < len(s.where); i++ {
			left = left.And(s.where[i])
		}

		Expressions := left
		//当前的  expressions 已经是个完整的结构，通过递归
		fmt.Println(Expressions)
		err := s.ReadExpressions(Expressions)
		if err != nil {
			return nil, err
		}
	}

	s.wb.WriteByte(';')
	fmt.Println(s.args, "s.args")
	return &Query{Sql: s.wb.String(), args: s.args}, nil
}
func (s *Selector[T]) ReadExpressions(expression Expression) error {
	switch typ := expression.(type) {
	case predicate:
		_, ok := typ.left.(predicate)
		if ok {
			s.wb.WriteByte('(')
		}
		if err := s.ReadExpressions(typ.left); err != nil {
			return err
		}
		if ok {
			s.wb.WriteByte(')')
		}

		s.wb.WriteByte(byte(32))
		s.wb.WriteString(typ.op)
		s.wb.WriteByte(byte(32))

		_, ok = typ.right.(predicate)
		if ok {
			s.wb.WriteByte('(')
		}
		if err := s.ReadExpressions(typ.right); err != nil {
			return err
		}
		if ok {
			s.wb.WriteByte(')')
		}
	case Values:
		s.wb.WriteByte('?')
		s.addarg(typ.val)
	case Column:
		s.wb.WriteByte('`')
		s.wb.WriteString(typ.column)
		s.wb.WriteByte('`')
	}
	return nil
}

func (s *Selector[T]) addarg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
}
