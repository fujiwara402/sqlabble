package stmt

import (
	"github.com/minodisk/sqlabble/keyword"
	"github.com/minodisk/sqlabble/token"
	"github.com/minodisk/sqlabble/tokenizer"
)

type JoinOperation struct {
	op  keyword.Operator
	ops []ComparisonOrLogicalOperation
}

func NewAnd(ops ...ComparisonOrLogicalOperation) JoinOperation {
	return JoinOperation{
		op:  keyword.And,
		ops: ops,
	}
}

func NewOr(ops ...ComparisonOrLogicalOperation) JoinOperation {
	return JoinOperation{
		op:  keyword.Or,
		ops: ops,
	}
}

func (o JoinOperation) nodeize() (tokenizer.Tokenizer, []interface{}) {
	ts := make(tokenizer.Tokenizers, len(o.ops))
	values := []interface{}{}
	for i, op := range o.ops {
		t, vals := op.nodeize()
		if _, ok := op.(JoinOperation); ok {
			t = tokenizer.NewParentheses(t)
		}
		ts[i] = t
		values = append(values, vals...)
	}
	return tokenizer.NewTokenizers(ts...).Prefix(
		token.Word(o.keyword()),
	), values
}

func (o JoinOperation) keyword() keyword.Operator {
	return o.op
}

type Not struct {
	operation ComparisonOrLogicalOperation
}

func NewNot(operation ComparisonOrLogicalOperation) Not {
	return Not{operation: operation}
}

func (o Not) nodeize() (tokenizer.Tokenizer, []interface{}) {
	middle, values := o.operation.nodeize()
	return tokenizer.NewParentheses(
		middle,
	).Prepend(
		token.Word(o.keyword()),
	), values
}

func (o Not) keyword() keyword.Operator {
	return keyword.Not
}

type ComparisonOperation struct {
	op     keyword.Operator
	column ValOrColOrFuncOrSub
	val    ValOrColOrFuncOrSub
}

func NewEq(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Eq,
		val: val,
	}
}

func NewNotEq(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.NotEq,
		val: val,
	}
}

func NewGt(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Gt,
		val: val,
	}
}

func NewGte(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Gte,
		val: val,
	}
}

func NewLt(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Lt,
		val: val,
	}
}

func NewLte(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Lte,
		val: val,
	}
}

func NewLike(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.Like,
		val: val,
	}
}

func NewRegExp(val ValOrColOrFuncOrSub) ComparisonOperation {
	return ComparisonOperation{
		op:  keyword.RegExp,
		val: val,
	}
}

func (o ComparisonOperation) nodeize() (tokenizer.Tokenizer, []interface{}) {
	t1, v1 := o.column.nodeize()
	t2, v2 := o.val.nodeize()
	return tokenizer.ConcatTokenizers(
		t1,
		t2,
		tokenizer.NewLine(
			token.Word(o.keyword()),
		),
	), append(v1, v2...)
}

func (o ComparisonOperation) keyword() keyword.Operator {
	return o.op
}

type Between struct {
	column   ValOrColOrFuncOrSub
	from, to ValOrColOrFuncOrSub
}

func NewBetween(from, to ValOrColOrFuncOrSub) Between {
	return Between{
		from: from,
		to:   to,
	}
}

func (o Between) nodeize() (tokenizer.Tokenizer, []interface{}) {
	t1, v1 := o.column.nodeize()
	t2, v2 := o.from.nodeize()
	t3, v3 := o.to.nodeize()
	return tokenizer.ConcatTokenizers(
		tokenizer.ConcatTokenizers(
			t1,
			t2,
			tokenizer.NewLine(
				token.Word(o.keyword()),
			),
		),
		t3,
		tokenizer.NewLine(
			token.Word(keyword.And),
		),
	), append(append(v1, v2...), v3...)
}

func (o Between) keyword() keyword.Operator {
	return keyword.Between
}

type ContainingOperation struct {
	op     keyword.Operator
	column ValOrColOrFuncOrSub
	vals   ValsOrSub
}

func NewIn(vals ValsOrSub) ContainingOperation {
	return ContainingOperation{
		op:   keyword.In,
		vals: vals,
	}
}

func NewNotIn(vals ValsOrSub) ContainingOperation {
	return ContainingOperation{
		op:   keyword.NotIn,
		vals: vals,
	}
}

func (o ContainingOperation) nodeize() (tokenizer.Tokenizer, []interface{}) {
	t1, v1 := o.column.nodeize()
	t2, v2 := o.vals.nodeize()
	return tokenizer.ConcatTokenizers(
		t1,
		t2,
		tokenizer.NewLine(
			token.Word(o.keyword()),
		),
	), append(v1, v2...)
}

func (o ContainingOperation) keyword() keyword.Operator {
	return o.op
}

type NullOperation struct {
	op     keyword.Operator
	column ColOrSub
}

func NewIsNull() NullOperation {
	return NullOperation{
		op: keyword.Is,
	}
}

func NewIsNotNull() NullOperation {
	return NullOperation{
		op: keyword.IsNot,
	}
}

func (o NullOperation) nodeize() (tokenizer.Tokenizer, []interface{}) {
	t1, v1 := o.column.nodeize()
	return tokenizer.ConcatTokenizers(
		t1,
		tokenizer.NewLine(token.Word(keyword.Null)),
		tokenizer.NewLine(
			token.Word(o.keyword()),
		),
	), v1
}

func (o NullOperation) keyword() keyword.Operator {
	return o.op
}
