package lexer

// LexerInterface defines the interface for tokenizers
type LexerInterface interface {
	NextToken() Token
}

// NewLexer creates a lexer instance
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:   input,
		posInfo: Position{Line: 1, Column: 0},
	}
	l.readChar()
	return l
}

type Lexer struct {
	input   string
	pos     int  // current position
	readPos int  // next position
	ch      rune // current char
	posInfo Position
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = rune(l.input[l.readPos])

		// Track position
		if l.ch == '\n' {
			l.posInfo.Line++
			l.posInfo.Column = 0
		} else {
			l.posInfo.Column++
		}
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '\\':
		tok = l.lexCommand()
	case '$':
		tok = l.lexMath()
	case '{':
		tok = Token{Type: TokenLBrace, Value: "{", Pos: l.posInfo}
		l.readChar()
	case '}':
		tok = Token{Type: TokenRBrace, Value: "}", Pos: l.posInfo}
		l.readChar()
	case '[':
		tok = l.lexOptionalArg()
	case '%':
		tok = l.lexText() // lexText handles comments
	case 0:
		tok = Token{Type: TokenEOF, Value: "", Pos: l.posInfo}
	default:
		tok = l.lexText()
	}

	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}
func (l *Lexer) lexCommand() Token {
	pos := l.posInfo
	value := ""

	// Skip the backslash
	l.readChar()

	// Read the command name (letters only for standard commands)
	for isLetter(l.ch) {
		value += string(l.ch)
		l.readChar()
	}

	// Handle special cases like \begin and \end
	if value == "begin" || value == "end" {
		// Skip whitespace before {
		l.skipWhitespace()

		if l.ch == '{' {
			// Look ahead to see if there's a matching closing brace
			tempPos := l.readPos
			braceCount := 1
			hasClosingBrace := false

			// Scan ahead to check for matching brace
			for tempPos < len(l.input) && braceCount > 0 {
				ch := rune(l.input[tempPos])
				if ch == '{' {
					braceCount++
				} else if ch == '}' {
					braceCount--
					if braceCount == 0 {
						hasClosingBrace = true
						break
					}
				}
				tempPos++
			}

			if hasClosingBrace {
				// Matching brace found, proceed normally
				l.readChar() // skip {
				envName := ""
				for l.ch != '}' && l.ch != 0 {
					envName += string(l.ch)
					l.readChar()
				}
				if l.ch == '}' {
					l.readChar() // skip }

					if value == "begin" {
						return Token{Type: TokenBeginEnv, Value: envName, Pos: pos}
					} else {
						return Token{Type: TokenEndEnv, Value: envName, Pos: pos}
					}
				}
			}
			// If no matching brace found, treat as regular command
			// Don't consume the { - let it be handled as a separate token
		}
	}

	return Token{Type: TokenCommand, Value: value, Pos: pos}
}
func (l *Lexer) lexMath() Token {
	pos := l.posInfo

	// Check for display math ($$)
	if l.readPos < len(l.input) && rune(l.input[l.readPos]) == '$' {
		// Display math $$...$$
		l.readChar() // skip first $
		l.readChar() // skip second $

		value := ""
		for {
			if l.ch == 0 {
				break // EOF
			}
			if l.ch == '$' && l.readPos < len(l.input) && rune(l.input[l.readPos]) == '$' {
				l.readChar() // skip first $
				l.readChar() // skip second $
				break
			}
			value += string(l.ch)
			l.readChar()
		}

		return Token{Type: TokenMathDisplay, Value: value, Pos: pos}
	} else {
		// Inline math $...$
		l.readChar() // skip $

		value := ""
		for l.ch != '$' && l.ch != 0 {
			value += string(l.ch)
			l.readChar()
		}
		if l.ch == '$' {
			l.readChar() // skip closing $
		} else {
			// EOF reached without closing $, but we'll still return the token
			// The parser can detect this if needed by checking for incomplete math
		}

		return Token{Type: TokenMathInline, Value: value, Pos: pos}
	}
}
func (l *Lexer) lexText() Token {
	pos := l.posInfo
	value := ""

	// Read text until we hit a special character
	for l.ch != 0 && l.ch != '\\' && l.ch != '$' && l.ch != '{' && l.ch != '}' && l.ch != '[' && l.ch != ']' && l.ch != '%' && l.ch != '\n' {
		value += string(l.ch)
		l.readChar()
	}

	// If we didn't read any characters, handle single special characters or EOF
	if value == "" {
		switch l.ch {
		case 0:
			return Token{Type: TokenEOF, Value: "", Pos: pos}
		case '}':
			l.readChar()
			return Token{Type: TokenRBrace, Value: "}", Pos: pos}
		case '%':
			// Handle comments
			l.readChar() // skip %
			comment := ""
			for l.ch != '\n' && l.ch != 0 {
				comment += string(l.ch)
				l.readChar()
			}
			return Token{Type: TokenComment, Value: comment, Pos: pos}
		case '\n':
			l.readChar()
			return Token{Type: TokenText, Value: "\n", Pos: pos}
		}
	}

	return Token{Type: TokenText, Value: value, Pos: pos}
}
func (l *Lexer) lexOptionalArg() Token {
	pos := l.posInfo
	l.readChar() // skip [

	value := ""
	for l.ch != ']' && l.ch != 0 {
		value += string(l.ch)
		l.readChar()
	}
	if l.ch == ']' {
		l.readChar() // skip ]
	}

	return Token{Type: TokenOptionalArg, Value: value, Pos: pos}
}

// Tokenize returns all tokens from the input as a slice
func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}

	for {
		token := l.NextToken()
		tokens = append(tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}

	return tokens
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
