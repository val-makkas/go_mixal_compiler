package main

import (
	"fmt"
)

// Recursive Descent Parser
type Parser struct {
	tokens   []Token // tokens apo ton lexer
	position int     // thesh token
	current  Token   // current token
}

func NewParser() *Parser {
	return &Parser{}
}

// main function pou metatrepei token se AST
func (p *Parser) Parse(tokens []Token) (*AST, error) {
	p.tokens = tokens
	p.position = 0

	if len(tokens) > 0 {
		p.current = tokens[0]
	} else {
		return nil, fmt.Errorf("no tokens to parse")
	}

	// Parsing ksekina apo PROGRAM -> METH-LIST | e
	methods, err := p.parseProgram()
	if err != nil {
		return nil, err
	}

	if !p.isAtEnd() && p.current.Type != TOK_EOF {
		return nil, p.error("unexpected token after end of program")
	}

	return &AST{Methods: methods}, nil
}

// PROGRAM -> METH-LIST | e
func (p *Parser) parseProgram() ([]Method, error) {
	var methods []Method

	// An oxi tokens h EOF return keno
	if p.isAtEnd() || p.current.Type == TOK_EOF {
		return methods, nil
	}

	// alliws synexizoume me thn anazhthsh
	// METH-LIST -> METH METH-LIST | METH
	for !p.isAtEnd() && p.current.Type == TOK_INT {
		method, err := p.parseMethod()
		if err != nil {
			return nil, err
		}
		methods = append(methods, method)
	}
	return methods, nil
}

// METH -> TYPE id '(' PARAMS ')' BODY
func (p *Parser) parseMethod() (Method, error) {
	startLine := p.current.Line

	// TYPE prepei na einai int
	if p.current.Type != TOK_INT {
		return Method{}, p.error(fmt.Sprintf("expected type 'int', got '%s'", p.current.Value))
	}
	returnType := p.current.Value
	p.advance()

	// id
	if p.current.Type != TOK_ID {
		return Method{}, p.error(fmt.Sprintf("expected method name, got '%s'", p.current.Value))
	}
	methodName := p.current.Value
	p.advance()

	// '('
	if p.current.Type != TOK_LPAREN {
		return Method{}, p.error(fmt.Sprintf("expected '(', got '%s'", p.current.Value))
	}
	p.advance()

	// PARAMS
	parameters, err := p.parseParameters()
	if err != nil {
		return Method{}, err
	}

	// ')'
	if p.current.Type != TOK_RPAREN {
		return Method{}, p.error(fmt.Sprintf("expected ')', got '%s'", p.current.Value))
	}
	p.advance()

	// BODY
	body, err := p.parseBody()
	if err != nil {
		return Method{}, err
	}
	return Method{
		ReturnType: returnType,
		Name:       methodName,
		Parameters: parameters,
		Body:       body,
		Line:       startLine,
	}, nil
}

// PARAMS -> FORMMALS | e
// FORMALS -> TYPE id (',' TYPE id)*
func (p *Parser) parseParameters() ([]Parameter, error) {
	var parameters []Parameter

	// oxi parametroi, epistrefw kenh lista
	if p.current.Type == TOK_RPAREN {
		return parameters, nil
	}

	// prwth parametros
	param, err := p.parseParameter()
	if err != nil {
		return nil, err
	}
	parameters = append(parameters, param)

	// ypoloipes parametroi
	for p.current.Type == TOK_COMMA {
		p.advance() // skip ','
		param, err := p.parseParameter()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}

	return parameters, nil
}

// TYPE id
func (p *Parser) parseParameter() (Parameter, error) {
	startLine := p.current.Line

	//TYPE
	if p.current.Type != TOK_INT {
		return Parameter{}, p.error(fmt.Sprintf("exprected int, got '%s'", p.current.Value))
	}
	paramType := p.current.Value
	p.advance()

	// id
	if p.current.Type != TOK_ID {
		return Parameter{}, p.error(fmt.Sprintf("expected parameter name, got '%s'", p.current.Value))
	}
	paramName := p.current.Value
	p.advance()

	return Parameter{
		Type: paramType,
		Name: paramName,
		Line: startLine,
	}, nil
}

