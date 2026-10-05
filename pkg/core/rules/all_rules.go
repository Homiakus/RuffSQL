package rules

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)
// RuleAL01: aliasing.table
type RuleAL01 struct{}

func (r *RuleAL01) Code() string        { return "AL01" }
func (r *RuleAL01) Name() string        { return "aliasing.table" }
func (r *RuleAL01) Description() string { return "Rule AL01" }
func (r *RuleAL01) Eval(ctx RuleContext) []LintViolation {
	return evalAL01(r, ctx)
}

// RuleAL02: aliasing.column
type RuleAL02 struct{}

func (r *RuleAL02) Code() string        { return "AL02" }
func (r *RuleAL02) Name() string        { return "aliasing.column" }
func (r *RuleAL02) Description() string { return "Rule AL02" }
func (r *RuleAL02) Eval(ctx RuleContext) []LintViolation {
	return evalAL02(r, ctx)
}

// RuleAL03: aliasing.expression
type RuleAL03 struct{}

func (r *RuleAL03) Code() string        { return "AL03" }
func (r *RuleAL03) Name() string        { return "aliasing.expression" }
func (r *RuleAL03) Description() string { return "Rule AL03" }
func (r *RuleAL03) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL04: aliasing.unique.table
type RuleAL04 struct{}

func (r *RuleAL04) Code() string        { return "AL04" }
func (r *RuleAL04) Name() string        { return "aliasing.unique.table" }
func (r *RuleAL04) Description() string { return "Rule AL04" }
func (r *RuleAL04) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL05: aliasing.unused
type RuleAL05 struct{}

func (r *RuleAL05) Code() string        { return "AL05" }
func (r *RuleAL05) Name() string        { return "aliasing.unused" }
func (r *RuleAL05) Description() string { return "Rule AL05" }
func (r *RuleAL05) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL06: aliasing.length
type RuleAL06 struct{}

func (r *RuleAL06) Code() string        { return "AL06" }
func (r *RuleAL06) Name() string        { return "aliasing.length" }
func (r *RuleAL06) Description() string { return "Rule AL06" }
func (r *RuleAL06) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL07: AL07
type RuleAL07 struct{}

func (r *RuleAL07) Code() string        { return "AL07" }
func (r *RuleAL07) Name() string        { return "AL07" }
func (r *RuleAL07) Description() string { return "Rule AL07" }
func (r *RuleAL07) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL08: aliasing.unique.column
type RuleAL08 struct{}

func (r *RuleAL08) Code() string        { return "AL08" }
func (r *RuleAL08) Name() string        { return "aliasing.unique.column" }
func (r *RuleAL08) Description() string { return "Rule AL08" }
func (r *RuleAL08) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAL09: AL09
type RuleAL09 struct{}

func (r *RuleAL09) Code() string        { return "AL09" }
func (r *RuleAL09) Name() string        { return "AL09" }
func (r *RuleAL09) Description() string { return "Rule AL09" }
func (r *RuleAL09) Eval(ctx RuleContext) []LintViolation {
	return evalAL09(r, ctx)
}

// RuleAL10: aliasing.required
type RuleAL10 struct{}

func (r *RuleAL10) Code() string        { return "AL10" }
func (r *RuleAL10) Name() string        { return "aliasing.required" }
func (r *RuleAL10) Description() string { return "Derived table must have an alias." }
func (r *RuleAL10) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAM01: ambiguous.distinct
type RuleAM01 struct{}

func (r *RuleAM01) Code() string        { return "AM01" }
func (r *RuleAM01) Name() string        { return "ambiguous.distinct" }
func (r *RuleAM01) Description() string { return "Rule AM01" }
func (r *RuleAM01) Eval(ctx RuleContext) []LintViolation {
	return evalAM01(r, ctx)
}

// RuleAM02: ambiguous.union
type RuleAM02 struct{}

func (r *RuleAM02) Code() string        { return "AM02" }
func (r *RuleAM02) Name() string        { return "ambiguous.union" }
func (r *RuleAM02) Description() string { return "Rule AM02" }
func (r *RuleAM02) Eval(ctx RuleContext) []LintViolation {
	return evalAM02(r, ctx)
}

