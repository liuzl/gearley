package earley3

import (
	"os"
	"testing"
)

func TestEarleyParse(t *testing.T) {
	SYM := NewRule("SYM", NewProduction(&Terminal{"a"}))
	OP := NewRule("OP", NewProduction(&Terminal{"+"}))
	EXPR := NewRule("EXPR", NewProduction(SYM))
	EXPR.add(NewProduction(EXPR, OP, EXPR))

	strs := []string{
		//"a",
		"a + a",
		"a + a + a",
		//"a + a + a + a",
		//"a + a + a + a + a",
		//"a + a + a + a + a + a",
		"a + a + a + a + a + a + a",
	}
	for _, text := range strs {
		p := NewParser(EXPR, text)
		trees := p.getTrees()
		t.Log("tree number:", len(*trees))
		for _, tree := range *trees {
			tree.Print(os.Stdout)
		}
		//(*trees)[0].Print(os.Stdout)
	}
}
