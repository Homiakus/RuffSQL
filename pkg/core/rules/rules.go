package rules

import (
	"fmt"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

type FixType string

const (
	FixReplace      FixType = "replace"
	FixDelete       FixType = "delete"
	FixCreateBefore FixType = "create_before"
	FixCreateAfter  FixType = "create_after"
)

// LintFix represents a concrete edit proposal for an AST segment
type LintFix struct {
	Type   FixType
	Anchor parser.BaseSegment
	Edits  []parser.BaseSegment
}

func ReplaceFix(anchor parser.BaseSegment, edits ...parser.BaseSegment) LintFix {
	return LintFix{
		Type:   FixReplace,
		Anchor: anchor,
		Edits:  edits,
	}
}

func DeleteFix(anchor parser.BaseSegment) LintFix {
	return LintFix{
		Type:   FixDelete,
		Anchor: anchor,
	}
}

func CreateBeforeFix(anchor parser.BaseSegment, edits ...parser.BaseSegment) LintFix {
	return LintFix{
		Type:   FixCreateBefore,
		Anchor: anchor,
		Edits:  edits,
	}
}

func CreateAfterFix(anchor parser.BaseSegment, edits ...parser.BaseSegment) LintFix {
	return LintFix{
		Type:   FixCreateAfter,
		Anchor: anchor,
		Edits:  edits,
	}
}

// LintViolation represents a rule violation with optional auto-fixes
type LintViolation struct {
	RuleCode    string             `json:"code"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Line        int                `json:"line"`
	Pos         int                `json:"pos"`
	Segment     parser.BaseSegment `json:"-"`
	Fixes       []LintFix          `json:"-"`
}

func (v LintViolation) String() string {
	return fmt.Sprintf("L: %3d | P: %3d | %-4s | %s", v.Line, v.Pos, v.RuleCode, v.Description)
}

// RuleContext passes contextual info during rule evaluation
type RuleContext struct {
	Segment parser.BaseSegment
	Parent  parser.BaseSegment
	Path    []parser.BaseSegment
}

// Rule is the interface every linting rule implements
type Rule interface {
	Code() string
	Name() string
	Description() string
	Eval(ctx RuleContext) []LintViolation
}
