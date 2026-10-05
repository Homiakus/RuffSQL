package rules

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// Group: AL (Aliasing AL01 - AL10)

// RuleAL01: Explicit aliasing of table reference (AS keyword required)
func evalAL01(r *RuleAL01, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "from_expression_element" {
		children := ctx.Segment.GetSegments()
		hasTable := false
		hasAs := false
		hasAlias := false
		var aliasSeg parser.BaseSegment

		for _, ch := range children {
			if ch.GetType() == "table_expression" {
				hasTable = true
			}
			if ch.GetType() == "keyword" && strings.ToUpper(ch.GetRaw()) == "AS" {
				hasAs = true
			}
			if ch.GetType() == "naked_identifier" && hasTable {
				hasAlias = true
				aliasSeg = ch
			}
		}

		// Implicit alias without 'AS'
		if hasTable && hasAlias && !hasAs {
			pos := aliasSeg.GetPosition()
			asSeg := parser.NewRawSegment("keyword", "AS", pos)
			spSeg := parser.NewRawSegment(string(lexer.TokenWhitespace), " ", pos)
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Use explicit 'AS' keyword for table alias '%s'.", aliasSeg.GetRaw()),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     aliasSeg,
				Fixes: []LintFix{
					CreateBeforeFix(aliasSeg, asSeg, spSeg),
				},
			}}
		}
	}
	return nil
}

// RuleAL02: Explicit aliasing of column reference (AS keyword required)
func evalAL02(r *RuleAL02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_clause_element" {
		children := ctx.Segment.GetSegments()
		hasExpr := false
		hasAs := false
		hasAlias := false
		var aliasSeg parser.BaseSegment

		for _, ch := range children {
			if ch.GetType() == "expression" {
				hasExpr = true
			}
			if ch.GetType() == "keyword" && strings.ToUpper(ch.GetRaw()) == "AS" {
				hasAs = true
			}
			if ch.GetType() == "naked_identifier" && hasExpr {
				hasAlias = true
				aliasSeg = ch
			}
		}

		if hasExpr && hasAlias && !hasAs {
			pos := aliasSeg.GetPosition()
			asSeg := parser.NewRawSegment("keyword", "AS", pos)
			spSeg := parser.NewRawSegment(string(lexer.TokenWhitespace), " ", pos)
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Use explicit 'AS' keyword for column alias '%s'.", aliasSeg.GetRaw()),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     aliasSeg,
				Fixes: []LintFix{
					CreateBeforeFix(aliasSeg, asSeg, spSeg),
				},
			}}
		}
	}
	return nil
}

// RuleAL09: Column should not be self-aliased (e.g. col AS col)
func evalAL09(r *RuleAL09, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_clause_element" {
		var colName string
		var aliasName string
		var aliasSeg parser.BaseSegment

		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "expression" {
				for _, sub := range ch.GetSegments() {
					if sub.GetType() == "column_reference" {
						colName = strings.TrimSpace(sub.GetRaw())
					}
				}
			}
			if ch.GetType() == "naked_identifier" {
				aliasName = strings.TrimSpace(ch.GetRaw())
				aliasSeg = ch
			}
		}

		if colName != "" && aliasName != "" && strings.EqualFold(colName, aliasName) {
			pos := aliasSeg.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Column '%s' is self-aliased as '%s'.", colName, aliasName),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     aliasSeg,
			}}
		}
	}
	return nil
}

// Group: ST (Structure ST01 - ST12)

// RuleST01: ELSE NULL in CASE expressions is redundant
func evalST01(r *RuleST01, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "case_expression" {
		raw := strings.ToUpper(ctx.Segment.GetRaw())
		if strings.Contains(raw, "ELSE NULL") {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Redundant 'ELSE NULL' in CASE expression.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleST07: Prefer ON over USING in JOIN clauses
func evalST07(r *RuleST07, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "join_clause" {
		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "keyword" && strings.ToUpper(ch.GetRaw()) == "USING" {
				pos := ch.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Found 'USING' statement. Expected only 'ON' statements in JOIN.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ch,
				}}
			}
		}
	}
	return nil
}

// RuleST12: Consecutive semicolons should not be used
func evalST12(r *RuleST12, ctx RuleContext) []LintViolation {
	if ctx.Parent != nil && ctx.Segment.GetType() == "statement_terminator" {
		siblings := ctx.Parent.GetSegments()
		for i, s := range siblings {
			if s == ctx.Segment && i+1 < len(siblings) {
				if siblings[i+1].GetType() == "statement_terminator" {
					pos := siblings[i+1].GetPosition()
					return []LintViolation{{
						RuleCode:    r.Code(),
						Name:        r.Name(),
						Description: "Consecutive semicolons detected.",
						Line:        pos.Line,
						Pos:         pos.Pos,
						Segment:     siblings[i+1],
						Fixes:       []LintFix{DeleteFix(siblings[i+1])},
					}}
				}
			}
		}
	}
	return nil
}
