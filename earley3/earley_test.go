package earley3

import (
	"testing"
)

func TestTermTypes(t *testing.T) {
	terms := []interface{}{
		&Terminal{value: "terminal"},
		&Rule{name: "Rule"},
	}
	for _, term := range terms {
		t.Log(term)
	}
}

func TestEarleyParse(t *testing.T) {
	SYM := NewRule("SYM", NewProduction(&Terminal{"a"}))
	OP := NewRule("OP", NewProduction(&Terminal{"+"}))
	EXPR := NewRule("EXPR", NewProduction(SYM))
	EXPR.add(NewProduction(EXPR, OP, EXPR))
	t.Log(SYM)
	t.Log(OP)
	t.Log(EXPR)

	strs := []string{
		//"a",
		"a + a",
		//"a + a + a",
		//"a + a + a + a",
		//"a + a + a + a + a",
		//"a + a + a + a + a + a",
		//"a + a + a + a + a + a + a",
	}
	for _, text := range strs {
		NewParser(EXPR, text)
		t.Log(text)
	}
}
