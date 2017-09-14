package gearley

import "fmt"

type symbol interface {
	// isTerminal indicates if the symbol is Terminal Symbol or Non Terminal Symbol.
	isTerminal() bool
	String() string
	isMatchingTerminal(rune) bool
	// s and input are slices of the full state set and the full input.
}

type terminal struct {
	value rune
}

func Terminal(r rune) terminal {
	return terminal{value: r}
}

func (t terminal) isTerminal() bool {
	return true
}

func (t terminal) String() string {
	return fmt.Sprintf("'%c'", t.value)
}

func (t terminal) isMatchingTerminal(r rune) bool {
	return r == t.value
}

type nonTerminal struct {
	name string
}

func NonTerminal(name string) nonTerminal {
	return nonTerminal{name: name}
}

func (n nonTerminal) isTerminal() bool {
	return false
}

func (n nonTerminal) String() string {
	return n.name
}

func (n nonTerminal) isMatchingTerminal(r rune) bool {
	return false
}
