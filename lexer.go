package gosqlparser

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
	tokenIdentifier
	tokenEnd
)

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
	for state := lexIdentifier; state != nil; {
		state = state(l)
	}

	close(l.tokens)
}

// produce sends the token.
func (l *lexer) produce(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.position]}
	l.start = l.position
}

func lexIdentifier(l *lexer) stateFunc {

	return nil
}
