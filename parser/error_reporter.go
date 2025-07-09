package parser

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ErrorReporter displays parse errors in Rust-style format with colors
type ErrorReporter struct {
	sourceCode string
	filename   string
	// Color functions for different severity levels
	warningColor  *color.Color
	errorColor    *color.Color
	fatalColor    *color.Color
	filenameColor *color.Color
	lineNumColor  *color.Color
	pointerColor  *color.Color
}

// NewErrorReporter creates a new error reporter with color support
func NewErrorReporter(sourceCode, filename string) *ErrorReporter {
	return &ErrorReporter{
		sourceCode:    sourceCode,
		filename:      filename,
		warningColor:  color.New(color.FgYellow, color.Bold),
		errorColor:    color.New(color.FgRed, color.Bold),
		fatalColor:    color.New(color.FgRed, color.Bold, color.BgWhite),
		filenameColor: color.New(color.FgCyan),
		lineNumColor:  color.New(color.FgBlue, color.Bold),
		pointerColor:  color.New(color.FgRed, color.Bold),
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

// printError displays a single error with code context and colorized pointer
func (r *ErrorReporter) printError(err ParseError, lines []string) {
	severityStr, colorFunc := r.getSeverityStringAndColor(err.Severity)

	// Print colored error message
	colorFunc.Printf("%s", severityStr)
	fmt.Printf(": %s\n", err.Message)

	// Print filename with color
	fmt.Print(" --> ")
	r.filenameColor.Printf("%s", r.filename)
	fmt.Printf(":%d:%d\n", err.Position.Line, err.Position.Column)

	// Ensure we have a valid line number
	if err.Position.Line > 0 && err.Position.Line <= len(lines) {
		line := lines[err.Position.Line-1]
		lineNumStr := fmt.Sprintf("%d", err.Position.Line)

		// Print the problematic line with colored line number
		r.lineNumColor.Printf("%s", lineNumStr)
		fmt.Printf(" | %s\n", line)

		// Print pointer to the column with color
		padding := max(err.Position.Column-1, 0)
		fmt.Printf("%s | %s",
			strings.Repeat(" ", len(lineNumStr)),
			strings.Repeat(" ", padding))
		r.pointerColor.Printf("^ here\n")
	}
	fmt.Println()
}

// getSeverityStringAndColor returns severity string and corresponding color function
func (r *ErrorReporter) getSeverityStringAndColor(s Severity) (string, *color.Color) {
	switch s {
	case Warning:
		return "warning", r.warningColor
	case Error:
		return "error", r.errorColor
	case Fatal:
		return "fatal", r.fatalColor
	default:
		return "error", r.errorColor
	}
}
