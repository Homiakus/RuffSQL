package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents classification of lexical elements
type TokenType string

const (
	TokenWhitespace         TokenType = "whitespace"
	TokenNewline            TokenType = "newline"
	TokenKeyword            TokenType = "keyword"
	TokenIdentifier         TokenType = "naked_identifier"
	TokenQuotedIdentifier   TokenType = "quoted_identifier"
	TokenNumericLiteral     TokenType = "numeric_literal"
	TokenStringLiteral      TokenType = "single_quote"
	TokenComparisonOperator TokenType = "comparison_operator"
	TokenBinaryOperator     TokenType = "binary_operator"
	TokenComma              TokenType = "comma"
	TokenDot                TokenType = "dot"
	TokenSemicolon          TokenType = "statement_terminator"
	TokenOpenBracket        TokenType = "start_bracket"
	TokenCloseBracket       TokenType = "end_bracket"
	TokenComment            TokenType = "comment"
	TokenMetaIndent         TokenType = "indent"
	TokenMetaDedent         TokenType = "dedent"
	TokenMetaEndOfFile      TokenType = "end_of_file"
)

// Position tracks 1-indexed line and column
type Position struct {
	Line int `json:"line"`
	Pos  int `json:"pos"`
}

// Token represents a single lexical unit
type Token struct {
	Type     TokenType `json:"type"`
	Raw      string    `json:"raw"`
	Position Position  `json:"position"`
}

func (t Token) String() string {
	return fmt.Sprintf("[L:%3d, P:%3d] %-25s: '%s'", t.Position.Line, t.Position.Pos, t.Type, t.Raw)
}

