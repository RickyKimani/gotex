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
			// Handle LaTeX commands like \alpha or \frac
			cmdStart := pos + 1
			pos++
			for pos < len(text) && isAlpha(text[pos]) {
				pos++
			}
			cmdName := text[cmdStart:pos]

			if symbol, exists := symbols.ConvertMathSymbol(cmdName); exists {
				nodes = append(nodes, &MathSymbol{Symbol: symbol, Command: cmdName, Position: tokenPos})
			} else {
				// Handle commands that might have arguments
				switch cmdName {
				case "frac":
					// Parse two arguments for fraction
					if pos < len(text) && text[pos] == '{' {
						numerator, newPos := p.parseBracedMathExpression(text, pos, tokenPos)
						pos = newPos

						if pos < len(text) && text[pos] == '{' {
							denominator, newPos := p.parseBracedMathExpression(text, pos, tokenPos)
							pos = newPos

							nodes = append(nodes, &MathFraction{
								Numerator:   numerator,
								Denominator: denominator,
								Position:    tokenPos,
							})
						} else {
							// Missing second argument, treat as regular command
							nodes = append(nodes, &Command{Name: cmdName, Position: tokenPos})
							if numerator != nil {
								nodes = append(nodes, numerator)
							}
						}
					} else {
						// Missing arguments, treat as regular command
						nodes = append(nodes, &Command{Name: cmdName, Position: tokenPos})
					}
				case "sqrt":
					// Parse one argument for square root
					if pos < len(text) && text[pos] == '{' {
						argument, newPos := p.parseBracedMathExpression(text, pos, tokenPos)
						pos = newPos

						nodes = append(nodes, &Command{
							Name:     cmdName,
							Args:     []Node{argument},
							Position: tokenPos,
						})
					} else {
						// Missing argument, treat as regular command
						nodes = append(nodes, &Command{Name: cmdName, Position: tokenPos})
					}
				default:
					// For other unknown commands, just store by name
					nodes = append(nodes, &Command{Name: cmdName, Position: tokenPos})
				}
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
		return &Group{Nodes: parsedContent, Position: tokenPos}, len(text)
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

	return &Group{Nodes: parsedContent, Position: tokenPos}, pos
}

// isMathSpecialChar checks for characters that have special meaning in our math parser.
func isMathSpecialChar(char byte) bool {
	return char == '^' || char == '_' || char == '\\'
}

// isAlpha checks if a character is a letter.
func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}
