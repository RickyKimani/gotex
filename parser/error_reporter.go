package parser

import (
	"fmt"
	"strings"
)

// ErrorReporter displays parse errors in Rust-style format
type ErrorReporter struct {
	sourceCode string
	filename   string
}

// NewErrorReporter creates a new error reporter
func NewErrorReporter(sourceCode, filename string) *ErrorReporter {
	return &ErrorReporter{
		sourceCode: sourceCode,
		filename:   filename,
	}
}

// ReportErrors displays all errors in Rust-style format
func (r *ErrorReporter) ReportErrors(errors []ParseError) {
	if len(errors) == 0 {
		return
	}

	fmt.Printf("\nParse errors found (%d):\n\n", len(errors))
	lines := strings.Split(r.sourceCode, "\n")

	for _, err := range errors {
		r.printError(err, lines)
	}
}

// printError displays a single error with code context and pointer
func (r *ErrorReporter) printError(err ParseError, lines []string) {
	severityStr := getSeverityString(err.Severity)

	fmt.Printf("%s: %s\n", severityStr, err.Message)
	fmt.Printf(" --> %s:%d:%d\n", r.filename, err.Position.Line, err.Position.Column)

	// Ensure we have a valid line number
	if err.Position.Line > 0 && err.Position.Line <= len(lines) {
		line := lines[err.Position.Line-1]
		lineNumStr := fmt.Sprintf("%d", err.Position.Line)

		// Print the problematic line
		fmt.Printf("%s | %s\n", lineNumStr, line)

		// Print pointer to the column
		padding := max(err.Position.Column - 1, 0)

		fmt.Printf("%s | %s^ here\n",
			strings.Repeat(" ", len(lineNumStr)),
			strings.Repeat(" ", padding))
	}
	fmt.Println()
}

// getSeverityString converts severity enum to display string
func getSeverityString(s Severity) string {
	switch s {
	case Warning:
		return "warning"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	default:
		return "error"
	}
}