// RuleAM03: ambiguous.order_by
type RuleAM03 struct{}

func (r *RuleAM03) Code() string        { return "AM03" }
func (r *RuleAM03) Name() string        { return "ambiguous.order_by" }
func (r *RuleAM03) Description() string { return "Rule AM03" }
func (r *RuleAM03) Eval(ctx RuleContext) []LintViolation {
	return evalAM03(r, ctx)
}

// RuleAM04: ambiguous.column_count
type RuleAM04 struct{}

func (r *RuleAM04) Code() string        { return "AM04" }
func (r *RuleAM04) Name() string        { return "ambiguous.column_count" }
func (r *RuleAM04) Description() string { return "Rule AM04" }
func (r *RuleAM04) Eval(ctx RuleContext) []LintViolation {
	return evalAM04(r, ctx)
}

// RuleAM05: ambiguous.join
type RuleAM05 struct{}

func (r *RuleAM05) Code() string        { return "AM05" }
func (r *RuleAM05) Name() string        { return "ambiguous.join" }
func (r *RuleAM05) Description() string { return "Rule AM05" }
func (r *RuleAM05) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAM06: ambiguous.column_references
type RuleAM06 struct{}

func (r *RuleAM06) Code() string        { return "AM06" }
func (r *RuleAM06) Name() string        { return "ambiguous.column_references" }
func (r *RuleAM06) Description() string { return "Rule AM06" }
func (r *RuleAM06) Eval(ctx RuleContext) []LintViolation {
	return evalAM06(r, ctx)
}

// RuleAM07: ambiguous.set_columns
type RuleAM07 struct{}

func (r *RuleAM07) Code() string        { return "AM07" }
func (r *RuleAM07) Name() string        { return "ambiguous.set_columns" }
func (r *RuleAM07) Description() string { return "Rule AM07" }
func (r *RuleAM07) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAM08: ambiguous.join_condition
type RuleAM08 struct{}

func (r *RuleAM08) Code() string        { return "AM08" }
func (r *RuleAM08) Name() string        { return "ambiguous.join_condition" }
func (r *RuleAM08) Description() string { return "Rule AM08" }
func (r *RuleAM08) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleAM09: ambiguous.order_by_limit
type RuleAM09 struct{}

func (r *RuleAM09) Code() string        { return "AM09" }
func (r *RuleAM09) Name() string        { return "ambiguous.order_by_limit" }
func (r *RuleAM09) Description() string { return "Rule AM09" }
func (r *RuleAM09) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleCP01: capitalisation.keywords
type RuleCP01 struct{}

func (r *RuleCP01) Code() string        { return "CP01" }
func (r *RuleCP01) Name() string        { return "capitalisation.keywords" }
func (r *RuleCP01) Description() string { return "Rule CP01" }
func (r *RuleCP01) Eval(ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "keyword" {
		raw := ctx.Segment.GetRaw()
		if raw != strings.ToUpper(raw) {
			pos := ctx.Segment.GetPosition()
			fixedRaw := parser.NewRawSegment("keyword", strings.ToUpper(raw), pos)
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Keyword '%s' must be uppercase ('%s').", raw, strings.ToUpper(raw)),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
				Fixes:       []LintFix{ReplaceFix(ctx.Segment, fixedRaw)},
			}}
		}
	}
	return nil
}

// RuleCP02: CP02
type RuleCP02 struct{}

func (r *RuleCP02) Code() string        { return "CP02" }
func (r *RuleCP02) Name() string        { return "CP02" }
func (r *RuleCP02) Description() string { return "Rule CP02" }
func (r *RuleCP02) Eval(ctx RuleContext) []LintViolation {
	return evalCP02(r, ctx)
}

// RuleCP03: capitalisation.functions
type RuleCP03 struct{}

func (r *RuleCP03) Code() string        { return "CP03" }
func (r *RuleCP03) Name() string        { return "capitalisation.functions" }
func (r *RuleCP03) Description() string { return "Rule CP03" }
func (r *RuleCP03) Eval(ctx RuleContext) []LintViolation {
	return evalCP03(r, ctx)
}

