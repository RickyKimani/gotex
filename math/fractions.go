package math

import "github.com/rickykimani/gotex/parser"

// renderFraction renders fractions
func (m *MathProcessor) renderFraction(frac *parser.MathFraction, x, y, fontSize float64) float64 {
	fracFontSize := fontSize * 0.8 // Slightly smaller for fractions

	// Calculate widths
	numWidth := m.calculateElementWidth(frac.Numerator, fracFontSize)
	denWidth := m.calculateElementWidth(frac.Denominator, fracFontSize)
	maxWidth := numWidth
	if denWidth > maxWidth {
		maxWidth = denWidth
	}

	// Center numerator and denominator
	numX := x + (maxWidth-numWidth)/2
	denX := x + (maxWidth-denWidth)/2

	// Position numerator above and denominator below the baseline
	numY := y - fontSize*0.4
	denY := y + fontSize*0.6

	// Render numerator and denominator
	m.renderMathElement(frac.Numerator, numX, numY, fracFontSize)
	m.renderMathElement(frac.Denominator, denX, denY, fracFontSize)

	// Draw fraction line
	lineY := y + fontSize*0.1
	m.generator.AddLine(x, lineY, x+maxWidth, lineY)

	return maxWidth
}
