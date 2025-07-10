package math

import (
	"math"

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
		for i, node := range n.Nodes {
			// Add spacing between elements if needed
			if i > 0 {
				totalWidth += m.getSpacing(n.Nodes[i-1], node)
			}
			totalWidth += m.calculateElementWidth(node, fontSize)
		}
		return totalWidth

	case *parser.MathFraction:
		// For fractions, return the width of the wider element
		numWidth := m.calculateElementWidth(n.Numerator, fontSize*0.7)
		denWidth := m.calculateElementWidth(n.Denominator, fontSize*0.7)
		return math.Max(numWidth, denWidth)

	case *parser.MathSuperscript:
		baseWidth := m.calculateElementWidth(n.Base, fontSize)
		expWidth := m.calculateElementWidth(n.Exponent, fontSize*0.7)
		return baseWidth + expWidth*0.8 // Superscript adds to width

	case *parser.MathSubscript:
		baseWidth := m.calculateElementWidth(n.Base, fontSize)
		subWidth := m.calculateElementWidth(n.Index, fontSize*0.7)
		return baseWidth + subWidth*0.8 // Subscript adds to width

	case *parser.MathSymbol:
		return m.generator.GetTextWidth(n.Symbol, fontSize, "normal")

	default:
		// Default approximation
		return fontSize * 0.5
	}
}
