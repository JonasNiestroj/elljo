package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/token"
	"elljo/compiler/js-parser/unistring"
)

func (self *Parser) ParseIdentifier() *ast.Identifier {
	identifier := &ast.Identifier{
		Name:  self.ParsedLiteral,
		Index: self.Index,
	}
	self.NextToken()
	return identifier
}

func (self *Parser) ParseExpression() ast.Expression {
	next := self.ParseAssignmentExpression
	left := next()

	if self.Token == token.COMMA {
		sequence := []ast.Expression{left}
		for {
			if self.Token != token.COMMA {
				break
			}
			self.NextToken()
			sequence = append(sequence, next())
		}
		return &ast.SequenceExpression{
			Sequence: sequence,
		}
	}
	return left
}

func (self *Parser) ParsePrimaryExpression() ast.Expression {
	literal, parsedLiteral := self.Literal, self.ParsedLiteral
	index := self.Index

	switch self.Token {
	case token.IDENTIFIER:
		self.NextToken()

		if len(literal) > 1 {
			isKeyword := token.StringIsKeyword(literal)
			if isKeyword != 0 {
				self.Error("Unexpected reserved keyword")
			}
		}
		return &ast.Identifier{
			Name:  parsedLiteral,
			Index: index,
		}
	case token.NULL:
		self.NextToken()
		return &ast.NullLiteral{
			Index:   index,
			Literal: literal,
		}
	case token.BOOLEAN:
		self.NextToken()
		value := false
		switch parsedLiteral {
		case "true":
			value = true
		case "false":
			value = false
		default:
			self.Error("Illegal boolean literal")
		}
		return &ast.BooleanLiteral{
			Index:   index,
			Literal: literal,
			Value:   value,
		}
	case token.STRING:
		self.NextToken()
		return &ast.StringLiteral{
			Index:   index,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.NUMBER:
		self.NextToken()
		value, err := ParseNumberLiteral(literal)
		if err != nil {
			self.Error(err.Error())
			value = 0
		}
		return &ast.NumberLiteral{
			Index:   index,
			Literal: literal,
			Value:   value,
		}
	case token.SLASH, token.QUOTIENT_ASSIGN:
		return self.ParseRegExpLiteral()
	case token.LEFT_BRACE:
		return self.ParseObjectLiteral()
	case token.LEFT_BRACKET:
		return self.ParseArrayLiteral()
	case token.LEFT_PARENTHESIS:
		self.ExpectToken(token.LEFT_PARENTHESIS)
		expression := self.ParseExpression()
		self.ExpectToken(token.RIGHT_PARENTHESIS)

		return expression
	case token.THIS:
		self.NextToken()
		return &ast.ThisExpression{
			Index: index,
		}
	case token.FUNCTION:
		return self.ParseFunction(false)
	case token.ARROW_FUNCTION:
		return self.ParseArrowFunction()
	case token.Ellipsis:
		return self.ParseSpreadElement()
	}

	self.ErrorUnexpectedToken(self.Token)
	self.NextStatement()
	return &ast.BadExpression{From: index, To: self.Index}
}

func (self *Parser) ParseSpreadElement() *ast.SpreadElement {
	start := self.ExpectToken(token.Ellipsis)
	return &ast.SpreadElement{
		Start:    start,
		Argument: self.ParseAssignmentExpression(),
	}
}

func (self *Parser) ParseRegExpLiteral() *ast.RegExpLiteral {
	offset := self.CharOffset - 1
	if self.Token == token.QUOTIENT_ASSIGN {
		offset -= 1
	}

	index := self.IndexOf(offset)

	pattern, _, err := self.ScanString(offset, false)
	endOffset := self.CharOffset

	if err == nil {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !IsLineTerminator(self.Char) && !IsLineWhiteSpace(self.Char) {
		self.NextToken()

		if self.Token == token.IDENTIFIER {
			flags = self.Literal
			self.NextToken()
			endOffset = self.CharOffset - 1
		}
	} else {
		self.NextToken()
	}

	literal := self.Template[offset:endOffset]

	return &ast.RegExpLiteral{
		Index:   index,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (self *Parser) ParseVariableDeclaration(declarationList *[]*ast.VariableExpression) ast.Expression {
	if self.Token != token.IDENTIFIER {
		index := self.ExpectToken(token.IDENTIFIER)
		self.NextStatement()
		return &ast.BadExpression{From: index, To: self.Index}
	}

	name := self.ParsedLiteral
	index := self.Index
	self.NextToken()
	node := &ast.VariableExpression{
		Name:  name,
		Index: index,
	}

	if declarationList != nil {
		*declarationList = append(*declarationList, node)
	}

	if self.Token == token.ASSIGN {
		self.NextToken()
		node.Initializer = self.ParseAssignmentExpression()
	}

	return node
}

func (self *Parser) ParseVariableDeclarationList(var_ int) []ast.Expression {
	var declarationList []*ast.VariableExpression
	var list []ast.Expression

	for {
		list = append(list, self.ParseVariableDeclaration(&declarationList))
		if self.Token != token.COMMA {
			break
		}
		self.NextToken()
	}

	self.Scope.Declare(&ast.VariableDeclaration{
		Var:  var_,
		List: declarationList,
	})

	return list
}

func (self *Parser) ParseObjectPropertyKey() (literal string, tkn token.Token, exp ast.Expression) {
	index, tkn, literal, parsedLiteral := self.Index, self.Token, self.Literal, self.ParsedLiteral
	var value ast.Expression
	self.NextToken()
	switch tkn {
	case token.IDENTIFIER:
		value = &ast.StringLiteral{
			Index:   index,
			Literal: literal,
			Value:   unistring.String(literal),
		}
	case token.NUMBER:
		num, err := ParseNumberLiteral(literal)
		if err != nil {
			self.Error(err.Error())
		} else {
			value = &ast.NumberLiteral{
				Index:   index,
				Literal: literal,
				Value:   num,
			}
		}
	case token.STRING:
		value = &ast.StringLiteral{
			Index:   index,
			Literal: literal,
			Value:   parsedLiteral,
		}
	case token.Ellipsis:

	default:
		if IsId(tkn) {
			value = &ast.StringLiteral{
				Index:   index,
				Literal: literal,
				Value:   unistring.String(literal),
			}
		}
	}

	return literal, tkn, value
}

func (self *Parser) ParseObjectProperty() ast.Property {
	literal, tkn, value := self.ParseObjectPropertyKey()
	if literal == "get" && self.Token != token.COLON {
		index := self.Index
		_, _, value = self.ParseObjectPropertyKey()
		parameterList := self.ParseFunctionParameterList()

		node := &ast.FunctionLiteral{
			Function:      index,
			ParameterList: parameterList,
		}

		self.ParseFunctionBlock(node)
		return &ast.ObjectProperty{
			Key:   value,
			Kind:  "get",
			Value: node,
		}
	} else if literal == "set" && self.Token != token.COLON {
		index := self.Index
		_, _, value := self.ParseObjectPropertyKey()
		parameterList := self.ParseFunctionParameterList()
		node := &ast.FunctionLiteral{
			Function:      index,
			ParameterList: parameterList,
		}
		self.ParseFunctionBlock(node)
		return &ast.ObjectProperty{
			Key:   value,
			Kind:  "set",
			Value: node,
		}
	} else if tkn == token.Ellipsis {
		index := self.Index
		return &ast.SpreadElement{
			Start:    index,
			Argument: self.ParseAssignmentExpression(),
		}
	}

	self.ExpectToken(token.COLON)
	return &ast.ObjectProperty{
		Key:   value,
		Kind:  "value",
		Value: self.ParseAssignmentExpression(),
	}
}

func (self *Parser) ParseObjectLiteral() ast.Expression {
	var value []ast.Property
	index0 := self.ExpectToken(token.LEFT_BRACE)
	for self.Token != token.RIGHT_BRACE && self.Token != token.EOF {
		property := self.ParseObjectProperty()
		value = append(value, property)
		if self.Token != token.RIGHT_BRACE {
			self.ExpectToken(token.COMMA)
		} else {
			break
		}
	}
	index1 := self.ExpectToken(token.RIGHT_BRACE)

	return &ast.ObjectLiteral{
		LeftBrace:  index0,
		RightBrace: index1,
		Value:      value,
	}
}

func (self *Parser) ParseArrayLiteral() ast.Expression {
	index0 := self.ExpectToken(token.LEFT_BRACKET)
	var value []ast.Expression
	for self.Token != token.RIGHT_BRACKET && self.Token != token.EOF {
		if self.Token == token.COMMA {
			self.NextToken()
			value = append(value, nil)
		}

		value = append(value, self.ParseAssignmentExpression())

		if self.Token != token.RIGHT_BRACKET {
			self.ExpectToken(token.COMMA)
		}
	}
	self.ExpectToken(token.RIGHT_BRACKET)

	return &ast.ArrayLiteral{
		LeftBracket:  index0,
		RightBracket: self.Index,
		Value:        value,
	}
}

func (self *Parser) ParseArgumentList() (argumentList []ast.Expression, index0, index1 int) {
	index0 = self.ExpectToken(token.LEFT_PARENTHESIS)
	if self.Token != token.RIGHT_PARENTHESIS {
		for {
			argumentList = append(argumentList, self.ParseAssignmentExpression())
			if self.Token != token.COMMA {
				break
			}
			self.NextToken()
		}
	}
	index1 = self.ExpectToken(token.RIGHT_PARENTHESIS)
	return
}

func (self *Parser) ParseCallExpression(left ast.Expression) ast.Expression {
	argumentList, index0, index1 := self.ParseArgumentList()
	return &ast.CallExpression{
		Callee:           left,
		LeftParanthesis:  index0,
		ArgumentList:     argumentList,
		RightParenthesis: index1,
	}
}

func (self *Parser) ParseDotMember(left ast.Expression) ast.Expression {
	period := self.ExpectToken(token.PERIOD)

	literal := self.ParsedLiteral
	index := self.Index

	if self.Token != token.IDENTIFIER && !IsId(self.Token) {
		self.ExpectToken(token.IDENTIFIER)
		self.NextStatement()
		return &ast.BadExpression{From: period, To: self.Index}
	}

	self.NextToken()

	return &ast.DotExpression{
		Left: left,
		Identifier: ast.Identifier{
			Index: index,
			Name:  literal,
		},
	}
}

func (self *Parser) ParseBracketMember(left ast.Expression) ast.Expression {
	index0 := self.ExpectToken(token.LEFT_BRACKET)
	member := self.ParseExpression()
	index1 := self.ExpectToken(token.RIGHT_BRACKET)
	return &ast.BracketExpression{
		LeftBracket:  index0,
		Left:         left,
		Member:       member,
		RightBracket: index1,
	}
}

func (self *Parser) ParseNewExpression() ast.Expression {
	index := self.ExpectToken(token.NEW)
	if self.Token == token.PERIOD {
		self.NextToken()
		prop := self.ParseIdentifier()
		if prop.Name == "target" {
			if !self.Scope.InFunction {
				self.Error("new.target expression is not allowed here")
			}
			return &ast.MetaProperty{
				Meta: &ast.Identifier{
					Name:  unistring.String(token.NEW.ToString()),
					Index: index,
				},
				Property: prop,
			}
		}
		self.ErrorUnexpectedToken(token.IDENTIFIER)
	}

	callee := self.ParseLeftHandSideExpression()
	node := &ast.NewExpression{
		New:    index,
		Callee: callee,
	}
	if self.Token == token.LEFT_PARENTHESIS {
		argumentList, index0, index1 := self.ParseArgumentList()
		node.ArgumentList = argumentList
		node.LeftParenthesis = index0
		node.RightParenthesis = index1
	}
	return node
}

func (self *Parser) ParseLeftHandSideExpression() ast.Expression {
	var left ast.Expression
	if self.Token == token.NEW {
		left = self.ParseNewExpression()
	} else {
		left = self.ParsePrimaryExpression()
	}

	for {
		if self.Token == token.PERIOD {
			left = self.ParseDotMember(left)
		} else if self.Token == token.LEFT_BRACKET {
			left = self.ParseBracketMember(left)
		} else {
			break
		}
	}
	return left
}

func (self *Parser) ParseLeftHandSideExpressionAllowCall() ast.Expression {
	allowIn := self.Scope.AllowIn
	self.Scope.AllowIn = true
	defer func() {
		self.Scope.AllowIn = allowIn
	}()
	if self.Token == token.LEFT_PARENTHESIS && self.LastToken == token.ASSIGN {
		self.NextToken()
		return self.ParseArrowFunction()
	}

	var left ast.Expression
	if self.Token == token.NEW {
		left = self.ParseNewExpression()
	} else {
		left = self.ParsePrimaryExpression()
	}

	for {
		if self.Token == token.PERIOD {
			left = self.ParseDotMember(left)
		} else if self.Token == token.LEFT_BRACKET {
			left = self.ParseBracketMember(left)
		} else if self.Token == token.LEFT_PARENTHESIS {
			left = self.ParseCallExpression(left)
		} else {
			break
		}
	}
	return left
}

func (self *Parser) ParsePostfixExpression() ast.Expression {
	operand := self.ParseLeftHandSideExpressionAllowCall()

	switch self.Token {
	case token.INCREMENT, token.DECREMENT:
		if self.ImplicitSemicolon {
			break
		}
		tkn := self.Token
		index := self.Index
		self.NextToken()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			self.Error("Invalid left-hand side in assigment")
			self.NextStatement()
			return &ast.BadExpression{From: index, To: self.Index}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Index:    index,
			Operand:  operand,
			Postfix:  true,
		}
	}
	return operand
}

func (self *Parser) ParseUnaryExpression() ast.Expression {
	switch self.Token {
	case token.PLUS, token.MINUS, token.NOT, token.BITWISE_NOT:
		fallthrough
	case token.DELETE, token.VOID, token.TYPEOF:
		tkn := self.Token
		index := self.Index
		self.NextToken()
		return &ast.UnaryExpression{
			Operator: tkn,
			Index:    index,
			Operand:  self.ParseUnaryExpression(),
		}
	case token.INCREMENT, token.DECREMENT:
		tkn := self.Token
		index := self.Index
		self.NextToken()
		operand := self.ParseUnaryExpression()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			self.Error("Invalid left-hand side in assignment")
			self.NextStatement()
			return &ast.BadExpression{From: index, To: self.Index}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Index:    index,
			Operand:  operand,
		}
	}

	return self.ParsePostfixExpression()
}

func (self *Parser) ParseMultiplicativeExpression() ast.Expression {
	next := self.ParseUnaryExpression
	left := next()

	for self.Token == token.MULTIPLY || self.Token == token.SLASH || self.Token == token.REMAINDER {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (self *Parser) ParseAdditiveExpression() ast.Expression {
	next := self.ParseMultiplicativeExpression
	left := next()

	for self.Token == token.PLUS || self.Token == token.MINUS {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseShiftExpression() ast.Expression {
	next := self.ParseAdditiveExpression
	left := next()

	for self.Token == token.SHIFT_LEFT || self.Token == token.SHIFT_RIGHT || self.Token == token.UNSIGNED_SHIFT_RIGHT {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseRelationalExpression() ast.Expression {
	next := self.ParseShiftExpression
	left := next()

	allowIn := self.Scope.AllowIn
	self.Scope.AllowIn = true
	defer func() {
		self.Scope.AllowIn = allowIn
	}()

	switch self.Token {
	case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUAL:
		tkn := self.Token
		self.NextToken()
		return &ast.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      self.ParseRelationalExpression(),
			Comparison: true,
		}
	case token.INSTANCEOF:
		tkn := self.Token
		self.NextToken()
		return &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    self.ParseRelationalExpression(),
		}
	case token.IN:
		if !allowIn {
			return left
		}
		tkn := self.Token
		self.NextToken()
		return &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    self.ParseRelationalExpression(),
		}
	}

	return left
}

func (self *Parser) ParseEqualityExpression() ast.Expression {
	next := self.ParseRelationalExpression
	left := next()

	for self.Token == token.EQUAL || self.Token == token.NOT_EQUAL || self.Token == token.STRICT_EQUAL || self.Token == token.STRICT_NOT_EQUAL {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      next(),
			Comparison: true,
		}
	}

	return left
}

func (self *Parser) ParseBitwiseAndExpression() ast.Expression {
	next := self.ParseEqualityExpression
	left := next()

	for self.Token == token.AND {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseBitwiseExclusiveOrExpression() ast.Expression {
	next := self.ParseBitwiseAndExpression
	left := next()

	for self.Token == token.EXCLUSIVE_OR {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseBitwiseOrExpression() ast.Expression {
	next := self.ParseBitwiseExclusiveOrExpression
	left := next()

	for self.Token == token.OR {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseLogicalAndExpression() ast.Expression {
	next := self.ParseBitwiseOrExpression
	left := next()

	for self.Token == token.LOGICAL_AND {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseLogicalOrExpression() ast.Expression {
	next := self.ParseLogicalAndExpression
	left := next()

	for self.Token == token.LOGICAL_OR {
		tkn := self.Token
		self.NextToken()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}
	return left
}

func (self *Parser) ParseConditionalExpression() ast.Expression {
	left := self.ParseLogicalOrExpression()

	if self.Token == token.QUESTION_MARK {
		self.NextToken()
		consequent := self.ParseAssignmentExpression()
		self.ExpectToken(token.COLON)
		return &ast.ConditionalExpression{
			Test:       left,
			Consequent: consequent,
			Alternate:  self.ParseAssignmentExpression(),
		}
	}

	return left
}

func (self *Parser) ParseAssignmentExpression() ast.Expression {
	left := self.ParseConditionalExpression()
	var operator token.Token
	if self.Token == token.ARROW_FUNCTION {
		return self.ParseArrowFunction()
	}

	switch self.Token {
	case token.ASSIGN:
		operator = self.Token
	case token.ADD_ASSIGN:
		operator = token.PLUS
	case token.SUBTRACT_ASSIGN:
		operator = token.MINUS
	case token.MULTIPLY_ASSIGN:
		operator = token.MULTIPLY
	case token.QUOTIENT_ASSIGN:
		operator = token.SLASH
	case token.REMAINDER_ASSIGN:
		operator = token.REMAINDER
	case token.AND_ASSIGN:
		operator = token.AND
	case token.OR_ASSIGN:
		operator = token.OR
	case token.EXCLUSIVE_OR_ASSIGN:
		operator = token.EXCLUSIVE_OR
	case token.SHIFT_LEFT_ASSIGN:
		operator = token.SHIFT_LEFT
	case token.SHIFT_RIGHT_ASSIGN:
		operator = token.SHIFT_RIGHT
	case token.UNSIGNED_SHIFT_RIGHT_ASSIGN:
		operator = token.UNSIGNED_SHIFT_RIGHT
	case token.LEFT_PARENTHESIS:
		operator = token.ARROW_FUNCTION
	}

	if operator != 0 {
		index := self.Index
		self.NextToken()
		switch left.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			self.Error("Invalid left-hand side in assignment")
			self.NextStatement()
			return &ast.BadExpression{From: index, To: self.Index}
		}

		right := self.ParseAssignmentExpression()

		return &ast.AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}
