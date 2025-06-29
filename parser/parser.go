package parser

import (
	"github.com/rickykimani/gotex/lexer"
)

// NewParser creates a new parser instance (alias for New for compatibility)
func NewParser(tokens []lexer.Token) *Parser {
	// Convert tokens to a simple lexer
	tokenLexer := &TokenLexer{tokens: tokens, pos: 0}
	return NewWithInterface(tokenLexer)
}

// TokenLexer is a simple lexer that works with pre-tokenized input
type TokenLexer struct {
	tokens []lexer.Token
	pos    int
}

func (tl *TokenLexer) NextToken() lexer.Token {
	if tl.pos >= len(tl.tokens) {
		return lexer.Token{Type: lexer.TokenEOF, Value: "", Pos: lexer.Position{}}
	}
	token := tl.tokens[tl.pos]
	tl.pos++
	return token
}

type Parser struct {
	lex       lexer.LexerInterface
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []ParseError // Collect errors during parsing
}

func New(l *lexer.Lexer) *Parser {
	return NewWithInterface(l)
}

func NewWithInterface(l lexer.LexerInterface) *Parser {
	p := &Parser{lex: l}
	// Read two tokens to initialize cur and peek
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

// Parse is an alias for ParseDocument that returns both document and errors
func (p *Parser) Parse() (*Document, []ParseError) {
	doc := p.ParseDocument()
	return doc, p.errors
}
