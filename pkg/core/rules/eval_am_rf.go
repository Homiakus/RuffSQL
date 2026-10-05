package rules

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// Group: AM (Ambiguous AM01 - AM09)

// RuleAM01: DISTINCT used with parentheses like DISTINCT(a)
func evalAM01(r *RuleAM01, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_clause" {
		raw := strings.ToUpper(ctx.Segment.GetRaw())
		if strings.Contains(raw, "DISTINCT(") {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "DISTINCT is not a function; remove parentheses: 'DISTINCT a' instead of 'DISTINCT(a)'.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleAM02: UNION [DISTINCT] should be specified instead of bare UNION
func evalAM02(r *RuleAM02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "keyword" && strings.ToUpper(ctx.Segment.GetRaw()) == "UNION" {
		if ctx.Parent != nil {
			siblings := ctx.Parent.GetSegments()
			for i, s := range siblings {
				if s == ctx.Segment {
					hasNextWord := false
					for j := i + 1; j < len(siblings); j++ {
						if siblings[j].GetType() == "keyword" {
							w := strings.ToUpper(siblings[j].GetRaw())
							if w == "ALL" || w == "DISTINCT" {
								hasNextWord = true
							}
							break
						}
					}
					if !hasNextWord {
						pos := ctx.Segment.GetPosition()
						return []LintViolation{{
							RuleCode:    r.Code(),
							Name:        r.Name(),
							Description: "UNION [DISTINCT] should be explicitly specified.",
							Line:        pos.Line,
							Pos:         pos.Pos,
							Segment:     ctx.Segment,
						}}
					}
				}
			}
		}
	}
	return nil
}

// RuleAM03: ORDER BY in subqueries without LIMIT is ambiguous
func evalAM03(r *RuleAM03, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "bracketed" {
		raw := strings.ToUpper(ctx.Segment.GetRaw())
		if strings.Contains(raw, "ORDER BY") && !strings.Contains(raw, "LIMIT") && !strings.Contains(raw, "TOP") {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: "ORDER BY in subquery without LIMIT has no predictable ordering effect.",
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// RuleAM04: Column count in UNION branches must match (heuristic)
func evalAM04(r *RuleAM04, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "file" {
		children := ctx.Segment.GetSegments()
		// Search for multiple select statements separated by UNION
		var selectElementCounts []int
		var selectPositions []parser.BaseSegment

		for _, ch := range children {
			if ch.GetType() == "select_statement" {
				count := 0
				for _, sub := range ch.GetSegments() {
					if sub.GetType() == "select_clause" {
						for _, elem := range sub.GetSegments() {
							if elem.GetType() == "select_clause_element" {
								count++
							}
						}
					}
				}
				if count > 0 {
					selectElementCounts = append(selectElementCounts, count)
					selectPositions = append(selectPositions, ch)
				}
			}
		}

		if len(selectElementCounts) >= 2 {
			first := selectElementCounts[0]
			for i := 1; i < len(selectElementCounts); i++ {
				if selectElementCounts[i] != first {
					pos := selectPositions[i].GetPosition()
					return []LintViolation{{
						RuleCode:    r.Code(),
						Name:        r.Name(),
						Description: fmt.Sprintf("Query has %d columns, but preceding branch has %d columns.", selectElementCounts[i], first),
						Line:        pos.Line,
						Pos:         pos.Pos,
						Segment:     selectPositions[i],
					}}
				}
			}
		}
	}
	return nil
}

// RuleAM06: Ambiguous column references (unqualified columns when multiple tables are joined)
func evalAM06(r *RuleAM06, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "column_reference" && !strings.Contains(ctx.Segment.GetRaw(), ".") && ctx.Segment.GetRaw() != "*" {
		// Check parent path for join_clause
		for _, p := range ctx.Path {
			if p.GetType() == "select_statement" {
				for _, sub := range p.GetSegments() {
					if sub.GetType() == "join_clause" {
						pos := ctx.Segment.GetPosition()
						return []LintViolation{{
							RuleCode:    r.Code(),
							Name:        r.Name(),
							Description: fmt.Sprintf("Ambiguous column reference '%s' in joined query.", ctx.Segment.GetRaw()),
							Line:        pos.Line,
							Pos:         pos.Pos,
							Segment:     ctx.Segment,
						}}
					}
				}
			}
		}
	}
	return nil
}

// Group: RF (References RF01 - RF07)

// RuleRF01: References must exist in FROM / JOIN clauses
func evalRF01(r *RuleRF01, ctx RuleContext) []LintViolation {
	return nil
}

// RuleRF02: References should be qualified when multiple tables are joined
func evalRF02(r *RuleRF02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "select_statement" {
		joinCount := 0
		var nakedCols []parser.BaseSegment

		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "join_clause" {
				joinCount++
			}
			if ch.GetType() == "select_clause" {
				for _, elem := range ch.GetSegments() {
					if elem.GetType() == "select_clause_element" {
						for _, sub := range elem.GetSegments() {
							if sub.GetType() == "expression" {
								for _, atom := range sub.GetSegments() {
									if atom.GetType() == "column_reference" && !strings.Contains(atom.GetRaw(), ".") && atom.GetRaw() != "*" {
										nakedCols = append(nakedCols, atom)
									}
								}
							}
						}
					}
				}
			}
		}

		if joinCount > 0 && len(nakedCols) > 0 {
			var violations []LintViolation
			for _, col := range nakedCols {
				pos := col.GetPosition()
				violations = append(violations, LintViolation{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: fmt.Sprintf("Unqualified reference '%s' when multiple tables are joined.", col.GetRaw()),
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     col,
				})
			}
			return violations
		}
	}
	return nil
}

// RuleRF04: Keywords used as identifiers should be quoted
func evalRF04(r *RuleRF04, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "naked_identifier" {
		raw := strings.ToUpper(ctx.Segment.GetRaw())
		// Common reserved keywords mistaken for column names
		if raw == "DATE" || raw == "USER" || raw == "ORDER" || raw == "GROUP" || raw == "TABLE" {
			pos := ctx.Segment.GetPosition()
			return []LintViolation{{
				RuleCode:    r.Code(),
				Name:        r.Name(),
				Description: fmt.Sprintf("Reserved keyword '%s' used as identifier; quote or rename.", ctx.Segment.GetRaw()),
				Line:        pos.Line,
				Pos:         pos.Pos,
				Segment:     ctx.Segment,
			}}
		}
	}
	return nil
}

// Group: TQ (Quotes and Types TQ01 - TQ04)

// RuleTQ01: Quoted identifiers that do not require quotes
func evalTQ01(r *RuleTQ01, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "quoted_identifier" {
		raw := ctx.Segment.GetRaw()
		if (strings.HasPrefix(raw, "\"") && strings.HasSuffix(raw, "\"")) ||
			(strings.HasPrefix(raw, "`") && strings.HasSuffix(raw, "`")) {
			inner := raw[1 : len(raw)-1]
			// If inner is simple identifier without spaces or special chars
			isSimple := true
			for _, ch := range inner {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
					isSimple = false
					break
				}
			}
			if isSimple && len(inner) > 0 {
				pos := ctx.Segment.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: fmt.Sprintf("Unnecessary quotes on identifier %s.", raw),
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ctx.Segment,
					Fixes: []LintFix{
						ReplaceFix(ctx.Segment, parser.NewRawSegment("naked_identifier", inner, pos)),
					},
				}}
			}
		}
	}
	return nil
}

// RuleTQ02: Unnecessary brackets around expressions
func evalTQ02(r *RuleTQ02, ctx RuleContext) []LintViolation {
	if ctx.Segment.GetType() == "bracketed" {
		for _, ch := range ctx.Segment.GetSegments() {
			if ch.GetType() == "bracketed" {
				pos := ch.GetPosition()
				return []LintViolation{{
					RuleCode:    r.Code(),
					Name:        r.Name(),
					Description: "Unnecessary duplicate parentheses found.",
					Line:        pos.Line,
					Pos:         pos.Pos,
					Segment:     ch,
				}}
			}
			if ch.GetType() == "expression" {
				for _, sub := range ch.GetSegments() {
					if sub.GetType() == "bracketed" {
						pos := sub.GetPosition()
						return []LintViolation{{
							RuleCode:    r.Code(),
							Name:        r.Name(),
							Description: "Unnecessary duplicate parentheses found.",
							Line:        pos.Line,
							Pos:         pos.Pos,
							Segment:     sub,
						}}
					}
				}
			}
		}
	}
	return nil
}

// Group: JJ (Jinja JJ01)

// RuleJJ01: Jinja tags should have consistent internal spacing
func evalJJ01(r *RuleJJ01, ctx RuleContext) []LintViolation {
	return nil
}
