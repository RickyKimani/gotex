package processor

import (
	"strings"

	"github.com/rickykimani/gotex/parser"
)

func (dp *DocumentProcessor) processNode(node parser.Node, style string) {
	switch n := node.(type) {
	case *parser.TextNode:
		// Handle single newline: line break
		if n.Value == "\n" {
			dp.newLine()
			dp.lastProcessedCommand = false
			return
		}
		// Check for paragraph break (double newline)
		if strings.Contains(n.Value, "\n\n") {
			paragraphs := strings.Split(n.Value, "\n\n")
			for i, para := range paragraphs {
				if i > 0 {
					dp.addParagraphBreak()
				}
				if strings.TrimSpace(para) != "" {
					dp.addText(para, style)
				}
			}
		} else {
			// Normal text or single newline embedded - addText handles internal spaces
			if strings.TrimSpace(n.Value) != "" {
				dp.addText(n.Value, style)
			} else if dp.lineHasContent && strings.Contains(n.Value, " ") {
				dp.addSpace(style)
			}
		}
		dp.lastProcessedCommand = false

	case *parser.Command:
		dp.processCommand(n, style)
		dp.lastProcessedCommand = true

	case *parser.Environment:
		dp.processEnvironment(n, style)
		dp.lastProcessedCommand = false

	case *parser.MathNode:
		dp.processMathNode(n, style)
		dp.lastProcessedCommand = false

	case *parser.Group:
		// Check if this is a styled text group (font command + text)
		if len(n.Nodes) >= 2 {
			if fontCmd, ok := n.Nodes[0].(*parser.Command); ok && fontCmd.Name == "font" {
				// This is a styled text group from macro expansion
				newStyle := style
				if len(fontCmd.Args) > 0 {
					fontArg := dp.extractText(fontCmd.Args[0])
					switch fontArg {
					case "bold":
						newStyle = "bold"
						if style == "italic" {
							newStyle = "bold-italic"
						}
					case "italic":
						newStyle = "italic"
						if style == "bold" {
							newStyle = "bold-italic"
						}
					}
				}
				// Process the remaining nodes with the new style
				for i := 1; i < len(n.Nodes); i++ {
					if i > 1 {
						// Add space between nodes if needed
						prevNode := n.Nodes[i-1]
						currNode := n.Nodes[i]
						if dp.shouldAddSpaceBetweenNodes(prevNode, currNode) && dp.lineHasContent {
							dp.addSpace(newStyle)
						}
					}
					dp.processNode(n.Nodes[i], newStyle)
				}
			} else {
				// Regular group processing
				dp.processNodes(n.Nodes, style)
			}
		} else {
			// Regular group processing
			dp.processNodes(n.Nodes, style)
		}
		dp.lastProcessedCommand = false
	}
}

func (dp *DocumentProcessor) processNodes(nodes []parser.Node, style string) {
	for i, node := range nodes {
		// Add space between nodes when appropriate
		if i > 0 {
			prevNode := nodes[i-1]
			needsSpace := dp.shouldAddSpaceBetweenNodes(prevNode, node)
			if needsSpace && dp.lineHasContent {
				dp.addSpace(style)
			}
		}
		dp.processNode(node, style)
	}
}
