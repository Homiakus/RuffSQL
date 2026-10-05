package linter

import (
	"strings"
	"testing"
)

func TestComprehensiveLinterRules(t *testing.T) {
	l, err := NewLinter("ansi")
	if err != nil {
		t.Fatalf("linter init error: %v", err)
	}

	tests := []struct {
		name         string
		sql          string
		expectedRule string
	}{
		{
			name:         "Rule CP01: lowercase keyword",
			sql:          "select 1;",
			expectedRule: "CP01",
		},
		{
			name:         "Rule CP03: lowercase function",
			sql:          "SELECT count(1) FROM t;",
			expectedRule: "CP03",
		},
		{
			name:         "Rule CP04: lowercase null",
			sql:          "SELECT null FROM t;",
			expectedRule: "CP04",
		},
		{
			name:         "Rule CV01: not equal !=",
			sql:          "SELECT a FROM t WHERE a != 1;",
			expectedRule: "CV01",
		},
		{
			name:         "Rule CV02: nvl or ifnull",
			sql:          "SELECT NVL(a, 0) FROM t;",
			expectedRule: "CV02",
		},
		{
			name:         "Rule CV04: count(0)",
			sql:          "SELECT COUNT(0) FROM t;",
			expectedRule: "CV04",
		},
		{
			name:         "Rule CV05: comparison with NULL",
			sql:          "SELECT a FROM t WHERE a = NULL;",
			expectedRule: "CV05",
		},
		{
			name:         "Rule CV06: missing statement terminator",
			sql:          "SELECT a FROM t",
			expectedRule: "CV06",
		},
		{
			name:         "Rule CV08: RIGHT JOIN instead of LEFT JOIN",
			sql:          "SELECT a FROM t RIGHT JOIN u ON t.id = u.id;",
			expectedRule: "CV08",
		},
		{
			name:         "Rule CV10: double quotes for string literals",
			sql:          `SELECT "hello" FROM t;`,
			expectedRule: "CV10",
		},
		{
			name:         "Rule CV11: convert instead of cast",
			sql:          "SELECT CONVERT(a) FROM t;",
			expectedRule: "CV11",
		},
		{
			name:         "Rule ST01: redundant ELSE NULL in CASE",
			sql:          "SELECT CASE WHEN a = 1 THEN 1 ELSE NULL END FROM t;",
			expectedRule: "ST01",
		},
		{
			name:         "Rule ST07: USING instead of ON in JOIN",
			sql:          "SELECT a FROM t JOIN u USING (id);",
			expectedRule: "ST07",
		},
		{
			name:         "Rule ST12: consecutive semicolons",
			sql:          "SELECT 1;;",
			expectedRule: "ST12",
		},
		{
			name:         "Rule AM01: bare distinct with parentheses",
			sql:          "SELECT DISTINCT(a) FROM t;",
			expectedRule: "AM01",
		},
		{
			name:         "Rule AM02: bare UNION without ALL or DISTINCT",
			sql:          "SELECT a FROM t UNION SELECT b FROM u;",
			expectedRule: "AM02",
		},
		{
			name:         "Rule TQ01: unnecessary quotes on identifier",
			sql:          "SELECT `my_col` FROM t;",
			expectedRule: "TQ01",
		},
		{
			name:         "Rule TQ02: duplicate redundant parentheses",
			sql:          "SELECT ((1 + 1)) FROM t;",
			expectedRule: "TQ02",
		},
		{
			name:         "Rule LT01: trailing whitespace",
			sql:          "SELECT 1 \nFROM t;",
			expectedRule: "LT01",
		},
		{
			name:         "Rule LT02: tab indentation",
			sql:          "SELECT 1\n\tFROM t;",
			expectedRule: "LT02",
		},
		{
			name:         "Rule LT13: starts with newline",
			sql:          "\nSELECT 1 FROM t;",
			expectedRule: "LT13",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := l.LintString(tc.sql, "test.sql")
			if err != nil {
				t.Fatalf("unexpected lint error: %v", err)
			}
			found := false
			for _, v := range res.Violations {
				if v.RuleCode == tc.expectedRule {
					found = true
					break
				}
			}
			if !found {
				var codes []string
				for _, v := range res.Violations {
					codes = append(codes, v.RuleCode)
				}
				t.Errorf("expected violation %s, got: %v", tc.expectedRule, codes)
			}
		})
	}
}

func TestLinterFixString(t *testing.T) {
	l, err := NewLinter("ansi")
	if err != nil {
		t.Fatalf("linter init error: %v", err)
	}

	badSQL := "select a, b FROM t WHERE a != 1;"
	fixedSQL, count, err := l.FixString(badSQL, 5)
	if err != nil {
		t.Fatalf("fix error: %v", err)
	}

	if count == 0 {
		t.Errorf("expected fixes to be applied")
	}

	// Keywords should now be uppercase
	if !strings.Contains(fixedSQL, "SELECT") {
		t.Errorf("expected 'select' to be fixed to 'SELECT', got: %s", fixedSQL)
	}
}