// RuleCP04: capitalisation.literals
type RuleCP04 struct{}

func (r *RuleCP04) Code() string        { return "CP04" }
func (r *RuleCP04) Name() string        { return "capitalisation.literals" }
func (r *RuleCP04) Description() string { return "Rule CP04" }
func (r *RuleCP04) Eval(ctx RuleContext) []LintViolation {
	return evalCP04(r, ctx)
}

// RuleCP05: capitalisation.types
type RuleCP05 struct{}

func (r *RuleCP05) Code() string        { return "CP05" }
func (r *RuleCP05) Name() string        { return "capitalisation.types" }
func (r *RuleCP05) Description() string { return "Rule CP05" }
func (r *RuleCP05) Eval(ctx RuleContext) []LintViolation {
	return evalCP05(r, ctx)
}

// RuleCV01: convention.not_equal
type RuleCV01 struct{}

func (r *RuleCV01) Code() string        { return "CV01" }
func (r *RuleCV01) Name() string        { return "convention.not_equal" }
func (r *RuleCV01) Description() string { return "Rule CV01" }
func (r *RuleCV01) Eval(ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "comparison_operator" && ctx.Segment.GetRaw() == "!=" {
		pos := ctx.Segment.GetPosition()
		return []LintViolation{{
			RuleCode:    r.Code(),
			Name:        r.Name(),
			Description: "Use '<>' instead of '!='.",
			Line:        pos.Line,
			Pos:         pos.Pos,
			Segment:     ctx.Segment,
			Fixes: []LintFix{
				ReplaceFix(ctx.Segment, parser.NewRawSegment("comparison_operator", "<>", pos)),
			},
		}}
	}
	return nil
}

// RuleCV02: convention.coalesce
type RuleCV02 struct{}

func (r *RuleCV02) Code() string        { return "CV02" }
func (r *RuleCV02) Name() string        { return "convention.coalesce" }
func (r *RuleCV02) Description() string { return "Rule CV02" }
func (r *RuleCV02) Eval(ctx RuleContext) []LintViolation {
	return evalCV02(r, ctx)
}

// RuleCV03: convention.select_trailing_comma
type RuleCV03 struct{}

func (r *RuleCV03) Code() string        { return "CV03" }
func (r *RuleCV03) Name() string        { return "convention.select_trailing_comma" }
func (r *RuleCV03) Description() string { return "Rule CV03" }
func (r *RuleCV03) Eval(ctx RuleContext) []LintViolation {
	return evalCV03(r, ctx)
}

// RuleCV04: convention.count_rows
type RuleCV04 struct{}

func (r *RuleCV04) Code() string        { return "CV04" }
func (r *RuleCV04) Name() string        { return "convention.count_rows" }
func (r *RuleCV04) Description() string { return "Rule CV04" }
func (r *RuleCV04) Eval(ctx RuleContext) []LintViolation {
	return evalCV04(r, ctx)
}

// RuleCV05: convention.is_null
type RuleCV05 struct{}

func (r *RuleCV05) Code() string        { return "CV05" }
func (r *RuleCV05) Name() string        { return "convention.is_null" }
func (r *RuleCV05) Description() string { return "Rule CV05" }
func (r *RuleCV05) Eval(ctx RuleContext) []LintViolation {
	return evalCV05(r, ctx)
}

// RuleCV06: convention.terminator
type RuleCV06 struct{}

func (r *RuleCV06) Code() string        { return "CV06" }
func (r *RuleCV06) Name() string        { return "convention.terminator" }
func (r *RuleCV06) Description() string { return "Rule CV06" }
func (r *RuleCV06) Eval(ctx RuleContext) []LintViolation {
	return evalCV06(r, ctx)
}

// RuleCV07: convention.statement_brackets
type RuleCV07 struct{}

func (r *RuleCV07) Code() string        { return "CV07" }
func (r *RuleCV07) Name() string        { return "convention.statement_brackets" }
func (r *RuleCV07) Description() string { return "Rule CV07" }
func (r *RuleCV07) Eval(ctx RuleContext) []LintViolation {
	return evalCV07(r, ctx)
}

