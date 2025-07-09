package typesetter

import (
	"github.com/rickykimani/gotex/fonts"
	"github.com/rickykimani/gotex/parser"
)

// Typesetter handles the main typesetting process
type Typesetter struct {
	DefaultFont *fonts.Font
	LineWidth   float64
	LineBreak   *LineBreak
	CurrentX    float64
	CurrentY    float64
	Lines       [][]*Box
}

// NewTypesetter creates a new typesetter instance
func NewTypesetter(lineWidth float64) *Typesetter {
	defaultFont := fonts.NewFont("Computer Modern", "serif", "normal", "normal", 12.0)

	return &Typesetter{
		DefaultFont: defaultFont,
		LineWidth:   lineWidth,
		LineBreak:   NewLineBreak(lineWidth),
		CurrentX:    0.0,
		CurrentY:    0.0,
		Lines:       make([][]*Box, 0),
	}
}

// TypesetDocument typesets a complete document AST
func (ts *Typesetter) TypesetDocument(doc *parser.Document) [][]*Box {
	ts.Lines = make([][]*Box, 0)

	for _, node := range doc.Body {
		boxes := ts.TypesetNode(node)
		lines := ts.LineBreak.BreakParagraph(boxes)
		ts.Lines = append(ts.Lines, lines...)
	}

	return ts.Lines
}

// TypesetNode typesets a single AST node
func (ts *Typesetter) TypesetNode(node parser.Node) []*Box {
	switch n := node.(type) {
	case *parser.TextNode:
		return ts.TypesetText(n.Value)
	case *parser.Command:
		return ts.TypesetCommand(n)
	case *parser.Environment:
		return ts.TypesetEnvironment(n)
	case *parser.MathNode:
		return ts.TypesetMathNode(n)
	case *parser.Group:
		return ts.TypesetGroup(n)
	default:
		// Unknown node type - return empty
		return []*Box{}
	}
}

// TypesetText converts text into a sequence of boxes
func (ts *Typesetter) TypesetText(text string) []*Box {
	boxes := make([]*Box, 0)

	for _, char := range text {
		charBox := NewCharBox(char, ts.DefaultFont)
		boxes = append(boxes, &charBox.Box)
	}

	return boxes
}

// TypesetCommand typesets a LaTeX command
func (ts *Typesetter) TypesetCommand(cmd *parser.Command) []*Box {
	boxes := make([]*Box, 0)

	switch cmd.Name {
	case "textbf":
		// Bold text - use modified font
		boldFont := fonts.NewFont(
			ts.DefaultFont.Name,
			ts.DefaultFont.Family,
			ts.DefaultFont.Style,
			"bold",
			ts.DefaultFont.Size,
		)

		// Typeset arguments with bold font
		for _, arg := range cmd.Args {
			argBoxes := ts.typesetNodeWithFont(arg, boldFont)
			boxes = append(boxes, argBoxes...)
		}

	case "textit":
		// Italic text - use modified font
		italicFont := fonts.NewFont(
			ts.DefaultFont.Name,
			ts.DefaultFont.Family,
			"italic",
			ts.DefaultFont.Weight,
			ts.DefaultFont.Size,
		)

		// Typeset arguments with italic font
		for _, arg := range cmd.Args {
			argBoxes := ts.typesetNodeWithFont(arg, italicFont)
			boxes = append(boxes, argBoxes...)
		}

	default:
		// Unknown command - render as text
		commandText := "\\" + cmd.Name
		for _, char := range commandText {
			charBox := NewCharBox(char, ts.DefaultFont)
			boxes = append(boxes, &charBox.Box)
		}
	}

	return boxes
}

// TypesetEnvironment typesets a LaTeX environment
func (ts *Typesetter) TypesetEnvironment(env *parser.Environment) []*Box {
	boxes := make([]*Box, 0)

	// For now, just typeset the content
	for _, node := range env.Body {
		nodeBoxes := ts.TypesetNode(node)
		boxes = append(boxes, nodeBoxes...)
	}

	return boxes
}

// TypesetMath typesets mathematical content
func (ts *Typesetter) TypesetMath(content string) []*Box {
	// Simplified math typesetting - just use italic font
	mathFont := fonts.NewFont(
		"Computer Modern Math",
		"serif",
		"italic",
		"normal",
		ts.DefaultFont.Size,
	)

	boxes := make([]*Box, 0)
	for _, char := range content {
		charBox := NewCharBox(char, mathFont)
		boxes = append(boxes, &charBox.Box)
	}

	return boxes
}

// TypesetMathNode typesets a mathematical node with its content nodes
func (ts *Typesetter) TypesetMathNode(mathNode *parser.MathNode) []*Box {
	// Simplified math typesetting - just use italic font
	mathFont := fonts.NewFont(
		"Computer Modern Math",
		"serif",
		"italic",
		"normal",
		ts.DefaultFont.Size,
	)

	boxes := make([]*Box, 0)

	// Process each content node
	for _, node := range mathNode.Content {
		nodeBoxes := ts.typesetNodeWithFont(node, mathFont)
		boxes = append(boxes, nodeBoxes...)
	}

	return boxes
}

// TypesetGroup typesets a grouped set of nodes
func (ts *Typesetter) TypesetGroup(group *parser.Group) []*Box {
	boxes := make([]*Box, 0)

	for _, node := range group.Nodes {
		nodeBoxes := ts.TypesetNode(node)
		boxes = append(boxes, nodeBoxes...)
	}

	return boxes
}

// typesetNodeWithFont typesets a node using a specific font
func (ts *Typesetter) typesetNodeWithFont(node parser.Node, font *fonts.Font) []*Box {
	oldFont := ts.DefaultFont
	ts.DefaultFont = font

	boxes := ts.TypesetNode(node)

	ts.DefaultFont = oldFont
	return boxes
}

// GetTotalHeight calculates the total height of all typeset lines
func (ts *Typesetter) GetTotalHeight() float64 {
	totalHeight := 0.0

	for _, line := range ts.Lines {
		lineHeight := 0.0
		for _, box := range line {
			if box.Height+box.Depth > lineHeight {
				lineHeight = box.Height + box.Depth
			}
		}
		totalHeight += lineHeight
	}

	return totalHeight
}

// GetTotalWidth calculates the maximum width of all lines
func (ts *Typesetter) GetTotalWidth() float64 {
	maxWidth := 0.0

	for _, line := range ts.Lines {
		lineWidth := 0.0
		for _, box := range line {
			lineWidth += box.Width
		}
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return maxWidth
}
