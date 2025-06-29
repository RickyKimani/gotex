package parser

import (
	"fmt"
	"testing"

	"github.com/rickykimani/gotex/lexer"
)

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int // Expected number of errors
	}{
		{
			name:     "Missing closing brace",
			input:    `\textbf{Bold text`,
			expected: 1,
		},
		{
			name:     "Unmatched environment",
			input:    `\begin{document} Some text`,
			expected: 1,
		},
		{
			name:     "Wrong environment end",
			input:    `\begin{document} Some text \end{article}`,
			expected: 2, // Wrong end + Missing proper end
		},
		{
			name:     "Multiple errors",
			input:    `\textbf{Bold \begin{center} text \end{wrong}`,
			expected: 3, // Wrong environment end + Missing environment end + Missing brace
		},
		{
			name:     "Valid document",
			input:    `\documentclass{article}\begin{document}Hello\end{document}`,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n=== Testing: %s ===\n", tc.name)
			fmt.Printf("Input: %s\n", tc.input)

			lexer := lexer.NewLexer(tc.input)
			parser := New(lexer)
			doc := parser.ParseDocument()

			errors := parser.GetErrors()
			fmt.Printf("Found %d errors:\n", len(errors))

			for i, err := range errors {
				fmt.Printf("  [%d] %s\n", i+1, err.Error())
			}

			if len(errors) != tc.expected {
				t.Errorf("Expected %d errors, got %d", tc.expected, len(errors))
			}

			// Show parsed AST even with errors
			fmt.Println("Parsed AST (with recovery):")
			printDocument(doc, 0)
		})
	}
}

func TestLaTeXLikeRecovery(t *testing.T) {
	// This simulates how LaTeX handles problematic input - it tries to continue
	input := `\documentclass{article}
\begin{document}

\section{Introduction}
This is a paragraph with \textbf{bold text.

% Missing closing brace above - LaTeX would warn but continue

\subsection{Math Examples}
Here's some math: $x + y = z

% Missing closing $ above

And here's display math:
$$E = mc^2$$

\begin{itemize}
\item First item
\item Second item
% Missing \end{itemize} - LaTeX would warn but try to close at document end

\end{document}`

	fmt.Println("=== LaTeX-like Error Recovery Test ===")
	fmt.Printf("Input:\n%s\n\n", input)

	lexer := lexer.NewLexer(input)
	parser := New(lexer)
	doc := parser.ParseDocument()

	fmt.Printf("Parser found %d errors:\n", len(parser.GetErrors()))
	for i, err := range parser.GetErrors() {
		fmt.Printf("  [%d] %s\n", i+1, err.Error())
	}

	fmt.Println("\nParsed AST (showing recovery):")
	printDocument(doc, 0)

	// Count successfully parsed nodes
	nodeCount := countNodes(doc)
	fmt.Printf("\nSuccessfully parsed %d nodes despite errors\n", nodeCount)
}

func countNodes(doc *Document) int {
	count := len(doc.Body)
	for _, node := range doc.Body {
		count += countNodesInNode(node)
	}
	return count
}

func countNodesInNode(node Node) int {
	switch n := node.(type) {
	case *Environment:
		count := 0
		for _, bodyNode := range n.Body {
			count += 1 + countNodesInNode(bodyNode)
		}
		return count
	case *Command:
		count := 0
		for _, arg := range n.Args {
			count += 1 + countNodesInNode(arg)
		}
		for _, opt := range n.Optional {
			count += 1 + countNodesInNode(opt)
		}
		return count
	default:
		return 0
	}
}
