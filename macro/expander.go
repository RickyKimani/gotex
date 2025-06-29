package macro

import "github.com/rickykimani/gotex/parser"

// Expander handles macro expansion with a macro store
type Expander struct {
	store *MacroStore
}

// NewExpander creates a new macro expander
func NewExpander(store *MacroStore) *Expander {
	return &Expander{store: store}
}

// ExpandDocument expands all macros in a document
func (e *Expander) ExpandDocument(doc *parser.Document) *parser.Document {
	expandedNodes := expandNodes(doc.Body, e.store)
	return parser.NewDocument(expandedNodes, doc.Position)
}

func Expand(node parser.Node, store *MacroStore) parser.Node {
	switch n := node.(type) {
	case *parser.Command:
		if macro, exists := store.Get(n.Name); exists {
			return expandMacro(macro, n, store)
		}
	case *parser.Environment:
		newStore := NewMacroStore(store) // Nested scope
		n.Body = expandNodes(n.Body, newStore)
	}
	return node
}

func expandMacro(macro *Macro, cmd *parser.Command, store *MacroStore) parser.Node {
	if len(cmd.Args) < macro.NumArgs {
		return cmd // Don't expand if not enough args
	}

	var expanded []parser.Node
	for _, node := range macro.Definitions {
		if arg, ok := node.(*parser.ArgumentPlaceholder); ok {
			if arg.Index < len(cmd.Args) {
				expanded = append(expanded, Expand(cmd.Args[arg.Index], store))
			}
		} else {
			expanded = append(expanded, Expand(node, store))
		}
	}

	if len(expanded) == 1 {
		return expanded[0]
	}
	return &parser.Group{Nodes: expanded}
}

func expandNodes(nodes []parser.Node, store *MacroStore) []parser.Node {
	var result []parser.Node
	for _, node := range nodes {
		expanded := Expand(node, store)
		result = append(result, expanded)
	}
	return result
}
