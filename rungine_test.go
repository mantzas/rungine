package rungine

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEval struct {
	success bool
	err     bool
	def     string
}

func (me mockEval) eval(facts Facts) (bool, string, error) {
	if me.err {
		return false, me.def, errors.New("TEST")
	}
	return me.success, me.def, nil
}

var evalSuccess = mockEval{success: true, def: "eval success", err: false}
var evalFailure = mockEval{success: false, def: "eval failure", err: false}
var evalError = mockEval{success: false, def: "eval failure", err: true}

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
	err := n.AppendDecisionRule(evalFailure.eval, &Node{})
	assert.NoError(t, err)
	err = n.AppendDecisionRule(evalFailure.eval, &Node{})
	assert.NoError(t, err)
	err = n.AppendDecisionRule(evalFailure.eval, &Node{})
	assert.NoError(t, err)
	fct := Facts{}
	res, aud, err := n.Eval(fct, nil)
	assert.Nil(t, res)
	assert.Empty(t, aud)
	assert.NoError(t, err)
}

func TestNode_Eval_Error(t *testing.T) {
	n := new(Node)
	err := n.AppendDecisionRule(evalFailure.eval, &Node{})
	assert.NoError(t, err)
	err = n.AppendDecisionRule(evalError.eval, &Node{})
	assert.NoError(t, err)
	err = n.AppendDecisionRule(evalFailure.eval, &Node{})
	assert.NoError(t, err)
	fct := Facts{}
	aud := []string{}
	res, aud, err := n.Eval(fct, aud)
	assert.Nil(t, res)
	assert.Empty(t, aud)
	assert.Error(t, err)
}

func TestNode_Eval(t *testing.T) {
	fct := Facts{"test": 123}
	res := Results{"res1": 123}
	aud := []string{"eval success", "eval success", "eval success", "eval success", "eval success", "res1:123"}
	root := createTree(res)
	actRes, actAud, err := root.Eval(fct, []string{})
	assert.Equal(t, res, actRes)
	assert.Equal(t, aud, actAud)
	assert.NoError(t, err)
}

func createTree(res Results) *Node {
	// final level
	final := new(Node)
	_ = final.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = final.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = final.AppendResultRule(evalSuccess.eval, res)

	// 3rd level
	third := new(Node)
	_ = third.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = third.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = third.AppendDecisionRule(evalSuccess.eval, final)

	// 2nd level
	second := new(Node)
	_ = second.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = second.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = second.AppendDecisionRule(evalSuccess.eval, third)

	// 1st level
	first := new(Node)
	_ = first.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = first.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = first.AppendDecisionRule(evalSuccess.eval, second)

	// root
	root := new(Node)
	_ = root.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = root.AppendDecisionRule(evalFailure.eval, &Node{})
	_ = root.AppendDecisionRule(evalSuccess.eval, first)
	return root
}

func Test_resultToString(t *testing.T) {
	res := Results{"res1": 123, "res2": 123}
	expected := "res1:123,res2:123"
	actual := resultToString(res)
	assert.Equal(t, actual, expected)
}
