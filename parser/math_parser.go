package parser

import (
	"github.com/rickykimani/gotex/lexer"
	"github.com/rickykimani/gotex/symbols"
)

// parseMath is the entry point for parsing a math node.
// It receives a token like $...$ and parses the content.
func (p *Parser) parseMath() *MathNode {
	isInline := p.curToken.Type == lexer.TokenMathInline
	rawContent := p.curToken.Value

	// The recursive math parser starts here.
	nodes, _ := p.parseMathExpression(rawContent, 0, p.curToken.Pos)

	return &MathNode{
		Inline:   isInline,
		Content:  nodes,
		Position: p.curToken.Pos,
	}
}

// parseMathExpression is the core of the math parser.
// It parses a string of math content and returns a list of nodes.
func (p *Parser) parseMathExpression(text string, startPos int, tokenPos lexer.Position) ([]Node, int) {
	pos := startPos
	var nodes []Node

	for pos < len(text) {
		ch := text[pos]

		switch ch {
		case '^', '_':
			var base Node
			if len(nodes) > 0 {
				base = nodes[len(nodes)-1]
				nodes = nodes[:len(nodes)-1] // Pop the base from the node list
			}

			// The script content parser will handle single chars or {...}
			scriptContent, newPos := p.parseScriptContent(text, pos+1, tokenPos)
			pos = newPos

			if ch == '^' {
				nodes = append(nodes, &MathSuperscript{Base: base, Exponent: scriptContent, Position: tokenPos})
			} else {
				nodes = append(nodes, &MathSubscript{Base: base, Index: scriptContent, Position: tokenPos})
			}

		case '\\':
			// Handle LaTeX commands like \alpha
			cmdStart := pos + 1
			pos++
			for pos < len(text) && isAlpha(text[pos]) {
				pos++
			}
			cmdName := text[cmdStart:pos]
			if symbol, exists := symbols.ConvertMathSymbol(cmdName); exists {
				nodes = append(nodes, &MathSymbol{Symbol: symbol, Command: cmdName, Position: tokenPos})
			} else {
				// For now, unknown commands are just stored by name.
				// A fuller implementation might parse arguments.
				nodes = append(nodes, &Command{Name: cmdName, Position: tokenPos})
			}

		default:
			// Regular text (e.g., "a", "1", "+")
			textStart := pos
			for pos < len(text) && !isMathSpecialChar(text[pos]) {
				pos++
			}
			value := text[textStart:pos]
			if value != "" {
				nodes = append(nodes, &TextNode{Value: value, Position: tokenPos})
			}
		}
	}

	return nodes, pos
}

// parseScriptContent determines if a script is a single character or a braced group.
func (p *Parser) parseScriptContent(text string, startPos int, tokenPos lexer.Position) (Node, int) {
	if startPos >= len(text) {
		return &TextNode{Value: "", Position: tokenPos}, startPos // Empty script
	}

	if text[startPos] == '{' {
		// Content is a braced expression, e.g., {n+1}
		return p.parseBracedMathExpression(text, startPos, tokenPos)
	}

	// Content is a single character, e.g., ^2
	return &TextNode{Value: string(text[startPos]), Position: tokenPos}, startPos + 1
}

// parseBracedMathExpression handles expressions inside {}.
// It recursively calls the main math parser on the content.
func (p *Parser) parseBracedMathExpression(text string, startPos int, tokenPos lexer.Position) (Node, int) {
	pos := startPos + 1 // Skip opening '{'
	braceCount := 1
	contentStart := pos

	// Find the matching closing brace to support nested braces
	for pos < len(text) && braceCount > 0 {
		switch text[pos] {
		case '{':
			braceCount++
		case '}':
			braceCount--
		}
		if braceCount > 0 {
			pos++
		}
	}

	if braceCount > 0 {
		// Unclosed brace, treat everything until the end as content.
		// This is an error recovery strategy.
		content := text[contentStart:]
		parsedContent, _ := p.parseMathExpression(content, 0, tokenPos)
		return &Group{Content: parsedContent, Position: tokenPos}, len(text)
	}

	// Recursively parse the content within the braces
	content := text[contentStart:pos]
	parsedContent, _ := p.parseMathExpression(content, 0, tokenPos)

	pos++ // Skip closing '}'

	// If the braced expression contained only one element, just return that element.
	// Otherwise, return a group.
	if len(parsedContent) == 1 {
		return parsedContent[0], pos
	}

	return &Group{Content: parsedContent, Position: tokenPos}, pos
}

// isMathSpecialChar checks for characters that have special meaning in our math parser.
func isMathSpecialChar(char byte) bool {
	return char == '^' || char == '_' || char == '\\'
}

// isAlpha checks if a character is a letter.
func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}
