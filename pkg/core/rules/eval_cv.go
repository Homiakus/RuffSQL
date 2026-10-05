package rules

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// Group: CV (Convention CV01 - CV13)

// RuleCV02: Use COALESCE instead of NVL / IFNULL
func evalCV02(r *RuleCV02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "function" {
		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "naked_identifier" {
				fnName := strings.ToUpper(ch.GetRaw())
				if fnName == "IFNULL" || fnName == "NVL" {
					pos := ch.GetPosition()
					fixedRaw := parser.NewRawSegment("naked_identifier", "COALESCE", pos)
					return []LintViolation{{
						RuleCode:    r.Code(),
						Name:        r.Name(),
						Description: fmt.Sprintf("Use 'COALESCE' instead of '%s'.", ch.GetRaw()),
						Line:        pos.Line,
						Pos:         pos.Pos,
						Segment:     ch,
						Fixes:       []LintFix{ReplaceFix(ch, fixedRaw)},
					}}
				}
			}
		}
	}
	return nil
}

// RuleCV03: Trailing comma in select clause (default: forbid trailing comma)
func evalCV03(r *RuleCV03, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_clause" {
		children := ctx.Segment.GetSegments()
		// Search backward for the last code token
		for i := len(children) - 1; i >= 0; i-- {
			ch := children[i]
			if ch.GetType() == string(lexer.TokenWhitespace) || ch.GetType() == string(lexer.TokenNewline) {
				continue
			}
			if ch.GetType() == "comma" {
				pos := ch.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Trailing comma in SELECT clause should not be used.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ch,
					Fixes:       []LintFix{DeleteFix(ch)},
				}}
			}
			break
		}
	}
	return nil
}

// RuleCV04: Count rows convention (prefer COUNT(1) or COUNT(*))
func evalCV04(r *RuleCV04, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "function" {
		children := ctx.Segment.GetSegments()
		if len(children) >= 2 && strings.ToUpper(children[0].GetRaw()) == "COUNT" {
			rawInside := strings.TrimSpace(children[1].GetRaw())
			// e.g. COUNT(0)
			if rawInside == "(0)" {
				pos := children[1].GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Use 'COUNT(*)' instead of 'COUNT(0)'.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     children[1],
				}}
			}
		}
	}
	return nil
}

// RuleCV05: Comparisons with NULL should use IS NULL / IS NOT NULL instead of = NULL or != NULL
func evalCV05(r *RuleCV05, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "expression" {
		var codeSegs []parser.BaseSegment
		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() != string(lexer.TokenWhitespace) && ch.GetType() != string(lexer.TokenNewline) {
				codeSegs = append(codeSegs, ch)
			}
		}
		for i := 0; i+2 < len(codeSegs); i++ {
			op := codeSegs[i+1]
			right := codeSegs[i+2]
			if (op.GetRaw() == "=" || op.GetRaw() == "!=" || op.GetRaw() == "<>") && strings.ToUpper(right.GetRaw()) == "NULL" {
				pos := op.GetPosition()
				desc := "Use 'IS NULL' instead of '= NULL'."
				if op.GetRaw() != "=" {
					desc = "Use 'IS NOT NULL' instead of '!= NULL' or '<> NULL'."
				}
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: desc,
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     op,
				}}
			}
		}
	}
	return nil
}

// RuleCV06: Statements must end with a semicolon delimiter
func evalCV06(r *RuleCV06, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "statement" {
		children := ctx.Segment.GetSegments()
		hasTerminator := false
		var lastSeg parser.BaseSegment
		for _, ch := range children {
			if ch.GetType() == "statement_terminator" {
				hasTerminator = true
			}
			if ch.GetType() != string(lexer.TokenWhitespace) && ch.GetType() != string(lexer.TokenNewline) {
				lastSeg = ch
			}
		}
		if !hasTerminator && lastSeg != nil {
			pos := lastSeg.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Statements must end with a semicolon delimiter.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     lastSeg,
				Fixes: []LintFix{
					CreateAfterFix(lastSeg, parser.NewRawSegment("statement_terminator", ";", pos)),
				},
			}}
		}
	}
	return nil
}

// RuleCV08: Disallow RIGHT OUTER JOIN, prefer LEFT OUTER JOIN
func evalCV08(r *RuleCV08, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "join_clause" {
		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "keyword" && strings.ToUpper(ch.GetRaw()) == "RIGHT" {
				pos := ch.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Use 'LEFT JOIN' instead of 'RIGHT JOIN'.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ch,
				}}
			}
		}
	}
	return nil
}

// RuleCV07: Statement brackets - redundant brackets around top-level queries
func evalCV07(r *RuleCV07, ctx RuleContext) []LintViolation {
	if ctx.Parent == nil && ctx.Segment.GetType() == "file" {
		children := ctx.Segment.GetSegments()
		// If the file consists solely of bracketed content
		var nonWs []parser.BaseSegment
		for _, ch := range children {
			if ch.GetType() != string(lexer.TokenWhitespace) && ch.GetType() != string(lexer.TokenNewline) {
				nonWs = append(nonWs, ch)
			}
		}
		if len(nonWs) == 1 && nonWs[0].GetType() == "bracketed" {
			pos := nonWs[0].GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Redundant outer brackets surrounding the entire statement.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     nonWs[0],
			}}
		}
	}
	return nil
}

// RuleCV09: Disallow blocked/deprecated words
func evalCV09(r *RuleCV09, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "keyword" {
		raw := strings.ToUpper(ctx.Segment.GetRaw())
		if raw == "GOTO" {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Disallowed statement keyword 'GOTO'.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleCV10: Quoted literals convention (prefer single quotes for strings)
func evalCV10(r *RuleCV10, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "literal" || ctx.Segment.GetType() == "quoted_identifier" {
		raw := ctx.Segment.GetRaw()
		if strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"") && len(raw) >= 2 {
			pos := ctx.Segment.GetPosition()
			singleQuoted := "'" + raw[1:len(raw)-1] + "'"
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Use single quotes for string literals instead of double quotes.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
				Fixes: []LintFix{
					ReplaceFix(ctx.Segment, parser.NewRawSegment("single_quote", singleQuoted, pos)),
				},
			}}
		}
	}
	return nil
}

// RuleCV11: Casting style (prefer CAST(...) or standard :: consistently)
func evalCV11(r *RuleCV11, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "function" {
		children := ctx.Segment.GetSegments()
		if len(children) > 0 && strings.ToUpper(children[0].GetRaw()) == "CONVERT" {
			pos := children[0].GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "Use 'CAST(...)' instead of vendor-specific 'CONVERT(...)'.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     children[0],
			}}
		}
	}
	return nil
}
