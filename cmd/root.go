package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	outputFile string
	scanMode   bool
)

var rootCmd = &cobra.Command{
	Use:   "gotex [input.tex]",
	Short: "GoTeX - A Go-based TeX to PDF compiler",
	Long: `GoTeX is a TeX to PDF compiler written in Go
and LaTeX-like compilation capabilities.

Examples:
  gotex document.tex                    # Compile document.tex to document.pdf
  gotex document.tex -o report.pdf     # Compile to custom output name
  gotex assignment                      # Compile assignment.tex (if it exists)
  gotex --scan                          # Scan current directory for .tex files`,
	Args: func(cmd *cobra.Command, args []string) error {
		// no arguments are needed for scan
		if scanMode {
			return nil
		}

		// Otherwise, exactly one argument is required
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one .tex file argument")
		}

		inputArg := args[0]

		// Check if the argument has .tex extension
		if strings.HasSuffix(strings.ToLower(inputArg), ".tex") {
			// Check if file exists
			if _, err := os.Stat(inputArg); os.IsNotExist(err) {
				return fmt.Errorf("file '%s' does not exist", inputArg)
			}
			return nil
		}

		// If no .tex extension, check for other extensions and reject them
		if strings.Contains(inputArg, ".") {
			ext := filepath.Ext(inputArg)
			return fmt.Errorf("unsupported file extension '%s'. GoTeX only processes .tex files", ext)
		}

		// No extension provided - try to find inputArg.tex
		texFilePath := inputArg + ".tex"
		if _, err := os.Stat(texFilePath); os.IsNotExist(err) {
			return fmt.Errorf("%s not found in current directory", texFilePath)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if scanMode {
			return scanTexFiles()
		}

		inputFile := args[0]

		// If input doesn't have .tex extension, append .tex
		if !strings.HasSuffix(strings.ToLower(inputFile), ".tex") {
			inputFile = inputFile + ".tex"
		}

		// Determine output file name
		if outputFile == "" {
			// Use base name of input file with .pdf extension
			base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
			outputFile = base + ".pdf"
		} else {
			// Clean up the output file name:
			// 1. Remove any existing .pdf extension to avoid double extensions
			// 2. Then add .pdf extension
			outputFile = strings.TrimSuffix(outputFile, ".pdf") + ".pdf"
		}

		return compileTeX(inputFile, outputFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output PDF file name (default: input basename + .pdf)")
	rootCmd.Flags().BoolVarP(&scanMode, "scan", "s", false, "Scan current directory for .tex files")
}

func scanTexFiles() error {
	fmt.Println("🔍 Scanning for .tex files in current directory...")
	fmt.Println()

	var texFiles []string

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Warning: error accessing %s: %v\n", path, err)
			return nil // Continue walking even if there's an error
		}

		// Skip hidden directories and files (but not the current directory)
		if path != "." && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for .tex files
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".tex") {
			texFiles = append(texFiles, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error scanning directory: %v", err)
	}

	if len(texFiles) == 0 {
		fmt.Println("No .tex files found in current directory.")
		return nil
	}

	fmt.Printf("Found %d .tex file(s):\n", len(texFiles))

	// First pass: collect file info and find max widths for alignment
	type fileInfo struct {
		path string
		size int64
		err  error
	}

	fileInfos := make([]fileInfo, len(texFiles))
	maxPathLen := 0
	maxSizeLen := 0

	for i, file := range texFiles {
		info, err := os.Stat(file)
		fileInfos[i] = fileInfo{path: file, err: err}

		if err == nil {
			fileInfos[i].size = info.Size()
			// Find max lengths for alignment
			if len(file) > maxPathLen {
				maxPathLen = len(file)
			}
			sizeStr := formatBytes(info.Size())
			if len(sizeStr) > maxSizeLen {
				maxSizeLen = len(sizeStr)
			}
		}
	}

	// Second pass: print with justified formatting
	for i, info := range fileInfos {
		if info.err != nil {
			fmt.Printf("  %2d. %-*s (error reading file info)\n", i+1, maxPathLen, info.path)
			continue
		}

		sizeStr := formatBytes(info.size)
		fmt.Printf("  %2d. %-*s %*s\n", i+1, maxPathLen, info.path, maxSizeLen, sizeStr)
	}

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"K", "M", "G"}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %sB", float64(bytes)/float64(div), units[exp])
}