// RuleCV08: convention.left_join
type RuleCV08 struct{}

func (r *RuleCV08) Code() string        { return "CV08" }
func (r *RuleCV08) Name() string        { return "convention.left_join" }
func (r *RuleCV08) Description() string { return "Rule CV08" }
func (r *RuleCV08) Eval(ctx RuleContext) []LintViolation {
	return evalCV08(r, ctx)
}

// RuleCV09: convention.blocked_words
type RuleCV09 struct{}

func (r *RuleCV09) Code() string        { return "CV09" }
func (r *RuleCV09) Name() string        { return "convention.blocked_words" }
func (r *RuleCV09) Description() string { return "Rule CV09" }
func (r *RuleCV09) Eval(ctx RuleContext) []LintViolation {
	return evalCV09(r, ctx)
}

// RuleCV10: convention.quoted_literals
type RuleCV10 struct{}

func (r *RuleCV10) Code() string        { return "CV10" }
func (r *RuleCV10) Name() string        { return "convention.quoted_literals" }
func (r *RuleCV10) Description() string { return "Rule CV10" }
func (r *RuleCV10) Eval(ctx RuleContext) []LintViolation {
	return evalCV10(r, ctx)
}

// RuleCV11: convention.casting_style
type RuleCV11 struct{}

func (r *RuleCV11) Code() string        { return "CV11" }
func (r *RuleCV11) Name() string        { return "convention.casting_style" }
func (r *RuleCV11) Description() string { return "Rule CV11" }
func (r *RuleCV11) Eval(ctx RuleContext) []LintViolation {
	return evalCV11(r, ctx)
}

// RuleCV12: convention.join_condition
type RuleCV12 struct{}

func (r *RuleCV12) Code() string        { return "CV12" }
func (r *RuleCV12) Name() string        { return "convention.join_condition" }
func (r *RuleCV12) Description() string { return "Rule CV12" }
func (r *RuleCV12) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleCV13: CV13
type RuleCV13 struct{}

func (r *RuleCV13) Code() string        { return "CV13" }
func (r *RuleCV13) Name() string        { return "CV13" }
func (r *RuleCV13) Description() string { return "Rule CV13" }
func (r *RuleCV13) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleJJ01: jinja.padding
type RuleJJ01 struct{}

func (r *RuleJJ01) Code() string        { return "JJ01" }
func (r *RuleJJ01) Name() string        { return "jinja.padding" }
func (r *RuleJJ01) Description() string { return "Rule JJ01" }
func (r *RuleJJ01) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT01: layout.spacing
type RuleLT01 struct{}

func (r *RuleLT01) Code() string        { return "LT01" }
func (r *RuleLT01) Name() string        { return "layout.spacing" }
func (r *RuleLT01) Description() string { return "Rule LT01" }
func (r *RuleLT01) Eval(ctx RuleContext) []LintViolation {
	if ctx.Parent != nil && ctx.Segment.GetType() == string(lexer.TokenWhitespace) {
		siblings := ctx.Parent.GetSegments()
		for i, s := range siblings {
			if s == ctx.Segment && i+1 < len(siblings) {
				if siblings[i+1].GetType() == string(lexer.TokenNewline) {
					pos := ctx.Segment.GetPosition()
					return []LintViolation{{
						RuleCode:    r.Code(),
						Name:        r.Name(),
						Description: "Trailing whitespace detected before newline.",
						Line:        pos.Line,
						Pos:         pos.Pos,
						Segment:     ctx.Segment,
						Fixes:       []LintFix{DeleteFix(ctx.Segment)},
					}}
				}
			}
		}
	}
	return nil
}

// RuleLT02: layout.indent
type RuleLT02 struct{}

func (r *RuleLT02) Code() string        { return "LT02" }
func (r *RuleLT02) Name() string        { return "layout.indent" }
func (r *RuleLT02) Description() string { return "Rule LT02" }
func (r *RuleLT02) Eval(ctx RuleContext) []LintViolation {
	return evalLT02(r, ctx)
}

// RuleLT03: layout.operators
type RuleLT03 struct{}

