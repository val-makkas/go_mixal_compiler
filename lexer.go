package main

import (
	"fmt"
	"regexp"
	"unicode"
)

type Lexer struct {
	input    string // kwdikas
	position int    // index
	line     int    // trexon grammh
	column   int    // trexon sthlh
}

func NewLexer() *Lexer {
	return &Lexer{
		line:   1,
		column: 1,
	}
}

// kyria methodos metatrepei input se tokens
func (l *Lexer) Tokenize(input string) ([]Token, error) {
	l.input = input
	l.position = 0
	l.line = 1
	l.column = 1

	var tokens []Token // lista apo tokens

	for l.position < len(l.input) {
		l.skipWhitespace()
		l.skipComments()

		if l.position >= len(l.input) {
			break
		}

		token, err := l.nextToken()
		if err != nil {
			return nil, err
		}

		if token.Type != TOK_EOF {
			tokens = append(tokens, token)
		}
	}

	tokens = append(tokens, Token{
		Type:   TOK_EOF,
		Value:  "",
		Line:   l.line,
		Column: l.column,
	})

	return tokens, nil
}

// diabazei to epomeno token
func (l *Lexer) nextToken() (Token, error) {
	if l.position >= len(l.input) {
		return Token{Type: TOK_EOF}, nil
	}

	// thesh tou token
	startLine := l.line
	startColumn := l.column
	ch := l.input[l.position]

	// literals
	switch ch {
	case '(':
		l.advance()
		return Token{Type: TOK_LPAREN, Value: "(", Line: startLine, Column: startColumn}, nil
	case ')':
		l.advance()
		return Token{Type: TOK_RPAREN, Value: ")", Line: startLine, Column: startColumn}, nil
	case '{':
		l.advance()
		return Token{Type: TOK_LBRACE, Value: "{", Line: startLine, Column: startColumn}, nil
	case '}':
		l.advance()
		return Token{Type: TOK_RBRACE, Value: "}", Line: startLine, Column: startColumn}, nil
	case ',':
		l.advance()
		return Token{Type: TOK_COMMA, Value: ",", Line: startLine, Column: startColumn}, nil
	case ';':
		l.advance()
		return Token{Type: TOK_SEMICOLON, Value: ";", Line: startLine, Column: startColumn}, nil
	case '+':
		l.advance()
		return Token{Type: TOK_PLUS, Value: "+", Line: startLine, Column: startColumn}, nil
	case '*':
		l.advance()
		return Token{Type: TOK_MULTIPLY, Value: "*", Line: startLine, Column: startColumn}, nil
	case '/':
		l.advance()
		return Token{Type: TOK_DIVIDE, Value: "/", Line: startLine, Column: startColumn}, nil

	case '-': // mporei na einai operator h sign gia arithmo
		return l.handleMinus()
	case '=': // mporei = h ==
		return l.handleEquals()
	case '>': // > h >=
		return l.handleGreaterThan()
	case '<': // < h <=
		return l.handleLessThan()
	case '!': // ! h !=
		return l.handleExclamation()
	}

	if l.isLetter(ch) {
		return l.readIdentifier()
	}

	if unicode.IsDigit(rune(ch)) {
		return l.readNumber()
	}

	return Token{
		Type:   TOK_ERROR,
		Value:  "",
		Line:   startLine,
		Column: startColumn,
	}, fmt.Errorf("unexpected character '%c' at line %d, column %d", ch, startLine, startColumn)
}

func (l *Lexer) handleMinus() (Token, error) {
	startLine := l.line
	startColumn := l.column

	l.advance()

	return Token{TOK_MINUS, "-", startLine, startColumn}, nil
}

func (l *Lexer) handleEquals() (Token, error) {
	startLine := l.line
	startColumn := l.column

	l.advance()

	// elegxoume an exei kai allo = meta
	if l.position < len(l.input) &&
		l.input[l.position] == '=' {
		l.advance()
		return Token{TOK_EQ, "==", startLine, startColumn}, nil
	}

	//alliws einai aplo =
	return Token{TOK_ASSIGN, "=", startLine, startColumn}, nil
}

