package math

import "github.com/rickykimani/gotex/parser"

// processDisplayMath renders display math expressions (centered)
func (mp *MathProcessor) processDisplayMath(content []parser.Node, _, y float64) float64 {
	// Treat as centered inline math
	totalWidth := mp.CalculateMathWidth(content)
	contentWidth := mp.generator.GetContentWidth()
	centerX := mp.generator.MarginLeft + (contentWidth-totalWidth)/2

	return mp.processInlineMath(content, centerX, y)
}
