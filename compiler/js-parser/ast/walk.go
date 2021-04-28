package ast

func Walk(node Node, nodeCallback func(node Node)) {
	switch node := node.(type) {
	case *Program:
		for _, b := range node.Body {
			nodeCallback(node)
			Walk(b, nodeCallback)
		}
	case *ExpressionStatement:
		nodeCallback(node)
		Walk(node.Expression, nodeCallback)
	case *DotExpression:
		nodeCallback(node)
		Walk(node.Left, nodeCallback)
	case *NewExpression:
		nodeCallback(node)
		Walk(node.Callee, nodeCallback)
	case *VariableExpression:
		nodeCallback(node)
	case *AssignExpression:
		nodeCallback(node)
	case *CallExpression:
		nodeCallback(node)
		Walk(node.Callee, nodeCallback)
		for _, e := range node.ArgumentList {
			Walk(e, nodeCallback)
		}
	case *FunctionLiteral:
		nodeCallback(node)
		Walk(node.Body, nodeCallback)
	case *Identifier:
		nodeCallback(node)
	case *ForStatement:
		nodeCallback(node)
		Walk(node.Body, nodeCallback)
	case *BracketExpression:
		nodeCallback(node)
		Walk(node.Member, nodeCallback)
	case *BlockStatement:
		nodeCallback(node)
		for _, b := range node.List {
			Walk(b, nodeCallback)
		}
	case *VariableStatement:
		nodeCallback(node)
		for _, b := range node.List {
			Walk(b, nodeCallback)
		}
	case *BinaryExpression:
		nodeCallback(node)
	case *SequenceExpression:
		nodeCallback(node)
		for _, b := range node.Sequence {
			Walk(b, nodeCallback)
		}
	case *UnaryExpression:
		nodeCallback(node)
		Walk(node.Operand, nodeCallback)
	case *ImportStatement:
		nodeCallback(node)
	case *IfStatement:
		nodeCallback(node)
		Walk(node.Consequent, nodeCallback)
	}
}