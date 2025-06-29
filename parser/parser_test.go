package parser

import (
	"fmt"
	"testing"

	"github.com/rickykimani/gotex/lexer"
)

func TestParser(t *testing.T) {
	input := `\documentclass[12pt]{article}
\usepackage{amsmath}

\begin{document}

\title{My Document}
\author{John Doe}

\section{Introduction}
Hello world! This is some text with $x + y = z$ inline math.

\begin{equation}
E = mc^2
\end{equation}

\textbf{Bold text} and \textit{italic text}.

\end{document}`

	lexer := lexer.NewLexer(input)
	parser := New(lexer)
	doc := parser.ParseDocument()

	fmt.Println("Parsing LaTeX document:")
	fmt.Println("Input:", input)
	fmt.Println("\nParsed AST:")
	printDocument(doc, 0)
}

func TestParserWithComments(t *testing.T) {
	input := `\documentclass{article}
% This is a comment at the top level
\begin{document}
% Comment inside environment
Hello world!
% Another comment
\textbf{Bold text} % Inline comment after command
% Final comment
\end{document}`

	lexer := lexer.NewLexer(input)
	parser := New(lexer)
	doc := parser.ParseDocument()

	fmt.Println("Parsing LaTeX document with comments:")
	fmt.Println("Input:", input)
	fmt.Println("\nParsed AST:")
	printDocument(doc, 0)

	// Count comment nodes
	commentCount := countComments(doc)
	fmt.Printf("\nTotal comment nodes found: %d\n", commentCount)
}

func countComments(doc *Document) int {
	count := 0
	for _, node := range doc.Body {
		count += countCommentsInNode(node)
	}
	return count
}

func countCommentsInNode(node Node) int {
	switch n := node.(type) {
	case *CommentNode:
		return 1
	case *Environment:
		count := 0
		for _, bodyNode := range n.Body {
			count += countCommentsInNode(bodyNode)
		}
		return count
	case *Command:
		count := 0
		for _, arg := range n.Args {
			count += countCommentsInNode(arg)
		}
		for _, opt := range n.Optional {
			count += countCommentsInNode(opt)
		}
		return count
	default:
		return 0
	}
}

func printDocument(doc *Document, indent int) {
	printIndent(indent)
	fmt.Printf("Document (Body: %d nodes)\n", len(doc.Body))

	for i, node := range doc.Body {
		printIndent(indent + 1)
		fmt.Printf("[%d] ", i)
		printNode(node, indent+1)
	}
}

func printNode(node Node, indent int) {
	switch n := node.(type) {
	case *Command:
		fmt.Printf("Command: \\%s\n", n.Name)
		if len(n.Optional) > 0 {
			printIndent(indent + 1)
			fmt.Printf("Optional args: %d\n", len(n.Optional))
			for _, opt := range n.Optional {
				printIndent(indent + 2)
				printNode(opt, indent+2)
			}
		}
		if len(n.Args) > 0 {
			printIndent(indent + 1)
			fmt.Printf("Required args: %d\n", len(n.Args))
			for _, arg := range n.Args {
				printIndent(indent + 2)
				printNode(arg, indent+2)
			}
		}
	case *Environment:
		fmt.Printf("Environment: %s (Body: %d nodes)\n", n.Name, len(n.Body))
		for _, bodyNode := range n.Body {
			printIndent(indent + 1)
			printNode(bodyNode, indent+1)
		}
	case *TextNode:
		fmt.Printf("Text: %q\n", n.Value)
	case *CommentNode:
		fmt.Printf("Comment: %q\n", n.Value)
	case *MathNode:
		if n.Inline {
			fmt.Printf("Math (inline): %d content nodes\n", len(n.Content))
		} else {
			fmt.Printf("Math (display): %d content nodes\n", len(n.Content))
		}
		for _, contentNode := range n.Content {
			printIndent(indent + 1)
			printNode(contentNode, indent+1)
		}
	default:
		fmt.Printf("Unknown node type\n")
	}
}

func printIndent(level int) {
	for i := 0; i < level; i++ {
		fmt.Print("  ")
	}
}