func (l *Lexer) handleLessThan() (Token, error) {
	startLine := l.line
	startColumn := l.column

	l.advance()

	// elegxoume an exei = meta
	if l.position < len(l.input) &&
		l.input[l.position] == '=' {
		l.advance()
		return Token{TOK_LE, "<=", startLine, startColumn}, nil
	}

	//alliws einai aplo <
	return Token{TOK_LT, "<", startLine, startColumn}, nil
}

func (l *Lexer) handleGreaterThan() (Token, error) {
	startLine := l.line
	startColumn := l.column

	l.advance()

	// elegxoume an exei = meta
	if l.position < len(l.input) &&
		l.input[l.position] == '=' {
		l.advance()
		return Token{TOK_GE, ">=", startLine, startColumn}, nil
	}

	//alliws einai aplo >
	return Token{TOK_GT, ">", startLine, startColumn}, nil
}

func (l *Lexer) handleExclamation() (Token, error) {
	startLine := l.line
	startColumn := l.column

	l.advance()

	// elegxoume an exei = meta
	if l.position < len(l.input) &&
		l.input[l.position] == '=' {
		l.advance()
		return Token{TOK_NE, "!=", startLine, startColumn}, nil
	}

	//alliws einai akyro
	return Token{TOK_ERROR, "", startLine, startColumn},
		fmt.Errorf("unexpected character '!' at line %d, column %d (did you mean '!='?)", startLine, startColumn)
}

// kanonas id = letter (letter | digit | '_")*
func (l *Lexer) readIdentifier() (Token, error) {
	start := l.position
	startLine := l.line
	startColumn := l.column

	// diabazw oso exei grammata, psifia, h _
	for l.position < len(l.input) {
		ch := l.input[l.position]
		if !l.isLetter(ch) && !unicode.IsDigit(rune(ch)) && ch != '_' {
			break
		}
		l.advance()
	}

	value := l.input[start:l.position]
	tokenType := TOK_ID

	// elegxw an einai desmeusmeno
	if keywordType, exists := keywords[value]; exists {
		tokenType = keywordType
	}

	return Token{tokenType, value, startLine, startColumn}, nil
}

// kanonas num = '-'? [1-9] digit*
func (l *Lexer) readNumber() (Token, error) {
	start := l.position
	startLine := l.line
	startColumn := l.column

	// an exei meion to diabazw
	if l.position < len(l.input) && l.input[l.position] == '-' {
		l.advance()
	}

	// airthmoi prepei na arxizoun apo [1-9] oxi apo 0
	if l.position < len(l.input) && l.input[l.position] == '0' {
		if l.position+1 < len(l.input) && unicode.IsDigit(rune(l.input[l.position+1])) {
			return Token{TOK_ERROR, "", startLine, startColumn},
				fmt.Errorf("invalid number format at line %d, column %d (numbers cannot have leading zeros)", startLine, startColumn)
		}
	}

	//diabazw psifia
	for l.position < len(l.input) && unicode.IsDigit(rune(l.input[l.position])) {
		l.advance()
	}

	value := l.input[start:l.position]

	// elsegxoume an einai egkyros arithmos
	matched, _ := regexp.MatchString(`^-?([1-9]\d*|0)$`, value)
	if !matched {
		return Token{TOK_ERROR, "", startLine, startColumn},
			fmt.Errorf("invalid number format '%s' at line %d, column %d", value, startLine, startColumn)
	}

	return Token{TOK_NUM, value, startLine, startColumn}, nil
}

// paraleipw kena, tabs kai newline
func (l *Lexer) skipWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(rune(l.input[l.position])) {
		l.advance()
	}
}

// paraleipw sxolia
func (l *Lexer) skipComments() {
	if l.position < len(l.input)-1 &&
		l.input[l.position] == '/' &&
		l.input[l.position+1] == '/' {
		// paraleipw mexri to newline
		for l.position < len(l.input) && l.input[l.position] != '\n' {
			l.advance()
		}
	}
}

// proxwraei ston epomeno xarakthra
func (l *Lexer) advance() {
	if l.position < len(l.input) &&
		l.input[l.position] == '\n' {
		l.line++     // nea grammh
		l.column = 1 // prwth stylh
	} else {
		l.column++ //apla epomenh stylh
	}
	l.position++
}

// elegxei an einai gramma
func (l *Lexer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
