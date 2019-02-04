package govaluate

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/mantzas/rungine"
)

// Expression definition.
type Expression struct {
	exp *govaluate.EvaluableExpression
}

// NewExpression constructor.
func NewExpression(exp string) (*Expression, error) {
	evExp, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		return nil, err
	}
	return &Expression{exp: evExp}, nil
}

// Eval evaluates the expression against the facts.
func (e *Expression) Eval(facts rungine.Facts) (bool, string, error) {
	val, err := e.exp.Evaluate(facts)
	if err != nil {
		return false, e.exp.String(), err
	}
	suc, ok := val.(bool)
	if !ok {
		return false, e.exp.String(), fmt.Errorf("expression returned non-bool result: %v", val)
	}
	return suc, e.exp.String(), nil
}
