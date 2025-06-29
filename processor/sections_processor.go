package processor

import (
	"fmt"
)

//TODO: Style sections better

func (dp *DocumentProcessor) addSection(text string) {
	// space before: (3.5ex + 1ex)
	ex := dp.fontSize * 0.5
	dp.newLine()
	dp.addVerticalSpace((3.5 + 1.0) * ex)

	// Increment section counter and reset subsection counter
	dp.sectionCounter++
	dp.subsectionCounter = 0
	dp.equationCounter = 0 // Reset equation counter for new section

	// Add section number prefix
	numberedText := fmt.Sprintf("%d %s", dp.sectionCounter, text)
	dp.generator.AddSection(numberedText, dp.generator.MarginLeft, dp.currentY, dp.fontSize)
	// space after: 2.3ex
	dp.addVerticalSpace(2.3 * ex)
	dp.lineHasContent = false // Reset line state
}

func (dp *DocumentProcessor) addSubsection(text string) {
	// space before: (3.25ex + 1ex)
	ex := dp.fontSize * 0.5
	dp.newLine()
	dp.addVerticalSpace((3.25 + 1.0) * ex)

	// Increment subsection counter
	dp.subsectionCounter++

	// Add subsection number prefix
	numberedText := fmt.Sprintf("%d.%d %s", dp.sectionCounter, dp.subsectionCounter, text)
	dp.generator.AddSubsection(numberedText, dp.generator.MarginLeft, dp.currentY, dp.fontSize)
	// space after: 1.5ex
	dp.addVerticalSpace(1.5 * ex)
	dp.lineHasContent = false // Reset line state
}

func (dp *DocumentProcessor) addSubsubsection(text string) {
	// space before: (3.25ex + 1ex)
	ex := dp.fontSize * 0.5
	dp.newLine()
	dp.addVerticalSpace((3.25 + 1.0) * ex)
	dp.generator.AddText(text, dp.generator.MarginLeft, dp.currentY, dp.fontSize*1.05, "bold")
	// space after: 1.5ex
	dp.addVerticalSpace(1.5 * ex)
	dp.lineHasContent = false // Reset line state
}
