package govaluate

import (
	"fmt"
	"testing"
	"time"

	"github.com/mantzas/rungine"
	"github.com/stretchr/testify/assert"
)

func TestNewExpression(t *testing.T) {
	type args struct {
		exp string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "failure", args: args{exp: ""}, wantErr: true},
		{name: "success", args: args{exp: "1 == 1"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewExpression(tt.args.exp)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestExpression_Eval(t *testing.T) {
	facts := rungine.Facts{
		"foo":                -1,
		"requests_made":      100,
		"requests_succeeded": 80,
		"http_response_body": "service is not ok",
		"val":                "T2",
		"date":               time.Now().Unix(),
		"channel":            make(chan bool),
	}
	type fields struct {
		exp string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name:    "no fact needed",
			fields:  fields{exp: "10 > 0"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "simple expression",
			fields:  fields{exp: "foo > 0"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "complex expression",
			fields:  fields{exp: "(requests_made * requests_succeeded / 100) >= 90"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "static comparisson",
			fields:  fields{exp: "http_response_body == 'service is ok'"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "date comparisson",
			fields:  fields{exp: "date > '2014-01-01 23:59:59'"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "in list",
			fields:  fields{exp: "val IN ('T1','T2','T3')"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "not a boolean expression",
			fields:  fields{exp: "requests_made - requests_succeeded"},
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid fact type",
			fields:  fields{exp: "channel > 0"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := NewExpression(tt.fields.exp)
			assert.NoError(t, err)
			success, actExp, err := exp.Eval(facts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, success)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, success)
			}
			assert.Equal(t, tt.fields.exp, actExp)
		})
	}
}

func TestExpression_Eval_Diamond(t *testing.T) {
	n := createDiamondTree()
	facts := rungine.Facts{
		"color":   "Z",
		"clarity": "I3",
		"cut":     "emerald",
		"weight":  11,
	}
	expResults := rungine.Results{"pricePerCarat": 5500}
	expAudit := []string{
		"color == 'Z'",
		"clarity == 'I3'",
		"cut == 'emerald'",
		"weight >= 9",
		"pricePerCarat:5500",
	}
	res, audit, err := n.Eval(facts, []string{})
	assert.NoError(t, err)
	assert.Equal(t, expResults, res)
	assert.Equal(t, expAudit, audit)
}

var res rungine.Results
var audit []string
var err error

func BenchmarkEvaluateDiamond(b *testing.B) {
	root := createDiamondTree()
	facts := rungine.Facts{
		"color":   "Z",
		"clarity": "I3",
		"cut":     "emerald",
		"weight":  11,
	}
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		res, audit, err = root.Eval(facts, nil)
	}
}

func createDiamondTree() *rungine.Node {

	// price per carat
	res1 := rungine.Results{"pricePerCarat": 1000}
	res2 := rungine.Results{"pricePerCarat": 1500}
	res3 := rungine.Results{"pricePerCarat": 2000}
	res4 := rungine.Results{"pricePerCarat": 2500}
	res5 := rungine.Results{"pricePerCarat": 3000}
	res6 := rungine.Results{"pricePerCarat": 3500}
	res7 := rungine.Results{"pricePerCarat": 4000}
	res8 := rungine.Results{"pricePerCarat": 4500}
	res9 := rungine.Results{"pricePerCarat": 5000}
	res10 := rungine.Results{"pricePerCarat": 5500}

	// final level (CARAT WEIGHT)
	w1, _ := NewExpression("weight < 1")
	w2, _ := NewExpression("(weight >= 1) && (weight < 2)")
	w3, _ := NewExpression("(weight >= 2) && (weight < 3)")
	w4, _ := NewExpression("(weight >= 3) && (weight < 4)")
	w5, _ := NewExpression("(weight >= 4) && (weight < 5)")
	w6, _ := NewExpression("(weight >= 5) && (weight < 6)")
	w7, _ := NewExpression("(weight >= 6) && (weight < 7)")
	w8, _ := NewExpression("(weight >= 7) && (weight < 8)")
	w9, _ := NewExpression("(weight >= 8) && (weight < 9)")
	w10, _ := NewExpression("weight >= 9")

	weightNode := new(rungine.Node)
	_ = weightNode.AppendResultRule(w1.Eval, res1)
	_ = weightNode.AppendResultRule(w2.Eval, res2)
	_ = weightNode.AppendResultRule(w3.Eval, res3)
	_ = weightNode.AppendResultRule(w4.Eval, res4)
	_ = weightNode.AppendResultRule(w5.Eval, res5)
	_ = weightNode.AppendResultRule(w6.Eval, res6)
	_ = weightNode.AppendResultRule(w7.Eval, res7)
	_ = weightNode.AppendResultRule(w8.Eval, res8)
	_ = weightNode.AppendResultRule(w9.Eval, res9)
	_ = weightNode.AppendResultRule(w10.Eval, res10)

	// 2nd level (CUT)
	cutNode := new(rungine.Node)
	cuts := []string{"marquise", "princess", "pear", "oval", "heart"}
	for _, c := range cuts {
		ct, _ := NewExpression(fmt.Sprintf("cut == '%s'", c))
		_ = cutNode.AppendDecisionRule(ct.Eval, &rungine.Node{})
	}
	cgEmerald, _ := NewExpression("cut == 'emerald'")
	_ = cutNode.AppendDecisionRule(cgEmerald.Eval, weightNode)

	// 1st level (CLARITY)
	clarityNode := new(rungine.Node)
	clarities := []string{"FL", "IF", "VVS1", "VVS2", "VS1", "VS2", "SI1", "SI2", "I1", "I2"}
	for _, cla := range clarities {
		clarity, _ := NewExpression(fmt.Sprintf("clarity == '%s'", cla))
		_ = clarityNode.AppendDecisionRule(clarity.Eval, &rungine.Node{})
	}

	clarityI3, _ := NewExpression("clarity == 'I3'")
	_ = clarityNode.AppendDecisionRule(clarityI3.Eval, cutNode)

	// root (COLOR D-to-Z)
	root := new(rungine.Node)
	colors := []string{"D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y"}

	for _, col := range colors {
		color, _ := NewExpression(fmt.Sprintf("color == '%s'", col))
		_ = root.AppendDecisionRule(color.Eval, &rungine.Node{})
	}
	colorZ, _ := NewExpression("color == 'Z'")
	_ = root.AppendDecisionRule(colorZ.Eval, clarityNode)

	return root
}

func evalFailure(facts rungine.Facts) (bool, string, error) {
	return false, "fail", nil
}
