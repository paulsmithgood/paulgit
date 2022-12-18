package model

import (
	"exercise/geektime/homework4/version1/internal/errs"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Model struct {
	//结构体名字
	TableName string
	//列名到字段定义映射
	FieldMap map[string]*Field
	//字段名到字段定义映射
	ColumnMap map[string]*Field
}

type Field struct {
	ColName string
	GoName  string
	Type    reflect.Type //字段类型
	Offset  uintptr      //偏移量
}

type TableName interface {
	TableName() string
}

type Option func(m *Model) error

type Regist interface {
	Get(val any) (*Model, error)
	Register(val any, options ...Option) (*Model, error)
}

type Register struct {
	Maps sync.Map
}

func Newregistry() Regist {
	return &Register{}
}

func (r *Register) Get(val any) (*Model, error) {
	fmt.Println("1", val)
	typ := reflect.TypeOf(val)
	m, ok := r.Maps.Load(typ)
	if ok {
		//到这里说明我们可以在 register maps 中找到
		return m.(*Model), nil
	}

	return r.Register(val)
}

func (r *Register) Register(val any, options ...Option) (*Model, error) {
	m, err := r.Parsemodel(val)
	if err != nil {
		return nil, err
	}
	for _, opt := range options {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}
	typ := reflect.TypeOf(val)
	r.Maps.Store(typ, m)
	return m, nil
}

func (r *Register) Parsemodel(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errs.ErrParseType
	}
	num := typ.NumField()

	//在这里 我们先初始化好 我们需要的 map
	FieldMap := make(map[string]*Field, num)
	ColumnMap := make(map[string]*Field, num)

	for i := 0; i < num; i++ {
		fd := typ.Field(i)
		tagmaps, err := r.ParseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		//	需要在 这个tagmaps 找对应的的 列 是我们设置的列
		columnname, _ := tagmaps["column"]
		if columnname == "" {
			columnname = underscoreName(fd.Name)
		}
		field := &Field{
			ColName: columnname, //到这里我们已经通过我们的tag 解析出来我们 需要的columnsname
			GoName:  fd.Name,
			Type:    fd.Type,
			Offset:  fd.Offset,
		} //到这里我们已经解析出来了我们需要的 field
		FieldMap[fd.Name] = field
		ColumnMap[columnname] = field
	}

	return &Model{
		TableName: underscoreName(typ.Name()),
		FieldMap:  FieldMap,
		ColumnMap: ColumnMap,
	}, nil
}

func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

func (r *Register) ParseTag(tag reflect.StructTag) (map[string]string, error) {
	//m := make(map[string]string, 0)
	tags := tag.Get("orm")
	fmt.Println(tags, "tags")
	if tags == "" {
		return map[string]string{}, nil
	}
	m := make(map[string]string, 1)
	//fmt.Println(tags)
	strings_list := strings.Split(tags, ",")

	for _, values := range strings_list {
		values_strings := strings.Split(values, "=")
		if len(values_strings) != 2 {
			return nil, errs.ErrInvalidTag
		}
		m[values_strings[0]] = values_strings[1]
	}
	return m, nil
}
