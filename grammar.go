package gearley

type grammar struct {
	rules []*rule
}

func Grammar(rules ...*rule) *grammar {
	return &grammar{rules: rules}
}
