package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/token"
)

func (self *Parser) ParseBlockStatement() *ast.BlockStatement {
	return &ast.BlockStatement{
		LeftBrace:  self.ExpectToken(token.LEFT_BRACE),
		List:       self.ParseStatementList(),
		RightBrace: self.ExpectToken(token.RIGHT_BRACE),
	}
}

func (self *Parser) ParseEmptyStatement() ast.Statement {
	return &ast.EmptyStatement{Semicolon: self.ExpectToken(token.SEMICOLON)}
}

func (self *Parser) ParseStatementList() (statements []ast.Statement) {
	for self.Token != token.RIGHT_BRACE && self.Token != token.EOF {
		statement := self.ParseStatement()
		statements = append(statements, statement)
	}
	return
}

func (self *Parser) ParseStatement() ast.Statement {
	if self.Token == token.EOF {
		self.ErrorUnexpectedToken(self.Token)
		return &ast.BadStatement{From: self.Index, To: self.Index + 1}
	}

	switch self.Token {
	case token.BREAK:
		return self.ParseBreakStatement()
	case token.CONTINUE:
		return self.ParseContinueStatement()
	case token.DEBUGGER:
		return self.ParseDebuggerStatement()
	case token.DO:
		return self.ParseDoWhileStatement()
	case token.FOR:
		return self.ParseForOrForInStatement()
	case token.FUNCTION:
		self.ParseFunction(true)
		return &ast.EmptyStatement{}
	case token.IF:
		return self.ParseIfStatement()
	case token.IMPORT:
		return self.ParseImportStatement()
	case token.LEFT_BRACE:
		return self.ParseBlockStatement()
	case token.RETURN:
		return self.ParseReturnStatement()
	case token.SEMICOLON:
		return self.ParseEmptyStatement()
	case token.SWITCH:
		return self.ParseSwitchStatement()
	case token.THROW:
		return self.ParseThrowStatement()
	case token.TRY:
		return self.ParseTryStatement()
	case token.WHILE:
		return self.ParseWhileStatement()
	case token.WITH:
		return self.ParseWithStatement()
	case token.VAR, token.LET, token.CONST:
		return self.ParseVariableStatement()
	}

	expression := self.ParseExpression()

	if identifier, isIdentifier := expression.(*ast.Identifier); isIdentifier && self.Token == token.COLON {
		colon := self.Index
		self.NextToken()
		label := identifier.Name
		for _, value := range self.Scope.Labels {
			if label == value {
				self.Error("Label '%s' already exists", label)
			}
		}
		self.Scope.Labels = append(self.Scope.Labels, label)
		statement := self.ParseStatement()
		self.Scope.Labels = self.Scope.Labels[:len(self.Scope.Labels)-1]
		return &ast.LabelledStatement{
			Label:     identifier,
			Colon:     colon,
			Statement: statement,
		}
	}
	self.OptionalSemicolon()

	return &ast.ExpressionStatement{
		Expression: expression,
	}
}

func (self *Parser) ParseTryStatement() ast.Statement {
	node := &ast.TryStatement{
		Try:  self.ExpectToken(token.TRY),
		Body: self.ParseBlockStatement(),
	}

	if self.Token == token.CATCH {
		catch := self.Index
		self.NextToken()
		self.ExpectToken(token.LEFT_PARENTHESIS)
		if self.Token != token.IDENTIFIER {
			self.ExpectToken(token.IDENTIFIER)
			self.NextStatement()
			return &ast.BadStatement{From: catch, To: self.Index}
		} else {
			identifier := self.ParseIdentifier()
			self.ExpectToken(token.RIGHT_PARENTHESIS)
			node.Catch = &ast.CatchStatement{
				Catch:     catch,
				Parameter: identifier,
				Body:      self.ParseBlockStatement(),
			}
		}
	}
	if self.Token == token.FINALLY {
		self.NextToken()
		node.Finally = self.ParseBlockStatement()
	}

	if node.Catch == nil && node.Finally == nil {
		self.Error("Missing catch or finally after try")
		return &ast.BadStatement{From: node.Try, To: node.Body.Index1()}
	}

	return node
}

