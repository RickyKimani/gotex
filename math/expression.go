package math

// MathElement represents a basic math element
type MathElement struct {
	Type     string  // "variable", "number", "operator", "symbol", "fraction", "superscript", "subscript"
	Content  string  // The actual content (x, +, 2, etc.)
	X        float64 // Position relative to math baseline
	Y        float64 // Position relative to math baseline
	Width    float64 // Width of this element
	Height   float64 // Height above baseline
	Depth    float64 // Depth below baseline
	FontSize float64 // Font size for this element
	Style    string  // Font style: "math-italic", "math-roman", "math-symbol"
}

// MathExpression represents a complete math expression
type MathExpression struct {
	Elements []*MathElement
	Width    float64 // Total width
	Height   float64 // Total height above baseline
	Depth    float64 // Total depth below baseline
	IsInline bool    // true for inline math, false for display math
	IsBuilt  bool    // true if layout has been calculated
}

// MathBox represents a positioned math element for rendering
type MathBox struct {
	Element   *MathElement
	AbsoluteX float64 // Absolute position in document
	AbsoluteY float64 // Absolute position in document
}

// NewMathExpression creates a new math expression
func NewMathExpression(inline bool) *MathExpression {
	return &MathExpression{
		Elements: make([]*MathElement, 0),
		IsInline: inline,
		IsBuilt:  false,
	}
}

// AddElement adds a math element to the expression
func (m *MathExpression) AddElement(element *MathElement) {
	m.Elements = append(m.Elements, element)
	m.IsBuilt = false // Mark as needing rebuild
}

// BuildLayout calculates positions and dimensions for all elements
func (m *MathExpression) BuildLayout() {
	if len(m.Elements) == 0 {
		return
	}

	currentX := 0.0
	maxHeight := 0.0
	maxDepth := 0.0

	for i, element := range m.Elements {
		// Set position
		element.X = currentX
		element.Y = 0.0 // Baseline for simple layout

		// Calculate spacing
		spacing := m.getElementSpacing(i)
		currentX += element.Width + spacing

		// Track max dimensions
		if element.Height > maxHeight {
			maxHeight = element.Height
		}
		if element.Depth > maxDepth {
			maxDepth = element.Depth
		}
	}

	m.Width = currentX
	m.Height = maxHeight
	m.Depth = maxDepth
	m.IsBuilt = true
}

// getElementSpacing returns the spacing after an element
func (m *MathExpression) getElementSpacing(index int) float64 {
	if index >= len(m.Elements)-1 {
		return 0.0 // No spacing after last element
	}

	current := m.Elements[index]
	next := m.Elements[index+1]

	// LaTeX spacing rules (simplified)
	switch current.Type {
	case "operator":
		if current.Content == "+" || current.Content == "-" {
			return 4.0 // Medium space around binary operators
		}
		return 2.0
	case "relation":
		if current.Content == "=" {
			return 5.0 // Thick space around relations
		}
		return 4.0
	case "variable", "number":
		if next.Type == "operator" || next.Type == "relation" {
			return 2.0 // Thin space before operators
		}
		return 1.0
	default:
		return 1.0 // Default thin space
	}
}

// GetBoundingBox returns the bounding box of the expression
func (m *MathExpression) GetBoundingBox() (width, height, depth float64) {
	if !m.IsBuilt {
		m.BuildLayout()
	}
	return m.Width, m.Height, m.Depth
}
