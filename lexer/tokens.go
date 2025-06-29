package lexer

type TokenType int

const (
	TokenText        TokenType = iota
	TokenCommand               // \command
	TokenBeginEnv              // \begin{env}
	TokenEndEnv                // \end{env}
	TokenMathInline            // $...$
	TokenMathDisplay           // $$...$$
	TokenLBrace                // {
	TokenRBrace                // }
	TokenOptionalArg           // [...]
	TokenComment               // %...
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
	Pos   Position // Line, Column information
}

type Position struct {
	Line   int
	Column int
}
