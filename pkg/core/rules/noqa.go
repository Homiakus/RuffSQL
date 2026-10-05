package rules

import (
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// NoqaMap stores disabled rules per line or file-wide
type NoqaMap struct {
	DisableAllLines map[int]bool
	DisabledRules   map[int]map[string]bool
	DisableFile     bool
}

func ExtractNoqa(root parser.BaseSegment) *NoqaMap {
	m := &NoqaMap{
		DisableAllLines: make(map[int]bool),
		DisabledRules:   make(map[int]map[string]bool),
	}

	var walk func(seg parser.BaseSegment)
	walk = func(seg parser.BaseSegment) {
		if seg.GetType() == "comment" {
			raw := seg.GetRaw()
			line := seg.GetPosition().Line
			parseNoqaComment(raw, line, m)
		}
		for _, ch := range seg.GetSegments() {
			walk(ch)
		}
	}

	walk(root)
	return m
}

func parseNoqaComment(comment string, line int, m *NoqaMap) {
	// e.g. -- noqa
	// e.g. -- noqa: disable=all
	// e.g. -- noqa: CP01,LT09
	lower := strings.ToLower(comment)
	idx := strings.Index(lower, "noqa")
	if idx == -1 {
		return
	}

	sub := strings.TrimSpace(comment[idx+4:])
	if sub == "" || strings.HasPrefix(sub, "\n") {
		// Just "-- noqa" disables all rules on this line
		m.DisableAllLines[line] = true
		return
	}

	if strings.HasPrefix(sub, ":") {
		content := strings.TrimSpace(sub[1:])
		if strings.HasPrefix(strings.ToLower(content), "disable=all") {
			m.DisableFile = true
			return
		}
		rules := strings.Split(content, ",")
		if m.DisabledRules[line] == nil {
			m.DisabledRules[line] = make(map[string]bool)
		}
		for _, r := range rules {
			code := strings.TrimSpace(strings.ToUpper(r))
			if code != "" {
				m.DisabledRules[line][code] = true
			}
		}
	}
}

// IsViolationIgnored checks if violation is suppressed by noqa
func (m *NoqaMap) IsViolationIgnored(v LintViolation) bool {
	if m.DisableFile {
		return true
	}
	if m.DisableAllLines[v.Line] {
		return true
	}
	if lineRules, ok := m.DisabledRules[v.Line]; ok {
		if lineRules[v.RuleCode] {
			return true
		}
	}
	return false
}
