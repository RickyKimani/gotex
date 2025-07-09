package math

import (
	"math"
	"slices"

	"github.com/rickykimani/gotex/parser"
)

// renderFraction renders fractions using improved positioning based on go-latex approach
func (m *MathProcessor) renderFraction(frac *parser.MathFraction, x, y, fontSize float64) float64 {
	// Use smaller font for numerator and denominator (similar to go-latex's shrink factor)
	fracFontSize := fontSize * 0.7

	// Calculate widths for proper centering
	numWidth := m.calculateElementWidth(frac.Numerator, fracFontSize)
	denWidth := m.calculateElementWidth(frac.Denominator, fracFontSize)
	maxWidth := math.Max(numWidth, denWidth)

	// Calculate heights to determine proper spacing for nested fractions
	numHeight := m.calculateElementHeight(frac.Numerator, fracFontSize)
	denHeight := m.calculateElementHeight(frac.Denominator, fracFontSize)

	// Add padding similar to go-latex (2 * thickness)
	thickness := fontSize * 0.05 // Approximate line thickness
	padding := 2 * thickness
	totalWidth := maxWidth + 2*padding

	// Center positions horizontally
	numX := x + padding + (maxWidth-numWidth)/2
	denX := x + padding + (maxWidth-denWidth)/2

	// Position fraction line at baseline (like equals sign middle)
	lineY := y

	// Calculate vertical spacing based on actual element heights
	// Use larger gaps for nested fractions to prevent overlap
	baseGap := fontSize * 0.2 // Increased base gap

	// Check if numerator contains a fraction for extra spacing
	numGap := baseGap + (numHeight * 0.2) // Increased multiplier for better spacing
	denGap := baseGap + (denHeight * 0.1) // Extra space for tall denominators

	// Add extra spacing if numerator is a fraction
	if m.isFraction(frac.Numerator) {
		numGap += fontSize * 0.15 // More additional gap for fraction numerators
	}

	// Position numerator above the line (higher Y value in this coordinate system)
	numY := lineY + numGap + (fracFontSize * 0.2)

	// Position denominator below the line (lower Y value in this coordinate system)
	denY := lineY - denGap - (fracFontSize * 0.8)

	// Render components
	m.renderMathElement(frac.Numerator, numX, numY, fracFontSize)   // Numerator
	m.renderMathElement(frac.Denominator, denX, denY, fracFontSize) // Denominator

	// Draw fraction line with proper thickness - ensure it covers the full width
	m.generator.AddLine(x+padding, lineY, x+padding+maxWidth, lineY)

	return totalWidth
}

// calculateElementHeight calculates the approximate height of a math element
func (m *MathProcessor) calculateElementHeight(node parser.Node, fontSize float64) float64 {
	if node == nil {
		return fontSize
	}

	switch n := node.(type) {
	case *parser.TextNode:
		return fontSize

	case *parser.MathFraction:
		// For nested fractions, calculate total height with proper spacing
		numHeight := m.calculateElementHeight(n.Numerator, fontSize*0.7)
		denHeight := m.calculateElementHeight(n.Denominator, fontSize*0.7)
		// Add extra spacing for nested fractions to prevent overlap
		return numHeight + denHeight + fontSize*0.6

	case *parser.MathSuperscript:
		baseHeight := m.calculateElementHeight(n.Base, fontSize)
		expHeight := m.calculateElementHeight(n.Exponent, fontSize*0.7)
		return baseHeight + expHeight*0.5 // Superscript adds to height

	case *parser.MathSubscript:
		baseHeight := m.calculateElementHeight(n.Base, fontSize)
		return baseHeight + fontSize*0.3 // Subscript adds some height

	case *parser.Group:
		maxHeight := fontSize
		for _, child := range n.Nodes {
			h := m.calculateElementHeight(child, fontSize)
			if h > maxHeight {
				maxHeight = h
			}
		}
		return maxHeight

	case *parser.Command:
		// For commands with arguments, check their heights
		maxHeight := fontSize
		for _, arg := range n.Args {
			h := m.calculateElementHeight(arg, fontSize)
			if h > maxHeight {
				maxHeight = h
			}
		}
		return maxHeight

	default:
		return fontSize
	}
}

// isFraction checks if a node contains a fraction
func (m *MathProcessor) isFraction(node parser.Node) bool {
	if node == nil {
		return false
	}

	switch node := node.(type) {
	case *parser.MathFraction:
		return true
	case *parser.Group:
		// Check if the group contains a fraction
		group := node
		return slices.ContainsFunc(group.Nodes, m.isFraction)
	default:
		return false
	}
}
