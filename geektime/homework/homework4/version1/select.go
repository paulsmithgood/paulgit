package version1

import (
	"exercise/geektime/homework4/version1/internal/errs"
	"exercise/geektime/homework4/version1/model"
	"strings"
)

var _ QueryBuild = &Selector[any]{}

type Seleable interface {
	selectable() // 这是一个标记接口
}

type Selector[T any] struct {
	sb      strings.Builder
	Table   string //通过 from 接口来的 ，如果是 空 那就代表使用的是 默认的表名
	where   []Predicate
	Args    []any
	db      *DB
	Molel   *model.Model
	columns []Seleable
	groupby []Columns
	having  []Predicate
	orderby []Orderby
	offset  int
	limit   int
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: strings.Builder{},
		db: db,
	}
}

func (s *Selector[T]) Select(selector ...Seleable) *Selector[T] {
	s.columns = selector
	return s
}
func (s *Selector[T]) From(name string) *Selector[T] {
	s.Table = name
	return s
}
func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	s.where = predicates
	return s
}
func (s *Selector[T]) GroupBy(columns ...Columns) *Selector[T] {
	s.groupby = columns
	return s
}
func (s *Selector[T]) Having(predicates ...Predicate) *Selector[T] {
	s.having = predicates
	return s
}
func (s *Selector[T]) OrderBy(orderbys ...Orderby) *Selector[T] {
	s.orderby = orderbys
	return s
}
func (s *Selector[T]) Offset(offset int) *Selector[T] {
	s.offset = offset
	return s
}
func (s *Selector[T]) Limit(limit int) *Selector[T] {
	s.limit = limit
	return s
}
func (s *Selector[T]) Build() (*Query, error) {
	var err error
	s.Molel, err = s.db.r.Get(*new(T))
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT ")
	if len(s.columns) > 0 {
		//在这里就说明有东西传进来了
		for index, value := range s.columns {
			if index > 0 {
				s.sb.WriteString(",")
			}
			switch typr := value.(type) {
			case Columns:
				//	是一个指定列
				s.sb.WriteByte('`')
				fd, ok := s.Molel.FieldMap[typr.name]
				if !ok {
					return nil, errs.ErrUnknowField
				}
				s.sb.WriteString(fd.ColName)
				s.sb.WriteByte('`')
				//	处理一下别名
				if typr.alias != "" {
					//	就 说明我们有 传别名的内容
					s.sb.WriteString(" AS ")
					s.sb.WriteByte('`')
					s.sb.WriteString(typr.alias)
					s.sb.WriteByte('`')
				}
			case Aggregate:
				s.sb.WriteString(typr.fn.String())
				s.sb.WriteByte('(')
				s.sb.WriteByte('`')
				fd, ok := s.Molel.FieldMap[typr.columns]
				if !ok {
					return nil, errs.ErrUnknowField
				}
				s.sb.WriteString(fd.ColName)
				s.sb.WriteByte('`')
				s.sb.WriteByte(')')

				if typr.alias != "" {
					//	就 说明我们有 传别名的内容
					s.sb.WriteString(" AS ")
					s.sb.WriteByte('`')
					s.sb.WriteString(typr.alias)
					s.sb.WriteByte('`')
				}
			case RawExpression:
				s.sb.WriteString(typr.exprs)
				if typr.alias != "" {
					//	就 说明我们有 传别名的内容
					s.sb.WriteString(" AS ")
					s.sb.WriteByte('`')
					s.sb.WriteString(typr.alias)
					s.sb.WriteByte('`')
				}
				if len(typr.args) > 0 {
					s.Addarg(typr.args...)
				}

			default:

			}

		}
	} else {
		s.sb.WriteString("*")
	}

	s.sb.WriteString(" FROM ")

	//用于 构造sql语句的 build 方法
	if s.Table == "" {
		//	需要去使用默认的字段
		s.sb.WriteByte('`')
		s.sb.WriteString(s.Molel.TableName)
		s.sb.WriteByte('`')
		//fmt.Println(reflect.TypeOf(new(T)).Name())
	} else {
		//s.sb.WriteByte('`')
		s.sb.WriteString(s.Table)
		//s.sb.WriteByte('`')
	}

	//首先先判断一下 里面是不是有where 条件

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		//	在这里需要先把所有的 predicate list 都先拼接成一个 整体
		new_left := s.where[0]
		for i := 1; i < len(s.where); i++ {
			new_left = new_left.And(s.where[i])
		}
		predicate_where := new_left //这里是一个 整体

		errx := s.buildexpression(predicate_where)
		if errx != nil {
			return nil, errx
		}
	}

	if len(s.groupby) > 0 {
		s.sb.WriteString(" GROUP BY ")
		for index, val := range s.groupby {
			if index > 0 {
				s.sb.WriteByte(',')
			}
			fd, ok := s.Molel.FieldMap[val.name]
			if !ok {
				return nil, errs.ErrUnknowField
			}
			s.sb.WriteByte('`')
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte('`')
		}
	}

	if len(s.having) > 0 {
		s.sb.WriteString(" HAVING ")
		new_leftX := s.having[0]
		for i := 1; i < len(s.having); i++ {
			new_leftX = new_leftX.And(s.having[i])
		}
		predicate_having := new_leftX //这里是一个 整体

		errx := s.buildexpression(predicate_having)
		if errx != nil {
			return nil, errx
		}
	}

	if len(s.orderby) > 0 {
		s.sb.WriteString(" ORDER BY ")
		for index, val := range s.orderby {
			if index > 0 {
				s.sb.WriteByte(',')
			}
			fd, ok := s.Molel.FieldMap[val.columns]
			if !ok {
				return nil, errs.ErrUnknowField
			}
			s.sb.WriteByte('`')
			s.sb.WriteString(fd.ColName)
			s.sb.WriteByte('`')
			s.sb.WriteByte(byte(32))
			s.sb.WriteString(val.direction)
		}
	}

	if s.limit > 0 {
		s.sb.WriteString(" LIMIT ?")
		s.Addarg(s.limit)
	}

	if s.offset > 0 {
		s.sb.WriteString(" OFFSET ?")
		s.Addarg(s.offset)
	}

	s.sb.WriteString(";")
	return &Query{SQL: s.sb.String(), Args: s.Args}, nil
}

