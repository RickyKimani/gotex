package parser

import (
	"fmt"
	"slices"

	"github.com/rickykimani/gotex/lexer"
)

func (p *Parser) parseEnvironment() *Environment {
	env := &Environment{
		Name:     p.curToken.Value,
		Body:     []Node{},
		Position: p.curToken.Pos,
	}

	// Check if this is a math environment that needs special parsing
	if p.isMathEnvironment(env.Name) {
		return p.parseMathEnvironment(env)
	}

	// Parse the environment body until we find \end{envname}
	found := false
	for p.peekToken.Type != lexer.TokenEOF {
		p.nextToken()

		// Check if this is the end of the environment
		if p.curToken.Type == lexer.TokenEndEnv {
			if p.curToken.Value == env.Name {
				found = true
				break
			} else {
				// Wrong environment name - report error but close the environment anyway
				p.addError(UnmatchedEnvironment,
					fmt.Sprintf("Found \\end{%s} but expected \\end{%s}",
						p.curToken.Value, env.Name), Error)
				// Close this environment anyway for error recovery
				found = true
				break
			}
		}

		// Parse the content inside the environment
		node := p.parseNode()
		if node != nil {
			env.Body = append(env.Body, node)
		}
	}

	// Check if we found the matching \end{}
	if !found {
		p.addError(UnmatchedEnvironment,
			fmt.Sprintf("Missing \\end{%s} for environment", env.Name), Error)
	} else {
		// Check for empty environments (warning)
		if len(env.Body) == 0 {
			p.addWarning(EmptyEnvironment,
				fmt.Sprintf("Empty environment: %s", env.Name))
		}
	}

	return env
}

// isMathEnvironment checks if an environment should be parsed as math content
func (p *Parser) isMathEnvironment(name string) bool {
	mathEnvironments := []string{"equation", "align", "gather", "multline", "split"}
	return slices.Contains(mathEnvironments, name)
}

// parseMathEnvironment parses math environments with special handling for braces
func (p *Parser) parseMathEnvironment(env *Environment) *Environment {
	// Parse the environment body with math-aware parsing
	found := false
	for p.peekToken.Type != lexer.TokenEOF {
		p.nextToken()

		// Check if this is the end of the environment
		if p.curToken.Type == lexer.TokenEndEnv {
			if p.curToken.Value == env.Name {
				found = true
				break
			} else {
				// Wrong environment name - report error but close the environment anyway
				p.addError(UnmatchedEnvironment,
					fmt.Sprintf("Found \\end{%s} but expected \\end{%s}",
						p.curToken.Value, env.Name), Error)
				// Close this environment anyway for error recovery
				found = true
				break
			}
		}

		// Parse the content inside the math environment with special brace handling
		node := p.parseMathEnvironmentNode()
		if node != nil {
			env.Body = append(env.Body, node)
		}
	}

	// Check if we found the matching \end{}
	if !found {
		p.addError(UnmatchedEnvironment,
			fmt.Sprintf("Missing \\end{%s} for environment", env.Name), Error)
	} else {
		// Check for empty environments (warning)
		if len(env.Body) == 0 {
			p.addWarning(EmptyEnvironment,
				fmt.Sprintf("Empty environment: %s", env.Name))
		}
	}

	return env
}

// parseMathEnvironmentNode parses nodes inside math environments with brace support
func (p *Parser) parseMathEnvironmentNode() Node {
	switch p.curToken.Type {
	case lexer.TokenCommand:
		return p.parseCommand()
	case lexer.TokenBeginEnv:
		return p.parseEnvironment()
	case lexer.TokenText:
		return p.parseText()
	case lexer.TokenMathInline, lexer.TokenMathDisplay:
		return p.parseMath()
	case lexer.TokenComment:
		return p.parseComment()
	case lexer.TokenLBrace:
		// In math environments, braces create groups
		return p.parseGroup()
	case lexer.TokenRBrace:
		// Unmatched closing brace - this is still a warning
		p.addWarning(UnexpectedToken, "unexpected '}' - no matching '{'")
		return nil
	default:
		// For other unexpected tokens
		return nil
	}
}

// parseGroup parses a group starting with { and ending with }
func (p *Parser) parseGroup() *Group {
	group := &Group{
		Nodes:    []Node{},
		Position: p.curToken.Pos,
	}

	// Parse until we find the matching }
	for p.peekToken.Type != lexer.TokenEOF {
		p.nextToken()

		if p.curToken.Type == lexer.TokenRBrace {
			// Found the matching closing brace
			break
		}

		// Parse the content inside the group using math environment parsing
		node := p.parseMathEnvironmentNode()
		if node != nil {
			group.Nodes = append(group.Nodes, node)
		}
	}

	// Check if we properly closed the group
	if p.curToken.Type != lexer.TokenRBrace {
		p.addError(UnmatchedBrace, "missing closing '}' for group", Error)
	}

	return group
}
