package linter

import (
	"testing"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

func BenchmarkLexer(b *testing.B) {
	sql := "SELECT user_id, user_name, email FROM users_table WHERE is_active = 1;"
	l := lexer.NewLexer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = l.Tokenize(sql)
	}
}

func BenchmarkParseAndLint(b *testing.B) {
	l, _ := NewLinter("ansi")
	sql := "SELECT user_id, user_name, email FROM users_table WHERE is_active = 1;"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = l.LintString(sql, "bench.sql")
	}
}
