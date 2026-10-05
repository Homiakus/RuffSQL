package parser

import (
	"testing"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

type mockDialect struct {
	segments map[string]Matchable
}

func (m *mockDialect) GetName() string { return "mock" }
func (m *mockDialect) GetSegment(name string) Matchable {
	return m.segments[name]
}

func TestParserParseSQL(t *testing.T) {
	d := &mockDialect{
		segments: make(map[string]Matchable),
	}
	// Define simple statement grammar matching any keyword
	d.segments["Statement"] = &mockMatchable{}

	p := NewParser(d)
	tree, err := p.ParseSQL("SELECT 1;")
	if err != nil {
		t.Fatalf("unexpected error parsing SQL: %v", err)
	}
	if tree == nil {
		t.Fatalf("expected non-nil AST tree")
	}
	if tree.GetType() != "file" {
		t.Errorf("expected root type 'file', got '%s'", tree.GetType())
	}
}

type mockMatchable struct{}

func (m *mockMatchable) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	if idx < len(segments) {
		return MatchSuccess([]BaseSegment{segments[idx]}, 1)
	}
	return MatchFailure()
}

func TestSegmentStringRepresentation(t *testing.T) {
	raw := NewRawSegment("keyword", "SELECT", lexer.Position{Line: 1, Pos: 1})
	str := raw.String(0)
	if str == "" {
		t.Errorf("expected non-empty string representation")
	}

	comp := NewCompositeSegment("select_statement", []BaseSegment{raw})
	if comp.GetRaw() != "SELECT" {
		t.Errorf("expected composite raw 'SELECT', got '%s'", comp.GetRaw())
	}
}
