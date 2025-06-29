package processor

// newLine moves to the next line with consistent spacing
func (dp *DocumentProcessor) newLine() {
	dp.currentY -= dp.lineHeight
	dp.checkNewPage()
	dp.lineHasContent = false
	dp.currentLineX = dp.generator.MarginLeft
	if dp.listLevel > 0 {
		dp.currentLineX += float64(dp.listLevel * 15)
	}
}
