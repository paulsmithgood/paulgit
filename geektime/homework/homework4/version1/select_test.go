//go:build PAULSTEST

package version1

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type PaulsModel struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  sql.NullString
}

func TestSelector_Build(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	testcase := []struct {
		name      string
		selector  *Selector[PaulsModel]
		wantquery *Query
		wanterror error
	}{
		{
			name:     "from diy",
			selector: NewSelector[PaulsModel](db).From("test_model"),
			wantquery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:     "default",
			selector: NewSelector[PaulsModel](db),
			wantquery: &Query{
				SQL: "SELECT * FROM `pauls_model`;",
			},
		},
		{
			name:     "where",
			selector: NewSelector[PaulsModel](db).Where(NewColumns("Id").Eq(8), NewColumns("Age").Eq(18)),
			wantquery: &Query{
				SQL:  "SELECT * FROM `pauls_model` WHERE `id` = ? AND `age` = ?;",
				Args: []any{8, 18},
			},
		},
		{
			name:     "where",
			selector: NewSelector[PaulsModel](db).Where(NewColumns("Id").Eq(8).And(NewColumns("Age").Eq(18))),
			wantquery: &Query{
				SQL:  "SELECT * FROM `pauls_model` WHERE `id` = ? AND `age` = ?;",
				Args: []any{8, 18},
			},
		},
		{
			name:     "specify column",
			selector: NewSelector[PaulsModel](db).Select(NewColumns("Id"), NewColumns("FirstName")),
			wantquery: &Query{
				SQL: "SELECT `id`,`first_name` FROM `pauls_model`;",
			},
		},
		{
			name:     "specify column",
			selector: NewSelector[PaulsModel](db).Select(NewColumns("Id").As("myid"), NewColumns("FirstName")),
			wantquery: &Query{
				SQL: "SELECT `id` AS `myid`,`first_name` FROM `pauls_model`;",
			},
		},
		{
			name:     "fn",
			selector: NewSelector[PaulsModel](db).Select(Sum("Age").As("sumage")),
			wantquery: &Query{
				SQL: "SELECT SUM(`age`) AS `sumage` FROM `pauls_model`;",
			},
		},
		{
			name:     "fn",
			selector: NewSelector[PaulsModel](db).Select(Count("Age")).Where(NewColumns("Age").Eq(18)),
			wantquery: &Query{
				SQL:  "SELECT COUNT(`age`) FROM `pauls_model` WHERE `age` = ?;",
				Args: []any{18},
			},
		},
		{
			name:     "raw",
			selector: NewSelector[PaulsModel](db).Select(NewRaw("ROW_NUMBER()OVER(ORDER BY ID)").As("index")),
			wantquery: &Query{
				SQL: "SELECT ROW_NUMBER()OVER(ORDER BY ID) AS `index` FROM `pauls_model`;",
			},
		},
		{
			name:     "raw",
			selector: NewSelector[PaulsModel](db).Select(NewRaw("age+?", 1).As("nextyear")),
			wantquery: &Query{
				SQL:  "SELECT age+? AS `nextyear` FROM `pauls_model`;",
				Args: []any{1},
			},
		},
		{
			name:     "where raw",
			selector: NewSelector[PaulsModel](db).Where(NewRaw("age=id+1").AsPredicate()),
			wantquery: &Query{
				SQL: "SELECT * FROM `pauls_model` WHERE age=id+1;",
			},
		},
		{
			name:     "group by",
			selector: NewSelector[PaulsModel](db).Select(Sum("Age").As("allage")).GroupBy(NewColumns("Age")),
			wantquery: &Query{
				SQL: "SELECT SUM(`age`) AS `allage` FROM `pauls_model` GROUP BY `age`;",
			},
		},
		{
			name:     "mu group by",
			selector: NewSelector[PaulsModel](db).GroupBy(NewColumns("Age"), NewColumns("FirstName")),
			wantquery: &Query{
				SQL: "SELECT * FROM `pauls_model` GROUP BY `age`,`first_name`;",
			},
		},
		{
			name:     "having",
			selector: NewSelector[PaulsModel](db).Having(Sum("Age").EQ(18)),
			wantquery: &Query{
				SQL:  "SELECT * FROM `pauls_model` HAVING SUM(`age`) = ?;",
				Args: []any{18},
			},
		},
		{
			name:     "order by",
			selector: NewSelector[PaulsModel](db).OrderBy(Asc("Id")),
			wantquery: &Query{
				SQL: "SELECT * FROM `pauls_model` ORDER BY `id` ASC;",
			},
		},
		{
			name:     "limit",
			selector: NewSelector[PaulsModel](db).Limit(8),
			wantquery: &Query{
				SQL:  "SELECT * FROM `pauls_model` LIMIT ?;",
				Args: []any{8},
			},
		},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.selector.Build()
			assert.NoError(t, tc.wanterror, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantquery, query)
		})
	}
}
