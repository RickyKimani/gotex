package math

import "github.com/rickykimani/gotex/parser"

// renderSubscript renders subscript (indices)
func (mp *MathProcessor) renderSubscript(sub *parser.MathSubscript, x, y, fontSize float64) float64 {
	// Render base
	var baseWidth float64
	if sub.Base != nil {
		baseWidth = mp.renderMathElement(sub.Base, x, y, fontSize)
	}

	// Render index (smaller and lowered)
	idxFontSize := fontSize * 0.7 // 70% of base size
	idxY := y - fontSize*0.3      // LOWER by 30%
	mp.renderMathElement(sub.Index, x+baseWidth, idxY, idxFontSize)

	// Calculate total width
	idxWidth := mp.calculateElementWidth(sub.Index, idxFontSize)
	return baseWidth + idxWidth
}
