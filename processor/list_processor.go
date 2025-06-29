package processor

import "fmt"

func (dp *DocumentProcessor) enterList(listType string) {
	dp.listLevel++
	dp.listType = append(dp.listType, listType)
	dp.listCounters = append(dp.listCounters, 0)
	dp.addVerticalSpace(5)
}

func (dp *DocumentProcessor) exitList() {
	if dp.listLevel > 0 {
		dp.listLevel--
		dp.listType = dp.listType[:len(dp.listType)-1]
		dp.listCounters = dp.listCounters[:len(dp.listCounters)-1]
		dp.newLine()              // Ensure proper line break before adding space
		dp.addVerticalSpace(15)   // Increased space after lists to prevent overlap
		dp.lineHasContent = false // Reset line state
	}
}

func (dp *DocumentProcessor) addListItem() {
	if dp.listLevel == 0 {
		return
	}

	dp.newLine()           // Ensure proper line positioning
	dp.addVerticalSpace(5) // Small space before item

	x := dp.generator.MarginLeft + float64((dp.listLevel-1)*15) // Use 15 instead of 20
	currentListType := dp.listType[dp.listLevel-1]

	if currentListType == "enumerate" {
		dp.listCounters[dp.listLevel-1]++
		number := dp.listCounters[dp.listLevel-1]
		dp.generator.AddText(fmt.Sprintf("%d.", number), x, dp.currentY, dp.fontSize, "normal")
	} else {
		dp.generator.AddText("• ", x, dp.currentY, dp.fontSize, "normal")
	}

	// The actual item text will be processed by subsequent nodes
	dp.currentX = x + 25     // Set indent for item content (reduced from 30 to 25)
	dp.lineHasContent = true // Mark that line has content
}