// BODY -> '{' DECLS STMTS '}'
func (p *Parser) parseBody() (Block, error) {
	// '{'
	if p.current.Type != TOK_LBRACE {
		return Block{}, p.error(fmt.Sprintf("expected '{', got '%s'", p.current.Value))
	}
	p.advance()

	//DECLS
	declarations, err := p.parseDeclarations()
	if err != nil {
		return Block{}, err
	}

	// STMTS
	statements, err := p.parseStatements()
	if err != nil {
		return Block{}, err
	}

	// '}'
	if p.current.Type != TOK_RBRACE {
		return Block{}, p.error(fmt.Sprintf("expected '}', got '%s'", p.current.Value))
	}
	p.advance()

	return Block{
		Declarations: declarations,
		Statements:   statements,
	}, nil
}

// DECLS -> DECL*
func (p *Parser) parseDeclarations() ([]Declaration, error) {
	var declarations []Declaration

	// oso exw int decls
	for p.current.Type == TOK_INT {
		decl, err := p.parseDeclaration()
		if err != nil {
			return nil, err
		}
		declarations = append(declarations, decl)
	}
	return declarations, nil
}

// DECL -> TYPE id VARS ';'
func (p *Parser) parseDeclaration() (Declaration, error) {
	startLine := p.current.Line

	// TYPE
	if p.current.Type != TOK_INT {
		return Declaration{}, p.error(fmt.Sprintf("expected type 'int', got '%s'", p.current.Value))
	}
	varType := p.current.Value
	p.advance()

	// prwth metablhth
	var variables []Variable
	variable, err := p.parseVariable()
	if err != nil {
		return Declaration{}, err
	}
	variables = append(variables, variable)

	//ypoloipes metavlhtes
	for p.current.Type == TOK_COMMA {
		p.advance() // skip ','

		variable, err := p.parseVariable()
		if err != nil {
			return Declaration{}, err
		}
		variables = append(variables, variable)
	}
	// ';'
	if p.current.Type != TOK_SEMICOLON {
		return Declaration{}, p.error(fmt.Sprintf("expected ';', got '%s'", p.current.Value))
	}
	p.advance()

	return Declaration{
		Type:      varType,
		Variables: variables,
		Line:      startLine,
	}, nil
}

// id | id '=' EXPR
func (p *Parser) parseVariable() (Variable, error) {
	// id
	if p.current.Type != TOK_ID {
		return Variable{}, p.error(fmt.Sprintf("expected variable name, got '%s'", p.current.Value))
	}
	varName := p.current.Value
	p.advance()

	var initialValue Expression

	// elegxw an exei timh
	if p.current.Type == TOK_ASSIGN {
		p.advance() // skip '='
		expr, err := p.parseExpression()
		if err != nil {
			return Variable{}, err
		}
		initialValue = expr
	} else {
		initialValue = nil // an den yparxei timh
	}

	return Variable{
		Name:         varName,
		InitialValue: initialValue,
	}, nil
}

