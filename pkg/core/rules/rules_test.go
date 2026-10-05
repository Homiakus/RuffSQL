package rules

import (
	"testing"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

func TestExtractNoqa(t *testing.T) {
	// Construct a dummy tree with comment
	commentSeg := parser.NewRawSegment("comment", "-- noqa: CP01", lexer.Position{Line: 2, Pos: 1})
	fileSeg := parser.NewCompositeSegment("file", []parser.BaseSegment{commentSeg})

	noqaMap := ExtractNoqa(fileSeg)
	if noqaMap.DisableFile {
		t.Fatalf("did not expect disableFile")
	}

	vCP01 := LintViolation{RuleCode: "CP01", Line: 2}
	vCP02 := LintViolation{RuleCode: "CP02", Line: 2}
	vLine3 := LintViolation{RuleCode: "CP01", Line: 3}

	if !noqaMap.IsViolationIgnored(vCP01) {
		t.Errorf("expected CP01 on line 2 to be ignored")
	}
	if noqaMap.IsViolationIgnored(vCP02) {
		t.Errorf("expected CP02 on line 2 NOT to be ignored")
	}
	if noqaMap.IsViolationIgnored(vLine3) {
		t.Errorf("expected CP01 on line 3 NOT to be ignored")
	}
}

func TestNoqaDisableAll(t *testing.T) {
	commentSeg := parser.NewRawSegment("comment", "-- noqa", lexer.Position{Line: 5, Pos: 1})
	fileSeg := parser.NewCompositeSegment("file", []parser.BaseSegment{commentSeg})

	noqaMap := ExtractNoqa(fileSeg)
	v := LintViolation{RuleCode: "LT01", Line: 5}
	if !noqaMap.IsViolationIgnored(v) {
		t.Errorf("expected line 5 to have all rules ignored")
	}
}

func TestNoqaDisableFile(t *testing.T) {
	commentSeg := parser.NewRawSegment("comment", "-- noqa: disable=all", lexer.Position{Line: 1, Pos: 1})
	fileSeg := parser.NewCompositeSegment("file", []parser.BaseSegment{commentSeg})

	noqaMap := ExtractNoqa(fileSeg)
	if !noqaMap.DisableFile {
		t.Errorf("expected disableFile to be true")
	}
	v := LintViolation{RuleCode: "CP01", Line: 42}
	if !noqaMap.IsViolationIgnored(v) {
		t.Errorf("expected any violation in file to be ignored")
	}
}

func TestRuleCP01(t *testing.T) {
	r := &RuleCP01{}
	keywordSeg := parser.NewRawSegment("keyword", "select", lexer.Position{Line: 1, Pos: 1})
	ctx := RuleContext{Segment: keywordSeg}
	violations := r.Eval(ctx)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleCode != "CP01" {
		t.Errorf("expected CP01, got %s", violations[0].RuleCode)
	}
	if len(violations[0].Fixes) != 1 {
		t.Fatalf("expected 1 fix")
	}
}

func TestRuleCV01(t *testing.T) {
	r := &RuleCV01{}
	opSeg := parser.NewRawSegment("comparison_operator", "!=", lexer.Position{Line: 1, Pos: 10})
	ctx := RuleContext{Segment: opSeg}
	violations := r.Eval(ctx)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for !=")
	}
	if violations[0].RuleCode != "CV01" {
		t.Errorf("expected CV01, got %s", violations[0].RuleCode)
	}
}

func TestRuleCV02(t *testing.T) {
	r := &RuleCV02{}
	fnName := parser.NewRawSegment("naked_identifier", "nvl", lexer.Position{Line: 1, Pos: 1})
	fnSeg := parser.NewCompositeSegment("function", []parser.BaseSegment{fnName})
	ctx := RuleContext{Segment: fnSeg}
	violations := r.Eval(ctx)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for NVL")
	}
	if violations[0].RuleCode != "CV02" {
		t.Errorf("expected CV02, got %s", violations[0].RuleCode)
	}
}

func TestRuleST01(t *testing.T) {
	r := &RuleST01{}
	caseSeg := parser.NewRawSegment("case_expression", "CASE WHEN a THEN 1 ELSE NULL END", lexer.Position{Line: 1, Pos: 1})
	ctx := RuleContext{Segment: caseSeg}
	violations := r.Eval(ctx)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for ELSE NULL")
	}
	if violations[0].RuleCode != "ST01" {
		t.Errorf("expected ST01, got %s", violations[0].RuleCode)
	}
}
