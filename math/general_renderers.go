package math

import (
	"strings"

	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/symbols"
)

// renderMathElement renders a single math element
func (m *MathProcessor) renderMathElement(node parser.Node, x, y, fontSize float64) float64 {
	// Safety check for nil nodes
	if node == nil {
		return 0
	}

	switch n := node.(type) {
	case *parser.TextNode:
		return m.renderMathText(n.Value, x, y, fontSize)

	case *parser.Command:
		return m.renderMathCommand(n, x, y, fontSize)

	case *parser.MathSymbol:
		return m.renderMathSymbol(n, x, y, fontSize)

	case *parser.MathSuperscript:
		return m.renderSuperscript(n, x, y, fontSize)

	case *parser.MathSubscript:
		return m.renderSubscript(n, x, y, fontSize)

	case *parser.MathFraction:
		return m.renderFraction(n, x, y, fontSize)

	case *parser.Group:
		return m.renderGroup(n, x, y, fontSize)

	default:
		// Fallback: treat as text
		if textNode, ok := node.(*parser.TextNode); ok {
			return m.renderMathText(textNode.Value, x, y, fontSize)
		}
		return 0
	}
}

// renderMathText renders mathematical text (variables, numbers, etc.)
func (m *MathProcessor) renderMathText(text string, x, y, fontSize float64) float64 {
	// Clean up the text
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// Choose appropriate font
	fontStyle := GetMathFont(text)

	// Render the text
	m.generator.AddText(text, x, y, fontSize, fontStyle)

	// Return width
	return m.generator.GetTextWidth(text, fontSize, fontStyle)
}

// renderMathCommand renders LaTeX math commands
func (m *MathProcessor) renderMathCommand(cmd *parser.Command, x, y, fontSize float64) float64 {
	// Check if it's a known math symbol
	if symbol, exists := symbols.ConvertMathSymbol(cmd.Name); exists {
		m.generator.AddText(symbol, x, y, fontSize, "normal")
		return m.generator.GetTextWidth(symbol, fontSize, "normal")
	}

	// Handle special commands
	switch cmd.Name {
	case "frac":
		if len(cmd.Args) >= 2 {
			// Create a temporary fraction node
			fracNode := &parser.MathFraction{
				Numerator:   cmd.Args[0],
				Denominator: cmd.Args[1],
				Position:    cmd.Position,
			}
			return m.renderFraction(fracNode, x, y, fontSize)
		}

	case "sqrt":
		if len(cmd.Args) >= 1 {
			return m.renderSquareRoot(cmd.Args[0], x, y, fontSize)
		}

	default:
		// For unknown commands, just render the arguments
		currentX := x
		for _, arg := range cmd.Args {
			width := m.renderMathElement(arg, currentX, y, fontSize)
			currentX += width
		}
		return currentX - x
	}

	return 0
}

// renderMathSymbol renders a math symbol
func (m *MathProcessor) renderMathSymbol(symbol *parser.MathSymbol, x, y, fontSize float64) float64 {
	// Check if the symbol is valid before rendering
	if symbol.Symbol == "" {
		return 0
	}

	// Use normal font for math symbols to ensure compatibility
	m.generator.AddText(symbol.Symbol, x, y, fontSize, "normal")
	return m.generator.GetTextWidth(symbol.Symbol, fontSize, "normal")
}

// renderGroup renders a group of math elements (like braced content)
func (m *MathProcessor) renderGroup(group *parser.Group, x, y, fontSize float64) float64 {
	currentX := x

	for i, node := range group.Nodes {
		// Add spacing between elements if needed
		if i > 0 {
			currentX += m.getSpacing(group.Nodes[i-1], node)
		}

		width := m.renderMathElement(node, currentX, y, fontSize)
		currentX += width
	}

	return currentX - x // Return total width
}
