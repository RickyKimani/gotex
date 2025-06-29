package math

import "github.com/rickykimani/gotex/parser"

// processInlineMath renders inline math expressions
func (m *MathProcessor) processInlineMath(content []parser.Node, x, y float64) float64 {
	currentX := x

	for i, node := range content {
		// Add spacing before operators and relations
		if i > 0 {
			currentX += m.getSpacing(content[i-1], node)
		}

		width := m.renderMathElement(node, currentX, y, m.fontSize)
		currentX += width
	}

	return currentX - x // Return total width
}
