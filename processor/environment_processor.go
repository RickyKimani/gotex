package processor

import (
	"fmt"
	"strings"

	"github.com/rickykimani/gotex/lexer"
	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/symbols"
)

//TODO: Modularize

func (dp *DocumentProcessor) processEnvironment(env *parser.Environment, style string) {
	switch env.Name {
	case "document":
		dp.processNodes(env.Body, style)

	case "itemize":
		dp.currentY += dp.lineHeight * 0.8
		dp.enterList("itemize")
		dp.processNodes(env.Body, style)
		dp.exitList()

	case "enumerate":
		dp.currentY += dp.lineHeight * 0.8
		dp.enterList("enumerate")
		dp.processNodes(env.Body, style)
		dp.exitList()

	case "center":
		// Add some space before and after
		dp.addVerticalSpace(10)
		dp.processNodes(env.Body, style)
		dp.addVerticalSpace(10)

	case "equation":
		dp.addVerticalSpace(dp.lineHeight * 0.5) // Add some space before
		// Increment equation counter
		dp.equationCounter++

		// Extract raw text from environment content and re-parse as math
		rawContent := dp.extractRawTextFromNodes(env.Body)
		if strings.TrimSpace(rawContent) != "" {
			// Calculate equation number string
			var equationNumber string
			if dp.sectionCounter > 0 {
				equationNumber = fmt.Sprintf("(%d.%d)", dp.sectionCounter, dp.equationCounter)
			} else {
				equationNumber = fmt.Sprintf("(%d)", dp.equationCounter)
			}

			// Re-parse as math content like inline math does
			mathNode := dp.parseMathContent(rawContent, false) // false = display math

			// Calculate actual math width for proper centering
			contentWidth := dp.generator.GetContentWidth()
			mathWidth := dp.mathProcessor.CalculateMathWidth(mathNode.Content)
			centerX := dp.generator.MarginLeft + (contentWidth-mathWidth)/2

			// Render the equation at center
			dp.mathProcessor.ProcessMathNode(mathNode, centerX, dp.currentY)

			// Add equation number on the right with proper positioning
			numberWidth := dp.generator.GetTextWidth(equationNumber, dp.fontSize, "normal")
			numberX := dp.generator.MarginLeft + contentWidth - numberWidth
			dp.generator.AddText(equationNumber, numberX, dp.currentY, dp.fontSize, "normal")
		}
		dp.newLine()                             // Ensure we're on a new line after the equation
		dp.addVerticalSpace(dp.lineHeight * 0.5) // Add some space after

	default:
		dp.processNodes(env.Body, style)
	}
}

// extractRawTextFromNodes extracts raw text from a slice of nodes, preserving LaTeX syntax
func (dp *DocumentProcessor) extractRawTextFromNodes(nodes []parser.Node) string {
	var result strings.Builder
	for _, node := range nodes {
		result.WriteString(dp.extractRawMathText(node))
	}
	return result.String()
}

// extractRawMathText extracts text while preserving LaTeX command syntax
func (dp *DocumentProcessor) extractRawMathText(node parser.Node) string {
	switch n := node.(type) {
	case *parser.TextNode:
		return n.Value
	case *parser.Group:
		var result strings.Builder
		for _, child := range n.Nodes {
			result.WriteString(dp.extractRawMathText(child))
		}
		return result.String()
	case *parser.Command:
		// Preserve the command with backslash for math processing
		var result strings.Builder
		result.WriteString("\\" + n.Name)
		// Add arguments if any
		for _, arg := range n.Args {
			result.WriteString("{")
			result.WriteString(dp.extractRawMathText(arg))
			result.WriteString("}")
		}
		return result.String()
	default:
		return ""
	}
}

// parseMathContent re-parses raw text as math content (similar to parser's parseMath)
func (dp *DocumentProcessor) parseMathContent(rawContent string, inline bool) *parser.MathNode {
	// Tokenize the math content
	tokens := lexer.NewLexer(rawContent).Tokenize()

	var content []parser.Node
	for _, tok := range tokens {
		switch tok.Type {
		case lexer.TokenCommand:
			// Check for math symbol commands in the central symbol table
			if symbol, exists := symbols.ConvertMathSymbol(tok.Value); exists {
				content = append(content, &parser.MathSymbol{Symbol: symbol, Command: tok.Value, Position: tok.Pos})
			} else {
				// Regular command node
				content = append(content, &parser.Command{Name: tok.Value, Args: []parser.Node{}, Optional: []parser.Node{}, Position: tok.Pos})
			}
		case lexer.TokenText:
			valRaw := tok.Value
			pos := 0
			for pos < len(valRaw) {
				// Find next ^ or _
				idx := strings.IndexAny(valRaw[pos:], "^_")
				if idx < 0 {
					// No special, append remaining text
					rest := valRaw[pos:]
					if strings.TrimSpace(rest) != "" {
						content = append(content, &parser.TextNode{Value: rest, Position: tok.Pos})
					}
					break
				}
				idx += pos
				// Append text before operator
				if idx > pos {
					seg := valRaw[pos:idx]
					if strings.TrimSpace(seg) != "" {
						content = append(content, &parser.TextNode{Value: seg, Position: tok.Pos})
					}
				}
				// Operator char
				op := valRaw[idx]
				// Following character(s) as exponent/subscript
				if idx+1 < len(valRaw) {
					expChar := string(valRaw[idx+1])
					// Base is last content
					if len(content) > 0 {
						base := content[len(content)-1]
						content = content[:len(content)-1]
						// Create superscript or subscript
						switch op {
						case '^':
							content = append(content, &parser.MathSuperscript{Base: base, Exponent: &parser.TextNode{Value: expChar, Position: tok.Pos}, Position: tok.Pos})
						case '_':
							content = append(content, &parser.MathSubscript{Base: base, Index: &parser.TextNode{Value: expChar, Position: tok.Pos}, Position: tok.Pos})
						}
					}
					pos = idx + 2
				} else {
					pos = idx + 1
				}
			}

		default:
			// ignore other tokens
		}
	}

	return &parser.MathNode{
		Inline:   inline,
		Content:  content,
		Position: lexer.Position{},
	}
}
