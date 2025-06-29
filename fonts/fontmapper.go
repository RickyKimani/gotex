package fonts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
)

// FontMapper handles mapping our font styles to TTF files and loading them into gopdf
type FontMapper struct {
	pdf      *gopdf.GoPdf
	fontPath string
	fonts    map[string]string // Map style ("normal", "bold", etc.) to font key ("dejavu-regular", "cm-bold")
	loaded   map[string]bool
}

// NewFontMapper creates a new font mapper for the given PDF instance
func NewFontMapper(pdf *gopdf.GoPdf, ttfDir string) *FontMapper {
	return &FontMapper{
		pdf:      pdf,
		fontPath: ttfDir,
		fonts: map[string]string{
			"normal":      "dejavu-regular",
			"regular":     "pagella-regular",
			"bold":        "pagella-bold",
			"italic":      "pagella-italic",
			"bold-italic": "pagella-bold-italic",
		},
		loaded: make(map[string]bool),
	}
}

// LoadFonts loads all required fonts into the PDF
func (fm *FontMapper) LoadFonts() error {
	//TODO: Find a computer-modern sans with all required symbols
	//TODO: and replace pagella and dejavu-sans
	fontFiles := map[string]string{
		"dejavu-regular":      filepath.Join(fm.fontPath, "dejavu-sans", "DejaVuSans.ttf"),
		"pagella-regular":     filepath.Join(fm.fontPath, "pagella", "texgyrepagella-regular.ttf"),
		"pagella-bold":        filepath.Join(fm.fontPath, "pagella", "texgyrepagella-bold.ttf"),
		"pagella-italic":      filepath.Join(fm.fontPath, "pagella", "texgyrepagella-italic.ttf"),
		"pagella-bold-italic": filepath.Join(fm.fontPath, "pagella", "texgyrepagella-bolditalic.ttf"),
	}

	for fontKey, fontPath := range fontFiles {
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			return fmt.Errorf("font file not found: %s", fontPath)
		}

		// Load font into PDF
		err := fm.pdf.AddTTFFont(fontKey, fontPath)
		if err != nil {
			return fmt.Errorf("failed to load font %s: %v", fontKey, err)
		}

		fm.loaded[fontKey] = true
	}

	return nil
}

// GetFontKey returns the appropriate font key for our style names
func (fm *FontMapper) GetFontKey(style string) string {
	if key, ok := fm.fonts[style]; ok {
		return key
	}
	return fm.fonts["regular"] // Default to pagella regular
}

// SetFont sets the font for the PDF based on style names
func (fm *FontMapper) SetFont(style string, size float64) error {
	fontKey := fm.GetFontKey(style)

	if !fm.loaded[fontKey] {
		return fmt.Errorf("font not loaded: %s", fontKey)
	}

	return fm.pdf.SetFont(fontKey, "", size)
}

// IsLoaded checks if a font style is loaded
func (fm *FontMapper) IsLoaded(style string) bool {
	fontKey := fm.GetFontKey(style)
	return fm.loaded[fontKey]
}
