package parser

import (
	"fmt"
	"testing"

	"github.com/rickykimani/gotex/lexer"
)

func TestBasicErrorHandling(t *testing.T) {
	// Test case 1: Missing closing brace
	fmt.Println("=== Test: Missing closing brace ===")
	input1 := `\textbf{Bold text without closing brace`

	lexer1 := lexer.NewLexer(input1)
	parser1 := New(lexer1)
	doc1 := parser1.ParseDocument()

	fmt.Printf("Input: %s\n", input1)
	fmt.Printf("Errors found: %d\n", len(parser1.GetErrors()))
	for _, err := range parser1.GetErrors() {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("Parsed AST:")
	printDocument(doc1, 0)

	// Test case 2: Unmatched environment
	fmt.Println("\n=== Test: Unmatched environment ===")
	input2 := `\begin{document}Hello world!` // Missing \end{document}

	lexer2 := lexer.NewLexer(input2)
	parser2 := New(lexer2)
	doc2 := parser2.ParseDocument()

	fmt.Printf("Input: %s\n", input2)
	fmt.Printf("Errors found: %d\n", len(parser2.GetErrors()))
	for _, err := range parser2.GetErrors() {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("Parsed AST:")
	printDocument(doc2, 0)

	// Test case 3: LaTeX-like recovery
	fmt.Println("\n=== Test: LaTeX-like error recovery ===")
	input3 := `\documentclass{article}
\begin{document}
\textbf{Bold text  % Missing closing brace
\textit{Italic text}  % This should still parse
\section{A Section}
Normal text continues here.
\end{document}`

	lexer3 := lexer.NewLexer(input3)
	parser3 := New(lexer3)
	doc3 := parser3.ParseDocument()

	fmt.Printf("Input:\n%s\n", input3)
	fmt.Printf("Errors found: %d\n", len(parser3.GetErrors()))
	for _, err := range parser3.GetErrors() {
		fmt.Printf("  - %s\n", err.Error())
	}

	fmt.Println("Parsed AST (showing recovery):")
	printDocument(doc3, 0)

	// Show that parsing continued despite errors
	nodeCount := countTotalNodes(doc3)
	fmt.Printf("\nSuccessfully parsed %d total nodes despite errors\n", nodeCount)
	fmt.Println("This demonstrates LaTeX-like error recovery!")
}

func countTotalNodes(doc *Document) int {
	count := len(doc.Body)
	for _, node := range doc.Body {
		count += countTotalNodesInNode(node)
	}
	return count
}

func countTotalNodesInNode(node Node) int {
	switch n := node.(type) {
	case *Environment:
		count := 0
		for _, bodyNode := range n.Body {
			count += 1 + countTotalNodesInNode(bodyNode)
		}
		return count
	case *Command:
		count := 0
		for _, arg := range n.Args {
			count += 1 + countTotalNodesInNode(arg)
		}
		for _, opt := range n.Optional {
			count += 1 + countTotalNodesInNode(opt)
		}
		return count
	default:
		return 0
	}
}
