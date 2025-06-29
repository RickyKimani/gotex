package processor

import "github.com/rickykimani/gotex/parser"

func (dp *DocumentProcessor) processCommand(cmd *parser.Command, style string) {
	switch cmd.Name {
	case "documentclass":
		// Skip document class
		//TODO: Implement different document classes

	case "usepackage":
		// Skip packages for now
		//TODO: Implement universally compatible package system

	case "title":
		if len(cmd.Args) > 0 {
			dp.title = dp.extractText(cmd.Args[0])
		}

	case "author":
		if len(cmd.Args) > 0 {
			dp.author = dp.extractText(cmd.Args[0])
		}

	case "date":
		if len(cmd.Args) > 0 {
			dp.date = dp.extractText(cmd.Args[0])
		}

	case "maketitle":
		dp.addTitle()

	case "section":
		if len(cmd.Args) > 0 {
			text := dp.extractText(cmd.Args[0])
			dp.addSection(text)
		}

	case "subsection":
		if len(cmd.Args) > 0 {
			text := dp.extractText(cmd.Args[0])
			dp.addSubsection(text)
		}

	case "subsubsection":
		if len(cmd.Args) > 0 {
			text := dp.extractText(cmd.Args[0])
			dp.addSubsubsection(text)
		}

	case "font":
		// Handle font changes from macro expansion
		// This command sets the style for subsequent text in the same group
		// We don't process anything here, at least not now

	case "textbf":
		if len(cmd.Args) > 0 {
			newStyle := "bold"
			if style == "italic" {
				newStyle = "bold-italic"
			}
			dp.processStyledText(cmd.Args[0], newStyle)
		}

	case "textit":
		if len(cmd.Args) > 0 {
			newStyle := "italic"
			if style == "bold" {
				newStyle = "bold-italic"
			}
			dp.processStyledText(cmd.Args[0], newStyle)
		}

	case "item":
		dp.addListItem()
		// Process the content that follows the \item command

	//Might be useless, will test later
	case "today":
		// Return current date in LaTeX format
		return // \today is handled in extractText

	default:
		// Process command arguments
		for _, arg := range cmd.Args {
			dp.processNode(arg, style)
		}
	}
}
