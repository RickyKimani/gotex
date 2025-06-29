package processor

import "github.com/rickykimani/gotex/parser"

func (dp *DocumentProcessor) processMathNode(mathNode *parser.MathNode, _ string) {
	// Process math content using the math processor
	mathWidth := dp.mathProcessor.ProcessMathNode(mathNode, dp.currentLineX, dp.currentY)

	// Update position after math
	dp.currentLineX += mathWidth
	dp.lineHasContent = true

	// Add a small space after inline math when followed by text
	if mathNode.Inline {
		dp.currentLineX += 3.0 // Small space after math
	}
}
