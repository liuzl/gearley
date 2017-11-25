package earley3

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTermTypes(t *testing.T) {
	terms := []interface{}{
		Terminal{Value: "terminal"},
		Rule{Name: "Rule"},
	}
	for _, term := range terms {
		t.Log(term.(ProductionTerm).Type())
	}
}

func TestTermEqual(t *testing.T) {
	term := Terminal{Value: "刘占亮"}
	type Case struct {
		other interface{}
		ret   bool
	}
	cases := []Case{
		{other: "刘占亮", ret: true},
		{other: Terminal{Value: "刘占亮"}, ret: true},
		{other: Terminal{Value: ""}, ret: false},
		{other: nil, ret: false},
		{other: "", ret: false},
	}

	for i, c := range cases {
		fmt.Println(reflect.TypeOf(c.other))
		if term.Equal(c.other) != c.ret {
			t.Error("case ", i, c, " do not passed")
		}
	}
}
