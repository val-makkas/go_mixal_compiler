package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Compiler struct {
	lexer    *Lexer
	parser   *Parser
	semantic *SemanticAnalyzer
	codegen  *CodeGenerator
	verbose  bool
}

func NewCompiler() *Compiler {
	return &Compiler{
		lexer:    NewLexer(),
		parser:   NewParser(),
		semantic: NewSemanticAnalyzer(),
		codegen:  NewCodeGenerator(),
		verbose:  true,
	}
}

func (c *Compiler) Compile(sourceFile string) error {
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	source := string(content)

	// LEKTIKH ANALYSH
	if c.verbose {
		fmt.Println("Phase 1: Lexical Analysis (Tokenization)")
		fmt.Println("-----------------------------------------")
	}

	tokens, err := c.lexer.Tokenize(source)
	if err != nil {
		return fmt.Errorf("lexical analysis failed: %w", err)
	}

	if c.verbose {
		fmt.Printf("Generated %d tokens\n\n", len(tokens))
		c.printFirstTokens(tokens, 10)
	}

	// SYNTAKTIKH ANALYSH
	if c.verbose {
		fmt.Println("Phase 2: Parsing (Syntax Analysis)")
		fmt.Println("-----------------------------------------")
	}

	ast, err := c.parser.Parse(tokens)
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}
	if c.verbose {
		fmt.Printf("Generated AST with %d methos\n\n", len(ast.Methods))
	}

	// SHMASIOLOGIKH ANALYSH
	if c.verbose {
		fmt.Println("Phase 3: Semantic Analysis")
		fmt.Println("-----------------------------------------")
	}

	symbolTables, err := c.semantic.Analyze(ast)
	if err != nil {
		return fmt.Errorf("semantic analysis failed: %w", err)
	}
	if c.verbose {
		fmt.Printf("Semantic analysis passed\n")
		fmt.Printf("   - Found %d methods\n", len(symbolTables))
		fmt.Printf("   - Main method: ✓\n\n")
	}

	// PARAGOGH KODIKA MIXAL
	if c.verbose {
		fmt.Println("Phase 4: Code Generation (MIXAL)")
		fmt.Println("-----------------------------------------")
	}

	mixalCode, err := c.codegen.Generate(ast, symbolTables)
	if err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}
	if c.verbose {
		fmt.Printf("Generated %d lines of MIXAL code\n", strings.Count(mixalCode, "\n")+1)
	}

	outputFile := c.getOutputFileName(sourceFile)
	if err := os.WriteFile(outputFile, []byte(mixalCode), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	if c.verbose {
		fmt.Printf("Output written to %s\n", outputFile)
	}
	return nil
}

func (c *Compiler) getOutputFileName(sourceFile string) string {
	ext := filepath.Ext(sourceFile)
	base := strings.TrimSuffix(filepath.Base(sourceFile), ext)
	return filepath.Join(filepath.Dir(sourceFile), base+".mixal")
}

func (c *Compiler) printFirstTokens(tokens []Token, count int) {
	if len(tokens) == 0 {
		fmt.Println("No tokens generated.")
		return
	}

	fmt.Printf("First %d tokens:\n", count)
	for i, token := range tokens {
		if i >= count {
			break
		}
		fmt.Printf("  %2d: %s\n", i+1, token)
	}
	fmt.Println()
}
