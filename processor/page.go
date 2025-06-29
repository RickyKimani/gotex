package processor

func (dp *DocumentProcessor) checkNewPage() {
	if dp.currentY < dp.generator.MarginBottom+50 {
		dp.generator.NewPage()
		dp.currentY = dp.generator.PageHeight - dp.generator.MarginTop
		dp.lineHasContent = false // Reset line state on new page
	}
}
