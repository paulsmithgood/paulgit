package model

import (
	"database/sql"
	"exercise/geektime/homework4/version1/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type PaulsModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  sql.NullString
}
type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

//func TestRegister_Get(t *testing.T) {
//	r := Newregistry()
//	u := &PaulsModel{}
//	model, err := r.Get(u)
//	fmt.Println(model, err)
//}

func TestRegistry_get(t *testing.T) {
	testCases := []struct {
		name      string
		val       any
		wantModel *Model
		wantErr   error
	}{
		{
			name: "test Model",
			val:  TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"Age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"LastName": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"last_name": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
			},
		},
		{
			// 指针
			name: "pointer",
			val:  &TestModel{},
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"Age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"LastName": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"last_name": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
			},
		},
		{
			// 多级指针
			name: "multiple pointer",
			// 因为 Go 编译器的原因，所以我们写成这样
			val: func() any {
				val := &TestModel{}
				return &val
			}(),
			wantModel: &Model{
				TableName: "test_model",
				FieldMap: map[string]*Field{
					"Id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"Age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"LastName": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type:    reflect.TypeOf(int64(0)),
						GoName:  "Id",
						Offset:  0,
					},
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
						Offset:  8,
					},
					"age": {
						ColName: "age",
						Type:    reflect.TypeOf(int8(0)),
						GoName:  "Age",
						Offset:  24,
					},
					"last_name": {
						ColName: "last_name",
						Type:    reflect.TypeOf(&sql.NullString{}),
						GoName:  "LastName",
						Offset:  32,
					},
				},
			},
		},
		{
			name:    "map",
			val:     map[string]string{},
			wantErr: errs.ErrParseType,
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errs.ErrParseType,
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errs.ErrParseType,
		},

		// 标签相关测试用例
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `orm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			wantModel: &Model{
				TableName: "column_tag",
				FieldMap: map[string]*Field{
					"ID": {
						ColName: "id",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "ID",
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "ID",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是传入一个空字符串，那么会用默认的名字
			name: "empty column",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type EmptyColumn struct {
					FirstName string `orm:"column="`
				}
				return &EmptyColumn{}
			}(),
			wantModel: &Model{
				TableName: "empty_column",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是没有赋值
			name: "invalid tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type InvalidTag struct {
					FirstName string `orm:"column"`
				}
				return &InvalidTag{}
			}(),
			wantErr: errs.ErrInvalidTag,
		},
		{
			// 如果用户设置了一些奇奇怪怪的内容，这部分内容我们会忽略掉
			name: "ignore tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type IgnoreTag struct {
					FirstName string `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &Model{
				TableName: "ignore_tag",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(""),
						GoName:  "FirstName",
					},
				},
			},
		},
	}

	r := &Register{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		// 我们这些用例就是为了确保
		// 在忘记 underscoreName 的行为特性之后
		// 可以从这里找回来
		// 比如说过了一段时间之后
		// 忘记了 ID 不能转化为 id
		// 那么这个测试能帮我们确定 ID 只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}