func (r *RuleLT03) Code() string        { return "LT03" }
func (r *RuleLT03) Name() string        { return "layout.operators" }
func (r *RuleLT03) Description() string { return "Rule LT03" }
func (r *RuleLT03) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT04: layout.commas
type RuleLT04 struct{}

func (r *RuleLT04) Code() string        { return "LT04" }
func (r *RuleLT04) Name() string        { return "layout.commas" }
func (r *RuleLT04) Description() string { return "Rule LT04" }
func (r *RuleLT04) Eval(ctx RuleContext) []LintViolation {
	return evalLT04(r, ctx)
}

// RuleLT05: layout.long_lines
type RuleLT05 struct{}

func (r *RuleLT05) Code() string        { return "LT05" }
func (r *RuleLT05) Name() string        { return "layout.long_lines" }
func (r *RuleLT05) Description() string { return "Rule LT05" }
func (r *RuleLT05) Eval(ctx RuleContext) []LintViolation {
	return evalLT05(r, ctx)
}

// RuleLT06: layout.functions
type RuleLT06 struct{}

func (r *RuleLT06) Code() string        { return "LT06" }
func (r *RuleLT06) Name() string        { return "layout.functions" }
func (r *RuleLT06) Description() string { return "Rule LT06" }
func (r *RuleLT06) Eval(ctx RuleContext) []LintViolation {
	return evalLT06(r, ctx)
}

// RuleLT07: layout.cte_bracket
type RuleLT07 struct{}

func (r *RuleLT07) Code() string        { return "LT07" }
func (r *RuleLT07) Name() string        { return "layout.cte_bracket" }
func (r *RuleLT07) Description() string { return "Rule LT07" }
func (r *RuleLT07) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT08: layout.cte_newline
type RuleLT08 struct{}

func (r *RuleLT08) Code() string        { return "LT08" }
func (r *RuleLT08) Name() string        { return "layout.cte_newline" }
func (r *RuleLT08) Description() string { return "Rule LT08" }
func (r *RuleLT08) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT09: LT09
type RuleLT09 struct{}

func (r *RuleLT09) Code() string        { return "LT09" }
func (r *RuleLT09) Name() string        { return "LT09" }
func (r *RuleLT09) Description() string { return "Rule LT09" }
func (r *RuleLT09) Eval(ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_clause" {
		var targets []parser.BaseSegment
		for _, child := range ctx.Segment.GetSegments() {
			if child.GetType() == "select_clause_element" {
				targets = append(targets, child)
			}
		}
		if len(targets) > 1 {
			firstLine := targets[0].GetPosition().Line
			allSameLine := true
			for _, t := range targets[1:] {
				if t.GetPosition().Line != firstLine {
					allSameLine = false
					break
				}
			}
			if allSameLine {
				pos := ctx.Segment.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Select targets should be on a new line unless there is only one select target.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ctx.Segment,
				}}
			}
		}
	}
	return nil
}

// RuleLT10: layout.select_modifiers
type RuleLT10 struct{}

func (r *RuleLT10) Code() string        { return "LT10" }
func (r *RuleLT10) Name() string        { return "layout.select_modifiers" }
func (r *RuleLT10) Description() string { return "Rule LT10" }
func (r *RuleLT10) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT11: layout.set_operators
type RuleLT11 struct{}

func (r *RuleLT11) Code() string        { return "LT11" }
func (r *RuleLT11) Name() string        { return "layout.set_operators" }
func (r *RuleLT11) Description() string { return "Rule LT11" }
func (r *RuleLT11) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleLT12: layout.end_of_file
type RuleLT12 struct{}

func (r *RuleLT12) Code() string        { return "LT12" }
func (r *RuleLT12) Name() string        { return "layout.end_of_file" }
func (r *RuleLT12) Description() string { return "Rule LT12" }
func (r *RuleLT12) Eval(ctx RuleContext) []LintViolation {
	return evalLT12(r, ctx)
}

// RuleLT13: layout.start_of_file
type RuleLT13 struct{}

