package main

// PROGRAM -> METH-LIST | e
type AST struct {
	Methods []Method // lista methodwn
}

// METH -> TYPE id '(' PARAMS ')' BODY
type Method struct {
	ReturnType string
	Name       string
	Parameters []Parameter
	Body       Block
	Line       int // errors
}

// FORMALS -> TYPE id
type Parameter struct {
	Type string
	Name string
	Line int // errors
}

// BODY -> '{' DECLS STMTS '}'
type Block struct {
	Declarations []Declaration // dhlwseis metablhtwn
	Statements   []Statement   // entoles
}

// DECL -> TYPE id VARS ';' | TYPE id '=' EXPR VARS ';'
type Declaration struct {
	Type      string
	Variables []Variable //lista metablhtwn
	Line      int        // grammh declare
}

// metablhth
type Variable struct {
	Name         string
	InitialValue Expression // arxikh timh (nil an den yparxei)
}

// interface entolwn
type Statement interface {
	statementNode()
}

// anathesh : ASSIGN -> LOCATION '=' EXPR
type Assignment struct {
	Variable   string     // onoma metablhths
	Expression Expression // ekfrash
	Line       int        // grammh
}

// return
type ReturnStatement struct {
	Expression Expression // express pou ginetai return
	Line       int
}

// if
type IfStatement struct {
	Condition Expression // synthiki
	ThenStmt  Statement  // true
	ElseStmt  Statement  // false
	Line      int
}

// while
type WhileStatement struct {
	Condition Expression // synthiki
	Body      Statement  // broxgos
	Line      int
}

// break
type BreakStatement struct {
	Line int
}

// block entolwn {}
type BlockStatement struct {
	Block Block
	Line  int
}

// interface ekfrasewn
type Expression interface {
	expressionNode()
}

// dyadikh ekfrash p.x. a + b , a == b ktlp
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	Line     int
}

// monadikh ekfrash p.x. (-x)
type UnaryExpression struct {
	Operator string
	Operand  Expression // ekfrash
	Line     int
}

// anafora se metablhth
type Identifier struct {
	Name string
	Line int
}

type NumberLiteral struct {
	Value string // p.x. "123", "4", "-5"
	Line  int
}

type BooleanLiteral struct {
	Value bool // true h false
	Line  int
}

// klhsh methodou
type MethodCall struct {
	Name      string
	Arguments []Expression
	Line      int
}

// interfaces
func (a *Assignment) statementNode()      {}
func (r *ReturnStatement) statementNode() {}
func (i *IfStatement) statementNode()     {}
func (w *WhileStatement) statementNode()  {}
func (b *BreakStatement) statementNode()  {}
func (b *BlockStatement) statementNode()  {}

func (b *BinaryExpression) expressionNode() {}
func (u *UnaryExpression) expressionNode()  {}
func (i *Identifier) expressionNode()       {}
func (n *NumberLiteral) expressionNode()    {}
func (b *BooleanLiteral) expressionNode()   {}
func (m *MethodCall) expressionNode()       {}
