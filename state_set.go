package gearley

import "fmt"

type stateSet struct {
	items   []*earleyItem
	itemSet map[earleyItem]bool
}

func newStateSet() *stateSet {
	return &stateSet{items: []*earleyItem{}, itemSet: make(map[earleyItem]bool)}
}

func (s *stateSet) String() string {
	return fmt.Sprint(s.items)
}

func (s *stateSet) length() int {
	return len(s.items)
}

func (s *stateSet) putItem(item *earleyItem) {
	if _, ok := s.itemSet[*item]; ok {
		return
	}
	s.itemSet[*item] = true
	s.items = append(s.items, item)
}

func (s *stateSet) findItemsToComplete(t nonTerminal) []*earleyItem {
	candidates := []*earleyItem{}
	for _, item := range s.items {
		switch c := item.rule.right[item.dot].(type) {
		case nonTerminal:
			if c.name == t.name {
				candidates = append(candidates, item)
			}
		default:
			continue
		}
	}
	return candidates
}
