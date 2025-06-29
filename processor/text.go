package processor

import (
	"strings"
	"time"

	"github.com/rickykimani/gotex/parser"
)

func (dp *DocumentProcessor) addText(text string, style string) {
	if text == "" {
		return
	}

	// Normalize newlines to spaces
	normalizedText := strings.ReplaceAll(text, "\n", " ")

	// Process text character by character to preserve exact spacing
	currentWord := ""

	for _, char := range normalizedText {
		if char == ' ' || char == '\t' {
			// Add the current word if we have one
			if currentWord != "" {
				dp.addWord(currentWord, style)
				currentWord = ""
			}

			// Add space if we have content on the line
			if dp.lineHasContent {
				dp.addSpace(style)
			}
		} else {
			currentWord += string(char)
		}
	}

	// Add the final word if we have one
	if currentWord != "" {
		dp.addWord(currentWord, style)
	}
}

func (dp *DocumentProcessor) extractText(node parser.Node) string {
	switch n := node.(type) {
	case *parser.TextNode:
		return n.Value
	case *parser.Group:
		var result strings.Builder
		for _, child := range n.Nodes {
			result.WriteString(dp.extractText(child))
		}
		return result.String()
	case *parser.Command:
		if n.Name == "textbf" && len(n.Args) > 0 {
			return dp.extractText(n.Args[0])
		}
		if n.Name == "textit" && len(n.Args) > 0 {
			return dp.extractText(n.Args[0])
		}
		if n.Name == "today" {
			return time.Now().Format("January 2, 2006") // Current date in LaTeX format
		}
		var result strings.Builder
		for _, arg := range n.Args {
			result.WriteString(dp.extractText(arg))
		}
		return result.String()
	default:
		return ""
	}
}
