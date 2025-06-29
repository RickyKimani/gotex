package math

import (
	"github.com/rickykimani/gotex/pdf"
)

// MathProcessor handles rendering of mathematical expressions
type MathProcessor struct {
	generator *pdf.Generator
	fontSize  float64
	spacing   MathSpacing
}

// NewMathProcessor creates a new math processor
func NewMathProcessor(generator *pdf.Generator, fontSize float64) *MathProcessor {
	return &MathProcessor{
		generator: generator,
		fontSize:  fontSize,
		spacing:   StandardMathSpacing(fontSize),
	}
}
