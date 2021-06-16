package ast

import (
	"elljo/compiler/js-parser/token"
	"elljo/compiler/js-parser/unistring"
)

type Node interface {
	Index0() int
	Index1() int
}

type (
	Expression interface {
		Node
	}

	ArrayLiteral struct {
		LeftBracket  int
		RightBracket int
		Value        []Expression
	}

	AssignExpression struct {
		Operator token.Token
		Left     Expression
		Right    Expression
	}

	BadExpression struct {
		From int
		To   int
	}

	BinaryExpression struct {
		Operator   token.Token
		Left       Expression
		Right      Expression
		Comparison bool
	}

	BooleanLiteral struct {
		Index   int
		Literal string
		Value   bool
	}

	BracketExpression struct {
		Left         Expression
		Member       Expression
		LeftBracket  int
		RightBracket int
	}

	CallExpression struct {
		Callee           Expression
		LeftParanthesis  int
		ArgumentList     []Expression
		RightParenthesis int
	}

	ConditionalExpression struct {
		Test       Expression
		Consequent Expression
		Alternate  Expression
	}

	DotExpression struct {
		Left       Expression
		Identifier Identifier
	}

	FunctionLiteral struct {
		Function        int
		Name            *Identifier
		ParameterList   *ParameterList
		Body            Statement
		Source          string
		DeclarationList []Declaration
	}

	Identifier struct {
		Name  unistring.String
		Index int
	}

	NewExpression struct {
		New              int
		Callee           Expression
		LeftParenthesis  int
		ArgumentList     []Expression
		RightParenthesis int
	}

	NullLiteral struct {
		Index   int
		Literal string
	}

	NumberLiteral struct {
		Index   int
		Literal string
		Value   interface{}
	}

	SpreadElement struct {
		Start    int
		Argument Expression
	}

	ObjectLiteral struct {
		LeftBrace  int
		RightBrace int
		Value      []Property
	}

	ParameterList struct {
		Opening int
		List    []*Identifier
		Closing int
	}

	ObjectProperty struct {
		Key   Expression
		Kind  string
		Value Expression
	}

	Property interface {
	}

	RegExpLiteral struct {
		Index   int
		Literal string
		Pattern string
		Flags   string
	}

	SequenceExpression struct {
		Sequence []Expression
	}

	StringLiteral struct {
		Index   int
		Literal string
		Value   unistring.String
	}

	ThisExpression struct {
		Index int
	}

	UnaryExpression struct {
		Operator token.Token
		Index    int
		Operand  Expression
		Postfix  bool
	}

	VariableExpression struct {
		Name        unistring.String
		Index       int
		Initializer Expression
	}

	MetaProperty struct {
		Meta     *Identifier
		Property *Identifier
		Index    int
	}

	Statement interface {
		Node
	}

	BadStatement struct {
		From int
		To   int
	}

	BlockStatement struct {
		LeftBrace  int
		List       []Statement
		RightBrace int
	}

	BranchStatement struct {
		Index int
		Token token.Token
		Label *Identifier
	}

	CaseStatement struct {
		Case       int
		Test       Expression
		Consequent []Statement
	}

	CatchStatement struct {
		Catch     int
		Parameter *Identifier
		Body      Statement
	}

	DebuggerStatement struct {
		Debugger int
	}

	DoWhileStatement struct {
		Do   int
		Test Expression
		Body Statement
	}

	EmptyStatement struct {
		Semicolon int
	}

	ExpressionStatement struct {
		Expression Expression
	}

	ForInStatement struct {
		For    int
		Into   Expression
		Source Expression
		Body   Statement
	}

	ForOfStatement struct {
		For    int
		Into   Expression
		Source Expression
		Body   Statement
	}

	ForStatement struct {
		For         int
		Initializer Expression
		Update      Expression
		Test        Expression
		Body        Statement
	}

	IfStatement struct {
		If         int
		Test       Expression
		Consequent Statement
		Alternate  Statement
	}

	ImportStatement struct {
		Source string
		Name   string
		From   int
		To     int
	}

	LabelledStatement struct {
		Label     *Identifier
		Colon     int
		Statement Statement
	}

	ReturnStatement struct {
		Return   int
		Argument Expression
	}

	SwitchStatement struct {
		Switch       int
		Discriminant Expression
		Default      int
		Body         []*CaseStatement
	}

	ThrowStatement struct {
		Throw    int
		Argument Expression
	}

	TryStatement struct {
		Try     int
		Body    Statement
		Catch   *CatchStatement
		Finally Statement
	}

	VariableStatement struct {
		Var  int
		List []Expression
	}

	WhileStatement struct {
		While int
		Test  Expression
		Body  Statement
	}

	WithStatement struct {
		With   int
		Object Expression
		Body   Statement
	}

	ExportStatement struct {
		Export    int
		Statement Statement
	}

	Declaration interface {
	}

	FunctionDeclaration struct {
		Function *FunctionLiteral
	}

	VariableDeclaration struct {
		Var  int
		List []*VariableExpression
	}
)

