package gearley

// Eerley parser.
// http://loup-vaillant.fr/tutorials/earley-parsing/recogniser

import (
	"fmt"
	"strings"
)

// state is the highest-level state of the parser.
type state []*stateSet

func (st *state) String() string {
	ss := make([]string, len(*st))
	for i, s := range *st {
		ss[i] = fmt.Sprint(i, " ", s)
	}
	return strings.Join(ss, "\n")
}

func (st *state) getAt(i int) *stateSet {
	return (*st)[i]
}

func (g *grammar) Parse(input string) {
	inputRunes := stringToRunes(input)
	st := initializeState(g, inputRunes)
	// the current index in the state 'st' that is being processed - S(stateIndex)
	stateIndex := 0
	// outter loop
	for stateIndex <= len(input) {
		set := st.getAt(stateIndex)
		fmt.Println("NOW SET S(", stateIndex, ")", set)
		i := 0
		// inner loop
		for i < set.length() {
			item := set.items[i]
			fmt.Println(i, item)
			i++

			fmt.Println("NOW item", item)

			if item.isCompleted() {
				fmt.Println("Completetion")
				originalSet := st.getAt(item.index)
				itemsToComplete := originalSet.findItemsToComplete(item.rule.left)
				fmt.Println("to complete: ", itemsToComplete)
				for _, itemToComplete := range itemsToComplete {
					nextItem := &earleyItem{
						rule:  itemToComplete.rule,
						dot:   itemToComplete.dot + 1,
						index: itemToComplete.index,
					}
					set.putItem(nextItem)
				}
				continue
			}
			if item.isNextMatchingTerminal(inputRunes[stateIndex]) {
				// Scan - the next symbol is Terminal and matches
				fmt.Println("Scan - terminal")
				nextItem := &earleyItem{
					rule:  item.rule,
					dot:   item.dot + 1,
					index: item.index,
				}
				// create next item
				// add it to the next stateSet
				fmt.Println("next item", nextItem)
				// TODO edge case when last stateIndex
				nextSet := st.getAt(stateIndex + 1)
				nextSet.putItem(nextItem)
				continue
			}
			if !item.getNext().isTerminal() {
				// Predict - the next symbol is Non Terminal
				nextSymbol := item.getNext().(nonTerminal)
				// Find all the rules for the symbol put those rules to the current set
				fmt.Println("Predict - NON TERMINAL")
				for _, r := range g.getRulesForSymbol(nextSymbol) {
					nextItem := &earleyItem{
						rule:  r,
						dot:   0,
						index: stateIndex,
					}
					set.putItem(nextItem)
				}
				continue
			}
		}

		fmt.Printf("S\n%v\n", st.String())
		stateIndex++
	}
}

func (g *grammar) getRulesForSymbol(s symbol) []*rule {
	found := []*rule{}
	for _, r := range g.rules {
		if r.left == s {
			found = append(found, r)
		}
	}
	return found
}

func initializeState(g *grammar, runes []rune) *state {
	sets := make([]*stateSet, len(runes)+1)
	for i := range sets {
		sets[i] = newStateSet()
	}
	sets[0] = newStateSetFromRules(g.rules)
	s := state(sets)
	return &s
}

func newStateSetFromRules(rules []*rule) *stateSet {
	items := make([]*earleyItem, len(rules))
	for i, r := range rules {
		items[i] = &earleyItem{rule: r, dot: 0, index: 0}
	}
	ss := newStateSet()
	ss.items = items
	return ss
}

func stringToRunes(input string) []rune {
	runes := []rune{}
	for _, r := range input {
		runes = append(runes, r)
	}
	return runes
}
