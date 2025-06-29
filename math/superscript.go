package math

import "github.com/rickykimani/gotex/parser"

// renderSuperscript renders superscript (exponents)
func (mp *MathProcessor) renderSuperscript(sup *parser.MathSuperscript, x, y, fontSize float64) float64 {
	// Render base
	var baseWidth float64
	if sup.Base != nil {
		baseWidth = mp.renderMathElement(sup.Base, x, y, fontSize)
	}

	// Render exponent (smaller and raised)
	expFontSize := fontSize * 0.7 // 70% of base size
	expY := y + fontSize*0.3      // RAISE by 30%
	mp.renderMathElement(sup.Exponent, x+baseWidth, expY, expFontSize)

	// Calculate total width
	expWidth := mp.calculateElementWidth(sup.Exponent, expFontSize)
	return baseWidth + expWidth
}
