package processor

import "github.com/rickykimani/gotex/parser"

// processStyledText processes a node with styling, ensuring text is split into words and spaces
func (dp *DocumentProcessor) processStyledText(node parser.Node, style string) {
	switch n := node.(type) {
	case *parser.TextNode:
		// Process the text with the given style, ensuring proper word separation
		dp.addText(n.Value, style)

	case *parser.Group:
		// Process each node in the group with the same style
		for i, child := range n.Nodes {
			// Add space between group nodes when appropriate
			if i > 0 {
				prevChild := n.Nodes[i-1]
				if dp.shouldAddSpaceBetweenNodes(prevChild, child) && dp.lineHasContent {
					dp.addSpace(style)
				}
			}
			dp.processStyledText(child, style)
		}

	case *parser.Command:
		// Handle nested commands with appropriate style combination
		switch n.Name {
		case "textbf":
			newStyle := "bold"
			switch style {
			case "italic":
				newStyle = "bold-italic"
			case "bold":
				newStyle = "bold" // Already bold
			}
			if len(n.Args) > 0 {
				dp.processStyledText(n.Args[0], newStyle)
			}
		case "textit":
			newStyle := "italic"
			switch style {
			case "bold":
				newStyle = "bold-italic"
			case "italic":
				newStyle = "italic" // Already italic
			}
			if len(n.Args) > 0 {
				dp.processStyledText(n.Args[0], newStyle)
			}
		default:
			// For other commands, process normally but preserve style
			dp.processCommand(n, style)
		}

	default:
		// For any other node type, process with the current style
		dp.processNode(node, style)
	}
}
