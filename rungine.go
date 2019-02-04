package rungine

import (
	"errors"
	"fmt"
	"strings"
)

// Facts definition.
type Facts map[string]interface{}

// Results definition.
type Results map[string]interface{}

// EvalFunc definition.
type EvalFunc func(facts Facts) (bool, string, error)

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
func (n *Node) Eval(facts Facts, audit []string) (Results, []string, error) {
	if audit == nil {
		audit = make([]string, 0)
	}
	for _, r := range n.rules {
		success, expr, err := r.Eval(facts)
		if err != nil {
			return nil, nil, err
		}
		if !success {
			continue
		}
		audit = append(audit, expr)
		if r.Next != nil {
			return r.Next.Eval(facts, audit)
		}
		audit = append(audit, resultToString(r.Res))
		return r.Res, audit, nil
	}
	return nil, nil, nil
}

func resultToString(res Results) string {
	sb := strings.Builder{}
	count := 0
	for k, v := range res {
		if count > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(k)
		sb.WriteRune(':')
		sb.WriteString(fmt.Sprint(v))
		count++
	}
	return sb.String()
}
