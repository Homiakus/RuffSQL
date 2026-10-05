package parser

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

// BaseSegment represents any node in the AST / CST (mimicking SQLFluff's BaseSegment)
type BaseSegment interface {
	GetType() string
	GetRaw() string
	GetPosition() lexer.Position
	GetSegments() []BaseSegment
	IsMeta() bool
	String(indent int) string
}

// RawSegment is a leaf node holding an exact lexical token
type RawSegment struct {
	Type     string         `json:"type"`
	Raw      string         `json:"raw"`
	Position lexer.Position `json:"position"`
	Meta     bool           `json:"meta"`
}

func NewRawSegment(t string, raw string, pos lexer.Position) *RawSegment {
	return &RawSegment{
		Type:     t,
		Raw:      raw,
		Position: pos,
		Meta:     false,
	}
}

func (s *RawSegment) GetType() string              { return s.Type }
func (s *RawSegment) GetRaw() string               { return s.Raw }
func (s *RawSegment) GetPosition() lexer.Position  { return s.Position }
func (s *RawSegment) GetSegments() []BaseSegment   { return nil }
func (s *RawSegment) IsMeta() bool                 { return s.Meta }
func (s *RawSegment) String(indent int) string {
	pad := strings.Repeat("    ", indent)
	if s.Meta {
		return fmt.Sprintf("[L:%3d, P:%3d] %s[META] %s:\n", s.Position.Line, s.Position.Pos, pad, s.Type)
	}
	return fmt.Sprintf("[L:%3d, P:%3d] %s%-30s '%s'\n", s.Position.Line, s.Position.Pos, pad, s.Type+":", s.Raw)
}

// MetaSegment represents zero-width formatting markers like indent / dedent / eof
type MetaSegment struct {
	Type     string         `json:"type"`
	Position lexer.Position `json:"position"`
}

func NewMetaSegment(t string, pos lexer.Position) *MetaSegment {
	return &MetaSegment{
		Type:     t,
		Position: pos,
	}
}

func (s *MetaSegment) GetType() string             { return s.Type }
func (s *MetaSegment) GetRaw() string              { return "" }
func (s *MetaSegment) GetPosition() lexer.Position { return s.Position }
func (s *MetaSegment) GetSegments() []BaseSegment  { return nil }
func (s *MetaSegment) IsMeta() bool                { return true }
func (s *MetaSegment) String(indent int) string {
	pad := strings.Repeat("    ", indent)
	return fmt.Sprintf("[L:%3d, P:%3d] %s[META] %s:\n", s.Position.Line, s.Position.Pos, pad, s.Type)
}

// CompositeSegment represents an interior branch node containing child segments
type CompositeSegment struct {
	Type     string         `json:"type"`
	Segments []BaseSegment  `json:"segments"`
	Position lexer.Position `json:"position"`
}

func NewCompositeSegment(t string, segs []BaseSegment) *CompositeSegment {
	pos := lexer.Position{Line: 1, Pos: 1}
	if len(segs) > 0 {
		pos = segs[0].GetPosition()
	}
	return &CompositeSegment{
		Type:     t,
		Segments: segs,
		Position: pos,
	}
}

func (s *CompositeSegment) GetType() string             { return s.Type }
func (s *CompositeSegment) GetSegments() []BaseSegment  { return s.Segments }
func (s *CompositeSegment) GetPosition() lexer.Position { return s.Position }
func (s *CompositeSegment) IsMeta() bool                { return false }
func (s *CompositeSegment) GetRaw() string {
	var b strings.Builder
	for _, child := range s.Segments {
		b.WriteString(child.GetRaw())
	}
	return b.String()
}

func (s *CompositeSegment) String(indent int) string {
	pad := strings.Repeat("    ", indent)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[L:%3d, P:%3d] %s%s:\n", s.Position.Line, s.Position.Pos, pad, s.Type))
	for _, child := range s.Segments {
		sb.WriteString(child.String(indent + 1))
	}
	return sb.String()
}
