package main

import (
	"errors"
	"fmt"
)

type Symbol struct {
	Name       string
	Type       string   // int
	Kind       string   // variable/parameter/method
	Offset     int      // thesh sto stack
	ParamCount int      // arithmos parametron (an einai methodos)
	ParamTypes []string // types twn parametron (an einai methodos)
	Line       int      // errors
}

// pinakas symbolwn gia ena scope
type SymbolTable struct {
	Symbols  map[string]*Symbol
	Name     string // onoma tou scope
	VarCount int    // arithmos metavlitwn sto scope
}

func NewSymbolTable(name string) *SymbolTable {
	return &SymbolTable{
		Symbols:  make(map[string]*Symbol),
		Name:     name,
		VarCount: 0,
	}
}

// add symbolo me duplicate check
func (st *SymbolTable) AddSymbol(symbol *Symbol) error {
	if _, exists := st.Symbols[symbol.Name]; exists {
		return fmt.Errorf("symbol '%s' already exists in scope '%s'", symbol.Name, st.Name)
	}

	st.Symbols[symbol.Name] = symbol

	// enhmerwsi offset kai varCount
	if symbol.Kind == "variable" || symbol.Kind == "parameter" {
		symbol.Offset = st.VarCount
		st.VarCount++
	}
	return nil
}

func (st *SymbolTable) Lookup(name string) (*Symbol, bool) {
	symbol, exists := st.Symbols[name]
	return symbol, exists
}

type SemanticAnalyzer struct {
	//Global scope gia methodous
	globalSymbols *SymbolTable

	//Symbol table gia kathe methodo
	methodTables map[string]*SymbolTable

	currentMethod string       // current methodos
	currentTable  *SymbolTable // current symbol table

	// Errors
	errors []string

	// Loop tracking
	loopDepth int
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		globalSymbols: NewSymbolTable("global"),
		methodTables:  make(map[string]*SymbolTable),
		errors:        []string{},
		loopDepth:     0,
	}
}

func (s *SemanticAnalyzer) Analyze(ast *AST) (map[string]*SymbolTable, error) {

	// Syllegw ta method signatures
	for _, method := range ast.Methods {
		if err := s.addMethodSignature(method); err != nil {
			s.errors = append(s.errors, err.Error())
		}
	}

	// elegxos main
	if _, exists := s.globalSymbols.Symbols["main"]; !exists {
		s.errors = append(s.errors, "error: no 'main' method found")
	}

	// method body analysis
	for _, method := range ast.Methods {
		if err := s.analyzeMethod(method); err != nil {
			s.errors = append(s.errors, err.Error())
		}
	}

	// errors
	if len(s.errors) > 0 {
		errorMsg := "Semantic errors found:\n"
		for _, err := range s.errors {
			errorMsg += fmt.Sprintf("- %s\n", err)
		}
		return nil, errors.New(errorMsg)
	}

	return s.methodTables, nil
}

func (s *SemanticAnalyzer) addMethodSignature(method Method) error {
	// overload checking
	paramTypes := make([]string, len(method.Parameters))
	for i, param := range method.Parameters {
		paramTypes[i] = param.Type
	}

	methodSymbol := &Symbol{
		Name:       method.Name,
		Type:       method.ReturnType,
		Kind:       "method",
		ParamCount: len(method.Parameters),
		ParamTypes: paramTypes,
		Line:       method.Line,
	}

	return s.globalSymbols.AddSymbol(methodSymbol)
}

func (s *SemanticAnalyzer) analyzeMethod(method Method) error {
	// dhmiourgw local symbol table gia th methodo
	methodTable := NewSymbolTable(method.Name)
	s.methodTables[method.Name] = methodTable

	// set current method and table
	s.currentMethod = method.Name
	s.currentTable = methodTable
	s.loopDepth = 0

	// add parametrwn sto method scope
	for _, param := range method.Parameters {
		paramSymbol := &Symbol{
			Name: param.Name,
			Type: param.Type,
			Kind: "parameter",
			Line: param.Line,
		}

		if err := methodTable.AddSymbol(paramSymbol); err != nil {
			return err
		}
	}

	return s.analyzeBlock(method.Body)
}

