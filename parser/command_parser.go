package parser

import "github.com/rickykimani/gotex/lexer"

func (p *Parser) parseCommand() *Command {
	cmd := &Command{
		Name:     p.curToken.Value,
		Args:     []Node{},
		Optional: []Node{},
		Position: p.curToken.Pos,
	}

	// Special handling for malformed \begin commands
	if p.curToken.Value == "begin" && p.peekToken.Type != lexer.TokenLBrace {
		p.addError(UnmatchedBrace, "missing '{' after \\begin", Error)
		return cmd
	}

	// Special handling for malformed \end commands
	if p.curToken.Value == "end" && p.peekToken.Type != lexer.TokenLBrace {
		p.addError(UnmatchedBrace, "missing '{' after \\end", Error)
		return cmd
	}

	// Look ahead for optional arguments [...]
	if p.peekToken.Type == lexer.TokenOptionalArg {
		p.nextToken() // move to optional arg token
		optArg := &TextNode{
			Value:    p.curToken.Value,
			Position: p.curToken.Pos,
		}
		cmd.Optional = append(cmd.Optional, optArg)
	}

	// Look ahead for required arguments {...}
	for p.peekToken.Type == lexer.TokenLBrace {
		p.nextToken()              // move to {
		bracePos := p.curToken.Pos // Save the position of the opening brace
		p.nextToken()              // move past {

		// Parse the argument content
		arg := p.parseArgumentWithStartPos(bracePos)
		if arg != nil {
			cmd.Args = append(cmd.Args, arg)
		}
	}

	return cmd
}
