package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CODE_START  = 1000 // program code
	VAR_START   = 2000 // storage metavlhtwn
	TEMP_START  = 3000 // storage temp
	STACK_START = 3500 // stack storage
)

type CodeGenerator struct {
	output         strings.Builder         // mixal code
	labelCounter   int                     // counter gia ta labels
	tempCounter    int                     // counter gia ta temp metavlhtes
	addressMap     map[string]int          // Var onoma -> memory address
	currentAddress int                     // current memory address
	breakLabels    []string                // stack gia ta break
	methodLabels   map[string]string       // Method onoma -> mixal label
	symbolTables   map[string]*SymbolTable // symbol tables gia kathe method
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		addressMap:     make(map[string]int),
		methodLabels:   make(map[string]string),
		currentAddress: VAR_START,
		labelCounter:   1,
		tempCounter:    1,
	}
}

func (c *CodeGenerator) Generate(ast *AST, symbolTables map[string]*SymbolTable) (string, error) {
	c.output.Reset()
	c.symbolTables = symbolTables

	// memory allocation
	if err := c.allocateMemory(symbolTables); err != nil {
		return "", fmt.Errorf("memory allocation error: %w", err)
	}

	// method label gen
	c.generateMethodLabels(ast)

	// main generation
	if err := c.generateMainProgram(ast); err != nil {
		return "", fmt.Errorf("main proccess generation error: %w", err)
	}

	// generation ypoloipwn methodwn
	if err := c.generateMethods(ast); err != nil {
		return "", fmt.Errorf("methods generation error: %w", err)
	}
	// telos programmatos
	c.generateFooter()

	return c.output.String(), nil
}

func (c *CodeGenerator) allocateMemory(symbolTables map[string]*SymbolTable) error {
	// allocate mnhmh gia parametrous
	for methodName, table := range symbolTables {
		paramCount := 0
		for _, symbol := range table.Symbols {
			if symbol.Kind == "parameter" {
				fullName := c.makeParameterName(methodName, paramCount)
				c.addressMap[fullName] = c.currentAddress
				c.currentAddress++
				paramCount++
			}
		}
	}

	// allocate mnhmh gia tis metavlites
	for methodName, table := range symbolTables {
		for varName, symbol := range table.Symbols {
			if symbol.Kind == "variable" {
				// dhmiourgia monadikou onomatos
				fullName := c.makeVariableName(methodName, varName)
				c.addressMap[fullName] = c.currentAddress
				c.currentAddress++
			}
		}
	}
	return nil
}

func (c *CodeGenerator) generateMethodLabels(ast *AST) {
	for _, method := range ast.Methods {
		if method.Name == "main" {
			c.methodLabels[method.Name] = "MAIN"
		} else {
			// metatroph se mixal label
			c.methodLabels[method.Name] = strings.ToUpper(method.Name)
		}
	}
}

func (c *CodeGenerator) generateMainProgram(ast *AST) error {
	// euresh main
	var mainMethod *Method
	for i, method := range ast.Methods {
		if method.Name == "main" {
			mainMethod = &ast.Methods[i]
			break
		}
	}

	if mainMethod == nil {
		return fmt.Errorf("main method not found")
	}

	if _, exists := c.symbolTables["main"]; !exists {
		return fmt.Errorf("symbol table for main method not found")
	}

	// mixal entry point
	c.output.WriteString("        ORIG  1000\n")
	c.output.WriteString("MAIN    NOP\n")

	// paragwgh body ths main
	if err := c.generateMethodBody(*mainMethod); err != nil {
		return fmt.Errorf("error generating main method body: %w", err)
	}

	c.output.WriteString("        HLT\n")

	return nil
}

func (c *CodeGenerator) generateMethods(ast *AST) error {
	for _, method := range ast.Methods {
		if method.Name != "main" {
			if err := c.generateMethod(method); err != nil {
				return fmt.Errorf("error generating method %s: %w", method.Name, err)
			}
		}
	}
	return nil
}

func (c *CodeGenerator) generateMethod(method Method) error {
	methodLabel := c.methodLabels[method.Name]
	exitLabel := methodLabel + "X"

	c.output.WriteString(fmt.Sprintf("%s    NOP\n", methodLabel))
	c.output.WriteString(fmt.Sprintf("        STJ   %s\n", exitLabel))

	if err := c.generateMethodBody(method); err != nil {
		return fmt.Errorf("error generating method body for %s: %w", method.Name, err)
	}

	c.output.WriteString(fmt.Sprintf("%s    JMP   *\n", exitLabel))
	return nil
}

