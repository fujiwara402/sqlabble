package sqlabble

import (
	"github.com/minodisk/sqlabble/generator"
	"github.com/minodisk/sqlabble/operator"
)

type joinOperation struct {
	op  operator.Operator
	ops []comparisonOrLogicalOperation
}

func newAnd(ops ...comparisonOrLogicalOperation) joinOperation {
	return joinOperation{
		op:  operator.And,
		ops: ops,
	}
}

func newOr(ops ...comparisonOrLogicalOperation) joinOperation {
	return joinOperation{
		op:  operator.Or,
		ops: ops,
	}
}

func (a joinOperation) node() generator.Node {
	ns := make([]generator.Node, len(a.ops))
	for i, op := range a.ops {
		ns[i] = op.node()
	}
	return generator.NewOperator(
		a.operator(),
		ns...,
	)
}

func (a joinOperation) operator() operator.Operator {
	return a.op
}

func (a joinOperation) operations() []comparisonOrLogicalOperation {
	return a.ops
}

type not struct {
	operation comparisonOrLogicalOperation
}

func newNot(operation comparisonOrLogicalOperation) not {
	return not{operation: operation}
}

func (n not) node() generator.Node {
	return generator.NewOpParenteses(
		n.operator(),
		generator.NewParentheses(n.operation.node()),
	)
}

func (n not) operator() operator.Operator {
	return operator.Not
}

func (n not) operations() []comparisonOrLogicalOperation {
	return []comparisonOrLogicalOperation{n.operation}
}

type comparisonOperation struct {
	op  operator.Operator
	col columnOrSubquery
	val interface{}
}

func newEq(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Eq,
		val: val,
	}
}

func newNotEq(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.NotEq,
		val: val,
	}
}

func newGt(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Gt,
		val: val,
	}
}

func newGte(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Gte,
		val: val,
	}
}

func newLt(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Lt,
		val: val,
	}
}

func newLte(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Lte,
		val: val,
	}
}

func newLike(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.Like,
		val: val,
	}
}

func newRegExp(val interface{}) comparisonOperation {
	return comparisonOperation{
		op:  operator.RegExp,
		val: val,
	}
}

func (c comparisonOperation) node() generator.Node {
	var n1, n2 generator.Node

	if c.col != nil {
		switch col := c.col.(type) {
		case column:
			n1 = col.expression()
		default:
			n1 = c.col.node()
		}
	}

	if c.val != nil {
		switch val := c.val.(type) {
		case sub:
			n2 = val.node()
		case Statement:
			n2 = newSub(val).node()
		default:
			n2 = generator.ValuesToExpression(val)
		}
	}

	op := generator.NewExpression(string(c.operator()))
	e1, ok1 := n1.(generator.Expression)
	e2, ok2 := n2.(generator.Expression)
	if ok1 && ok2 {
		return generator.JoinExpressions(e1, op, e2)
	}
	if ok1 {
		return generator.NewParallelNodes(e1.Append(op), n2)
	}
	if ok2 {
		return generator.NewParallelNodes(n1, e2.Prepend(op))
	}
	return generator.NewParallelNodes(n1, op, n2)
}

func (c comparisonOperation) operator() operator.Operator {
	return c.op
}

type between struct {
	col      column
	from, to interface{}
}

func newBetween(col column, from, to interface{}) between {
	return between{
		col:  col,
		from: from,
		to:   to,
	}
}

func (b between) node() generator.Node {
	return generator.JoinExpressions(
		generator.NewExpression(b.col.name),
		generator.NewExpression(string(b.operator())),
		generator.ValuesToExpression(b.from),
		generator.NewExpression(string(operator.And)),
		generator.ValuesToExpression(b.to),
	)
}

func (b between) operator() operator.Operator {
	return operator.Between
}

type containingOperation struct {
	op   operator.Operator
	col  column
	vals []interface{}
}

func newIn(col column, vals ...interface{}) containingOperation {
	return containingOperation{
		op:   operator.In,
		col:  col,
		vals: vals,
	}
}

func newNotIn(col column, vals ...interface{}) containingOperation {
	return containingOperation{
		op:   operator.NotIn,
		col:  col,
		vals: vals,
	}
}

func (o containingOperation) node() generator.Node {
	return generator.JoinExpressions(
		o.col.expression(),
		generator.NewExpression(string(o.operator())),
		generator.ValuesToExpression(o.vals...).
			WrapSQL("(", ")"),
	)
}

func (o containingOperation) operator() operator.Operator {
	return o.op
}

type nullyOperation struct {
	op  operator.Operator
	col column
}

func newIsNull(col column) nullyOperation {
	return nullyOperation{
		op:  operator.IsNull,
		col: col,
	}
}

func newIsNotNull(col column) nullyOperation {
	return nullyOperation{
		op:  operator.IsNotNull,
		col: col,
	}
}

func (o nullyOperation) node() generator.Node {
	return generator.JoinExpressions(
		o.col.expression(),
		generator.NewExpression(string(o.operator())),
	)
}

func (o nullyOperation) operator() operator.Operator {
	return o.op
}