func (r *RuleLT13) Code() string        { return "LT13" }
func (r *RuleLT13) Name() string        { return "layout.start_of_file" }
func (r *RuleLT13) Description() string { return "Rule LT13" }
func (r *RuleLT13) Eval(ctx RuleContext) []LintViolation {
	return evalLT13(r, ctx)
}

// RuleLT14: layout.keyword_newline
type RuleLT14 struct{}

func (r *RuleLT14) Code() string        { return "LT14" }
func (r *RuleLT14) Name() string        { return "layout.keyword_newline" }
func (r *RuleLT14) Description() string { return "Rule LT14" }
func (r *RuleLT14) Eval(ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "where_clause" {
		pos := ctx.Segment.GetPosition()
		if pos.Pos > 1 {
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "The 'WHERE' keyword should always start a new line.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
				Fixes: []LintFix{
					CreateBeforeFix(ctx.Segment, parser.NewRawSegment(string(lexer.TokenNewline), "\n", pos)),
				},
			}}
		}
	}
	return nil
}

// RuleLT15: LT15
type RuleLT15 struct{}

func (r *RuleLT15) Code() string        { return "LT15" }
func (r *RuleLT15) Name() string        { return "LT15" }
func (r *RuleLT15) Description() string { return "Rule LT15" }
func (r *RuleLT15) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleOR01: oracle.empty_batch
type RuleOR01 struct{}

func (r *RuleOR01) Code() string        { return "OR01" }
func (r *RuleOR01) Name() string        { return "oracle.empty_batch" }
func (r *RuleOR01) Description() string { return "Rule OR01" }
func (r *RuleOR01) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RulePG01: postgres.excessive_locks
type RulePG01 struct{}

func (r *RulePG01) Code() string        { return "PG01" }
func (r *RulePG01) Name() string        { return "postgres.excessive_locks" }
func (r *RulePG01) Description() string { return "Rule PG01" }
func (r *RulePG01) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RulePG02: postgres.not_valid_foreign_key
type RulePG02 struct{}

func (r *RulePG02) Code() string        { return "PG02" }
func (r *RulePG02) Name() string        { return "postgres.not_valid_foreign_key" }
func (r *RulePG02) Description() string { return "Rule PG02" }
func (r *RulePG02) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF01: references.from
type RuleRF01 struct{}

func (r *RuleRF01) Code() string        { return "RF01" }
func (r *RuleRF01) Name() string        { return "references.from" }
func (r *RuleRF01) Description() string { return "Rule RF01" }
func (r *RuleRF01) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF02: references.qualification
type RuleRF02 struct{}

func (r *RuleRF02) Code() string        { return "RF02" }
func (r *RuleRF02) Name() string        { return "references.qualification" }
func (r *RuleRF02) Description() string { return "Rule RF02" }
func (r *RuleRF02) Eval(ctx RuleContext) []LintViolation {
	return evalRF02(r, ctx)
}

// RuleRF03: references.consistent
type RuleRF03 struct{}

func (r *RuleRF03) Code() string        { return "RF03" }
func (r *RuleRF03) Name() string        { return "references.consistent" }
func (r *RuleRF03) Description() string { return "Rule RF03" }
func (r *RuleRF03) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF04: references.keywords
type RuleRF04 struct{}

func (r *RuleRF04) Code() string        { return "RF04" }
func (r *RuleRF04) Name() string        { return "references.keywords" }
func (r *RuleRF04) Description() string { return "Rule RF04" }
func (r *RuleRF04) Eval(ctx RuleContext) []LintViolation {
	return evalRF04(r, ctx)
}

// RuleRF05: references.special_chars
type RuleRF05 struct{}

func (r *RuleRF05) Code() string        { return "RF05" }
func (r *RuleRF05) Name() string        { return "references.special_chars" }
func (r *RuleRF05) Description() string { return "Rule RF05" }
func (r *RuleRF05) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF06: RF06
type RuleRF06 struct{}

func (r *RuleRF06) Code() string        { return "RF06" }
func (r *RuleRF06) Name() string        { return "RF06" }
func (r *RuleRF06) Description() string { return "Rule RF06" }
func (r *RuleRF06) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF07: references.window_alias
type RuleRF07 struct{}

