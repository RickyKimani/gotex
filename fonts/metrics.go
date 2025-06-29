package fonts

// NewFontMetrics creates default font metrics for a given font size
func NewFontMetrics(size float64) *FontMetrics {
	return &FontMetrics{
		Height:    size * 1.2, // 120% of font size
		Ascent:    size * 0.8, // 80% above baseline
		Descent:   size * 0.2, // 20% below baseline
		Leading:   size * 0.2, // 20% for line spacing
		XHeight:   size * 0.5, // 50% for x-height
		CapHeight: size * 0.7, // 70% for capital height
	}
}

// GetMetric returns a specific metric value by name
func (fm *FontMetrics) GetMetric(name string) float64 {
	switch name {
	case "height":
		return fm.Height
	case "ascent":
		return fm.Ascent
	case "descent":
		return fm.Descent
	case "leading":
		return fm.Leading
	case "xheight":
		return fm.XHeight
	case "capheight":
		return fm.CapHeight
	default:
		return 0.0
	}
}

// ScaleMetrics scales all metrics by a given factor
func (fm *FontMetrics) ScaleMetrics(factor float64) *FontMetrics {
	return &FontMetrics{
		Height:    fm.Height * factor,
		Ascent:    fm.Ascent * factor,
		Descent:   fm.Descent * factor,
		Leading:   fm.Leading * factor,
		XHeight:   fm.XHeight * factor,
		CapHeight: fm.CapHeight * factor,
	}
}
