package processor

//TODO: Automate with \maketitle

func (dp *DocumentProcessor) addTitle() {
	if dp.title != "" {
		dp.addVerticalSpace(40) // More space before title
		dp.generator.AddTitle(dp.title, dp.currentX, dp.currentY, dp.fontSize)
		dp.currentY -= dp.lineHeight * 2.5 // More space after title to account for larger font
	}

	if dp.author != "" {
		dp.addVerticalSpace(15) // Slightly more space before author
		dp.generator.AddTextWithAlignment(dp.author, dp.currentX, dp.currentY, dp.fontSize, "normal", "center")
		dp.currentY -= dp.lineHeight
	}

	if dp.date != "" {
		dp.addVerticalSpace(8) // Slightly more space before date
		dp.generator.AddTextWithAlignment(dp.date, dp.currentX, dp.currentY, dp.fontSize, "normal", "center")
		dp.currentY -= dp.lineHeight * 2
	}

	dp.addVerticalSpace(40) // More space after title block
}
