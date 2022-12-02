package version2

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	testcase := []struct {
		name string

		build QueryBuild

		query *Query
		error error
	}{
		{
			name:  "no from ",
			build: &Selector[TestModel]{},
			query: &Query{
				Sql:  "SELECT * FROM `TestMode`;",
				args: nil,
			},
		},
		{
			name:  " from ",
			build: (&Selector[TestModel]{}).From("`temp_table`"),
			query: &Query{
				Sql:  "SELECT * FROM `temp_table`;",
				args: nil,
			},
		},
		{
			name:  " empty from ",
			build: (&Selector[TestModel]{}).From("`db`.`temp_data`"),
			query: &Query{
				Sql:  "SELECT * FROM `db`.`temp_data`;",
				args: nil,
			},
		},
		{
			name:  " where ",
			build: (&Selector[TestModel]{}).Where(NewColumn("Age").Eq(18)),
			query: &Query{
				Sql:  "SELECT * FROM `TestMode` WHERE `Age` = ?;",
				args: []any{18},
			},
		},
		{
			name:  " not ",
			build: (&Selector[TestModel]{}).Where(Not(NewColumn("Age").Eq(18))),
			query: &Query{
				Sql:  "SELECT * FROM `TestMode` WHERE  NOT (`Age` = ?);",
				args: []any{18},
			},
		},
		{
			name:  " and ",
			build: (&Selector[TestModel]{}).Where(NewColumn("Age").Eq(18).And(NewColumn("FirstName").Eq("zh"))),
			query: &Query{
				Sql:  "SELECT * FROM `TestMode` WHERE (`Age` = ?) AND (`FirstName` = ?);",
				args: []any{18, "zh"},
			},
		},
		{
			name:  " or ",
			build: (&Selector[TestModel]{}).Where(NewColumn("Age").Eq(18).Or(NewColumn("FirstName").Eq("zh"))),
			query: &Query{
				Sql:  "SELECT * FROM `TestMode` WHERE (`Age` = ?) OR (`FirstName` = ?);",
				args: []any{18, "zh"},
			},
		},
		//{
		//	name:  " or ",
		//	build: (&Selector[TestMode]{}).Where(NewColumn("Age").Eq(18).Or(NewColumn("ccc").Eq("zh"))),
		//	query: &Query{
		//		Sql: "SELECT * FROM `TestMode` WHERE (`Age` = ?) OR (`first_name` = ?);",
		//	},
		//	error: errors.New("未知字段"),
		//},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			new_query, err := tc.build.Build()
			assert.Equal(t, tc.error, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.query, new_query)
		})
	}
}
