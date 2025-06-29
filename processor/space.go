package processor

import (
	"strings"

	"github.com/rickykimani/gotex/parser"
)

func (dp *DocumentProcessor) addSpace(style string) {
	spaceWidth := dp.calculateTextWidth(" ", style)

	// Check if adding space would exceed the line width
	if dp.currentLineX+spaceWidth > dp.generator.PageWidth-dp.generator.MarginRight {
		dp.newLine()
		return // Don't add space at the beginning of a new line
	}

	// Add space by moving position forward
	dp.currentLineX += spaceWidth
}

func (dp *DocumentProcessor) addVerticalSpace(space float64) {
	dp.currentY -= space
	dp.checkNewPage()
}

// shouldAddSpaceBetweenNodes determines if we need to add space between two adjacent nodes
func (dp *DocumentProcessor) shouldAddSpaceBetweenNodes(prev, curr parser.Node) bool {
	// Check if the previous node ended with a space or if current starts with space
	prevEndsWithSpace := false
	currStartsWithSpace := false

	if prevText, ok := prev.(*parser.TextNode); ok {
		prevEndsWithSpace = strings.HasSuffix(prevText.Value, " ") || strings.HasSuffix(prevText.Value, "\n") || strings.HasSuffix(prevText.Value, "\t")
	}

	if currText, ok := curr.(*parser.TextNode); ok {
		currStartsWithSpace = strings.HasPrefix(currText.Value, " ") || strings.HasPrefix(currText.Value, "\n") || strings.HasPrefix(currText.Value, "\t")
	}

	// If neither has explicit spacing, we typically need to add space between nodes
	// Exception: don't add space if nodes are commands or other non-text elements
	if !prevEndsWithSpace && !currStartsWithSpace {
		// Add space between text nodes or between text and commands
		_, prevIsText := prev.(*parser.TextNode)
		_, currIsText := curr.(*parser.TextNode)
		_, prevIsCommand := prev.(*parser.Command)
		_, currIsCommand := curr.(*parser.Command)

		// Add space if we have text-to-text, text-to-command, or command-to-text transitions
		return (prevIsText && currIsText) || (prevIsText && currIsCommand) || (prevIsCommand && currIsText)
	}

	return false
}
