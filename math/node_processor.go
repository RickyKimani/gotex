package math

import "github.com/rickykimani/gotex/parser"

// ProcessMathNode renders a math node to PDF
func (m *MathProcessor) ProcessMathNode(node *parser.MathNode, x, y float64) float64 {
	if node.Inline {
		return m.processInlineMath(node.Content, x, y)
	} else {
		return m.processDisplayMath(node.Content, x, y)
	}
}
