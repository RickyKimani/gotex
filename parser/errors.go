package parser

import (
	"fmt"

	"github.com/rickykimani/gotex/lexer"
)

// ErrorType represents different kinds of parsing errors
type ErrorType int

const (
	UnmatchedBrace ErrorType = iota
	UnmatchedEnvironment
	MissingArgument
	UnexpectedToken
	UnmatchedMath
	UnexpectedEOF
)

// ParseError represents a parsing error with recovery information
type ParseError struct {
	Type     ErrorType
	Message  string
	Position lexer.Position
	Expected string
	Got      string
	Severity Severity
}

type Severity int

const (
	Warning Severity = iota
	Error
	Fatal
)

func (e ParseError) Error() string {
	severityStr := map[Severity]string{
		Warning: "Warning",
		Error:   "Error",
		Fatal:   "Fatal",
	}
	return fmt.Sprintf("%s at line %d, col %d: %s",
		severityStr[e.Severity], e.Position.Line, e.Position.Column, e.Message)
}
