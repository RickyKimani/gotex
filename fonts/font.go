package fonts

// Font represents a font with its metrics and properties
type Font struct {
	Name    string
	Family  string
	Style   string
	Weight  string
	Size    float64
	Metrics *FontMetrics
}

// FontMetrics holds detailed font measurement information
type FontMetrics struct {
	Height    float64 // Total height of the font
	Ascent    float64 // Distance from baseline to top
	Descent   float64 // Distance from baseline to bottom
	Leading   float64 // Line spacing
	XHeight   float64 // Height of lowercase 'x' (For non ascending characters)
	CapHeight float64 // Height of capital letters
}

// NewFont creates a new font with default metrics
func NewFont(name, family, style, weight string, size float64) *Font {
	return &Font{
		Name:    name,
		Family:  family,
		Style:   style,
		Weight:  weight,
		Size:    size,
		Metrics: NewFontMetrics(size),
	}
}

// GetCharWidth returns the width of a character in this font
func (f *Font) GetCharWidth(char rune) float64 {
	//TODO: Automate font width calculation
	baseWidth := f.Size * 0.6

	switch char {
	case ' ':
		return baseWidth * 0.5
	case 'i', 'l', 'I', 'j', 't', 'f':
		return baseWidth * 0.3
	case 'm', 'w', 'M', 'W':
		return baseWidth * 1.6
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'X', 'Y', 'Z':
		return baseWidth * 1.2
	case '.', ',', ':', ';', '!', '|':
		return baseWidth * 0.3
	case '-', '=', '+':
		return baseWidth * 0.7
	case '(', ')', '[', ']', '{', '}':
		return baseWidth * 0.4
	case '"', '\'':
		return baseWidth * 0.3
	default:
		return baseWidth
	}
}

// GetStringWidth calculates the total width of a string in this font
func (f *Font) GetStringWidth(text string) float64 {
	width := 0.0
	for _, char := range text {
		width += f.GetCharWidth(char)
	}
	return width
}
