package lexer

import (
	"testing"
)

func TestLexerBasic(t *testing.T) {
	sql := "SELECT col1, col2 FROM tbl WHERE id = 123;"
	l := NewLexer()
	tokens, err := l.Tokenize(sql)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatalf("expected tokens, got none")
	}

	expectedTypes := []TokenType{
		TokenKeyword,
		TokenWhitespace,
		TokenIdentifier,
		TokenComma,
		TokenWhitespace,
		TokenIdentifier,
		TokenWhitespace,
		TokenKeyword,
		TokenWhitespace,
		TokenIdentifier,
		TokenWhitespace,
		TokenKeyword,
		TokenWhitespace,
		TokenIdentifier,
		TokenWhitespace,
		TokenComparisonOperator,
		TokenWhitespace,
		TokenNumericLiteral,
		TokenSemicolon,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, exp := range expectedTypes {
		if tokens[i].Type != exp {
			t.Errorf("token %d: expected %s, got %s (%s)", i, exp, tokens[i].Type, tokens[i].Raw)
		}
	}
}
