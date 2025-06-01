# Go MIXAL Compiler

A compiler written in Go that translates a simple procedural programming language to MIXAL assembly code for Donald Knuth's MIX computer architecture.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Language Specification](#language-specification)
- [Usage](#usage)
- [Testing](#testing)

## Overview

This project implements a complete compiler pipeline that transforms source code written in a C-like procedural language into MIXAL (MIX Assembly Language) code. The compiler demonstrates fundamental compiler construction principles including lexical analysis, parsing, semantic analysis, and code generation.

The target architecture is the MIX computer, a hypothetical computer designed by Donald Knuth for educational purposes in "The Art of Computer Programming" series.

## Features

- **Complete Compilation Pipeline**: Lexer → Parser → Semantic Analyzer → Code Generator
- **C-like Syntax**: Familiar C-style syntax with functions and control structures
- **Type System**: Static typing with integer types
- **Procedural Programming**: Function definitions with parameters and local variables
- **Control Flow**: Support for if-else statements, while loops, and function calls
- **Expression Evaluation**: Arithmetic expressions with proper operator precedence
- **Error Handling**: Comprehensive error reporting for lexical, syntactic, and semantic errors
- **MIXAL Output**: Generates optimized MIXAL assembly code

## Architecture

The compiler is structured into several key components:

### 1. Lexical Analyzer (`lexer.go`)
- Tokenizes source code into meaningful symbols
- Handles keywords, identifiers, operators, and literals
- Provides position tracking for error reporting

### 2. Parser (`parser.go`)
- Implements recursive descent parsing
- Builds Abstract Syntax Tree (AST) from tokens
- Handles operator precedence and associativity

### 3. Abstract Syntax Tree (`ast.go`)
- Defines node types for program structure
- Represents classes, methods, statements, and expressions
- Supports visitor pattern for tree traversal

### 4. Semantic Analyzer (`semantic.go`)
- Type checking and validation
- Symbol table management
- Scope resolution and variable binding

### 5. Code Generator (`codegen.go`)
- Translates AST to MIXAL assembly
- Handles memory allocation and addressing
- Implements method calling conventions

## Language Specification

### Grammar

```
Program     ::= MethodDecl*
MethodDecl  ::= 'int' IDENTIFIER '(' ParamList? ')' Block
ParamList   ::= 'int' IDENTIFIER (',' 'int' IDENTIFIER)*
Block       ::= '{' Declaration* Statement* '}'
Declaration ::= 'int' IDENTIFIER ('=' Expression)? ';'
Statement   ::= Assignment | IfStmt | WhileStmt | ReturnStmt | BreakStmt | Block ';'
Assignment  ::= IDENTIFIER '=' Expression ';'
IfStmt      ::= 'if' '(' Expression ')' Statement ('else' Statement)?
WhileStmt   ::= 'while' '(' Expression ')' Statement
BreakStmt   ::= 'break' ';'
ReturnStmt  ::= 'return' Expression? ';'
Expression  ::= RelExpr
RelExpr     ::= AddExpr (('<' | '>' | '<=' | '>=' | '==' | '!=') AddExpr)*
AddExpr     ::= MulExpr (('+' | '-') MulExpr)*
MulExpr     ::= UnaryExpr (('*' | '/') UnaryExpr)*
UnaryExpr   ::= '-' UnaryExpr | PrimaryExpr
PrimaryExpr ::= INTEGER | IDENTIFIER | MethodCall | '(' Expression ')'
MethodCall  ::= IDENTIFIER '(' ArgList? ')'
ArgList     ::= Expression (',' Expression)*
```

### Language Features

- **Functions**: C-style function definitions with parameters and return values
- **Types**: Only `int` type is supported
- **Variables**: Local variable declarations with optional initialization
- **Control Flow**: `if-else`, `while` loops, `break` statements
- **Expressions**: Arithmetic and relational expressions with standard operator precedence
- **Function Calls**: Functions can call other functions with parameters

### Lexical Elements

- **Keywords**: `int`, `if`, `else`, `while`, `return`, `break`, `true`, `false`
- **Operators**: `+`, `-`, `*`, `/`, `=`, `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Delimiters**: `(`, `)`, `{`, `}`, `;`, `,`
- **Literals**: Integer constants (positive integers, no leading zeros)
- **Identifiers**: Alphanumeric sequences starting with letter or underscore

## Installation

### Prerequisites

- Go 1.18 or higher
- Git

### Clone and Build

```bash
git clone <repository-url>
cd go_mixal_compiler
go build -o mixal_compiler
```

## Usage

### Basic Compilation

```bash
./mixal_compiler_go input.txt
```

The compiler automatically generates an output file with the same name but `.mixal` extension.

### Command Line

```bash
./mixal_compiler_go <input-file>
```

**Arguments:**
- `input-file`: Source code file to compile (output file is automatically generated)

## Testing

### Running Tests

The `examples/` directory contains test cases:

```bash
# Test simple arithmetic
./mixal_compiler_go examples/simple.txt

# Test complex expressions
./mixal_compiler_go examples/complex.txt

# Test multiplication
./mixal_compiler_go examples/multiply.txt

# Test functions without parameters
./mixal_compiler_go examples/noparams.txt
```

### Verification

To verify the generated MIXAL code:

1. Use a MIX emulator (like the [online MIX emulator](https://www.mix-emulator.org/))
2. Load the generated `.mixal` file
3. Run and check the output
4. Compare with expected results