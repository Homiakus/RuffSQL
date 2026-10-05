package templater

import (
	"strings"
	)

// JinjaTagType indicates the kind of template delimiter
type JinjaTagType string

const (
	JinjaExpression JinjaTagType = "expression" // {{ ... }}
	JinjaBlock      JinjaTagType = "block"      // {% ... %}
	JinjaComment    JinjaTagType = "comment"    // {# ... #}
)

// JinjaSlice records the position of templated code in the raw input
type JinjaSlice struct {
	Type     JinjaTagType
	Raw      string
	StartPos int
	EndPos   int
}

// TemplatedResult holds the sanitized SQL (for parsing) and the slices mapping back
type TemplatedResult struct {
	RawSQL       string
	SanitizedSQL string
	Slices       []JinjaSlice
}

// SanitizeJinja replaces Jinja expressions and blocks with valid placeholder SQL tokens
// while maintaining identical byte length and newline structure so (Line, Pos) stay 100% accurate.
func SanitizeJinja(raw string) *TemplatedResult {
	var slices []JinjaSlice
	var sb strings.Builder
	runes := []rune(raw)
	n := len(runes)
	i := 0

	for i < n {
		// 1. Jinja Expression: {{ ... }}
		if runes[i] == '{' && i+1 < n && runes[i+1] == '{' {
			start := i
			i += 2
			for i+1 < n && !(runes[i] == '}' && runes[i+1] == '}') {
				i++
			}
			if i+1 < n {
				i += 2
			}
			content := string(runes[start:i])
			slices = append(slices, JinjaSlice{
				Type:     JinjaExpression,
				Raw:      content,
				StartPos: start,
				EndPos:   i,
			})

			// Replace with a dummy valid identifier like `dbt_ref_placeholder` padded with whitespace
			placeholder := "dbt_expr"
			if len(content) > len(placeholder) {
				sb.WriteString(placeholder)
				sb.WriteString(strings.Repeat(" ", len(content)-len(placeholder)))
			} else {
				sb.WriteString(content)
			}
			continue
		}

		// 2. Jinja Block: {% ... %}
		if runes[i] == '{' && i+1 < n && runes[i+1] == '%' {
			start := i
			i += 2
			for i+1 < n && !(runes[i] == '%' && runes[i+1] == '}') {
				i++
			}
			if i+1 < n {
				i += 2
			}
			content := string(runes[start:i])
			slices = append(slices, JinjaSlice{
				Type:     JinjaBlock,
				Raw:      content,
				StartPos: start,
				EndPos:   i,
			})

			// Blocks often wrap clauses, replace with whitespace/newlines preserving line counts
			for _, r := range content {
				if r == '\n' {
					sb.WriteRune('\n')
				} else if r == '\r' {
					sb.WriteRune('\r')
				} else {
					sb.WriteRune(' ')
				}
			}
			continue
		}

		// 3. Jinja Comment: {# ... #}
		if runes[i] == '{' && i+1 < n && runes[i+1] == '#' {
			start := i
			i += 2
			for i+1 < n && !(runes[i] == '#' && runes[i+1] == '}') {
				i++
			}
			if i+1 < n {
				i += 2
			}
			content := string(runes[start:i])
			slices = append(slices, JinjaSlice{
				Type:     JinjaComment,
				Raw:      content,
				StartPos: start,
				EndPos:   i,
			})

			for _, r := range content {
				if r == '\n' {
					sb.WriteRune('\n')
				} else if r == '\r' {
					sb.WriteRune('\r')
				} else {
					sb.WriteRune(' ')
				}
			}
			continue
		}

		sb.WriteRune(runes[i])
		i++
	}

	return &TemplatedResult{
		RawSQL:       raw,
		SanitizedSQL: sb.String(),
		Slices:       slices,
	}
}

// IsTemplated checks if the SQL string contains Jinja delimiters
func IsTemplated(sql string) bool {
	return strings.Contains(sql, "{{") || strings.Contains(sql, "{%") || strings.Contains(sql, "{#")
}