type Program struct {
	Body            []Statement
	DeclarationList []Declaration
}

func (self *ArrayLiteral) Index0() int          { return self.LeftBracket }
func (self *AssignExpression) Index0() int      { return self.Left.Index0() }
func (self *BadExpression) Index0() int         { return self.From }
func (self *BinaryExpression) Index0() int      { return self.Left.Index0() }
func (self *BooleanLiteral) Index0() int        { return self.Index }
func (self *BracketExpression) Index0() int     { return self.Left.Index0() }
func (self *CallExpression) Index0() int        { return self.Callee.Index0() }
func (self *ConditionalExpression) Index0() int { return self.Test.Index0() }
func (self *DotExpression) Index0() int         { return self.Left.Index0() }
func (self *FunctionLiteral) Index0() int       { return self.Function }
func (self *Identifier) Index0() int            { return self.Index }
func (self *NewExpression) Index0() int         { return self.New }
func (self *NullLiteral) Index0() int           { return self.Index }
func (self *NumberLiteral) Index0() int         { return self.Index }
func (self *ObjectLiteral) Index0() int         { return self.LeftBrace }
func (self *SpreadElement) Index0() int         { return self.Start }
func (self *RegExpLiteral) Index0() int         { return self.Index }
func (self *SequenceExpression) Index0() int    { return self.Sequence[0].Index0() }
func (self *StringLiteral) Index0() int         { return self.Index }
func (self *ThisExpression) Index0() int        { return self.Index }
func (self *UnaryExpression) Index0() int       { return self.Index }
func (self *VariableExpression) Index0() int    { return self.Index }
func (self *MetaProperty) Index0() int          { return self.Index }

func (self *BadStatement) Index0() int        { return self.From }
func (self *BlockStatement) Index0() int      { return self.LeftBrace }
func (self *BranchStatement) Index0() int     { return self.Index }
func (self *CaseStatement) Index0() int       { return self.Case }
func (self *CatchStatement) Index0() int      { return self.Catch }
func (self *DebuggerStatement) Index0() int   { return self.Debugger }
func (self *DoWhileStatement) Index0() int    { return self.Do }
func (self *EmptyStatement) Index0() int      { return self.Semicolon }
func (self *ExpressionStatement) Index0() int { return self.Expression.Index0() }
func (self *ForInStatement) Index0() int      { return self.For }
func (self *ForOfStatement) Index0() int      { return self.For }
func (self *ForStatement) Index0() int        { return self.For }
func (self *IfStatement) Index0() int         { return self.If }
func (self *ImportStatement) Index0() int     { return self.From }
func (self *LabelledStatement) Index0() int   { return self.Label.Index0() }
func (self *Program) Index0() int             { return self.Body[0].Index0() }
func (self *ReturnStatement) Index0() int     { return self.Return }
func (self *SwitchStatement) Index0() int     { return self.Switch }
func (self *ThrowStatement) Index0() int      { return self.Throw }
func (self *TryStatement) Index0() int        { return self.Try }
func (self *VariableStatement) Index0() int   { return self.Var }
func (self *WhileStatement) Index0() int      { return self.While }
func (self *WithStatement) Index0() int       { return self.With }
func (self *ExportStatement) Index0() int     { return self.Export }