// STMTS -> STMT*
func (p *Parser) parseStatements() ([]Statement, error) {
	var statements []Statement

	// oso exw statements
	for !p.isAtEnd() && p.current.Type != TOK_RBRACE {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		// an stmt den einai nil to prostheto sto slice
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements, nil
}

// STMT -> ASSIGN ';' | return EXPR ';' | if '(' EXPR ')' STMT else STMT
//
//	| while '(' EXPR ')' STMT | break ';' | BLOCK | ';'

func (p *Parser) parseStatement() (Statement, error) {
	switch p.current.Type {
	case TOK_RETURN:
		return p.parseReturnStatement()
	case TOK_IF:
		return p.parseIfStatement()
	case TOK_WHILE:
		return p.parseWhileStatement()
	case TOK_BREAK:
		return p.parseBreakStatement()
	case TOK_LBRACE:
		return p.parseBlockStatement()
	case TOK_SEMICOLON:
		p.advance() // skip ';'
		return nil, nil
	case TOK_ID:
		return p.parseAssignmentStatement()
	default:
		return nil, p.error(fmt.Sprintf("unexpected token '%s' at line %d", p.current.Value, p.current.Line))
	}
}

// return EXPR ';'
func (p *Parser) parseReturnStatement() (Statement, error) {
	startLine := p.current.Line
	p.advance() // skip 'return'

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.current.Type != TOK_SEMICOLON {
		return nil, p.error(fmt.Sprintf("expected ';' after return expression, got '%s'", p.current.Value))
	}
	p.advance() // skip ';'

	return &ReturnStatement{
		Expression: expr,
		Line:       startLine,
	}, nil
}

// if '(' EXPR ')' STMT else STMT
func (p *Parser) parseIfStatement() (Statement, error) {
	startLine := p.current.Line
	p.advance()

	// '('
	if p.current.Type != TOK_LPAREN {
		return nil, p.error(fmt.Sprintf("expected '(', got '%s'", p.current.Value))
	}
	p.advance()

	// EXPR
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// ')'
	if p.current.Type != TOK_RPAREN {
		return nil, p.error(fmt.Sprintf("expected ')', got '%s'", p.current.Value))
	}
	p.advance()

	// STMT (then branch)
	thenStmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	var elseStmt Statement
	// else (ean yparxei)
	if p.current.Type == TOK_ELSE {
		p.advance() // skip 'else'

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		elseStmt = stmt
	}

	return &IfStatement{
		Condition: condition,
		ThenStmt:  thenStmt,
		ElseStmt:  elseStmt,
		Line:      startLine,
	}, nil
}

// while '(' EXPR ')' STMT
func (p *Parser) parseWhileStatement() (Statement, error) {
	startLine := p.current.Line
	p.advance()

	// '('
	if p.current.Type != TOK_LPAREN {
		return nil, p.error(fmt.Sprintf("expected '(', got '%s'", p.current.Value))
	}
	p.advance()

	// EXPR
	condition, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// ')'
	if p.current.Type != TOK_RPAREN {
		return nil, p.error(fmt.Sprintf("expected ')', got '%s'", p.current.Value))
	}
	p.advance()

	// STMT
	body, err := p.parseStatement()
	if err != nil {
		return nil, err
	}

	return &WhileStatement{
		Condition: condition,
		Body:      body,
		Line:      startLine,
	}, nil
}

// break ';'
func (p *Parser) parseBreakStatement() (Statement, error) {
	startLine := p.current.Line
	p.advance()

	if p.current.Type != TOK_SEMICOLON {
		return nil, p.error(fmt.Sprintf("expected ';' after break, got '%s'", p.current.Value))
	}
	p.advance()

	return &BreakStatement{Line: startLine}, nil
}

// '{' STMTS '}'
func (p *Parser) parseBlockStatement() (Statement, error) {
	startLine := p.current.Line
	p.advance() // skip '{'

	// STMTS
	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	// '}'
	if p.current.Type != TOK_RBRACE {
		return nil, p.error(fmt.Sprintf("expected '}', got '%s'", p.current.Value))
	}
	p.advance()

	return &BlockStatement{
		Block: Block{
			Declarations: []Declaration{}, // ta blocks den exoun declarations
			Statements:   statements,
		},
		Line: startLine,
	}, nil
}

// ASSIGN -> LOCATION '=' EXPR ';'
func (p *Parser) parseAssignmentStatement() (Statement, error) {
	startLine := p.current.Line

	// LOCATION
	if p.current.Type != TOK_ID {
		return nil, p.error(fmt.Sprintf("expected variable name, got '%s'", p.current.Value))
	}
	varName := p.current.Value
	p.advance()

	// '='
	if p.current.Type != TOK_ASSIGN {
		return nil, p.error(fmt.Sprintf("expected '=', got '%s'", p.current.Value))
	}
	p.advance()

	// EXPR
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// ';'
	if p.current.Type != TOK_SEMICOLON {
		return nil, p.error(fmt.Sprintf("expected ';' after assignment, got '%s'", p.current.Value))
	}
	p.advance()

	return &Assignment{
		Variable:   varName,
		Expression: expr,
		Line:       startLine,
	}, nil
}

// EXPR -> ADD-EXPR RELOP ADD-EXPR | ADD-EXPR
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseRelationalExpression()
}

// EXPR -> ADD-EXPR RELOP ADD-EXPR | ADD-EXPR
// RELOP -> '==' | '!=' | '<' | '<=' | '>' | '>='
func (p *Parser) parseRelationalExpression() (Expression, error) {
	left, err := p.parseAddExpression()
	if err != nil {
		return nil, err
	}

	// elegxw an exei sxesiako telesth
	if p.isRelationalOperator() {
		operator := p.current.Value
		line := p.current.Line
		p.advance() // skip relational operator

		right, err := p.parseAddExpression()
		if err != nil {
			return nil, err
		}

		return &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
			Line:     line,
		}, nil
	}

	return left, nil
}

// ADD-EXPR -> ADD-EXPR ADDOP TERM | TERM
// ADDOP -> '+' | '-'
func (p *Parser) parseAddExpression() (Expression, error) {
	left, err := p.parseMultiplyExpression()
	if err != nil {
		return nil, err
	}

	// a + b + c = (a + b) + c
	for p.isAddOperator() {
		operator := p.current.Value
		line := p.current.Line
		p.advance() // skip operator

		right, err := p.parseMultiplyExpression()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
			Line:     line,
		}
	}
	return left, nil
}