func (s *Selector[T]) buildexpression(predicate_where Expression) error {
	//	传入我们的  predicate_where 并且解析出来
	switch typr := predicate_where.(type) {
	case Predicate:
		_, ok := typr.left.(Predicate)
		if ok {

		}
		err := s.buildexpression(typr.left)
		if err != nil {
			return err
		}
		if typr.op == "" {
			//可能是 没有 op的，比如 raw
			return nil
		}
		s.sb.WriteByte(byte(32))
		s.sb.WriteString(typr.op.String())
		s.sb.WriteByte(byte(32))

		_, ok = typr.right.(Predicate)
		if ok {

		}
		err = s.buildexpression(typr.right)
		if err != nil {
			return err
		}
	case Aggregate:
		s.sb.WriteString(typr.fn.String())
		s.sb.WriteString("(`")
		fd, ok := s.Molel.FieldMap[typr.columns]
		if !ok {
			return errs.ErrUnknowField
		}
		s.sb.WriteString(fd.ColName)
		s.sb.WriteString("`)")

	case Columns:
		s.sb.WriteByte('`')
		fd, ok := s.Molel.FieldMap[typr.name]
		if !ok {
			return errs.ErrUnknowField
		}
		s.sb.WriteString(fd.ColName)
		s.sb.WriteByte('`')
	case Values:
		s.sb.WriteString("?")
		s.Addarg(typr.val)
	case RawExpression:
		s.sb.WriteString(typr.exprs)
		if len(typr.args) > 0 {
			s.Addarg(typr.args...)
		}
	}
	return nil
}

func (s *Selector[T]) Addarg(val ...any) {
	if s.Args == nil {
		s.Args = make([]any, 0, 8)
	}
	s.Args = append(s.Args, val...)
}
