package processor

func (dp *DocumentProcessor) addWord(word string, style string) {
	if word == "" {
		return
	}

	wordWidth := dp.calculateTextWidth(word, style)

	if dp.lineHasContent && (dp.currentLineX+wordWidth > dp.generator.PageWidth-dp.generator.MarginRight) {
		dp.newLine()
	}

	dp.generator.AddText(word, dp.currentLineX, dp.currentY, dp.fontSize, style)
	dp.currentLineX += wordWidth
	dp.lineHasContent = true
}
