package earley3

import (
	"errors"
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
	Equal(interface{}) bool
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
	if reflect.DeepEqual(self, other) {
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

func (self Rule) Equal(other interface{}) bool {
	return false
}

func (self Rule) String() string {
	return self.Name
}

func (self Rule) Type() string {
	return "Non-Terminal"
}
