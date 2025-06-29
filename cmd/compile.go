package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rickykimani/gotex/lexer"
	"github.com/rickykimani/gotex/macro"
	"github.com/rickykimani/gotex/parser"
	"github.com/rickykimani/gotex/pdf"
	"github.com/rickykimani/gotex/processor"
)

// compileTeX compiles a .tex file to PDF
func compileTeX(inputFile, outputFile string) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	//TODO: Remove debug logs
	fmt.Printf("=== GoTeX Compiler ===\n")
	fmt.Printf("Input: %s\n", inputFile)
	fmt.Printf("Output: %s\n\n", outputFile)

	// Tokenize
	fmt.Println("Step 1: Tokenizing...")
	l := lexer.NewLexer(string(content))
	tokens := l.Tokenize()
	fmt.Printf("Generated %d tokens\n", len(tokens))

	// Parse
	fmt.Println("\nStep 2: Parsing...")
	p := parser.NewParser(tokens)
	doc, errors := p.Parse()

	// Show any parse errors using Rust-style error reporting
	if len(errors) > 0 {
		reporter := parser.NewErrorReporter(string(content), inputFile)
		reporter.ReportErrors(errors)
	}

	// Expand macros
	fmt.Println("Step 3: Expanding macros...")
	store := macro.NewMacroStore(nil)
	store.AddBuiltins()

	expander := macro.NewExpander(store)
	expandedDoc := expander.ExpandDocument(doc)
	fmt.Printf("Document expanded: %d nodes\n", len(expandedDoc.Nodes))

	// Generate PDF
	fmt.Println("\nStep 4: Generating PDF...")

	// Get TTF directory path
	execDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %v", err)
	}
	ttfDir := filepath.Join(filepath.Dir(execDir), "ttf")

	// Check if ttf directory exists, if not use relative path
	if _, err := os.Stat(ttfDir); os.IsNotExist(err) {
		ttfDir = "ttf"
	}

	generator, err := pdf.NewGenerator(ttfDir)
	if err != nil {
		return fmt.Errorf("error creating PDF generator: %v", err)
	}

	// Process the document and generate PDF
	docProcessor := processor.NewDocumentProcessor(generator)
	docProcessor.ProcessDocument(expandedDoc.Nodes)

	err = generator.GeneratePDF(outputFile)
	if err != nil {
		return fmt.Errorf("error generating PDF: %v", err)
	}

	fmt.Printf("\n✅ PDF generated successfully: %s\n", outputFile)
	fmt.Printf("📄 Pages: %d\n", generator.GetPageCount())

	return nil
}
