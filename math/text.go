package math

import (
	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/symbols"
)

// extractTextContent extracts text content from a node for spacing analysis
func (m *MathProcessor) extractTextContent(node parser.Node) string {
	switch n := node.(type) {
	case *parser.TextNode:
		return n.Value
	case *parser.Command:
		if symbol, exists := symbols.ConvertMathSymbol(n.Name); exists {
			return symbol
		}
		return n.Name
	default:
		return ""
	}
}
