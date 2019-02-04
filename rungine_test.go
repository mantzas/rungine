package rungine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEval struct {
	success bool
	def     string
}

func (me mockEval) eval(facts Facts) EvalResult {
	return EvalResult{
		Definition: me.def,
		Success:    me.success,
	}
}

var evalSuccess = mockEval{success: true, def: "eval success"}
var evalFailure = mockEval{success: false, def: "eval failure"}

func TestNode_AppendDecisionRule(t *testing.T) {
	type args struct {
		eval EvalFunc
		next *Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "missing eval function", args: args{eval: nil, next: new(Node)}, wantErr: true},
		{name: "missing next node", args: args{eval: evalSuccess.eval, next: nil}, wantErr: true},
		{name: "success", args: args{eval: evalSuccess.eval, next: new(Node)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := new(Node)
			err := n.AppendDecisionRule(tt.args.eval, tt.args.next)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_AppendResultRule(t *testing.T) {
	type args struct {
		eval EvalFunc
		res  Results
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "missing eval function", args: args{eval: nil, res: Results{}}, wantErr: true},
		{name: "missing result", args: args{eval: evalSuccess.eval, res: nil}, wantErr: true},
		{name: "success", args: args{eval: evalSuccess.eval, res: Results{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := new(Node)
			err := n.AppendResultRule(tt.args.eval, tt.args.res)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNode_Eval_NoMatchingRules(t *testing.T) {
	n := new(Node)
	n.AppendDecisionRule(evalFailure.eval, &Node{})
	n.AppendDecisionRule(evalFailure.eval, &Node{})
	n.AppendDecisionRule(evalFailure.eval, &Node{})
	fct := Facts{}
	aud := []string{}
	res, aud := n.Eval(fct, aud)
	assert.Nil(t, res)
	assert.Empty(t, aud)
}

func TestNode_Eval(t *testing.T) {
	fct := Facts{"test": 123}
	res := Results{"res1": 123}
	aud := []string{"eval success", "eval success", "eval success", "eval success", "eval success"}
	root := createTree(res)
	actRes, actAud := root.Eval(fct, []string{})
	assert.Equal(t, res, actRes)
	assert.Equal(t, aud, actAud)
}

func createTree(res Results) *Node {
	// final level
	final := new(Node)
	final.AppendDecisionRule(evalFailure.eval, &Node{})
	final.AppendDecisionRule(evalFailure.eval, &Node{})
	final.AppendResultRule(evalSuccess.eval, res)

	// 3rd level
	third := new(Node)
	third.AppendDecisionRule(evalFailure.eval, &Node{})
	third.AppendDecisionRule(evalFailure.eval, &Node{})
	third.AppendDecisionRule(evalSuccess.eval, final)

	// 2nd level
	second := new(Node)
	second.AppendDecisionRule(evalFailure.eval, &Node{})
	second.AppendDecisionRule(evalFailure.eval, &Node{})
	second.AppendDecisionRule(evalSuccess.eval, third)

	// 1st level
	first := new(Node)
	first.AppendDecisionRule(evalFailure.eval, &Node{})
	first.AppendDecisionRule(evalFailure.eval, &Node{})
	first.AppendDecisionRule(evalSuccess.eval, second)

	// root
	root := new(Node)
	root.AppendDecisionRule(evalFailure.eval, &Node{})
	root.AppendDecisionRule(evalFailure.eval, &Node{})
	root.AppendDecisionRule(evalSuccess.eval, first)
	return root
}
