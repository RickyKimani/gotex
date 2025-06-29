package processor

func (dp *DocumentProcessor) addParagraphBreak() {
	if dp.lineHasContent {
		dp.newLine()
	}
	dp.addVerticalSpace(dp.lineHeight * 0.8)
	dp.lineHasContent = false
	dp.currentLineX = dp.generator.MarginLeft
}
