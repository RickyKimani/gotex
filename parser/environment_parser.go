package parser

import (
	"fmt"

	"github.com/rickykimani/gotex/lexer"
)

func (p *Parser) parseEnvironment() *Environment {
	env := &Environment{
		Name:     p.curToken.Value,
		Body:     []Node{},
		Position: p.curToken.Pos,
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
	}

	// Sync the Content field with Body for compatibility
	env.Content = env.Body

	return env
}
