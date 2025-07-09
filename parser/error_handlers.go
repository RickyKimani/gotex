package parser

import "github.com/rickykimani/gotex/lexer"

func (p *Parser) addError(errType ErrorType, message string, severity Severity) {
	err := ParseError{
		Type:     errType,
		Message:  message,
		Position: p.curToken.Pos,
		Severity: severity,
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) addErrorAtPosition(errType ErrorType, message string, severity Severity, pos lexer.Position) {
	err := ParseError{
		Type:     errType,
		Message:  message,
		Position: pos,
		Severity: severity,
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) GetErrors() []ParseError {
	return p.errors
}

func (p *Parser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *Parser) HasFatalErrors() bool {
	for _, err := range p.errors {
		if err.Severity == Fatal {
			return true
		}
	}
	return false
}

// addWarning is a convenience function for adding warnings
func (p *Parser) addWarning(errType ErrorType, message string) {
	p.addError(errType, message, Warning)
}