func (c *CodeGenerator) generateMethodBody(method Method) error {
	// paragwgh dhlwsewn metavlhtwn
	for _, decl := range method.Body.Declarations {
		if err := c.generateDeclaration(decl, method.Name); err != nil {
			return fmt.Errorf("error generating declaration: %w", err)
		}
	}

	// paragwgh entolwn
	for _, stmt := range method.Body.Statements {
		if err := c.generateStatement(stmt, method.Name); err != nil {
			return fmt.Errorf("error generating statement: %w", err)
		}
	}

	return nil
}

func (c *CodeGenerator) generateDeclaration(decl Declaration, methodName string) error {
	for _, variable := range decl.Variables {
		if variable.InitialValue != nil {
			// arxikopoihsh: var = initialValue
			if err := c.generateExpression(variable.InitialValue, methodName); err != nil {
				return fmt.Errorf("error generating initial value for variable %s: %w", variable.Name, err)
			}

			// apothikeush apotelesmatos
			varName := c.makeVariableName(methodName, variable.Name)
			varAddr := c.addressMap[varName]
			c.output.WriteString(fmt.Sprintf("        STA   %d\n", varAddr))
		}
	}
	return nil
}

func (c *CodeGenerator) generateStatement(stmt Statement, methodName string) error {
	switch s := stmt.(type) {
	case *ReturnStatement:
		return c.generateReturnStatement(s, methodName)
	case *Assignment:
		return c.generateAssignment(s, methodName)
	case *IfStatement:
		return c.generateIfStatement(s, methodName)
	case *WhileStatement:
		return c.generateWhileStatement(s, methodName)
	case *BreakStatement:
		return c.generateBreakStatement(s)
	case *BlockStatement:
		return c.generateBlock(s.Block, methodName)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

func (c *CodeGenerator) generateReturnStatement(stmt *ReturnStatement, methodName string) error {
	//kanonikh ekfrash return
	if err := c.generateExpression(stmt.Expression, methodName); err != nil {
		return fmt.Errorf("error generating return value: %w", err)
	}

	return nil
}

func (c *CodeGenerator) generateAssignment(stmt *Assignment, methodName string) error {
	// deksia pleura
	if err := c.generateExpression(stmt.Expression, methodName); err != nil {
		return err
	}

	// apothikeush apotelesmatos
	var varAddr int
	var found bool

	// ✅ FIX: First check if it's a variable
	varName := c.makeVariableName(methodName, stmt.Variable)
	if addr, exists := c.addressMap[varName]; exists {
		varAddr = addr
		found = true
	} else {
		// ✅ FIX: If not a variable, check if it's a parameter
		symbolTable := c.symbolTables[methodName]
		paramAddr := c.findParameterByName(methodName, stmt.Variable, symbolTable)
		if paramAddr != -1 {
			varAddr = paramAddr
			found = true
		}
	}

	if !found {
		return fmt.Errorf("variable or parameter '%s' not found in method '%s'", stmt.Variable, methodName)
	}

	c.output.WriteString(fmt.Sprintf("        STA   %d\n", varAddr))

	return nil
}

func (c *CodeGenerator) generateIfStatement(stmt *IfStatement, methodName string) error {
	elseLabel := c.newLabel("ELSE")
	endifLabel := c.newLabel("ENDIF")

	// paragwgh synthikhs
	if err := c.generateExpression(stmt.Condition, methodName); err != nil {
		return fmt.Errorf("error generating if condition: %w", err)
	}

	// sygkrish rA me 0
	c.output.WriteString("        CMPA   =0=\n")

	if stmt.ElseStmt != nil {
		// goto sto else an false
		c.output.WriteString(fmt.Sprintf("        JE    %s\n", elseLabel))
	} else {
		// goto end an den yparxei else
		c.output.WriteString(fmt.Sprintf("        JE    %s\n", endifLabel))
	}

	// paragwgh then
	if err := c.generateStatement(stmt.ThenStmt, methodName); err != nil {
		return fmt.Errorf("error generating if then statement: %w", err)
	}

	if stmt.ElseStmt != nil {
		// paraleipw else
		c.output.WriteString(fmt.Sprintf("        JMP   %s\n", endifLabel))

		// etiketa else
		c.output.WriteString(fmt.Sprintf("%s    NOP\n", elseLabel))

		// ekselesh else
		if err := c.generateStatement(stmt.ElseStmt, methodName); err != nil {
			return fmt.Errorf("error generating if else statement: %w", err)
		}
	}

	c.output.WriteString(fmt.Sprintf("%s    NOP\n", endifLabel))

	return nil
}

func (c *CodeGenerator) generateWhileStatement(stmt *WhileStatement, methodName string) error {
	loopLabel := c.newLabel("LOOP")
	endLabel := c.newLabel("ENDLOOP")

	// append emfoleyumena break labels
	c.breakLabels = append(c.breakLabels, endLabel)

	c.output.WriteString(fmt.Sprintf("%s    NOP\n", loopLabel))

	// paragwgh synthikhs
	if err := c.generateExpression(stmt.Condition, methodName); err != nil {
		return fmt.Errorf("error generating while condition: %w", err)
	}

	// eksodos an false
	c.output.WriteString("        CMPA   =0=\n")
	c.output.WriteString(fmt.Sprintf("        JE    %s\n", endLabel))

	// paragwgh body
	if err := c.generateStatement(stmt.Body, methodName); err != nil {
		return fmt.Errorf("error generating while body: %w", err)
	}

	// goto elegxo synthikhs
	c.output.WriteString(fmt.Sprintf("        JMP   %s\n", loopLabel))

	// telos loop
	c.output.WriteString(fmt.Sprintf("%s    NOP\n", endLabel))

	// afairesh break label apo stack
	c.breakLabels = c.breakLabels[:len(c.breakLabels)-1]

	return nil
}

func (c *CodeGenerator) generateBreakStatement(_ *BreakStatement) error {
	if len(c.breakLabels) == 0 {
		return fmt.Errorf("break statement outside of loop")
	}

	// goto plhsiestoro brongxo
	breakLabel := c.breakLabels[len(c.breakLabels)-1]
	c.output.WriteString(fmt.Sprintf("        JMP   %s\n", breakLabel))

	return nil
}

func (c *CodeGenerator) generateBlock(block Block, methodName string) error {
	// paragwgh dhlwsewn
	for _, decl := range block.Declarations {
		if err := c.generateDeclaration(decl, methodName); err != nil {
			return fmt.Errorf("error generating block declaration: %w", err)
		}
	}

	// paragwgh entolwn
	for _, stmt := range block.Statements {
		if err := c.generateStatement(stmt, methodName); err != nil {
			return fmt.Errorf("error generating block statement: %w", err)
		}
	}

	return nil
}

func (c *CodeGenerator) generateExpression(expr Expression, methodName string) error {
	symbolTable := c.symbolTables[methodName]

	switch e := expr.(type) {
	case *NumberLiteral:
		value, err := strconv.Atoi(e.Value)
		if err != nil {
			return fmt.Errorf("invalid number literal '%s': %w", e.Value, err)
		}
		c.output.WriteString(fmt.Sprintf("        LDA   =%d=\n", value))
		return nil

	case *BooleanLiteral:
		value := 0
		if e.Value {
			value = 1
		}
		c.output.WriteString(fmt.Sprintf("        LDA   =%d=\n", value))
		return nil

	case *Identifier:
		varName := c.makeVariableName(methodName, e.Name)
		if varAddr, exists := c.addressMap[varName]; exists {
			c.output.WriteString(fmt.Sprintf("        LDA   %d\n", varAddr))
			return nil
		}
		// xrhshmopoiw symbol table gia parametrous
		paramAddr := c.findParameterByName(methodName, e.Name, symbolTable)
		if paramAddr != -1 {
			c.output.WriteString(fmt.Sprintf("        LDA   %d\n", paramAddr))
			return nil
		}
		return fmt.Errorf("undefined variable or parameter '%s' in method '%s'", e.Name, methodName)

	case *BinaryExpression:
		return c.generateBinaryExpression(e, methodName)

	case *UnaryExpression:
		return c.generateUnaryExpression(e, methodName)

	case *MethodCall:
		return c.generateMethodCall(e, methodName)

	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (c *CodeGenerator) generateBinaryExpression(expr *BinaryExpression, methodName string) error {
	if leftIdent, ok := expr.Left.(*Identifier); ok {
		if rightIdent, ok := expr.Right.(*Identifier); ok {
			symbolTable := c.symbolTables[methodName]

			leftAddr := c.findParameterByName(methodName, leftIdent.Name, symbolTable)
			if leftAddr == -1 {
				// an oxi parameter, psaxnoume gia variable
				varName := c.makeVariableName(methodName, leftIdent.Name)
				if addr, exists := c.addressMap[varName]; exists {
					leftAddr = addr
				}
			}

			rightAddr := c.findParameterByName(methodName, rightIdent.Name, symbolTable)
			if rightAddr == -1 {
				// an oxi parameter, psaxnoume gia variable
				varName := c.makeVariableName(methodName, rightIdent.Name)
				if addr, exists := c.addressMap[varName]; exists {
					rightAddr = addr
				}
			}

			if leftAddr != -1 && rightAddr != -1 {
				c.output.WriteString(fmt.Sprintf("        LDA   %d\n", leftAddr))

				switch expr.Operator {
				case "+":
					c.output.WriteString(fmt.Sprintf("        ADD   %d\n", rightAddr))
				case "-":
					c.output.WriteString(fmt.Sprintf("        SUB   %d\n", rightAddr))
				case "*":
					c.output.WriteString(fmt.Sprintf("        MUL   %d\n", rightAddr))
				case "/":
					c.output.WriteString(fmt.Sprintf("        DIV   %d\n", rightAddr))
				case "==", "!=", "<", "<=", ">", ">=":
					return c.generateComparison(expr.Operator, rightAddr)
				}
				return nil
			}
		}
	}

	// aristerh pleura
	if err := c.generateExpression(expr.Left, methodName); err != nil {
		return fmt.Errorf("error generating left expression: %w", err)
	}

	// apothikeush aristerou apotelesmatos proswrina
	leftTemp := c.allocateTemp()
	c.output.WriteString(fmt.Sprintf("        STA   %d\n", leftTemp))

	// deksia pleura
	if err := c.generateExpression(expr.Right, methodName); err != nil {
		return fmt.Errorf("error generating right expression: %w", err)
	}

	// apothikeush deksiou apotelesmatos proswrina
	rightTemp := c.allocateTemp()
	c.output.WriteString(fmt.Sprintf("        STA   %d\n", rightTemp))

	// fortwsh aristerou apotelesmatos sto rA
	c.output.WriteString(fmt.Sprintf("        LDA   %d\n", leftTemp))

	// praksh
	switch expr.Operator {
	case "+":
		c.output.WriteString(fmt.Sprintf("        ADD   %d\n", rightTemp))
	case "-":
		c.output.WriteString(fmt.Sprintf("        SUB   %d\n", rightTemp))
	case "*":
		c.output.WriteString(fmt.Sprintf("        MUL   %d\n", rightTemp))
	case "/":
		c.output.WriteString(fmt.Sprintf("        DIV   %d\n", rightTemp))
	case "==", "!=", "<", "<=", ">", ">=":
		return c.generateComparison(expr.Operator, rightTemp)
	default:
		return fmt.Errorf("unsupported operator: %s", expr.Operator)
	}

	return nil
}

func (c *CodeGenerator) generateComparison(op string, rightAddr int) error {
	trueLabel := c.newLabel("TRUE")
	endLabel := c.newLabel("ENDCMP")

	c.output.WriteString(fmt.Sprintf("        CMPA   %d\n", rightAddr))

	// goto vash apotelesmatos
	switch op {
	case "==":
		c.output.WriteString(fmt.Sprintf("        JE    %s\n", trueLabel))
	case "!=":
		c.output.WriteString(fmt.Sprintf("        JNE   %s\n", trueLabel))
	case "<":
		c.output.WriteString(fmt.Sprintf("        JL    %s\n", trueLabel))
	case "<=":
		c.output.WriteString(fmt.Sprintf("        JLE   %s\n", trueLabel))
	case ">":
		c.output.WriteString(fmt.Sprintf("        JG    %s\n", trueLabel))
	case ">=":
		c.output.WriteString(fmt.Sprintf("        JGE   %s\n", trueLabel))
	}

	// false: fortwsh 0
	c.output.WriteString("        LDA   =0=\n")
	c.output.WriteString(fmt.Sprintf("        JMP   %s\n", endLabel))

	// true: fortwsh 1
	c.output.WriteString(fmt.Sprintf("%s    LDA   =1=\n", trueLabel))
	c.output.WriteString(fmt.Sprintf("%s    NOP\n", endLabel))

	return nil
}

func (c *CodeGenerator) generateUnaryExpression(expr *UnaryExpression, methodName string) error {
	// telesths
	if err := c.generateExpression(expr.Operand, methodName); err != nil {
		return fmt.Errorf("error generating unary expression: %w", err)
	}

	switch expr.Operator {
	case "-":
		// arithmitikh arnhsh ( -x = 0 - x )
		tempAddr := c.allocateTemp()
		c.output.WriteString(fmt.Sprintf("        STA   %d\n", tempAddr))
		c.output.WriteString("        LDA   =0=\n")
		c.output.WriteString(fmt.Sprintf("        SUB   %d\n", tempAddr))

	case "!":
		// logikh arnhsh
		// an == 0 tote 1 alliws pali 0
		trueLabel := c.newLabel("TRUE")
		endLabel := c.newLabel("ENDNOT")

		c.output.WriteString("        CMPA   =0=\n")
		c.output.WriteString(fmt.Sprintf("        JE    %s\n", trueLabel))

		// != 0 tote apotelesma 0
		c.output.WriteString("        LDA   =0=\n")
		c.output.WriteString(fmt.Sprintf("        JMP   %s\n", endLabel))

		// == 0 tote apotelesma 1
		c.output.WriteString(fmt.Sprintf("%s    LDA   =1=\n", trueLabel))
		c.output.WriteString(fmt.Sprintf("%s    NOP\n", endLabel))
	}
	return nil
}

func (c *CodeGenerator) generateMethodCall(expr *MethodCall, methodName string) error {
	for i, arg := range expr.Arguments {
		if err := c.generateExpression(arg, methodName); err != nil {
			return err
		}

		paramAddr := c.getParameterAddress(expr.Name, i)
		c.output.WriteString(fmt.Sprintf("        STA   %d\n", paramAddr))
	}

	methodLabel := c.methodLabels[expr.Name]
	c.output.WriteString(fmt.Sprintf("        JMP   %s\n", methodLabel))

	return nil
}

/* func (c *CodeGenerator) generateVariableStorage() {
	allAddresses := make(map[int]bool)
	for _, address := range c.addressMap {
		allAddresses[address] = true
	}
	for i := 1; i < c.tempCounter; i++ {
		tempAddr := TEMP_START + i
		allAddresses[tempAddr] = true
	}

	addresses := make([]int, 0, len(allAddresses))
	for addr := range allAddresses {
		addresses = append(addresses, addr)
	}

	sort.Ints(addresses)

	for _, address := range addresses {
		c.output.WriteString(fmt.Sprintf("        ORIG  %d\n", address))
		c.output.WriteString("        CON   0\n")
	}
} */

func (c *CodeGenerator) generateFooter() {
	c.output.WriteString("        END   MAIN")
}

// HELPERS

func (c *CodeGenerator) newLabel(prefix string) string {
	label := fmt.Sprintf("%s%d", prefix, c.labelCounter)
	c.labelCounter++
	return label
}

func (c *CodeGenerator) allocateTemp() int {
	tempAddr := TEMP_START + c.tempCounter
	c.tempCounter++

	return tempAddr
}

func (c *CodeGenerator) makeParameterName(methodName string, index int) string {
	return fmt.Sprintf("%s_param_%d", methodName, index)
}

func (c *CodeGenerator) makeVariableName(methodName, varName string) string {
	return fmt.Sprintf("%s_%s", methodName, varName)
}

func (c *CodeGenerator) findParameterByName(methodName, paramName string, symbolTable *SymbolTable) int {
	if symbol, exists := symbolTable.Symbols[paramName]; exists && symbol.Kind == "parameter" {
		paramKey := c.makeParameterName(methodName, symbol.Offset)
		if addr, exists := c.addressMap[paramKey]; exists {
			return addr
		}
	}
	return -1
}

func (c *CodeGenerator) getParameterAddress(methodName string, index int) int {
	paramName := c.makeParameterName(methodName, index)
	if addr, exists := c.addressMap[paramName]; exists {
		return addr
	}

	addr := c.currentAddress
	c.addressMap[paramName] = addr
	c.currentAddress++
	return addr
}
