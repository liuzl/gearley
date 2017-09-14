package gearley

import (
	"fmt"
	"strings"
)

type rule struct {
	left  nonTerminal
	right []symbol
}

func (r *rule) length() int {
	return len(r.right)
}

func (r *rule) String() string {
	rightStrings := make([]string, len(r.right))
	for i, s := range r.right {
		rightStrings[i] = s.String()
	}
	return fmt.Sprintf("%v -> %v", r.left.String(), strings.Join(rightStrings, " "))
}

func Rule(t nonTerminal, symbols ...symbol) *rule {
	return &rule{left: t, right: symbols}
}
