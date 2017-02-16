package operator

// Operator is reserved keyword in SQL statement.
type Operator string

// Logical operators.
const (
	And Operator = "AND"
	Or  Operator = "OR"
	Not Operator = "NOT"

	All    Operator = "ALL"
	Any    Operator = "ANY"
	Exists Operator = "EXISTS"
	Unique Operator = "UNIQUE"

	Like   Operator = "LIKE"
	RegExp Operator = "REGEXP"

	Between Operator = "BETWEEN"

	In    Operator = "IN"
	NotIn Operator = "NOT IN"

	IsNull    Operator = "IS NULL"
	IsNotNull Operator = "IS NOT NULL"
)

// Comparison operators.
const (
	Eq    Operator = "="
	NotEq Operator = "!="
	Gt    Operator = ">"
	Gte   Operator = ">="
	Lt    Operator = "<"
	Lte   Operator = "<="
)

// Arithmetic operators.
const (
	Add Operator = "+"
	Sub Operator = "-"
	Mul Operator = "*"
	Div Operator = "/"
	Mod Operator = "%"
)

// Alias operators.
const (
	As Operator = "AS"
)