package linter

import (
	"fmt"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/templater"
	"os"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/config"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/dialects"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/rules"
)

type Linter struct {
	Dialect *dialects.BaseDialect
	Rules   []rules.Rule
	Config  *config.FluffConfig
}

func NewLinter(dialectName string) (*Linter, error) {
	d, err := dialects.GetDialect(dialectName)
	if err != nil {
		return nil, err
	}
	return &Linter{
		Dialect: d,
		Rules:   rules.GetAllRules(),
		Config:  config.DefaultConfig(),
	}, nil
}

func NewLinterWithConfig(cfg *config.FluffConfig) (*Linter, error) {
	dialectName := cfg.Dialect
	if dialectName == "" {
		dialectName = "ansi"
	}
	d, err := dialects.GetDialect(dialectName)
	if err != nil {
		return nil, err
	}
	return &Linter{
		Dialect: d,
		Rules:   rules.GetAllRules(),
		Config:  cfg,
	}, nil
}

type LintResult struct {
	FilePath   string
	ParsedTree parser.BaseSegment
	Violations []rules.LintViolation
}

func (l *Linter) LintString(sql string, filePath string) (*LintResult, error) {
	cleanSQL := sql
	if templater.IsTemplated(sql) {
		cleanSQL = templater.SanitizeJinja(sql).SanitizedSQL
	}
	p := parser.NewParser(l.Dialect)
	tree, err := p.ParseSQL(cleanSQL)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var rawViolations []rules.LintViolation
	l.crawlAndLint(tree, nil, nil, &rawViolations)

	// Extract and apply inline -- noqa directives
	noqaMap := rules.ExtractNoqa(tree)
	var filteredViolations []rules.LintViolation

	for _, v := range rawViolations {
		// Filter by config (enabled/disabled rules)
		if l.Config != nil && !l.Config.IsRuleEnabled(v.RuleCode) {
			continue
		}
		// Filter by inline -- noqa
		if noqaMap.IsViolationIgnored(v) {
			continue
		}
		filteredViolations = append(filteredViolations, v)
	}

	return &LintResult{
		FilePath:   filePath,
		ParsedTree: tree,
		Violations: filteredViolations,
	}, nil
}

func (l *Linter) LintFile(path string) (*LintResult, error) {
	if l.Config != nil && l.Config.IsPathIgnored(path) {
		return &LintResult{FilePath: path}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return l.LintString(string(data), path)
}

func (l *Linter) FixString(sql string, maxIterations int) (string, int, error) {
	current := sql
	totalFixed := 0

	for iter := 0; iter < maxIterations; iter++ {
		res, err := l.LintString(current, "<fix>")
		if err != nil {
			return current, totalFixed, err
		}

		var allFixes []rules.LintFix
		for _, v := range res.Violations {
			allFixes = append(allFixes, v.Fixes...)
		}

		if len(allFixes) == 0 {
			break
		}

		newSQL, count := applyFixes(res.ParsedTree, allFixes)
		if count == 0 || newSQL == current {
			break
		}
		totalFixed += count
		current = newSQL
	}

	return current, totalFixed, nil
}

func applyFixes(root parser.BaseSegment, fixes []rules.LintFix) (string, int) {
	replaceMap := make(map[parser.BaseSegment][]parser.BaseSegment)
	deleteMap := make(map[parser.BaseSegment]bool)
	createBeforeMap := make(map[parser.BaseSegment][]parser.BaseSegment)

	count := 0
	for _, fix := range fixes {
		switch fix.Type {
		case rules.FixReplace:
			replaceMap[fix.Anchor] = fix.Edits
			count++
		case rules.FixDelete:
			deleteMap[fix.Anchor] = true
			count++
		case rules.FixCreateBefore:
			createBeforeMap[fix.Anchor] = append(createBeforeMap[fix.Anchor], fix.Edits...)
			count++
		}
	}

	var sb strings.Builder
	var serialize func(seg parser.BaseSegment)
	serialize = func(seg parser.BaseSegment) {
		if deletes, ok := deleteMap[seg]; ok && deletes {
			return
		}
		if befores, ok := createBeforeMap[seg]; ok {
			for _, b := range befores {
				sb.WriteString(b.GetRaw())
			}
		}
		if replacements, ok := replaceMap[seg]; ok {
			for _, r := range replacements {
				sb.WriteString(r.GetRaw())
			}
			return
		}

		children := seg.GetSegments()
		if len(children) == 0 {
			sb.WriteString(seg.GetRaw())
		} else {
			for _, ch := range children {
				serialize(ch)
			}
		}
	}

	serialize(root)
	return sb.String(), count
}

func (l *Linter) crawlAndLint(seg parser.BaseSegment, parent parser.BaseSegment, path []parser.BaseSegment, out *[]rules.LintViolation) {
	ctx := rules.RuleContext{
		Segment: seg,
		Parent:  parent,
		Path:    path,
	}

	for _, r := range l.Rules {
		v := r.Eval(ctx)
		if len(v) > 0 {
			*out = append(*out, v...)
		}
	}

	currentPath := append(path, seg)
	for _, child := range seg.GetSegments() {
		l.crawlAndLint(child, seg, currentPath, out)
	}
}