func (self *Parser) ParseImportStatement() ast.Statement {
	from := self.Index
	self.ExpectToken(token.IMPORT)
	name := self.Literal
	self.ExpectToken(token.IDENTIFIER)
	self.ExpectToken(token.IMPORTFROM)
	source := self.Literal
	self.ExpectToken(token.STRING)

	return &ast.ImportStatement{
		From:   from,
		To:     self.Index,
		Name:   name,
		Source: source,
	}
}

func (self *Parser) ParseFunctionParameterList() *ast.ParameterList {
	opening := self.ExpectToken(token.LEFT_PARENTHESIS)
	var list []*ast.Identifier
	for self.Token != token.RIGHT_PARENTHESIS && self.Token != token.EOF {
		if self.Token != token.IDENTIFIER {
			self.ExpectToken(token.IDENTIFIER)
		} else {
			list = append(list, self.ParseIdentifier())
		}
		if self.Token != token.RIGHT_PARENTHESIS {
			self.ExpectToken(token.COMMA)
		}
	}
	closing := self.ExpectToken(token.RIGHT_PARENTHESIS)

	return &ast.ParameterList{
		Opening: opening,
		List:    list,
		Closing: closing,
	}
}

func (self *Parser) ParseParameterList() (list []string) {
	for self.Token != token.EOF {
		if self.Token != token.IDENTIFIER {
			self.ExpectToken(token.IDENTIFIER)
		}
		list = append(list, self.Literal)
		self.NextToken()
		if self.Token != token.EOF {
			self.ExpectToken(token.COMMA)
		}
	}
	return
}

func (self *Parser) ParseFunction(declaration bool) *ast.FunctionLiteral {
	node := &ast.FunctionLiteral{
		Function: self.ExpectToken(token.FUNCTION),
	}

	var name *ast.Identifier
	if self.Token == token.IDENTIFIER {
		name = self.ParseIdentifier()
		if declaration {
			self.Scope.Declare(&ast.FunctionDeclaration{
				Function: node,
			})
		}
	} else if declaration {
		self.ExpectToken(token.IDENTIFIER)
	}
	node.Name = name
	node.ParameterList = self.ParseFunctionParameterList()
	self.ParseFunctionBlock(node)
	node.Source = self.Slice(node.Index0(), node.Index1())
	return node
}

func (self *Parser) ParseArrowFunction() *ast.FunctionLiteral {
	node := &ast.FunctionLiteral{}
	if self.Token != token.RIGHT_PARENTHESIS {
		node.ParameterList = self.ParseFunctionParameterList()
	}
	self.NextToken()
	self.ExpectToken(token.ARROW_FUNCTION)
	self.ParseFunctionBlock(node)
	node.Source = self.Slice(node.Index0(), node.Index1())
	return node
}

func (self *Parser) ParseFunctionBlock(node *ast.FunctionLiteral) {
	self.OpenScope()
	inFunction := self.Scope.InFunction
	self.Scope.InFunction = true
	defer func() {
		self.Scope.InFunction = inFunction
		self.CloseScope()
	}()
	node.Body = self.ParseBlockStatement()
	node.DeclarationList = self.Scope.DeclarationList
}

func (self *Parser) ParseDebuggerStatement() ast.Statement {
	index := self.ExpectToken(token.DEBUGGER)

	node := &ast.DebuggerStatement{
		Debugger: index,
	}

	self.Semicolon()
	return node
}

func (self *Parser) ParseReturnStatement() ast.Statement {
	index := self.ExpectToken(token.RETURN)

	if !self.Scope.InFunction {
		self.Error("Illegal return statement")
		self.NextStatement()
		return &ast.BadStatement{From: index, To: self.Index}
	}

	node := &ast.ReturnStatement{
		Return: index,
	}

	if !self.ImplicitSemicolon && self.Token != token.SEMICOLON && self.Token != token.RIGHT_BRACE && self.Token != token.EOF {
		node.Argument = self.ParseExpression()
	}

	self.Semicolon()
	return node
}

