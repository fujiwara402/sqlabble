package stmt

import (
	"github.com/minodisk/sqlabble/keyword"
	"github.com/minodisk/sqlabble/token"
	"github.com/minodisk/sqlabble/tokenizer"
)

type Where struct {
	prev      Prever
	operation ComparisonOrLogicalOperation
}

func NewWhere(operation ComparisonOrLogicalOperation) Where {
	return Where{
		operation: operation,
	}
}

func (w Where) GroupBy(column Column, columns ...Column) GroupBy {
	g := NewGroupBy(column, columns...)
	g.prev = w
	return g
}

func (w Where) OrderBy(orders ...Order) OrderBy {
	o := NewOrderBy(orders...)
	o.prev = w
	return o
}

func (w Where) Limit(count int) Limit {
	l := NewLimit(count)
	l.prev = w
	return l
}

func (w Where) nodeize() (tokenizer.Tokenizer, []interface{}) {
	return nodeizePrevs(w)
}

func (w Where) nodeizeSelf() (tokenizer.Tokenizer, []interface{}) {
	middle, values := w.operation.nodeize()
	return tokenizer.NewContainer(
		tokenizer.NewLine(token.Word(keyword.Where)),
	).SetMiddle(
		middle,
	), values
}

func (w Where) previous() Prever {
	return w.prev
}
