package parser

import "github.com/rickykimani/gotex/lexer"

// parseNode determines what type of node to parse based on the current token
func (p *Parser) parseNode() Node {
	switch p.curToken.Type {
	case lexer.TokenCommand:
		return p.parseCommand()
	case lexer.TokenBeginEnv:
		return p.parseEnvironment()
	case lexer.TokenText:
		return p.parseText()
	case lexer.TokenMathInline, lexer.TokenMathDisplay:
		return p.parseMath()
	case lexer.TokenComment:
		return p.parseComment()
	case lexer.TokenRBrace:
		// Unmatched closing brace at document level - this is a warning since we can recover
		p.addWarning(UnexpectedToken, "unexpected '}' - no matching '{'")
		return nil
	case lexer.TokenLBrace:
		// Unmatched opening brace at document level - this is an error since it affects parsing
		p.addError(UnexpectedToken, "unexpected '{' - not part of a command", Error)
		return nil
	default:
		// For other unexpected tokens
		return nil
	}
}