func (self *Parser) ParseThrowStatement() ast.Statement {
	index := self.ExpectToken(token.THROW)

	if self.ImplicitSemicolon {
		if self.Char == -1 {
			self.Error("Unexpected end of input")
		} else {
			self.Error("Illegal newline after throw")
		}
		self.NextStatement()
		return &ast.BadStatement{From: index, To: self.Index}
	}
	node := &ast.ThrowStatement{
		Argument: self.ParseExpression(),
	}

	self.Semicolon()

	return node
}

func (self *Parser) ParseSwitchStatement() ast.Statement {
	self.ExpectToken(token.SWITCH)
	self.ExpectToken(token.LEFT_PARENTHESIS)
	node := &ast.SwitchStatement{
		Discriminant: self.ParseExpression(),
		Default:      -1,
	}
	self.ExpectToken(token.RIGHT_PARENTHESIS)
	self.ExpectToken(token.LEFT_BRACE)
	inSwitch := self.Scope.InSwitch
	self.Scope.InSwitch = true
	defer func() {
		self.Scope.InSwitch = inSwitch
	}()

	for index := 0; self.Token != token.EOF; index++ {
		if self.Token == token.RIGHT_BRACE {
			self.NextToken()
			break
		}

		clause := self.ParseCaseStatement()
		if clause.Test == nil {
			if node.Default != -1 {
				self.Error("Already saw a default in switch")
			}
			node.Default = index
		}
		node.Body = append(node.Body, clause)
	}
	return node
}

func (self *Parser) ParseWithStatement() ast.Statement {
	self.ExpectToken(token.WITH)
	self.ExpectToken(token.LEFT_PARENTHESIS)
	node := &ast.WithStatement{
		Object: self.ParseExpression(),
	}
	self.ExpectToken(token.RIGHT_PARENTHESIS)
	node.Body = self.ParseStatement()

	return node
}

func (self *Parser) ParseCaseStatement() *ast.CaseStatement {
	node := &ast.CaseStatement{
		Case: self.Index,
	}
	if self.Token == token.DEFAULT {
		self.NextToken()
	} else {
		self.ExpectToken(token.CASE)
		node.Test = self.ParseExpression()
	}
	self.ExpectToken(token.COLON)

	for {
		if self.Token == token.EOF || self.Token == token.RIGHT_BRACE || self.Token == token.CASE || self.Token == token.DEFAULT {
			break
		}
		node.Consequent = append(node.Consequent, self.ParseStatement())
	}
	return node
}

func (self *Parser) ParseIterationStatement() ast.Statement {
	inIteration := self.Scope.InIteration
	self.Scope.InIteration = true
	defer func() {
		self.Scope.InIteration = inIteration
	}()
	return self.ParseStatement()
}

func (self *Parser) ParseForIn(index int, into ast.Expression) *ast.ForInStatement {
	source := self.ParseExpression()
	self.ExpectToken(token.RIGHT_PARENTHESIS)

	return &ast.ForInStatement{
		For:    index,
		Into:   into,
		Source: source,
		Body:   self.ParseIterationStatement(),
	}
}

func (self *Parser) ParseForOf(index int, into ast.Expression) *ast.ForOfStatement {
	source := self.ParseExpression()
	self.ExpectToken(token.RIGHT_PARENTHESIS)

	return &ast.ForOfStatement{
		For:    index,
		Into:   into,
		Source: source,
		Body:   self.ParseIterationStatement(),
	}
}

func (self *Parser) ParseFor(index int, initializer ast.Expression) *ast.ForStatement {
	var test, update ast.Expression
	if self.Token != token.SEMICOLON {
		test = self.ParseExpression()
	}
	self.ExpectToken(token.SEMICOLON)

	if self.Token != token.RIGHT_PARENTHESIS {
		update = self.ParseExpression()
	}
	self.ExpectToken(token.RIGHT_PARENTHESIS)

	return &ast.ForStatement{
		For:         index,
		Initializer: initializer,
		Test:        test,
		Update:      update,
		Body:        self.ParseIterationStatement(),
	}
}

