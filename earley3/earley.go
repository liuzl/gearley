package earley3

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
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

// Epsilon transition: an empty production
var Epsilon = Production{}

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
	states []*TableState
}

/*
 * only insert a state if it is not already contained in the list of states. return the
 * inserted state, or the pre-existing one.
 */
func (self *TableColumn) insert(state *TableState) *TableState {
	for _, s := range self.states {
		if *state == *s {
			return s
		}
	}
	self.states = append(self.states, state)
	state.endCol = self
	return self.get(self.size() - 1)
}

func (self *TableColumn) size() int {
	return len(self.states)
}

func (self *TableColumn) get(index int) *TableState {
	return self.states[index]
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

func NewParser(startRule *Rule, text string) *Parser {
	tokens := strings.Fields(text)
	parser := &Parser{}
	parser.columns = append(parser.columns, &TableColumn{index: 0, token: ""})
	for i, token := range tokens {
		parser.columns = append(parser.columns,
			&TableColumn{index: i + 1, token: token})
	}
	parser.finalState = parser.parse(startRule)
	// TODO
	return parser
}

// this is the name of the special "gamma" rule added by the algorithm
// (this is unicode for 'LATIN SMALL LETTER GAMMA')
const GAMMA_RULE = "\u0263" // "\u0194"

/*
 * the Earley algorithm's core: add gamma rule, fill up table, and check if the gamma rule
 * span from the first column to the last one. return the final gamma state, or null,
 * if the parse failed.
 */
func (self *Parser) parse(startRule *Rule) *TableState {
	self.columns[0].states = append(self.columns[0].states,
		&TableState{
			name:       GAMMA_RULE,
			production: NewProductionFromTerms(startRule),
			dotIndex:   0,
			startCol:   self.columns[0],
		})
	for i, col := range self.columns {
		for j := 0; j < len(col.states); j++ {
			state := col.states[j]
			if state.isCompleted() {
				self.complete(col, state)
			} else {
				var term interface{} = state.getNextTerm()
				switch term.(type) {
				case *Rule:
					self.predict(col, term.(*Rule))
				case *Terminal:
					if i+1 < len(self.columns) {
						self.scan(self.columns[i+1], state, term.(*Terminal).Value)
					}
				}
			}
		}
		self.handleEpsilons(col)
		// DEBUG -- uncomment to print the table during parsing, column after column
		col.Print(os.Stdout, false)
	}

	// find end state (return nil if not found)
	lastCol := self.columns[len(self.columns)-1]
	for i := 0; i < len(lastCol.states); i++ {
		if lastCol.states[i].name == GAMMA_RULE && lastCol.states[i].isCompleted() {
			return lastCol.states[i]
		}
	}
	return nil
}

/*
 * Earley scan
 */
func (self *Parser) scan(col *TableColumn, state *TableState, token string) {
	if token == col.token {
		col.insert(&TableState{
			name:       state.name,
			production: state.production,
			dotIndex:   state.dotIndex + 1,
			startCol:   state.startCol,
		})
	}
}

/*
 * Earley predict. returns true if the table has been changed, false otherwise
 */
func (self *Parser) predict(col *TableColumn, rule *Rule) bool {
	changed := false
	for _, prod := range rule.Productions {
		st := &TableState{name: rule.Name, production: prod, dotIndex: 0, startCol: col}
		st2 := col.insert(st)
		changed = changed || (st == st2)
	}
	return changed
}

/*
 * Earley complete. returns true if the table has been changed, false otherwise
 */
func (self *Parser) complete(col *TableColumn, state *TableState) bool {
	changed := false
	for _, st := range state.startCol.states {
		var term interface{} = st.getNextTerm()
		if r, ok := term.(*Rule); ok && r.Name == state.name {
			st := &TableState{name: r.Name, production: st.production, dotIndex: st.dotIndex + 1, startCol: st.startCol}
			st2 := col.insert(st)
			changed = changed || (st == st2)
		}
	}
	return changed
}

/*
 * call predict() and complete() for as long as the table keeps changing (may only happen
 * if we've got epsilon transitions)
 */
func (self *Parser) handleEpsilons(col *TableColumn) {
	changed := true
	for changed {
		changed = false
		for _, state := range col.states {
			var term interface{} = state.getNextTerm()
			if r, ok := term.(*Rule); ok {
				changed = changed || self.predict(col, r)
			}
			if state.isCompleted() {
				changed = changed || self.complete(col, state)
			}
		}
	}
}

/*
 * return all parse trees (forest). the forest is simply a list of root nodes, each
 * representing a possible parse tree. a node is contains a value and the node's children,
 * and supports pretty-printing
 */
func (self *Parser) getTrees() []*Node {
	return self.buildTrees(self.finalState)
}

/*
 * this is a bit "magical" -- i wrote the code that extracts a single parse tree,
 * and with some help from a colleague (non-student) we managed to make it return all
 * parse trees.
 *
 * how it works: suppose we're trying to match [X -> Y Z W]. we go from finish-to-start,
 * e.g., first we'll try to match W in X.encCol. let this matching state be M1. next we'll
 * try to match Z in M1.startCol. let this matching state be M2. and finally, we'll try to
 * match Y in M2.startCol, which must also start at X.startCol. let this matching state be
 * M3.
 *
 * if we matched M1, M2 and M3, then we've found a parsing for X:
 * X->
 *    Y -> M3
 *    Z -> M2
 *    W -> M1
 *
 */
func (self *Parser) buildTrees(state *TableState) []*Node {
	return self.buildTreesHelper(&[]*Node{}, state, len(state.production.Rules)-1, state.endCol)
}

func (self *Parser) buildTreesHelper(children *[]*Node, state *TableState, ruleIndex int, endCol *TableColumn) []*Node {
	var outputs []*Node
	var startCol *TableColumn
	if ruleIndex < 0 {
		// this is the base-case for the recursion (we matched the entire rule)
		outputs = append(outputs, &Node{value: state, children: *children})
		return outputs
	} else if ruleIndex == 0 {
		// if this is the first rule
		startCol = state.startCol
	}
	rule := state.production.Rules[ruleIndex]

	for _, st := range state.endCol.states {
		if st == state {
			// this prevents an endless recursion: since the states are filled in order of
			// completion, we know that X cannot depend on state Y that comes after it X
			// in chronological order
			break
		}
		if !st.isCompleted() || st.name != rule.Name {
			// this state is out of the question -- either not completed or does not match
			// the name
			continue
		}
		if startCol != nil && st.startCol != startCol {
			// if startCol isn't nil, this state must span from startCol to endCol
			continue
		}
		// okay, so `st` matches -- now we need to create a tree for every possible sub-match
		for _, subTree := range self.buildTrees(st) {
			// in python: children2 = [subTree] + children
			children2 := []*Node{}
			children2 = append(children2, subTree)
			children2 = append(children2, *children...)
			// now try all options
			for _, node := range self.buildTreesHelper(&children2, state, ruleIndex-1, st.startCol) {
				outputs = append(outputs, node)
			}
		}
	}
	return outputs
}