func (self *ArrayLiteral) Index1() int          { return self.RightBracket }
func (self *AssignExpression) Index1() int      { return self.Right.Index1() }
func (self *BadExpression) Index1() int         { return self.To }
func (self *BinaryExpression) Index1() int      { return self.Right.Index1() }
func (self *BooleanLiteral) Index1() int        { return self.Index + len(self.Literal) }
func (self *BracketExpression) Index1() int     { return self.RightBracket + 1 }
func (self *CallExpression) Index1() int        { return self.RightParenthesis + 1 }
func (self *ConditionalExpression) Index1() int { return self.Test.Index1() }
func (self *DotExpression) Index1() int         { return self.Identifier.Index1() }
func (self *FunctionLiteral) Index1() int       { return self.Body.Index1() }
func (self *Identifier) Index1() int            { return self.Index + len(self.Name) }
func (self *NewExpression) Index1() int         { return self.RightParenthesis + 1 }
func (self *NullLiteral) Index1() int           { return self.Index + 4 }
func (self *NumberLiteral) Index1() int         { return self.Index + len(self.Literal) }
func (self *ObjectLiteral) Index1() int         { return self.RightBrace }
func (self *SpreadElement) Index1() int         { return self.Argument.Index1() }
func (self *RegExpLiteral) Index1() int         { return self.Index + len(self.Literal) }
func (self *SequenceExpression) Index1() int    { return self.Sequence[0].Index1() }
func (self *StringLiteral) Index1() int         { return self.Index + len(self.Literal) }
func (self *ThisExpression) Index1() int        { return self.Index }
func (self *UnaryExpression) Index1() int {
	if self.Postfix {
		return self.Operand.Index1() + 2 // ++ --
	}
	return self.Operand.Index1()
}
func (self *VariableExpression) Index1() int {
	if self.Initializer == nil {
		return self.Index + len(self.Name) + 1
	}
	return self.Initializer.Index1()
}
func (self *MetaProperty) Index1() int {
	return self.Property.Index1()
}

func (self *BadStatement) Index1() int        { return self.To }
func (self *BlockStatement) Index1() int      { return self.RightBrace + 1 }
func (self *BranchStatement) Index1() int     { return self.Index }
func (self *CaseStatement) Index1() int       { return self.Consequent[len(self.Consequent)-1].Index1() }
func (self *CatchStatement) Index1() int      { return self.Body.Index1() }
func (self *DebuggerStatement) Index1() int   { return self.Debugger + 8 }
func (self *DoWhileStatement) Index1() int    { return self.Test.Index1() }
func (self *EmptyStatement) Index1() int      { return self.Semicolon + 1 }
func (self *ExpressionStatement) Index1() int { return self.Expression.Index1() }
func (self *ForInStatement) Index1() int      { return self.Body.Index1() }
func (self *ForOfStatement) Index1() int      { return self.Body.Index1() }
func (self *ForStatement) Index1() int        { return self.Body.Index1() }
func (self *IfStatement) Index1() int {
	if self.Alternate != nil {
		return self.Alternate.Index1()
	}
	return self.Consequent.Index1()
}
func (self *ImportStatement) Index1() int   { return self.To }
func (self *LabelledStatement) Index1() int { return self.Colon + 1 }
func (self *Program) Index1() int           { return self.Body[len(self.Body)-1].Index1() }
func (self *ReturnStatement) Index1() int   { return self.Return }
func (self *SwitchStatement) Index1() int   { return self.Body[len(self.Body)-1].Index1() }
func (self *ThrowStatement) Index1() int    { return self.Throw }
func (self *TryStatement) Index1() int      { return self.Try }
func (self *VariableStatement) Index1() int { return self.List[len(self.List)-1].Index1() }
func (self *WhileStatement) Index1() int    { return self.Body.Index1() }
func (self *WithStatement) Index1() int     { return self.Body.Index1() }
func (self *ExportStatement) Index1() int   { return self.Statement.Index1() }
