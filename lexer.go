package gosqlparser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// lexer tokenizes the input and produces tokens.
type lexer struct {
	input string // raw SQL query

	start    int // start position of the current token
	position int // current position
	width    int // width of last rune read from input

	tokens chan token
}

// token is an entity produced by tokenizer for parser that represents a smaller typed piece
// of input string.
type token struct {
	t     tokenType
	value string
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenSpace
	tokenIdentifier
	tokenDot
	tokenField
	tokenEnd

	// tokenEq  // "=="
	// tokenLt  // "<"
	// tokenLte // "<="
	// tokenGt  // ">"
	// tokenGte // ">="

	tokenSelect // "SELECT"
	tokenFrom   // "FROM"
	tokenWhere  // "WHERE"
	tokenLimit  // "LIMIT"
	tokenAlias  // "AS"
)

const (
	keywordSelect = "SELECT"
	keywordFrom   = "FROM"
	keywordWhere  = "WHERE"
	keywordLimit  = "LIMIT"
	keywordAlias  = "AS"
)

var keywords = map[string]tokenType{
	keywordSelect: tokenSelect,
	keywordFrom:   tokenFrom,
	keywordWhere:  tokenWhere,
	keywordAlias:  tokenAlias,
	keywordLimit:  tokenLimit,
}

const end = -1

type stateFunc func(*lexer) stateFunc

func lex(input string) <-chan token {
	l := newLexer(input)

	go l.run()

	return l.tokens
}

// newLexer returns an instance of the new lexer.
func newLexer(input string) *lexer {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}

	return l
}

func (l *lexer) run() {
	for state := lexInit; state != nil; {
		state = state(l)
	}

	close(l.tokens)
}

// produce sends the token.
func (l *lexer) produce(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.position]}
	l.start = l.position
}

func (l *lexer) next() rune {
	if l.position >= len(l.input) {
		l.width = 0

		return end
	}

	r, w := utf8.DecodeRuneInString(l.input[l.position:])

	l.width = w
	l.position += w

	return r
}

func (l *lexer) revert() {
	l.position -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.revert()

	return r
}

func lexInit(l *lexer) stateFunc {
	r := l.next()

	switch true {
	case isAlphaNumeric(r):
		return lexIdentifier
	case unicode.IsSpace(r):
		l.produce(tokenSpace)

		return lexInit
	case r == end:

		l.produce(tokenEnd)
		return nil
	}

	// TODO: resolve
	panic("unreachable")
}

func lexIdentifier(l *lexer) stateFunc {
	r := l.next()

	if isAlphaNumeric(r) {
		// advance
		return lexIdentifier
	}

	l.revert()

	word := l.input[l.start:l.position]
	if t, ok := keywords[strings.ToUpper(word)]; ok {
		l.produce(t)
	} else {
		l.produce(tokenIdentifier)
	}

	return lexInit
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
