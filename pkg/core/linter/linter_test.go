package linter

import (
	"testing"
)

func TestLinterRules(t *testing.T) {
	l, err := NewLinter("ansi")
	if err != nil {
		t.Fatalf("linter init error: %v", err)
	}

	badSQL := "select a, b FROM t WHERE id = 1;"
	res, err := l.LintString(badSQL, "inline.sql")
	if err != nil {
		t.Fatalf("lint error: %v", err)
	}

	ruleMap := make(map[string]bool)
	for _, v := range res.Violations {
		ruleMap[v.RuleCode] = true
	}

	// Should violate CP01 (lowercase 'select')
	if !ruleMap["CP01"] {
		t.Errorf("expected CP01 violation for lowercase 'select'")
	}
	// Should violate LT09 (multiple targets on one line)
	if !ruleMap["LT09"] {
		t.Errorf("expected LT09 violation for multiple targets on one line")
	}
	// Should violate LT14 (WHERE not on newline)
	if !ruleMap["LT14"] {
		t.Errorf("expected LT14 violation for WHERE not on newline")
	}
}
