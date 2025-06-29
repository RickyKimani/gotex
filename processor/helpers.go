package processor

// calculateTextWidth calculates the actual width of text using the PDF generator
func (dp *DocumentProcessor) calculateTextWidth(text string, style string) float64 {
	return dp.generator.GetTextWidth(text, dp.fontSize, style)
}
