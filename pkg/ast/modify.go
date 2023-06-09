package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			node.Statements[i] = Modify(stmt, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Right = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Index = Modify(node.Index, modifier).(Expression)
	case *IfExpression:
		node.Condition = Modify(node.Condition, modifier).(Expression)
		node.Consequence = Modify(node.Consequence, modifier).(*BlockStatement)
		if node.Alternative != nil {
			node.Alternative = Modify(node.Alternative, modifier).(*BlockStatement)
		}
	case *BlockStatement:
		for i, stmt := range node.Statements {
			node.Statements[i] = Modify(stmt, modifier).(Statement)
		}
	case *ReturnStatement:
		node.ReturnValue = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i, param := range node.Parameters {
			node.Parameters[i] = Modify(param, modifier).(*Identifier)
		}
		node.Body = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, elem := range node.Elements {
			node.Elements[i] = Modify(elem, modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			newKey := Modify(key, modifier).(Expression)
			newValue := Modify(value, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Pairs = newPairs
	}
	return modifier(node)
}
