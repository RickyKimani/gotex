package math

import (
	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/symbols"
)

// calculateMathWidth calculates the total width of a math expression
func (m *MathProcessor) CalculateMathWidth(content []parser.Node) float64 {
	totalWidth := 0.0

	for i, node := range content {
		if i > 0 {
			totalWidth += m.getSpacing(content[i-1], node)
		}
		totalWidth += m.calculateElementWidth(node, m.fontSize)
	}

	return totalWidth
}

// calculateElementWidth calculates the width of a single math element
func (m *MathProcessor) calculateElementWidth(node parser.Node, fontSize float64) float64 {
	// Safety check for nil nodes
	if node == nil {
		return 0
	}

	switch n := node.(type) {
	case *parser.TextNode:
		fontStyle := GetMathFont(n.Value)
		return m.generator.GetTextWidth(n.Value, fontSize, fontStyle)

	case *parser.Command:
		if symbol, exists := symbols.ConvertMathSymbol(n.Name); exists {
			return m.generator.GetTextWidth(symbol, fontSize, "normal")
		}
		// For other commands, approximate based on arguments
		totalWidth := 0.0
		for _, arg := range n.Args {
			totalWidth += m.calculateElementWidth(arg, fontSize)
		}
		return totalWidth

	case *parser.Group:
		totalWidth := 0.0
		for _, node := range n.Content {
			totalWidth += m.calculateElementWidth(node, fontSize)
		}
		return totalWidth

	default:
		// Default approximation
		return fontSize * 0.5
	}
}