func (r *RuleRF07) Code() string        { return "RF07" }
func (r *RuleRF07) Name() string        { return "references.window_alias" }
func (r *RuleRF07) Description() string { return "Rule RF07" }
func (r *RuleRF07) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST01: structure.else_null
type RuleST01 struct{}

func (r *RuleST01) Code() string        { return "ST01" }
func (r *RuleST01) Name() string        { return "structure.else_null" }
func (r *RuleST01) Description() string { return "Rule ST01" }
func (r *RuleST01) Eval(ctx RuleContext) []LintViolation {
	return evalST01(r, ctx)
}

// RuleST02: structure.simple_case
type RuleST02 struct{}

func (r *RuleST02) Code() string        { return "ST02" }
func (r *RuleST02) Name() string        { return "structure.simple_case" }
func (r *RuleST02) Description() string { return "Rule ST02" }
func (r *RuleST02) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST03: structure.unused_cte
type RuleST03 struct{}

func (r *RuleST03) Code() string        { return "ST03" }
func (r *RuleST03) Name() string        { return "structure.unused_cte" }
func (r *RuleST03) Description() string { return "Rule ST03" }
func (r *RuleST03) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST04: structure.nested_case
type RuleST04 struct{}

func (r *RuleST04) Code() string        { return "ST04" }
func (r *RuleST04) Name() string        { return "structure.nested_case" }
func (r *RuleST04) Description() string { return "Rule ST04" }
func (r *RuleST04) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST05: structure.subquery
type RuleST05 struct{}

func (r *RuleST05) Code() string        { return "ST05" }
func (r *RuleST05) Name() string        { return "structure.subquery" }
func (r *RuleST05) Description() string { return "Rule ST05" }
func (r *RuleST05) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST06: structure.column_order
type RuleST06 struct{}

func (r *RuleST06) Code() string        { return "ST06" }
func (r *RuleST06) Name() string        { return "structure.column_order" }
func (r *RuleST06) Description() string { return "Rule ST06" }
func (r *RuleST06) Eval(ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "star" {
		pos := ctx.Segment.GetPosition()
		return []LintViolation{{
			RuleCode:    r.Code(),
			Name:        r.Name(),
			Description: "Avoid using 'SELECT *'.",
			Line:        pos.Line,
			Pos:         pos.Pos,
			Segment:     ctx.Segment,
		}}
	}
	return nil
}

// RuleST07: structure.using
type RuleST07 struct{}

func (r *RuleST07) Code() string        { return "ST07" }
func (r *RuleST07) Name() string        { return "structure.using" }
func (r *RuleST07) Description() string { return "Rule ST07" }
func (r *RuleST07) Eval(ctx RuleContext) []LintViolation {
	return evalST07(r, ctx)
}

// RuleST08: structure.distinct
type RuleST08 struct{}

func (r *RuleST08) Code() string        { return "ST08" }
func (r *RuleST08) Name() string        { return "structure.distinct" }
func (r *RuleST08) Description() string { return "Rule ST08" }
func (r *RuleST08) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST09: ST09
type RuleST09 struct{}

func (r *RuleST09) Code() string        { return "ST09" }
func (r *RuleST09) Name() string        { return "ST09" }
func (r *RuleST09) Description() string { return "Rule ST09" }
func (r *RuleST09) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST10: structure.constant_expression
type RuleST10 struct{}

func (r *RuleST10) Code() string        { return "ST10" }
func (r *RuleST10) Name() string        { return "structure.constant_expression" }
func (r *RuleST10) Description() string { return "Rule ST10" }
func (r *RuleST10) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST11: ST11
type RuleST11 struct{}

func (r *RuleST11) Code() string        { return "ST11" }
func (r *RuleST11) Name() string        { return "ST11" }
func (r *RuleST11) Description() string { return "Rule ST11" }
func (r *RuleST11) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleST12: structure.consecutive_semicolons
type RuleST12 struct{}

