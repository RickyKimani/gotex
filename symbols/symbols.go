package symbols

import "strings"

//TODO: Add more symbols if possible

// MathSymbolTable maps LaTeX math commands to Unicode symbols for parser lookups
var MathSymbolTable = map[string]string{
	// Greek letters (lowercase)
	"alpha":   "α",
	"beta":    "β",
	"gamma":   "γ",
	"delta":   "δ",
	"epsilon": "ε",
	"zeta":    "ζ",
	"eta":     "η",
	"theta":   "θ",
	"iota":    "ι",
	"kappa":   "κ",
	"lambda":  "λ",
	"mu":      "μ",
	"nu":      "ν",
	"xi":      "ξ",
	"pi":      "π",
	"rho":     "ρ",
	"sigma":   "σ",
	"tau":     "τ",
	"upsilon": "υ",
	"phi":     "φ",
	"chi":     "χ",
	"psi":     "ψ",
	"omega":   "ω",

	// Greek letters (uppercase)
	"Alpha":   "Α",
	"Beta":    "Β",
	"Gamma":   "Γ",
	"Delta":   "Δ",
	"Epsilon": "Ε",
	"Zeta":    "Ζ",
	"Eta":     "Η",
	"Theta":   "Θ",
	"Iota":    "Ι",
	"Kappa":   "Κ",
	"Lambda":  "Λ",
	"Mu":      "Μ",
	"Nu":      "Ν",
	"Xi":      "Ξ",
	"Pi":      "Π",
	"Rho":     "Ρ",
	"Sigma":   "Σ",
	"Tau":     "Τ",
	"Upsilon": "Υ",
	"Phi":     "Φ",
	"Chi":     "Χ",
	"Psi":     "Ψ",
	"Omega":   "Ω",

	// Common math symbols
	"sum":     "∑",
	"prod":    "∏",
	"int":     "∫",
	"oint":    "∮",
	"partial": "∂",
	"infty":   "∞",
	"nabla":   "∇",
	"pm":      "±",
	"mp":      "∓",
	"times":   "×",
	"div":     "÷",
	"cdot":    "·",
	"bullet":  "•",

	// Relations
	"leq":    "≤",
	"geq":    "≥",
	"neq":    "≠",
	"equiv":  "≡",
	"approx": "≈",
	"sim":    "∼",
	"simeq":  "≃",
	"cong":   "≅",
	"propto": "∝",

	// Set theory
	"in":       "∈",
	"notin":    "∉",
	"subset":   "⊂",
	"supset":   "⊃",
	"subseteq": "⊆",
	"supseteq": "⊇",
	"cup":      "∪",
	"cap":      "∩",
	"emptyset": "∅",

	// Logic
	"forall":  "∀",
	"exists":  "∃",
	"neg":     "¬",
	"land":    "∧",
	"lor":     "∨",
	"implies": "⇒",
	"iff":     "⇔",

	// Arrows
	"leftarrow":      "←",
	"rightarrow":     "→",
	"uparrow":        "↑",
	"downarrow":      "↓",
	"leftrightarrow": "↔",
	"Leftarrow":      "⇐",
	"Rightarrow":     "⇒",
	"Uparrow":        "⇑",
	"Downarrow":      "⇓",
	"Leftrightarrow": "⇔",
}

// ConvertMathSymbol checks if a command corresponds to a known math symbol
func ConvertMathSymbol(command string) (string, bool) {
	cmd := strings.TrimPrefix(command, "\\")
	symbol, exists := MathSymbolTable[cmd]
	return symbol, exists
}

// IsMathSymbol returns true if a command is a known math symbol
func IsMathSymbol(command string) bool {
	cmd := strings.TrimPrefix(command, "\\")
	_, exists := MathSymbolTable[cmd]
	return exists
}
