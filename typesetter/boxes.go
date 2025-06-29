package typesetter

import (
	"github.com/rickykimani/gotex/fonts"
	"github.com/rickykimani/gotex/parser"
)

// Box represents a typeset element with dimensions and position
type Box struct {
	Width   float64
	Height  float64
	Depth   float64
	X       float64
	Y       float64
	Content any // Could be text, image, etc.
	Font    *fonts.Font
}

// CharBox represents a single character box
type CharBox struct {
	Box
	Char rune
}

// TextBox represents a sequence of characters
type TextBox struct {
	Box
	Text string
}

// HBox represents a horizontal box containing other boxes
type HBox struct {
	Box
	Children []*Box
}

// VBox represents a vertical box containing other boxes
type VBox struct {
	Box
	Children []*Box
}

// NewCharBox creates a new character box
func NewCharBox(char rune, font *fonts.Font) *CharBox {
	width := font.GetCharWidth(char)
	return &CharBox{
		Box: Box{
			Width:  width,
			Height: font.Metrics.Ascent,
			Depth:  font.Metrics.Descent,
			Font:   font,
		},
		Char: char,
	}
}

// NewTextBox creates a new text box
func NewTextBox(text string, font *fonts.Font) *TextBox {
	width := font.GetStringWidth(text)
	return &TextBox{
		Box: Box{
			Width:  width,
			Height: font.Metrics.Ascent,
			Depth:  font.Metrics.Descent,
			Font:   font,
		},
		Text: text,
	}
}

// NewHBox creates a new horizontal box
func NewHBox() *HBox {
	return &HBox{
		Box:      Box{},
		Children: make([]*Box, 0),
	}
}

// NewVBox creates a new vertical box
func NewVBox() *VBox {
	return &VBox{
		Box:      Box{},
		Children: make([]*Box, 0),
	}
}

// AddChild adds a child box to an HBox and updates dimensions
func (hb *HBox) AddChild(child *Box) {
	hb.Children = append(hb.Children, child)
	hb.Width += child.Width
	if child.Height > hb.Height {
		hb.Height = child.Height
	}
	if child.Depth > hb.Depth {
		hb.Depth = child.Depth
	}
}

// AddChild adds a child box to a VBox and updates dimensions
func (vb *VBox) AddChild(child *Box) {
	vb.Children = append(vb.Children, child)
	vb.Height += child.Height + child.Depth
	if child.Width > vb.Width {
		vb.Width = child.Width
	}
}

// ConvertASTToBoxes converts AST nodes to boxes (placeholder implementation)
func ConvertASTToBoxes(node parser.Node, font *fonts.Font) *Box {
	switch n := node.(type) {
	case *parser.TextNode:
		textBox := NewTextBox(n.Value, font)
		return &textBox.Box
	case *parser.Command:
		// For now, just treat commands as text -> should change soon
		textBox := NewTextBox("\\"+n.Name, font)
		return &textBox.Box
	default:
		// Default empty box
		return &Box{}
	}
}