func (self *Parser) ParseForOrForInStatement() ast.Statement {
	index := self.ExpectToken(token.FOR)
	self.ExpectToken(token.LEFT_PARENTHESIS)

	var left []ast.Expression

	forIn := false
	forOf := false
	if self.Token != token.SEMICOLON {
		allowIn := self.Scope.AllowIn
		self.Scope.AllowIn = false
		if self.Token == token.VAR || self.Token == token.CONST || self.Token == token.LET {
			var_ := self.Index
			self.NextToken()
			list := self.ParseVariableDeclarationList(var_)
			if len(list) == 1 {
				if self.Token == token.IN {
					self.NextToken()
					forIn = true
				} else if self.Token == token.IDENTIFIER {
					if self.Literal == "of" {
						self.NextToken()
						forOf = true
					}
				}
			}
			left = list
		} else {
			left = append(left, self.ParseExpression())
			if self.Token == token.IN {
				self.NextToken()
				forIn = true
			} else if self.Token == token.IDENTIFIER {
				if self.Literal == "of" {
					self.NextToken()
					forOf = true
				}
			}
		}
		self.Scope.AllowIn = allowIn
	}

	if forIn || forOf {
		switch left[0].(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression, *ast.VariableExpression:
		default:
			self.Error("Invalid left-hand side in for-in or for-of")
			self.NextStatement()
			return &ast.BadStatement{From: index, To: self.Index}
		}
		if forIn {
			return self.ParseFor(index, left[0])
		}
		return self.ParseForOf(index, left[0])
	}

	self.ExpectToken(token.SEMICOLON)
	return self.ParseFor(index, &ast.SequenceExpression{Sequence: left})
}

func (self *Parser) ParseVariableStatement() *ast.VariableStatement {
	index := self.ExpectToken(token.VAR, token.LET, token.CONST)
	list := self.ParseVariableDeclarationList(index)
	self.Semicolon()
	return &ast.VariableStatement{
		Var:  index,
		List: list,
	}
}

func (self *Parser) ParseDoWhileStatement() ast.Statement {
	inIteration := self.Scope.InIteration
	self.Scope.InIteration = true
	defer func() {
		self.Scope.InIteration = inIteration
	}()

	self.ExpectToken(token.DO)
	node := &ast.DoWhileStatement{}
	if self.Token == token.LEFT_BRACE {
		node.Body = self.ParseBlockStatement()
	} else {
		node.Body = self.ParseStatement()
	}

	self.ExpectToken(token.WHILE)
	self.ExpectToken(token.LEFT_PARENTHESIS)
	node.Test = self.ParseExpression()
	self.ExpectToken(token.RIGHT_PARENTHESIS)
	if self.Token == token.SEMICOLON {
		self.NextToken()
	}
	return node
}

func (self *Parser) ParseWhileStatement() ast.Statement {
	self.ExpectToken(token.WHILE)
	self.ExpectToken(token.LEFT_PARENTHESIS)
	node := &ast.WhileStatement{
		Test: self.ParseExpression(),
	}
	self.ExpectToken(token.RIGHT_PARENTHESIS)
	node.Body = self.ParseIterationStatement()

	return node
}

func (self *Parser) ParseIfStatement() ast.Statement {
	self.ExpectToken(token.IF)
	self.ExpectToken(token.LEFT_PARENTHESIS)
	node := &ast.IfStatement{
		Test: self.ParseExpression(),
	}
	self.ExpectToken(token.RIGHT_PARENTHESIS)

	if self.Token == token.LEFT_BRACE {
		node.Consequent = self.ParseBlockStatement()
	} else {
		node.Consequent = self.ParseStatement()
	}

	if self.Token == token.ELSE {
		self.NextToken()
		node.Alternate = self.ParseStatement()
	}
	return node
}