// Lexer breaks raw SQL string into tokens
type Lexer struct{}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Tokenize(raw string) ([]Token, error) {
	var tokens []Token
	runes := []rune(raw)
	n := len(runes)

	line := 1
	pos := 1
	i := 0

	for i < n {
		startLine := line
		startPos := pos
		r := runes[i]

		// 1. Newline
		if r == '\n' {
			tokens = append(tokens, Token{
				Type:     TokenNewline,
				Raw:      "\n",
				Position: Position{Line: startLine, Pos: startPos},
			})
			line++
			pos = 1
			i++
			continue
		} else if r == '\r' {
			rawStr := "\r"
			if i+1 < n && runes[i+1] == '\n' {
				rawStr = "\r\n"
				i++
			}
			tokens = append(tokens, Token{
				Type:     TokenNewline,
				Raw:      rawStr,
				Position: Position{Line: startLine, Pos: startPos},
			})
			line++
			pos = 1
			i++
			continue
		}

		// 2. Whitespace
		if unicode.IsSpace(r) {
			start := i
			for i < n && unicode.IsSpace(runes[i]) && runes[i] != '\n' && runes[i] != '\r' {
				i++
				pos++
			}
			tokens = append(tokens, Token{
				Type:     TokenWhitespace,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// 3. Comments (-- single line or /* */ multi-line)
		if r == '-' && i+1 < n && runes[i+1] == '-' {
			start := i
			for i < n && runes[i] != '\n' && runes[i] != '\r' {
				i++
				pos++
			}
			tokens = append(tokens, Token{
				Type:     TokenComment,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}
		if r == '/' && i+1 < n && runes[i+1] == '*' {
			start := i
			i += 2
			pos += 2
			for i+1 < n && !(runes[i] == '*' && runes[i+1] == '/') {
				if runes[i] == '\n' {
					line++
					pos = 1
				} else {
					pos++
				}
				i++
			}
			if i+1 < n {
				i += 2
				pos += 2
			}
			tokens = append(tokens, Token{
				Type:     TokenComment,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// 4. Strings: '...'
		if r == '\'' {
			start := i
			i++
			pos++
			for i < n {
				if runes[i] == '\'' {
					if i+1 < n && runes[i+1] == '\'' {
						// escaped quote
						i += 2
						pos += 2
						continue
					}
					i++
					pos++
					break
				}
				if runes[i] == '\n' {
					line++
					pos = 1
				} else {
					pos++
				}
				i++
			}
			tokens = append(tokens, Token{
				Type:     TokenStringLiteral,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// 5. Quoted Identifiers: "..." or `...`
		if r == '"' || r == '`' {
			quote := r
			start := i
			i++
			pos++
			for i < n && runes[i] != quote {
				if runes[i] == '\n' {
					line++
					pos = 1
				} else {
					pos++
				}
				i++
			}
			if i < n && runes[i] == quote {
				i++
				pos++
			}
			tokens = append(tokens, Token{
				Type:     TokenQuotedIdentifier,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// 6. Numeric literals
		if unicode.IsDigit(r) {
			start := i
			hasDot := false
			for i < n && (unicode.IsDigit(runes[i]) || (!hasDot && runes[i] == '.' && i+1 < n && unicode.IsDigit(runes[i+1]))) {
				if runes[i] == '.' {
					hasDot = true
				}
				i++
				pos++
			}
			tokens = append(tokens, Token{
				Type:     TokenNumericLiteral,
				Raw:      string(runes[start:i]),
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// 7. Symbols & Operators
		// Postgres :: and -> / ->>
		if r == ':' && i+1 < n && runes[i+1] == ':' {
			tokens = append(tokens, Token{Type: TokenBinaryOperator, Raw: "::", Position: Position{Line: startLine, Pos: startPos}})
			i += 2
			pos += 2
			continue
		}
		if r == '-' && i+1 < n && runes[i+1] == '>' {
			if i+2 < n && runes[i+2] == '>' {
				tokens = append(tokens, Token{Type: TokenBinaryOperator, Raw: "->>", Position: Position{Line: startLine, Pos: startPos}})
				i += 3
				pos += 3
				continue
			}
			tokens = append(tokens, Token{Type: TokenBinaryOperator, Raw: "->", Position: Position{Line: startLine, Pos: startPos}})
			i += 2
			pos += 2
			continue
		}
		if r == ',' {
			tokens = append(tokens, Token{Type: TokenComma, Raw: ",", Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}
		if r == ';' {
			tokens = append(tokens, Token{Type: TokenSemicolon, Raw: ";", Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}
		if r == '.' {
			tokens = append(tokens, Token{Type: TokenDot, Raw: ".", Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}
		if r == '(' {
			tokens = append(tokens, Token{Type: TokenOpenBracket, Raw: "(", Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}
		if r == ')' {
			tokens = append(tokens, Token{Type: TokenCloseBracket, Raw: ")", Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}

		// Multi-character and comparison operators
		if (r == '<' || r == '>' || r == '=' || r == '!') && i+1 < n {
			next := runes[i+1]
			if (r == '<' && (next == '=' || next == '>')) ||
				(r == '>' && next == '=') ||
				(r == '!' && next == '=') {
				op := string(runes[i : i+2])
				tokens = append(tokens, Token{Type: TokenComparisonOperator, Raw: op, Position: Position{Line: startLine, Pos: startPos}})
				i += 2
				pos += 2
				continue
			}
		}
		if r == '=' || r == '<' || r == '>' {
			tokens = append(tokens, Token{Type: TokenComparisonOperator, Raw: string(r), Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}
		if r == '+' || r == '-' || r == '*' || r == '/' || r == '%' || r == '|' {
			if r == '|' && i+1 < n && runes[i+1] == '|' {
				tokens = append(tokens, Token{Type: TokenBinaryOperator, Raw: "||", Position: Position{Line: startLine, Pos: startPos}})
				i += 2
				pos += 2
				continue
			}
			tokens = append(tokens, Token{Type: TokenBinaryOperator, Raw: string(r), Position: Position{Line: startLine, Pos: startPos}})
			i++
			pos++
			continue
		}

		// 8. Word tokens (Keywords or Naked Identifiers)
		if unicode.IsLetter(r) || r == '_' {
			start := i
			for i < n && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
				pos++
			}
			word := string(runes[start:i])
			tokenType := TokenIdentifier
			if IsStandardKeyword(word) {
				tokenType = TokenKeyword
			}
			tokens = append(tokens, Token{
				Type:     tokenType,
				Raw:      word,
				Position: Position{Line: startLine, Pos: startPos},
			})
			continue
		}

		// Default fallback for any other char
		tokens = append(tokens, Token{
			Type:     TokenType(fmt.Sprintf("char_%c", r)),
			Raw:      string(r),
			Position: Position{Line: startLine, Pos: startPos},
		})
		i++
		pos++
	}

	return tokens, nil
}

// Standard SQL keywords
var standardKeywords = map[string]struct{}{
	"SELECT": {}, "FROM": {}, "WHERE": {}, "JOIN": {}, "INNER": {}, "LEFT": {}, "RIGHT": {}, "FULL": {},
	"CROSS": {}, "ON": {}, "GROUP": {}, "BY": {}, "HAVING": {}, "ORDER": {}, "LIMIT": {}, "OFFSET": {},
	"INSERT": {}, "INTO": {}, "VALUES": {}, "UPDATE": {}, "SET": {}, "DELETE": {}, "CREATE": {}, "TABLE": {},
	"DROP": {}, "ALTER": {}, "AND": {}, "OR": {}, "NOT": {}, "IN": {}, "IS": {}, "NULL": {}, "AS": {},
	"DISTINCT": {}, "ALL": {}, "UNION": {}, "CASE": {}, "WHEN": {}, "THEN": {}, "ELSE": {}, "END": {},
	"LIKE": {}, "BETWEEN": {}, "EXISTS": {}, "ASC": {}, "DESC": {}, "WITH": {}, "OVER": {}, "PARTITION": {}, "CONFLICT": {}, "DO": {}, "NOTHING": {}, "DUPLICATE": {}, "KEY": {}, "QUALIFY": {}, "PROCEDURE": {}, "FUNCTION": {}, "TRIGGER": {}, "CALL": {}, "EXECUTE": {}, "EXEC": {}, "DECLARE": {}, "RETURN": {}, "RETURNS": {}, "BEGIN": {}, "LANGUAGE": {}, "BEFORE": {}, "AFTER": {}, "INSTEAD": {}, "OF": {}, "EACH": {}, "ROW": {}, "INOUT": {},
}

func IsStandardKeyword(word string) bool {
	_, ok := standardKeywords[strings.ToUpper(word)]
	return ok
}
