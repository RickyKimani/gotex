package processor

import (
	"github.com/rickykimani/gotex/math"
	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/pdf"
)

type DocumentProcessor struct {
	generator     *pdf.Generator
	mathProcessor *math.MathProcessor
	currentY      float64
	currentX      float64
	currentLineX  float64
	lineHeight    float64
	fontSize      float64

	// Document metadata
	title  string
	author string
	date   string

	// List state
	listLevel    int
	listType     []string // "enumerate" or "itemize"
	listCounters []int

	// Section counters
	sectionCounter    int
	subsectionCounter int
	equationCounter   int

	// Line state tracking
	lineHasContent bool // Track if current line has content

	// Track if we just processed a command for spacing logic
	lastProcessedCommand bool
}

func NewDocumentProcessor(generator *pdf.Generator) *DocumentProcessor {
	generator.NewPage() // Ensure we have at least one page

	fontSize := 12.0

	return &DocumentProcessor{
		generator:            generator,
		mathProcessor:        math.NewMathProcessor(generator, fontSize),
		currentY:             generator.PageHeight - generator.MarginTop,
		currentX:             0,
		currentLineX:         generator.MarginLeft,
		lineHeight:           20.0,
		fontSize:             fontSize,
		listLevel:            0,
		listType:             make([]string, 0),
		listCounters:         make([]int, 0),
		sectionCounter:       0,
		subsectionCounter:    0,
		equationCounter:      0,
		lineHasContent:       false,
		lastProcessedCommand: false,
	}
}

func (dp *DocumentProcessor) ProcessDocument(nodes []parser.Node) {
	dp.processNodes(nodes, "normal")
}
