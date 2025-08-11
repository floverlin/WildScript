package lexer

import (
	"strings"
)

type Tokenizer interface {
	NextToken() Token
}

type Lexer struct {
	input   []byte
	pos     int
	readPos int
	ch      byte
	line    int
	column  int
}

func New(input []byte) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' ||
		l.ch == '\n' ||
		l.ch == '\t' ||
		l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	var token Token

	l.skipWhitespace()

	if t, ok := dual[string([]byte{l.ch, l.peekChar()})]; ok {
		token = newToken(t, string(t), l.line, l.column)
		l.readChar()
	} else if t, ok := mono[l.ch]; ok {
		token = newToken(t, string(l.ch), l.line, l.column)
	} else if l.ch == '#' {
		l.skipComment()
		return l.NextToken()
	} else if l.ch == '"' {
		return l.readString()
	} else if isDigit(l.ch) {
		return l.readNumber()
	} else if isLetter(l.ch) {
		return l.readIdentifier()
	} else if l.ch == 0 {
		token = newToken(EOF, "", l.line, l.column)
	} else {
		token = newToken(ILLEGAL, string(l.ch), l.line, l.column)
	}

	l.readChar()
	return token
}

func (l *Lexer) readString() Token {
	line, col := l.line, l.column
	l.readChar()
	var sb strings.Builder
	for l.ch != '"' {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			case 'n':
				sb.WriteByte('\n')
			}
		} else {
			sb.WriteByte(l.ch)
		}

		if l.peekChar() == 0 {
			l.readChar()
			return newToken(ILLEGAL, sb.String(), line, col)
		}

		l.readChar()
	}
	l.readChar()
	return newToken(STRING, sb.String(), line, col)
}

func (l *Lexer) readNumber() Token {
	line, col := l.line, l.column
	start := l.pos
	var dots int
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			dots++
		}
		l.readChar()
	}
	if dots > 1 {
		return newToken(ILLEGAL, string(l.input[start:l.pos]), line, col)
	}
	return newToken(NUMBER, string(l.input[start:l.pos]), line, col)
}

func (l *Lexer) readIdentifier() Token {
	line, col := l.line, l.column
	start := l.pos
	for isDigit(l.ch) || isLetter(l.ch) || l.ch == '_' {
		l.readChar()
	}
	literal := string(l.input[start:l.pos])
	tokenType := lookupIdent(literal)
	return newToken(tokenType, literal, line, col)
}

func lookupIdent(ident string) TokenType {
	if identType, ok := specialIdents[ident]; ok {
		return identType
	}
	return IDENT
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
