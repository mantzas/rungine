package rungine

import "errors"

// Facts definition.
type Facts map[string]interface{}

// Results definition.
type Results map[string]interface{}

// EvalResult defines the results of the evaluation.
type EvalResult struct {
	Definition string
	Success    bool
}

// EvalFunc definition.
type EvalFunc func(facts Facts) EvalResult

// Rule definition.
type Rule struct {
	Eval EvalFunc
	Next *Node
	Res  Results
}

// Node definition.
type Node struct {
	rules []*Rule
}

// AppendDecisionRule that links the previous node to the next forming a tree.
func (n *Node) AppendDecisionRule(eval EvalFunc, next *Node) error {
	if eval == nil {
		return errors.New("evaluation func is nil")
	}
	if next == nil {
		return errors.New("next is nil")
	}
	n.rules = append(n.rules, &Rule{Eval: eval, Next: next})
	return nil
}

// AppendResultRule that forms tha leaf of the tree containing the predefined result.
func (n *Node) AppendResultRule(eval EvalFunc, res Results) error {
	if eval == nil {
		return errors.New("evaluation func is nil")
	}
	if res == nil {
		return errors.New("res is nil")
	}
	n.rules = append(n.rules, &Rule{Eval: eval, Res: res})
	return nil
}

// Eval returns the result of the evaluation of the decision tree.
func (n *Node) Eval(facts Facts, audit []string) (Results, []string) {
	for _, r := range n.rules {
		re := r.Eval(facts)
		if !re.Success {
			continue
		}
		audit = append(audit, re.Definition)
		if r.Next != nil {
			return r.Next.Eval(facts, audit)
		}
		return r.Res, audit
	}
	return nil, audit
}
