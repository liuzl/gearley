package earley3

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

/*
 * Terminology
 * ==========
 * Consider the following context-free rule:
 *
 *     X -> A B C | A hello
 *
 * We say rule 'X' has two __production__: "A B C" and "A hello".
 * Each production is made of __production terms__, which can be either
 * __terminals__ (in our case, "hello") or __rules__ (non-terminals, such as "A", "B", and "C")
 */

/*
 * an abstract notation of the elements that can be placed within production
 */
type ProductionTerm interface {
	String() string
	Type() string
}

/*
 * Represents a terminal element in a production
 */
type Terminal struct {
	Value string
}

func (self Terminal) String() string {
	return self.Value
}

func (self Terminal) Type() string {
	return "Terminal"
}

func (self Terminal) Equal(other interface{}) bool {
	if self == other { // no need using reflect.DeepEqual
		return true
	}
	if s, ok := other.(string); ok && s == self.Value {
		return true
	}
	return false
}

/*
 * Represents a production of the rule.
 */
type Production struct {
	Terms []ProductionTerm
	Rules []Rule
}

func NewProductionFromTerms(terms ...ProductionTerm) *Production {
	prod := &Production{Terms: terms}
	prod.getRules()
	return prod
}

func NewProduction(terms ...interface{}) (*Production, error) {
	prod := &Production{}
	for _, term := range terms {
		switch term.(type) {
		case string:
			prod.Terms = append(prod.Terms, Terminal{Value: term.(string)})
		case ProductionTerm:
			prod.Terms = append(prod.Terms, term.(ProductionTerm))
		case *ProductionTerm:
			prod.Terms = append(prod.Terms, *term.(*ProductionTerm))
		default:
			return nil, errors.New("Term must be ProductionTerm or string, not " + reflect.TypeOf(term).String())
		}
	}
	return prod, nil
}

func (self *Production) size() int {
	return len(self.Terms)
}

func (self *Production) get(index int) *ProductionTerm {
	return &self.Terms[index]
}

func (self *Production) getRules() {
	self.Rules = nil
	for _, term := range self.Terms {
		switch term.(type) {
		case Rule:
			self.Rules = append(self.Rules, term.(Rule))
		}
	}
}

func (self Production) String() string {
	s := ""
	for i, term := range self.Terms {
		switch term.(type) {
		case Terminal:
			s += term.String()
		case Rule:
			s += term.(Rule).Name
		}
		if i != self.size()-1 {
			s += " "
		}
	}
	return s
}

/*
 * A CFG rule. Since CFG rules can be self-referential, more productions may be added
 * to them after construction. For example:
 *
 * Grammar:
 *    SYM -> a
 *    OP -> + | -
 *    EXPR -> SYM | EXPR OP EXPR
 *
 * In Java:
 *     Rule SYM = new Rule("SYM", new Production("a"));
 *     Rule OP = new Rule("OP", new Production("+"), new Production("-"));
 *     Rule EXPR = new Rule("EXPR", new Production(SYM));
 *     EXPR.add(new Production(EXPR, OP, EXPR));            // needs to reference EXPR
 *
 */

type Rule struct {
	Name        string
	Productions []*Production
}

func NewRule(name string, productions ...*Production) *Rule {
	return &Rule{Name: name, Productions: productions}
}

func (self *Rule) add(productions ...*Production) {
	self.Productions = append(self.Productions, productions...)
}

func (self *Rule) size() int {
	return len(self.Productions)
}

func (self *Rule) get(index int) *Production {
	return self.Productions[index]
}

func (self Rule) String() string {
	s := self.Name + " -> "
	for i, prod := range self.Productions {
		s += prod.String()
		if i != self.size()-1 {
			s += " | "
		}
	}
	return s
}

func (self Rule) Type() string {
	return "Non-Terminal"
}

/*
 * Represents a state in the Earley parsing table. A state has its rule's name,
 * the rule's production, dot-location, and starting- and ending-column in the parsing
 * table
 */
type TableState struct {
	name       string
	production *Production
	dotIndex   int
	startCol   *TableColumn
	endCol     *TableColumn
}

func (self *TableState) isCompleted() bool {
	return self.dotIndex >= self.production.size()
}

func (self *TableState) getNextTerm() *ProductionTerm {
	if self.isCompleted() {
		return nil
	}
	return self.production.get(self.dotIndex)
}

func (self TableState) String() string {
	s := ""
	for i, term := range self.production.Terms {
		if i == self.dotIndex {
			s += "\u00B7"
		}
		switch term.(type) {
		case Terminal:
			s += term.String()
		case Rule:
			s += term.(Rule).Name
		}
		s += " "
	}
	if self.dotIndex == self.production.size() {
		s += "\u00B7"
	}
	return fmt.Sprintf("%s -> %s [%d-%d]",
		self.name, s, self.startCol.index, self.endCol.index)
}

/*
 * Represents a column in the Earley parsing table
 */
type TableColumn struct {
	token  string
	index  int
	states []TableState
}

/*
 * only insert a state if it is not already contained in the list of states. return the
 * inserted state, or the pre-existing one.
 */
func (self *TableColumn) insert(state TableState) *TableState {
	for _, s := range self.states {
		if state == s {
			return &s
		}
	}
	self.states = append(self.states, state)
	return self.get(self.size() - 1)
}

func (self *TableColumn) size() int {
	return len(self.states)
}

func (self *TableColumn) get(index int) *TableState {
	return &self.states[index]
}

func (self *TableColumn) Print(out *os.File, showUncompleted bool) {
	fmt.Fprintf(out, "[%d] '%s'\n", self.index, self.token)
	fmt.Fprintln(out, "=======================================")
	for _, s := range self.states {
		if !s.isCompleted() && !showUncompleted {
			continue
		}
		fmt.Fprintln(out, s)
	}
	fmt.Fprintln(out)
}

/*
 * A generic tree node
 */
type Node struct {
	value    interface{}
	children []*Node
}

func (self *Node) Print(out *os.File) {
	self.PrintLevel(out, 0)
}

func (self *Node) PrintLevel(out *os.File, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += " "
	}
	fmt.Fprintf(out, "%s%v", indentation, self.value)
	for _, child := range self.children {
		child.PrintLevel(out, level+1)
	}
}

/*
 * The Earley Parser.
 *
 * Usage:
 *
 *   var p *Parser = NewParser(StartRule, "my space-delimited statement")
 *   for _, tree := range p.getTrees() {
 *     tree.Print(os.Stdout)
 *   }
 *
 */
type Parser struct {
	columns    []*TableColumn
	finalState *TableState
}