func (r *RuleST12) Code() string        { return "ST12" }
func (r *RuleST12) Name() string        { return "structure.consecutive_semicolons" }
func (r *RuleST12) Description() string { return "Rule ST12" }
func (r *RuleST12) Eval(ctx RuleContext) []LintViolation {
	return evalST12(r, ctx)
}

// RuleTQ01: tsql.sp_prefix
type RuleTQ01 struct{}

func (r *RuleTQ01) Code() string        { return "TQ01" }
func (r *RuleTQ01) Name() string        { return "tsql.sp_prefix" }
func (r *RuleTQ01) Description() string { return "Rule TQ01" }
func (r *RuleTQ01) Eval(ctx RuleContext) []LintViolation {
	return evalTQ01(r, ctx)
}

// RuleTQ02: tsql.procedure_begin_end
type RuleTQ02 struct{}

func (r *RuleTQ02) Code() string        { return "TQ02" }
func (r *RuleTQ02) Name() string        { return "tsql.procedure_begin_end" }
func (r *RuleTQ02) Description() string { return "Rule TQ02" }
func (r *RuleTQ02) Eval(ctx RuleContext) []LintViolation {
	return evalTQ02(r, ctx)
}

// RuleTQ03: tsql.empty_batch
type RuleTQ03 struct{}

func (r *RuleTQ03) Code() string        { return "TQ03" }
func (r *RuleTQ03) Name() string        { return "tsql.empty_batch" }
func (r *RuleTQ03) Description() string { return "Rule TQ03" }
func (r *RuleTQ03) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// RuleTQ04: tsql.prefer_as_alias
type RuleTQ04 struct{}

func (r *RuleTQ04) Code() string        { return "TQ04" }
func (r *RuleTQ04) Name() string        { return "tsql.prefer_as_alias" }
func (r *RuleTQ04) Description() string { return "Rule TQ04" }
func (r *RuleTQ04) Eval(ctx RuleContext) []LintViolation {
	return nil
}

// GetAllRules returns all 79 registered rules
func GetAllRules() []Rule {
	return []Rule{
		&RuleAL01{},
		&RuleAL02{},
		&RuleAL03{},
		&RuleAL04{},
		&RuleAL05{},
		&RuleAL06{},
		&RuleAL07{},
		&RuleAL08{},
		&RuleAL09{},
		&RuleAL10{},
		&RuleAM01{},
		&RuleAM02{},
		&RuleAM03{},
		&RuleAM04{},
		&RuleAM05{},
		&RuleAM06{},
		&RuleAM07{},
		&RuleAM08{},
		&RuleAM09{},
		&RuleCP01{},
		&RuleCP02{},
		&RuleCP03{},
		&RuleCP04{},
		&RuleCP05{},
		&RuleCV01{},
		&RuleCV02{},
		&RuleCV03{},
		&RuleCV04{},
		&RuleCV05{},
		&RuleCV06{},
		&RuleCV07{},
		&RuleCV08{},
		&RuleCV09{},
		&RuleCV10{},
		&RuleCV11{},
		&RuleCV12{},
		&RuleCV13{},
		&RuleJJ01{},
		&RuleLT01{},
		&RuleLT02{},
		&RuleLT03{},
		&RuleLT04{},
		&RuleLT05{},
		&RuleLT06{},
		&RuleLT07{},
		&RuleLT08{},
		&RuleLT09{},
		&RuleLT10{},
		&RuleLT11{},
		&RuleLT12{},
		&RuleLT13{},
		&RuleLT14{},
		&RuleLT15{},
		&RuleOR01{},
		&RulePG01{},
		&RulePG02{},
		&RuleRF01{},
		&RuleRF02{},
		&RuleRF03{},
		&RuleRF04{},
		&RuleRF05{},
		&RuleRF06{},
		&RuleRF07{},
		&RuleST01{},
		&RuleST02{},
		&RuleST03{},
		&RuleST04{},
		&RuleST05{},
		&RuleST06{},
		&RuleST07{},
		&RuleST08{},
		&RuleST09{},
		&RuleST10{},
		&RuleST11{},
		&RuleST12{},
		&RuleTQ01{},
		&RuleTQ02{},
		&RuleTQ03{},
		&RuleTQ04{},
	}
}
