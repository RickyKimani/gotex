package typesetter

// LineBreak handles line breaking algorithms similar to TeX
type LineBreak struct {
	LineWidth   float64
	Tolerance   int
	Penalties   map[string]int
	MinDemerits int
}

// BreakPoint represents a potential line break point
type BreakPoint struct {
	Position int
	Demerits int
	Previous *BreakPoint
	Width    float64
	Stretch  float64
	Shrink   float64
}

// NewLineBreak creates a new line breaking instance
func NewLineBreak(lineWidth float64) *LineBreak {
	return &LineBreak{
		LineWidth:   lineWidth,
		Tolerance:   200,
		MinDemerits: 10000,
		Penalties: map[string]int{
			"hyphen":      50,
			"line":        10,
			"widow":       150,
			"club":        150,
			"broken-math": 100,
		},
	}
}

//TODO: Implement full Knuth-Plass algo

// BreakParagraph breaks a paragraph into lines using a simplified TeX algorithm
func (lb *LineBreak) BreakParagraph(boxes []*Box) [][]*Box {
	if len(boxes) == 0 {
		return [][]*Box{}
	}

	breakPoints := lb.findBreakPoints(boxes)
	return lb.constructLines(boxes, breakPoints)
}

// findBreakPoints identifies potential line break points
func (lb *LineBreak) findBreakPoints(boxes []*Box) []*BreakPoint {
	breakPoints := []*BreakPoint{}

	// Add start point
	breakPoints = append(breakPoints, &BreakPoint{
		Position: 0,
		Demerits: 0,
		Width:    0,
	})

	currentWidth := 0.0

	for i, box := range boxes {
		currentWidth += box.Width

		// Check if we can break at this position
		if lb.canBreakAt(i, boxes) {
			demerits := lb.calculateDemerits(currentWidth, lb.LineWidth)

			breakPoint := &BreakPoint{
				Position: i + 1,
				Demerits: demerits,
				Width:    currentWidth,
			}

			breakPoints = append(breakPoints, breakPoint)

			// Reset width for next line calculation
			currentWidth = 0.0
		}
	}

	return breakPoints
}

// canBreakAt determines if a line break is allowed at this position
func (lb *LineBreak) canBreakAt(position int, boxes []*Box) bool {
	if position >= len(boxes)-1 {
		return true // Can always break at end
	}

	// Simplified: allow breaks at spaces or after certain characters
	box := boxes[position]
	if textBox, ok := box.Content.(*TextBox); ok {
		text := textBox.Text
		if len(text) > 0 {
			lastChar := text[len(text)-1]
			return lastChar == ' ' || lastChar == '-'
		}
	}

	return false
}

// calculateDemerits calculates the penalty for breaking at this point
func (lb *LineBreak) calculateDemerits(actualWidth, targetWidth float64) int {
	ratio := actualWidth / targetWidth

	if ratio > 1.0 {
		// Line is too long - high penalty
		excess := ratio - 1.0
		return int(excess * 1000)
	} else if ratio < 0.8 {
		// Line is too short - moderate penalty
		shortage := 0.8 - ratio
		return int(shortage * 500)
	}

	// Good fit - low penalty
	return 10
}

// constructLines builds the final line structure from break points
func (lb *LineBreak) constructLines(boxes []*Box, breakPoints []*BreakPoint) [][]*Box {
	lines := [][]*Box{}

	if len(breakPoints) < 2 {
		// Single line
		return [][]*Box{boxes}
	}

	for i := 0; i < len(breakPoints)-1; i++ {
		start := breakPoints[i].Position
		end := breakPoints[i+1].Position

		if end > len(boxes) {
			end = len(boxes)
		}

		if start < end {
			line := make([]*Box, end-start)
			copy(line, boxes[start:end])
			lines = append(lines, line)
		}
	}

	return lines
}

// JustifyLine adjusts spacing in a line to fit the target width
func (lb *LineBreak) JustifyLine(line []*Box, targetWidth float64) {
	if len(line) == 0 {
		return
	}

	totalWidth := 0.0
	spaceCount := 0

	// Calculate current width and count spaces
	for _, box := range line {
		totalWidth += box.Width
		if textBox, ok := box.Content.(*TextBox); ok {
			for _, char := range textBox.Text {
				if char == ' ' {
					spaceCount++
				}
			}
		}
	}

	if spaceCount == 0 || totalWidth >= targetWidth {
		return // Can't justify
	}

	// Distribute extra space among spaces
	extraSpace := targetWidth - totalWidth
	spaceAdjustment := extraSpace / float64(spaceCount)

	// Apply adjustments (simplified - in practice would modify actual spacing)
	_ = spaceAdjustment // TODO: implement actual spacing adjustment
}
