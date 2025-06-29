package math

import "github.com/rickykimani/gotex/parser"

//TODO: Implement different radicals

//TODO: Make sure the radical covers the width of operand (\radical{})

// renderSquareRoot renders square root symbols
func (mp *MathProcessor) renderSquareRoot(arg parser.Node, x, y, fontSize float64) float64 {
	// Render the square root symbol
	mp.generator.AddText("√", x, y, fontSize, "normal")
	symbolWidth := mp.generator.GetTextWidth("√", fontSize, "normal")

	// Render the argument
	argWidth := mp.renderMathElement(arg, x+symbolWidth, y, fontSize)

	return symbolWidth + argWidth
}
