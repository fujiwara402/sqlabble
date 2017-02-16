package generator

import "strings"

// Context is a container for storing the state to be output
// in the process of building a query.
type Context struct {
	prefix, indent      string
	breaking, flatSets  bool
	head                string
	depth, bracketDepth int
}

func newContext(o Options) Context {
	return Context{
		prefix:       o.Prefix,
		indent:       o.Indent,
		breaking:     o.Prefix != "" || o.Indent != "",
		flatSets:     o.FlatSets,
		head:         "",
		depth:        0,
		bracketDepth: 0,
	}
}

func (c Context) currentHead() string {
	return c.head
}

func (c Context) clearHead() Context {
	c.head = ""
	return c
}

func (c Context) setHead(head string) Context {
	c.head = head
	return c
}

func (c Context) isBreaking() bool {
	return c.breaking
}

func (c Context) pre() string {
	return c.prefix + strings.Repeat(c.indent, c.depth)
}

func (c Context) incDepth() Context {
	c.depth++
	return c
}

func (c Context) clearParenthesesDepth() Context {
	c.bracketDepth = 0
	return c
}

func (c Context) incParenthesesDepth() Context {
	c.bracketDepth++
	return c
}

func (c Context) isTopParentheses() bool {
	return c.bracketDepth == 0
}

func (c Context) setFlatSet(flat bool) Context {
	c.flatSets = flat
	return c
}

func (c Context) join(sqls ...string) string {
	ss := []string{}
	for _, sql := range sqls {
		if sql != "" {
			ss = append(ss, sql)
		}
	}

	if c.isBreaking() {
		return strings.Join(ss, "")
	}
	return strings.Join(ss, " ")
}
