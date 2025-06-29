package math

import "github.com/rickykimani/gotex/parser"

// getSpacing returns appropriate spacing between two math elements
func (m *MathProcessor) getSpacing(prev, curr parser.Node) float64 {
	// Extract text content for analysis
	prevText := m.extractTextContent(prev)
	currText := m.extractTextContent(curr)

	// Apply spacing rules
	if IsOperator(prevText) || IsOperator(currText) {
		return m.spacing.BeforeOperator
	}

	if IsRelation(prevText) || IsRelation(currText) {
		return m.spacing.BeforeRelation
	}

	// Default spacing for other elements
	return m.fontSize * 0.1
}
