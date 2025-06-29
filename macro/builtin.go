package macro

import "github.com/rickykimani/gotex/parser"

type Macro struct {
	Name         string
	NumArgs      int
	Definitions  []parser.Node
	IsExpandable bool
}

type MacroStore struct {
	macros map[string]*Macro
	parent *MacroStore
}

func NewMacroStore(parent *MacroStore) *MacroStore {
	return &MacroStore{
		macros: make(map[string]*Macro),
		parent: parent,
	}
}

func (ms *MacroStore) AddBuiltins() {
	ms.macros["textbf"] = &Macro{
		Name:    "textbf",
		NumArgs: 1,
		Definitions: []parser.Node{
			&parser.Command{
				Name: "font",
				Args: []parser.Node{
					&parser.TextNode{Value: "bold"},
				},
			},
			// The argument will be inserted here
			&parser.ArgumentPlaceholder{Index: 0},
		},
	}

	ms.macros["textit"] = &Macro{
		Name:    "textit",
		NumArgs: 1,
		Definitions: []parser.Node{
			&parser.Command{
				Name: "font",
				Args: []parser.Node{
					&parser.TextNode{Value: "italic"},
				},
			},
			// The argument will be inserted here
			&parser.ArgumentPlaceholder{Index: 0},
		},
	}
}

func (ms *MacroStore) Get(name string) (*Macro, bool) {
	// Look in current scope first
	if macro, exists := ms.macros[name]; exists {
		return macro, true
	}

	// Look in parent scopes
	if ms.parent != nil {
		return ms.parent.Get(name)
	}

	return nil, false
}

func (ms *MacroStore) Set(name string, macro *Macro) {
	ms.macros[name] = macro
}
