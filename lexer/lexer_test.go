package lexer

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `\documentclass{article}
\begin{document}
Hello $x + y = z$ world!
\textbf{Bold text}
% This is a comment
$$\sum_{i=1}^{n} i = \frac{n(n+1)}{2}$$
\end{document}`

	lexer := NewLexer(input)

	fmt.Println("Tokenizing LaTeX input:")
	fmt.Println("Input:", input)
	fmt.Println("\nTokens:")

	for {
		token := lexer.NextToken()
		fmt.Printf("Type: %-15s Value: %-20s Pos: Line %d, Col %d\n",
			tokenTypeToString(token.Type),
			fmt.Sprintf("'%s'", token.Value),
			token.Pos.Line,
			token.Pos.Column)

		if token.Type == TokenEOF {
			break
		}
	}
}

func tokenTypeToString(tokenType TokenType) string {
	switch tokenType {
	case TokenText:
		return "TEXT"
	case TokenCommand:
		return "COMMAND"
	case TokenBeginEnv:
		return "BEGIN_ENV"
	case TokenEndEnv:
		return "END_ENV"
	case TokenMathInline:
		return "MATH_INLINE"
	case TokenMathDisplay:
		return "MATH_DISPLAY"
	case TokenLBrace:
		return "LBRACE"
	case TokenRBrace:
		return "RBRACE"
	case TokenOptionalArg:
		return "OPTIONAL_ARG"
	case TokenComment:
		return "COMMENT"
	case TokenEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
