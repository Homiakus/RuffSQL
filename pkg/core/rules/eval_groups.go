package rules

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// Implement detailed evaluation for Capitalisation group (CP01-CP05)

// RuleCP02: Unquoted identifiers must be uppercase (or consistently cased)
func evalCP02(r *RuleCP02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "naked_identifier" {
		raw := ctx.Segment.GetRaw()
		if raw != strings.ToUpper(raw) && raw != strings.ToLower(raw) {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Unquoted identifier '%s' has mixed case.", raw),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleCP03: Function names must be consistently uppercase
func evalCP03(r *RuleCP03, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "function" {
		for _, child := range ctx.Segment.GetSegments() {
			if child.GetType() == "naked_identifier" {
				raw := child.GetRaw()
				if raw != strings.ToUpper(raw) {
					pos := child.GetPosition()
					fixedRaw := parser.NewRawSegment("naked_identifier", strings.ToUpper(raw), pos)
					return []LintViolation{{
						RuleCode:    r.Code(),
						Name:        r.Name(),
						Description: fmt.Sprintf("Function name '%s' must be uppercase ('%s').", raw, strings.ToUpper(raw)),
						Line:        pos.Line,
						Pos:         pos.Pos,
						Segment:     child,
						Fixes:       []LintFix{ReplaceFix(child, fixedRaw)},
					}}
				}
			}
		}
	}
	return nil
}

// RuleCP04: Boolean and NULL literals must be uppercase
func evalCP04(r *RuleCP04, ctx RuleContext) []LintViolation {
	rawUpper := strings.ToUpper(ctx.Segment.GetRaw())
	if rawUpper == "NULL" || rawUpper == "TRUE" || rawUpper == "FALSE" {
		if ctx.Segment.GetRaw() != rawUpper {
			pos := ctx.Segment.GetPosition()
			fixedRaw := parser.NewRawSegment(ctx.Segment.GetType(), rawUpper, pos)
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Literal '%s' must be uppercase ('%s').", ctx.Segment.GetRaw(), rawUpper),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
				Fixes:       []LintFix{ReplaceFix(ctx.Segment, fixedRaw)},
			}}
		}
	}
	return nil
}

// RuleCP05: Datatypes must be consistently uppercase
func evalCP05(r *RuleCP05, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "data_type" || ctx.Segment.GetType() == "primitive_type" {
		raw := ctx.Segment.GetRaw()
		if raw != strings.ToUpper(raw) {
			pos := ctx.Segment.GetPosition()
			fixedRaw := parser.NewRawSegment(ctx.Segment.GetType(), strings.ToUpper(raw), pos)
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Data type '%s' must be uppercase ('%s').", raw, strings.ToUpper(raw)),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
				Fixes:       []LintFix{ReplaceFix(ctx.Segment, fixedRaw)},
			}}
		}
	}
	return nil
}

// Implement detailed evaluation for Layout group (LT02 - LT15)

// RuleLT04: Leading/Trailing commas enforcement (default trailing)
func evalLT04(r *RuleLT04, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "comma" {
		// If comma is the first token on a new line (leading comma anti-pattern)
		if ctx.Segment.GetPosition().Pos == 1 {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Commas should not be leading on a new line.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleLT05: Line length limit (e.g. 120 chars)
func evalLT05(r *RuleLT05, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetPosition().Pos > 120 {
		pos := ctx.Segment.GetPosition()
		return []LintViolation{{
			RuleCode:    r.Code(),
			Name:        r.Name(),
			Description: fmt.Sprintf("Line exceeds maximum length of 120 characters (at position %d).", pos.Pos),
			Line:        pos.Line,
			Pos:         pos.Pos,
			Segment:     ctx.Segment,
		}}
	}
	return nil
}

// RuleLT06: Function name and open parenthesis spacing
func evalLT06(r *RuleLT06, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "function" {
		children := ctx.Segment.GetSegments()
		for i := 0; i+1 < len(children); i++ {
			if children[i].GetType() == "naked_identifier" && children[i+1].GetType() == string(lexer.TokenWhitespace) {
				pos := children[i+1].GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "No whitespace allowed between function name and opening bracket.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     children[i+1],
					Fixes:       []LintFix{DeleteFix(children[i+1])},
				}}
			}
		}
	}
	return nil
}

// RuleLT01: No trailing whitespace
func evalLT01(r *RuleLT01, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == string(lexer.TokenWhitespace) && ctx.Parent != nil {
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

// RuleLT02: Indentation must follow consistent tab/space indentation
func evalLT02(r *RuleLT02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == string(lexer.TokenWhitespace) && ctx.Segment.GetPosition().Pos == 1 {
		raw := ctx.Segment.GetRaw()
		if strings.Contains(raw, "\t") {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Tabs used for indentation; spaces required.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleLT12: File must end with a single newline
func evalLT12(r *RuleLT12, ctx RuleContext) []LintViolation {
	if ctx.Parent == nil && ctx.Segment.GetType() == "file" {
		children := ctx.Segment.GetSegments()
		if len(children) > 0 {
			last := children[len(children)-1]
			if last.GetType() != string(lexer.TokenNewline) {
				pos := last.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Files must end with a trailing newline.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     last,
					Fixes: []LintFix{
						CreateAfterFix(last, parser.NewRawSegment(string(lexer.TokenNewline), "\n", pos)),
					},
				}}
			}
		}
	}
	return nil
}

// RuleLT13: Files must not begin with blank lines or newlines
func evalLT13(r *RuleLT13, ctx RuleContext) []LintViolation {
	if ctx.Parent == nil && ctx.Segment.GetType() == "file" {
		children := ctx.Segment.GetSegments()
		if len(children) > 0 && children[0].GetType() == string(lexer.TokenNewline) {
			pos := children[0].GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Files must not begin with a newline.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     children[0],
				Fixes:       []LintFix{DeleteFix(children[0])},
			}}
		}
	}
	return nil
}