func (s *SemanticAnalyzer) analyzeBlock(block Block) error {
	// declarations apo vars (panta prwtes)
	for _, decl := range block.Declarations {
		if err := s.analyzeDeclaration(decl); err != nil {
			return err
		}
	}

	// statement exec
	for _, stmt := range block.Statements {
		if err := s.analyzeStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (s *SemanticAnalyzer) analyzeDeclaration(decl Declaration) error {
	for _, variable := range decl.Variables {
		varSymbol := &Symbol{
			Name: variable.Name,
			Type: decl.Type, // mono int
			Kind: "variable",
			Line: decl.Line,
		}

		// prosthiki sto scope
		if err := s.currentTable.AddSymbol(varSymbol); err != nil {
			return err
		}

		// elegxos gia init
		if variable.InitialValue != nil {
			exprType, err := s.analyzeExpression(variable.InitialValue)
			if err != nil {
				return err
			}

			// type compatibility
			if exprType != "int" {
				return fmt.Errorf("type mismatch in variable initialization at line %d: expected int, got %s",
					decl.Line, exprType)
			}

		}
	}
	return nil
}

// statement switch
func (s *SemanticAnalyzer) analyzeStatement(stmt Statement) error {
	switch stmt := stmt.(type) {
	case *ReturnStatement:
		return s.analyzeReturnStatement(stmt)
	case *IfStatement:
		return s.analyzeIfStatement(stmt)
	case *WhileStatement:
		return s.analyzeWhileStatement(stmt)
	case *BreakStatement:
		return s.analyzeBreakStatement(stmt)
	case *BlockStatement:
		return s.analyzeBlock(stmt.Block)
	case *Assignment:
		return s.analyzeAssignment(stmt)
	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

// expression switch
func (s *SemanticAnalyzer) analyzeExpression(expr Expression) (string, error) {
	switch e := expr.(type) {
	case *NumberLiteral:
		return "int", nil
	case *BooleanLiteral:
		return "int", nil // true = 1 , false = 0
	case *Identifier:
		return s.analyzeIdentifier(e)
	case *BinaryExpression:
		return s.analyzeBinaryExpression(e)
	case *UnaryExpression:
		return s.analyzeUnaryExpression(e)
	case *MethodCall:
		return s.analyzeMethodCall(e)
	default:
		return "", fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (s *SemanticAnalyzer) analyzeIdentifier(expr *Identifier) (string, error) {
	symbol, exists := s.currentTable.Lookup(expr.Name)
	if !exists {
		return "", fmt.Errorf("undefined identifier '%s' at line %d", expr.Name, expr.Line)
	}

	return symbol.Type, nil
}

func (s *SemanticAnalyzer) analyzeBinaryExpression(expr *BinaryExpression) (string, error) {
	// aristerh pleura
	leftType, err := s.analyzeExpression(expr.Left)
	if err != nil {
		return "", err
	}

	// dexia pleura
	rightType, err := s.analyzeExpression(expr.Right)
	if err != nil {
		return "", err
	}

	// type compatibility
	if leftType != "int" || rightType != "int" {
		return "", fmt.Errorf("type mismatch in binary expression at line %d: expected int, got %s and %s",
			expr.Line, leftType, rightType)
	}

	return "int", nil
}

func (s *SemanticAnalyzer) analyzeUnaryExpression(expr *UnaryExpression) (string, error) {
	operandType, err := s.analyzeExpression(expr.Operand)
	if err != nil {
		return "", err
	}

	if operandType != "int" {
		return "", fmt.Errorf("type mismatch in unary expression at line %d: expected int, got %s",
			expr.Line, operandType)
	}
	return "int", nil
}

func (s *SemanticAnalyzer) analyzeMethodCall(expr *MethodCall) (string, error) {
	// anazhthsh methodou global scope
	methodSymbol, exists := s.globalSymbols.Lookup(expr.Name)
	if !exists {
		return "", fmt.Errorf("undefined method '%s' at line %d", expr.Name, expr.Line)
	}

	// elegxos parametron
	if len(expr.Arguments) != methodSymbol.ParamCount {
		return "", fmt.Errorf("method '%s' called with wrong number of arguments at line %d: expected %d, got %d",
			expr.Name, expr.Line, methodSymbol.ParamCount, len(expr.Arguments))
	}

	// elegxos typwn parametron
	for i, arg := range expr.Arguments {
		argType, err := s.analyzeExpression(arg)
		if err != nil {
			return "", err
		}

		if argType != methodSymbol.ParamTypes[i] {
			return "", fmt.Errorf("type mismatch in argument %d of method '%s' at line %d: expected %s, got %s",
				i+1, expr.Name, expr.Line, methodSymbol.ParamTypes[i], argType)
		}
	}

	return methodSymbol.Type, nil
}

func (s *SemanticAnalyzer) analyzeReturnStatement(stmt *ReturnStatement) error {
	// elegxos return
	exprType, err := s.analyzeExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// elegxos an to type tou return tairiazei me to method type
	methodSymbol, exists := s.globalSymbols.Symbols[s.currentMethod]
	if !exists {
		return fmt.Errorf("return statement outside of method at line %d", stmt.Line)
	}

	if exprType != methodSymbol.Type {
		return fmt.Errorf("type mismatch in return statement at line %d: expected %s, got %s",
			stmt.Line, methodSymbol.Type, exprType)
	}
	return nil
}

func (s *SemanticAnalyzer) analyzeIfStatement(stmt *IfStatement) error {
	_, err := s.analyzeExpression(stmt.Condition)
	if err != nil {
		return err
	}

	// analysh to then
	if err := s.analyzeStatement(stmt.ThenStmt); err != nil {
		return err
	}

	// analysh to else
	if stmt.ElseStmt != nil {
		if err := s.analyzeStatement(stmt.ElseStmt); err != nil {
			return err
		}
	}

	return nil
}

func (s *SemanticAnalyzer) analyzeWhileStatement(stmt *WhileStatement) error {
	// elegxos condition
	_, err := s.analyzeExpression(stmt.Condition)
	if err != nil {
		return err
	}

	// ++ sto loop depth gia break validation
	s.loopDepth++

	// analysh to body
	err = s.analyzeStatement(stmt.Body)

	s.loopDepth--

	return err
}

func (s *SemanticAnalyzer) analyzeBreakStatement(stmt *BreakStatement) error {
	// elegxos an eimaste se loop
	if s.loopDepth == 0 {
		return fmt.Errorf("break statement outside of loop at line %d", stmt.Line)
	}
	return nil
}

func (s *SemanticAnalyzer) analyzeAssignment(stmt *Assignment) error {
	// elegxos an einai dlwmenh h metavlhti
	symbol, exists := s.currentTable.Lookup(stmt.Variable)
	if !exists {
		return fmt.Errorf("undefined variable '%s' at line %d", stmt.Variable, stmt.Line)
	}

	// elegxos an einai metavlhti h parametros kai oxi methodos
	if symbol.Kind == "method" {
		return fmt.Errorf("cannot assign to method '%s' at line %d", stmt.Variable, stmt.Line)
	}

	// elegxos tou assigned expression
	exprType, err := s.analyzeExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// type compatibility
	if exprType != symbol.Type {
		return fmt.Errorf("type mismatch in assignment to '%s' at line %d: expected %s, got %s",
			symbol.Name, stmt.Line, symbol.Type, exprType)
	}
	return nil
}
