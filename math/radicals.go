package math

import "github.com/rickykimani/gotex/parser"

//TODO: Implement different radicals

// renderSquareRoot renders square root symbols with proper vinculum
func (mp *MathProcessor) renderSquareRoot(arg parser.Node, x, y, fontSize float64) float64 {
	// Render the square root symbol
	mp.generator.AddText("√", x, y, fontSize, "normal")
	symbolWidth := mp.generator.GetTextWidth("√", fontSize, "normal")

	// Calculate the width of the argument to know how long the vinculum should be
	argWidth := mp.calculateElementWidth(arg, fontSize)

	// Render the argument
	mp.renderMathElement(arg, x+symbolWidth, y, fontSize)

	// Draw the vinculum (horizontal line) above the radicand
	// Position it at the top of the radical symbol (above the baseline)
	//TODO: Implement precise font height and multipliers to avoid hard coded values
	vinculumY := y + fontSize*0.8
	vinculumStartX := x + symbolWidth
	vinculumEndX := x + symbolWidth + argWidth

	mp.generator.AddLine(vinculumStartX, vinculumY, vinculumEndX, vinculumY)

	return symbolWidth + argWidth
}
