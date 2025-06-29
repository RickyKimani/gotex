package parser

import "github.com/rickykimani/gotex/lexer"

func (p *Parser) parseArgumentWithStartPos(startPos lexer.Position) Node {
	// Parse content inside braces until we hit the closing brace
	var nodes []Node

	for p.curToken.Type != lexer.TokenRBrace && p.curToken.Type != lexer.TokenEOF {
		node := p.parseNode()
		if node != nil {
			nodes = append(nodes, node)
		}
		p.nextToken()
	}

	// Check if we found the closing brace
	if p.curToken.Type == lexer.TokenEOF {
		// Report error at the starting position, not current position
		p.addErrorAtPosition(UnmatchedBrace, "Missing closing brace '}' for argument", Error, startPos)
		// Return what we parsed so far
	}

	// If we only have one text node, return it directly
	// Otherwise, we might need a compound node structure
	if len(nodes) == 1 {
		return nodes[0]
	} else if len(nodes) > 1 {
		// If there are multiple nodes, return a Group to preserve structure
		return &Group{Nodes: nodes, Position: startPos}
	}

	return nil
}
