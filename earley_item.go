package gearley

import (
	"fmt"
	"strings"
)

const FLAT_DOT = "\u25CF"

type earleyItem struct {
	rule  *rule
	dot   int
	index int
}

func (t *earleyItem) String() string {
	rightStrings := make([]string, len(t.rule.right))
	for i, r := range t.rule.right {
		rightStrings[i] = r.String()
	}
	return fmt.Sprintf("%v -> %v%v%v (%d)",
		t.rule.left.String(),
		strings.Join(rightStrings[0:t.dot], " "),
		FLAT_DOT,
		strings.Join(rightStrings[t.dot:], " "),
		t.index,
	)
}

func (t *earleyItem) isCompleted() bool {
	return t.dot == t.rule.length()
}

func (t *earleyItem) getSymbolAt(i int) symbol {
	return t.rule.right[i]
}

func (t *earleyItem) getNext() symbol {
	return t.rule.right[t.dot]
}

func (t *earleyItem) isNextMatchingTerminal(nextRune rune) bool {
	return t.getNext().isMatchingTerminal(nextRune)
}