// TERM -> TERM MULOP FACTOR | FACTOR
// MULOP -> '*' | '/'
func (p *Parser) parseMultiplyExpression() (Expression, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	// a * b * c = (a * b) * c
	for p.isMultiplyOperator() {
		operator := p.current.Value
		line := p.current.Line
		p.advance()

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		left = &BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
			Line:     line,
		}
	}
	return left, nil
}

// xeirizetai ta vasika stoixeia twn expression
// FACTOR -> '(' EXPR ')' | LOCATION | num | treu | false | METHOD '(' ACTUALS ')'
func (p *Parser) parseFactor() (Expression, error) {
	switch p.current.Type {
	case TOK_LPAREN:
		// '(' EXPR ')'
		p.advance() // skip '('

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if p.current.Type != TOK_RPAREN {
			return nil, p.error(fmt.Sprintf("expected ')', got '%s'", p.current.Value))
		}
		p.advance()

		return expr, nil

	case TOK_NUM:
		// num
		value := p.current.Value
		line := p.current.Line
		p.advance()
		return &NumberLiteral{
			Value: value,
			Line:  line,
		}, nil

	case TOK_TRUE:
		// true
		line := p.current.Line
		p.advance()
		return &BooleanLiteral{
			Value: true,
			Line:  line,
		}, nil

	case TOK_FALSE:
		// false
		line := p.current.Line
		p.advance()
		return &BooleanLiteral{
			Value: false,
			Line:  line,
		}, nil

	case TOK_ID:
		// LOCATION | METHOD '(' ACTUALS ')'
		name := p.current.Value
		line := p.current.Line
		p.advance()

		// an einai klhsh methodou
		if p.current.Type == TOK_LPAREN {
			// METHOD '(' ACTUALS ')'
			p.advance()

			arguments, err := p.parseActuals()
			if err != nil {
				return nil, err
			}

			if p.current.Type != TOK_RPAREN {
				return nil, p.error(fmt.Sprintf("expected ')', got '%s'", p.current.Value))
			}
			p.advance()

			return &MethodCall{
				Name:      name,
				Arguments: arguments,
				Line:      line,
			}, nil
		}

		// alliws einai identifier
		return &Identifier{
			Name: name,
			Line: line,
		}, nil

	case TOK_MINUS:
		// monadikh ekfrash p.x. -x
		line := p.current.Line
		p.advance()

		operand, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		return &UnaryExpression{
			Operator: "-",
			Operand:  operand,
			Line:     line,
		}, nil

	default:
		return nil, p.error(fmt.Sprintf("unexpected token in expression: '%s'", p.current.Value))
	}
}

// ACTUALS -> EXPR ARGS | e
// ARGS -> ',' EXPR ARGS | e
func (p *Parser) parseActuals() ([]Expression, error) {
	var arguments []Expression

	// an den exw args
	if p.current.Type == TOK_RPAREN {
		return arguments, nil
	}

	// prwth ekfrash
	arg, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	arguments = append(arguments, arg)

	// ypoloipes an yparxoun
	for p.current.Type == TOK_COMMA {
		p.advance()

		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}
	return arguments, nil
}

// HELPERS
func (p *Parser) isRelationalOperator() bool {
	return p.current.Type == TOK_LT || p.current.Type == TOK_LE ||
		p.current.Type == TOK_GT || p.current.Type == TOK_GE ||
		p.current.Type == TOK_EQ || p.current.Type == TOK_NE
}

func (p *Parser) isAddOperator() bool {
	return p.current.Type == TOK_PLUS || p.current.Type == TOK_MINUS
}

func (p *Parser) isMultiplyOperator() bool {
	return p.current.Type == TOK_MULTIPLY || p.current.Type == TOK_DIVIDE
}
func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.position++
		if p.position < len(p.tokens) {
			p.current = p.tokens[p.position]
		} else {
			p.current = Token{Type: TOK_EOF, Value: "", Line: -1}
		}
	}
}

func (p *Parser) isAtEnd() bool {
	return p.position >= len(p.tokens) || p.current.Type == TOK_EOF
}

func (p *Parser) error(message string) error {
	return fmt.Errorf("syntax error at line %d, column %d: %s",
		p.current.Line, p.current.Column, message)
}
