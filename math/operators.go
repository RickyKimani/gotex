package math

import "slices"

// GetMathFont returns the appropriate font style for math content
func GetMathFont(content string) string {
	// Use normal for now for testing
	return "normal"
}

// MathSpacing defines spacing rules for mathematical typography
type MathSpacing struct {
	BeforeOperator float64 // space before +, -, *, etc.
	AfterOperator  float64 // space after +, -, *, etc.
	BeforeRelation float64 // space before =, <, >, etc.
	AfterRelation  float64 // space after =, <, >, etc.
	BeforeFunction float64 // space before sin, cos, log, etc.
	AfterFunction  float64 // space after sin, cos, log, etc.
}

// StandardMathSpacing returns standard mathematical spacing rules
func StandardMathSpacing(fontSize float64) MathSpacing {
	baseUnit := fontSize * 0.1 // 10% of font size as base spacing unit
	return MathSpacing{
		BeforeOperator: baseUnit * 2, // medium space
		AfterOperator:  baseUnit * 2,
		BeforeRelation: baseUnit * 3, // thick space
		AfterRelation:  baseUnit * 3,
		BeforeFunction: baseUnit * 1, // thin space
		AfterFunction:  0,            // no space after function names
	}
}

// IsOperator checks if a string is a mathematical operator
func IsOperator(s string) bool {
	operators := []string{"+", "-", "*", "×", "÷", "·", "±", "∓"}
	return slices.Contains(operators, s)
}

// IsRelation checks if a string is a mathematical relation
func IsRelation(s string) bool {
	relations := []string{"=", "<", ">", "≤", "≥", "≠", "≡", "≈", "∼", "≃", "≅", "∝"}
	return slices.Contains(relations, s)
}
