package pdf

import (
	"fmt"

	"github.com/rickykimani/gotex/fonts"
	"github.com/signintech/gopdf"
)

// Generator handles PDF generation using gopdf with the font(s) set in fonts dir
type Generator struct {
	pdf        *gopdf.GoPdf
	fontMapper *fonts.FontMapper

	// Page dimensions (points)
	PageWidth    float64
	PageHeight   float64
	MarginTop    float64
	MarginRight  float64
	MarginBottom float64
	MarginLeft   float64

	// State tracking
	CurrentPage int
	pageCount   int
}

// NewGenerator creates a new PDF generator with gopdf backend
func NewGenerator(ttfDir string) (*Generator, error) {
	pdf := &gopdf.GoPdf{}

	// Configure page settings
	pageConfig := gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
		Unit:     gopdf.UnitPT,
	}

	pdf.Start(pageConfig)

	// Initialize font mapper
	fontMapper := fonts.NewFontMapper(pdf, ttfDir)

	//TODO: Automate inch to pt with gopdf tools

	generator := &Generator{
		pdf:          pdf,
		fontMapper:   fontMapper,
		PageWidth:    pageConfig.PageSize.W,
		PageHeight:   pageConfig.PageSize.H,
		MarginTop:    72.0, // 1 inch = 72 points
		MarginRight:  72.0, // 1 inch = 72 points
		MarginBottom: 72.0, // 1 inch = 72 points
		MarginLeft:   72.0, // 1 inch = 72 points
		CurrentPage:  0,
		pageCount:    0,
	}

	// Load Computer Modern fonts
	if err := fontMapper.LoadFonts(); err != nil {
		return nil, fmt.Errorf("failed to load fonts: %v", err)
	}

	return generator, nil
}

// NewPage creates a new page
func (g *Generator) NewPage() {
	g.pdf.AddPage()
	g.CurrentPage++
	g.pageCount++
}

// AddText adds text at the specified position with the given style
func (g *Generator) AddText(text string, x, y, fontSize float64, style string) {
	if text == "" {
		return
	}

	// Set font for the style
	if err := g.fontMapper.SetFont(style, fontSize); err != nil {
		// Fallback to regular font if style not available
		g.fontMapper.SetFont("normal", fontSize)
	}

	// Convert our coordinate system (top-left origin) to PDF coordinate system (bottom-left origin)
	pdfY := g.PageHeight - y

	g.pdf.SetX(x)
	g.pdf.SetY(pdfY)
	g.pdf.Text(text)
}

// AddTextWithAlignment adds text with specified alignment
func (g *Generator) AddTextWithAlignment(text string, x, y, fontSize float64, style, alignment string) {
	if text == "" {
		return
	}

	// Set font for the style
	if err := g.fontMapper.SetFont(style, fontSize); err != nil {
		g.fontMapper.SetFont("normal", fontSize)
	}

	// Calculate alignment offset
	textWidth, _ := g.pdf.MeasureTextWidth(text)

	var alignedX float64
	switch alignment {
	case "center":
		alignedX = (g.PageWidth - textWidth) / 2
	case "right":
		alignedX = g.PageWidth - g.MarginRight - textWidth
	default: // left
		alignedX = x
	}

	// Convert coordinate system
	pdfY := g.PageHeight - y

	g.pdf.SetX(alignedX)
	g.pdf.SetY(pdfY)
	g.pdf.Text(text)
}

// AddTitle adds a title with larger font size
func (g *Generator) AddTitle(text string, x, y, baseFontSize float64) {
	titleSize := baseFontSize * 1.8
	g.AddTextWithAlignment(text, x, y, titleSize, "normal", "center")
}

// AddSection adds a section heading
func (g *Generator) AddSection(text string, x, y, baseFontSize float64) {
	sectionSize := baseFontSize * 1.4
	g.AddText(text, x, y, sectionSize, "bold")
}

// AddSubsection adds a subsection heading
func (g *Generator) AddSubsection(text string, x, y, baseFontSize float64) {
	subsectionSize := baseFontSize * 1.2
	g.AddText(text, x, y, subsectionSize, "bold")
}

// AddLine adds a line from (x1,y1) to (x2,y2)
func (g *Generator) AddLine(x1, y1, x2, y2 float64) {
	// Convert coordinate system
	pdfY1 := g.PageHeight - y1
	pdfY2 := g.PageHeight - y2

	// Set stroke properties to make the line visible
	g.pdf.SetStrokeColor(0, 0, 0) // Black
	g.pdf.SetLineWidth(0.5)       // Thin line
	g.pdf.Line(x1, pdfY1, x2, pdfY2)
}

// Legacy compatibility methods that map old font names to new styles

// AddBoldText adds bold text (legacy compatibility)
func (g *Generator) AddBoldText(text string, x, y, fontSize float64) {
	g.AddText(text, x, y, fontSize, "bold")
}

// AddItalicText adds italic text (legacy compatibility)
func (g *Generator) AddItalicText(text string, x, y, fontSize float64) {
	g.AddText(text, x, y, fontSize, "italic")
}

// AddBoldItalicText adds bold italic text (legacy compatibility)
func (g *Generator) AddBoldItalicText(text string, x, y, fontSize float64) {
	g.AddText(text, x, y, fontSize, "bold-italic")
}

// GetContentWidth returns the available content width
func (g *Generator) GetContentWidth() float64 {
	return g.PageWidth - g.MarginLeft - g.MarginRight
}

// GetContentHeight returns the available content height
func (g *Generator) GetContentHeight() float64 {
	return g.PageHeight - g.MarginTop - g.MarginBottom
}

// GetTextWidth returns the width of text in the given style
func (g *Generator) GetTextWidth(text string, fontSize float64, style string) float64 {
	// Set font temporarily to measure
	if err := g.fontMapper.SetFont(style, fontSize); err != nil {
		g.fontMapper.SetFont("normal", fontSize)
	}

	width, _ := g.pdf.MeasureTextWidth(text)
	return width
}

// GetPageCount returns the current number of pages
func (g *Generator) GetPageCount() int {
	return g.pageCount
}

// GeneratePDF writes the PDF to a file
func (g *Generator) GeneratePDF(filename string) error {
	return g.pdf.WritePdf(filename)
}