func (self *Parser) ParseSourceElement() ast.Statement {
	return self.ParseStatement()
}

func (self *Parser) ParseSourceElements() []ast.Statement {
	body := []ast.Statement(nil)

	for {
		if self.Token != token.STRING {
			break
		}
		body = append(body, self.ParseSourceElement())
	}
	for self.Token != token.EOF && self.Index < len(self.Template) {
		body = append(body, self.ParseSourceElement())
	}
	return body
}

func (self *Parser) ParseProgram() *ast.Program {
	self.OpenScope()
	defer self.CloseScope()
	prg := &ast.Program{
		Body:            self.ParseSourceElements(),
		DeclarationList: self.Scope.DeclarationList,
	}
	return prg
}

func (self *Parser) ParseBreakStatement() ast.Statement {
	index := self.ExpectToken(token.BREAK)
	semicolon := self.ImplicitSemicolon
	if self.Token == token.SEMICOLON {
		semicolon = true
		self.NextToken()
	}

	if semicolon || self.Token == token.RIGHT_BRACE {
		self.ImplicitSemicolon = false
		if !self.Scope.InIteration && !self.Scope.InSwitch {
			return self.CreateBadStatement(index, "Illegal break statement")
		}
		return &ast.BranchStatement{
			Index: index,
			Token: token.BREAK,
		}
	}
	if self.Token == token.IDENTIFIER {
		identifier := self.ParseIdentifier()
		if !self.Scope.HasLabel(identifier.Name) {
			self.Error("Undefined label '%s'", identifier.Name)
			return &ast.BadStatement{From: index, To: identifier.Index1()}
		}
		self.Semicolon()
		return &ast.BranchStatement{
			Index: index,
			Token: token.BREAK,
			Label: identifier,
		}
	}
	self.ExpectToken(token.IDENTIFIER)
	return self.CreateBadStatement(index, "Illegal break statement")
}

func (self *Parser) CreateBadStatement(index int, message string) ast.Statement {
	self.Error(message)
	self.NextStatement()
	return &ast.BadStatement{From: index, To: self.Index}
}

func (self *Parser) ParseContinueStatement() ast.Statement {
	index := self.ExpectToken(token.CONTINUE)
	semicolon := self.ImplicitSemicolon
	if self.Token == token.SEMICOLON {
		semicolon = true
		self.NextToken()
	}

	if semicolon || self.Token == token.RIGHT_BRACE {
		self.ImplicitSemicolon = true
		if !self.Scope.InIteration {
			return self.CreateBadStatement(index, "Illegal continue statement")
		}
		return &ast.BranchStatement{
			Index: index,
			Token: token.CONTINUE,
		}
	}

	if self.Token == token.IDENTIFIER {
		identifier := self.ParseIdentifier()
		if !self.Scope.HasLabel(identifier.Name) {
			self.Error("Undefined label '%s'", identifier.Name)
			return &ast.BadStatement{From: index, To: identifier.Index1()}
		}
		if !self.Scope.InIteration {
			return self.CreateBadStatement(index, "Illegal continue statement")
		}
		self.Semicolon()
		return &ast.BranchStatement{
			Index: index,
			Token: token.CONTINUE,
			Label: identifier,
		}
	}
	self.ExpectToken(token.IDENTIFIER)
	return self.CreateBadStatement(index, "Illegal continue statement")
}

func (self *Parser) NextStatement() {
	for {
		switch self.Token {
		case token.BREAK, token.CONTINUE, token.FOR, token.IF, token.RETURN,
			token.SWITCH, token.VAR, token.DO, token.TRY, token.WITH, token.WHILE,
			token.THROW, token.CATCH, token.FINALLY, token.IMPORT:
			if self.Index == self.Recover.Index && self.Recover.Count < 10 {
				self.Recover.Count++
				return
			}
			if self.Index > self.Recover.Index {
				self.Recover.Index = self.Index
				self.Recover.Count = 0
				return
			}
		case token.EOF:
			return
		}
		self.NextToken()
	}
}
