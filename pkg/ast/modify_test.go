package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression {
		return &IntegerLiteral{Value: 1}
	}
	two := func() Expression {
		return &IntegerLiteral{Value: 2}
	}

	trunOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{one(), two()},
		{&Program{Statements: []Statement{&ExpressionStatement{Expression: one()}}}, &Program{Statements: []Statement{&ExpressionStatement{Expression: two()}}}},
		{&InfixExpression{Left: one(), Operator: "+", Right: one()}, &InfixExpression{Left: two(), Operator: "+", Right: two()}},
		{&PrefixExpression{Operator: "-", Right: one()}, &PrefixExpression{Operator: "-", Right: two()}},
		{&IndexExpression{Left: one(), Index: one()}, &IndexExpression{Left: two(), Index: two()}},
		{&IfExpression{Condition: one(), Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}}, Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}}}, &IfExpression{Condition: two(), Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}}, Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}}}},
		{&ReturnStatement{ReturnValue: one()}, &ReturnStatement{ReturnValue: two()}},
		{&LetStatement{Value: one()}, &LetStatement{Value: two()}},
		{&FunctionLiteral{Parameters: []*Identifier{}, Body: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}}}, &FunctionLiteral{Parameters: []*Identifier{}, Body: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}}}},
		{&ArrayLiteral{Elements: []Expression{one()}}, &ArrayLiteral{Elements: []Expression{two()}}},
	}

	for _, tt := range tests {
		modified := Modify(tt.input, trunOneIntoTwo)
		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. want=%q, got=%q", tt.expected, modified)
		}
	}

	// DeepEqual doesn't work for map
	hashLiteral := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, trunOneIntoTwo)

	for k, v := range hashLiteral.Pairs {
		k := k.(*IntegerLiteral)
		if k.Value != 2 {
			t.Errorf("key not modified. want=2, got=%d", k.Value)
		}
		v := v.(*IntegerLiteral)
		if v.Value != 2 {
			t.Errorf("value not modified. want=2, got=%d", v.Value)
		}
	}
}
