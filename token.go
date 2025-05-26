package main

import "fmt"

type TokenType int

const (
	//literals
	TOK_ID TokenType = iota
	TOK_NUM
	TOK_TRUE
	TOK_FALSE

	// deysmemenes lekseis
	TOK_INT
	TOK_RETURN
	TOK_IF
	TOK_ELSE
	TOK_WHILE
	TOK_BREAK

	// operatos
	TOK_ASSIGN   // =
	TOK_PLUS     // +
	TOK_MINUS    // -
	TOK_MULTIPLY // *
	TOK_DIVIDE   // /

	// relational ops
	TOK_LT // <
	TOK_LE // <=
	TOK_GT // >
	TOK_GE // >=
	TOK_EQ // ==
	TOK_NE // !=

	//diaxwristika
	TOK_LPAREN    // (
	TOK_RPAREN    // )
	TOK_LBRACE    // {
	TOK_RBRACE    // }
	TOK_COMMA     // ,
	TOK_SEMICOLON // ;

	TOK_EOF   // EOF
	TOK_ERROR // ERROR
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// DEBUG
func (t Token) String() string {
	return fmt.Sprintf("<%s: '%s' at %d:%d>",
		tokenTypeNames[t.Type], t.Value, t.Line, t.Column)
}

var tokenTypeNames = map[TokenType]string{
	TOK_ID:        "ID",
	TOK_NUM:       "NUM",
	TOK_TRUE:      "TRUE",
	TOK_FALSE:     "FALSE",
	TOK_INT:       "INT",
	TOK_RETURN:    "RETURN",
	TOK_IF:        "IF",
	TOK_ELSE:      "ELSE",
	TOK_WHILE:     "WHILE",
	TOK_BREAK:     "BREAK",
	TOK_ASSIGN:    "ASSIGN",
	TOK_PLUS:      "PLUS",
	TOK_MINUS:     "MINUS",
	TOK_MULTIPLY:  "MULTIPLY",
	TOK_DIVIDE:    "DIVIDE",
	TOK_LT:        "LT",
	TOK_LE:        "LE",
	TOK_GT:        "GT",
	TOK_GE:        "GE",
	TOK_EQ:        "EQ",
	TOK_NE:        "NE",
	TOK_LPAREN:    "LPAREN",
	TOK_RPAREN:    "RPAREN",
	TOK_LBRACE:    "LBRACE",
	TOK_RBRACE:    "RBRACE",
	TOK_COMMA:     "COMMA",
	TOK_SEMICOLON: "SEMICOLON",
	TOK_EOF:       "EOF",
	TOK_ERROR:     "ERROR",
}

var keywords = map[string]TokenType{
	"int":    TOK_INT,
	"return": TOK_RETURN,
	"if":     TOK_IF,
	"else":   TOK_ELSE,
	"while":  TOK_WHILE,
	"break":  TOK_BREAK,
	"true":   TOK_TRUE,
	"false":  TOK_FALSE,
}
