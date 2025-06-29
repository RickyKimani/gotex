package parser

import "github.com/rickykimani/gotex/lexer"

type Node interface {
	Pos() lexer.Position
}

type Command struct {
	Name     string
	Args     []Node
	Optional []Node
	Position lexer.Position
}

type Environment struct {
	Name     string
	Body     []Node
	Content  []Node // Alias for Body for compatibility
	Position lexer.Position
}

type TextNode struct {
	Value    string
	Position lexer.Position
}

type MathNode struct {
	Inline   bool
	Content  []Node
	Position lexer.Position
}

type Document struct {
	Body     []Node
	Nodes    []Node // Alias for Body for compatibility
	Position lexer.Position
}

type CommentNode struct {
	Value    string
	Position lexer.Position
}

// ErrorRecoveryNode represents content that couldn't be parsed correctly
type ErrorRecoveryNode struct {
	Error    ParseError
	RawText  string // The problematic text
	Position lexer.Position
}

type ArgumentPlaceholder struct {
	Index    int // Which argument (0, 1, 2, ...)
	Position lexer.Position
}

type Group struct {
	Nodes    []Node
	Content  []Node // Alias for Nodes for compatibility
	Position lexer.Position
}

// Enhanced math node types for better math processing
type MathSymbol struct {
	Symbol   string // The actual symbol (α, β, ∑, etc.)
	Command  string // The LaTeX command (\alpha, \beta, \sum, etc.)
	Position lexer.Position
}

type MathOperator struct {
	Operator string // +, -, *, /, =, <, >, etc.
	Position lexer.Position
}

type MathSuperscript struct {
	Base     Node
	Exponent Node
	Position lexer.Position
}

type MathSubscript struct {
	Base     Node
	Index    Node
	Position lexer.Position
}

type MathFraction struct {
	Numerator   Node
	Denominator Node
	Position    lexer.Position
}

// Implement the Node interface for all AST nodes
func (c *Command) Pos() lexer.Position             { return c.Position }
func (e *Environment) Pos() lexer.Position         { return e.Position }
func (t *TextNode) Pos() lexer.Position            { return t.Position }
func (m *MathNode) Pos() lexer.Position            { return m.Position }
func (d *Document) Pos() lexer.Position            { return d.Position }
func (c *CommentNode) Pos() lexer.Position         { return c.Position }
func (e *ErrorRecoveryNode) Pos() lexer.Position   { return e.Position }
func (a *ArgumentPlaceholder) Pos() lexer.Position { return a.Position }
func (g *Group) Pos() lexer.Position               { return g.Position }
func (m *MathSymbol) Pos() lexer.Position          { return m.Position }
func (m *MathOperator) Pos() lexer.Position        { return m.Position }
func (m *MathSuperscript) Pos() lexer.Position     { return m.Position }
func (m *MathSubscript) Pos() lexer.Position       { return m.Position }
func (m *MathFraction) Pos() lexer.Position        { return m.Position }

// NewDocument creates a new document with synchronized fields
func NewDocument(nodes []Node, pos lexer.Position) *Document {
	return &Document{
		Body:     nodes,
		Nodes:    nodes, // Keep in sync
		Position: pos,
	}
}

// NewEnvironment creates a new environment with synchronized fields
func NewEnvironment(name string, content []Node, pos lexer.Position) *Environment {
	return &Environment{
		Name:     name,
		Body:     content,
		Content:  content, // Keep in sync
		Position: pos,
	}
}

// NewGroup creates a new group with synchronized fields
func NewGroup(nodes []Node, pos lexer.Position) *Group {
	return &Group{
		Nodes:    nodes,
		Content:  nodes, // Keep in sync
		Position: pos,
	}
}
