package statement_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/minodisk/sqlabble/internal/diff"
	"github.com/minodisk/sqlabble/statement"
)

// func TestSetType(t *testing.T) {
// 	for _, c := range []interface{}{
// 		statement.Set{},
// 	} {
// 		t.Run(fmt.Sprintf("Type %T", c), func(t *testing.T) {
// 			if _, ok := c.(statement.ClauseNode); !ok {
// 				t.Errorf("%T should implement statement.Clause", c)
// 			}
// 		})
// 	}
// }

func TestSet(t *testing.T) {
	for i, c := range []struct {
		statement statement.Statement
		sql       string
		sqlIndent string
		values    []interface{}
	}{
		{
			statement.NewSet(
				statement.NewColumn("foo").Assign(100),
			),
			"SET foo = ?",
			`> SET
>   foo = ?
`,
			[]interface{}{
				100,
			},
		},
		{
			statement.NewSet(
				statement.NewColumn("foo").Assign(100),
				statement.NewColumn("bar").Assign(200),
			),
			"SET foo = ?, bar = ?",
			`> SET
>   foo = ?
>   , bar = ?
`,
			[]interface{}{
				100,
				200,
			},
		},
		{
			statement.NewSet(
				statement.NewColumn("foo").Assign(100),
				statement.NewColumn("bar").Assign(200),
				statement.NewColumn("baz").Assign(300),
			),
			"SET foo = ?, bar = ?, baz = ?",
			`> SET
>   foo = ?
>   , bar = ?
>   , baz = ?
`,
			[]interface{}{
				100,
				200,
				300,
			},
		},
		{
			statement.NewSet(
				statement.NewColumn("foo").Assign(100),
				statement.NewColumn("bar").Assign(200),
				statement.NewColumn("baz").Assign(300),
			).Where(
				statement.NewColumn("qux").Lte(400),
			),
			"SET foo = ?, bar = ?, baz = ? WHERE qux <= ?",
			`> SET
>   foo = ?
>   , bar = ?
>   , baz = ?
> WHERE
>   qux <= ?
`,
			[]interface{}{
				100,
				200,
				300,
				400,
			},
		},
	} {
		t.Run(fmt.Sprintf("%d Build", i), func(t *testing.T) {
			sql, values := b.Build(c.statement)
			if sql != c.sql {
				t.Error(diff.SQL(sql, c.sql))
			}
			if !reflect.DeepEqual(values, c.values) {
				t.Error(diff.Values(values, c.values))
			}
		})
		t.Run(fmt.Sprintf("%d BuildIndent", i), func(t *testing.T) {
			sql, values := bi.Build(c.statement)
			if sql != c.sqlIndent {
				t.Error(diff.SQL(sql, c.sqlIndent))
			}
			if !reflect.DeepEqual(values, c.values) {
				t.Error(diff.Values(values, c.values))
			}
		})
	}
}
