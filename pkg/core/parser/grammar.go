package parser

import (
	"fmt"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

// MatchResult represents the outcome of matching a grammar against tokens
type MatchResult struct {
	MatchedSegments []BaseSegment
	MatchedLength   int
	Success         bool
}

func MatchSuccess(segs []BaseSegment, length int) MatchResult {
	return MatchResult{
		MatchedSegments: segs,
		MatchedLength:   length,
		Success:         true,
	}
}

func MatchFailure() MatchResult {
	return MatchResult{
		MatchedSegments: nil,
		MatchedLength:   0,
		Success:         false,
	}
}

// ParseContext passes parser configuration and dialect registry during recursive parsing
type ParseContext struct {
	Dialect Dialect
}

// Dialect interface defines segment resolution and grammar overrides
type Dialect interface {
	GetName() string
	GetSegment(name string) Matchable
}

// Matchable is the core interface for both Segments and Grammars (Ref, Sequence, OneOf, etc.)
type Matchable interface {
	Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult
}

// StringParser matches a specific keyword or literal text
type StringParser struct {
	Template  string
	TokenType string
}

func NewStringParser(template string, tokenType string) *StringParser {
	return &StringParser{
		Template:  template,
		TokenType: tokenType,
	}
}

func (p *StringParser) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	// Skip leading whitespace / comments for code match if needed, but in SQLFluff
	// raw segment comparison tests the current token
	if idx >= len(segments) {
		return MatchFailure()
	}
	seg := segments[idx]
	if strings.EqualFold(seg.GetRaw(), p.Template) {
		matched := NewRawSegment(p.TokenType, seg.GetRaw(), seg.GetPosition())
		return MatchSuccess([]BaseSegment{matched}, 1)
	}
	return MatchFailure()
}

// Ref references another named segment in the current Dialect
type Ref struct {
	Name string
}

func NewRef(name string) *Ref {
	return &Ref{Name: name}
}

func (r *Ref) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	target := ctx.Dialect.GetSegment(r.Name)
	if target == nil {
		panic(fmt.Sprintf("Reference not found in dialect: %s", r.Name))
	}
	return target.Match(segments, idx, ctx)
}

// Sequence matches a sequence of matchables in order
type Sequence struct {
	Elements []Matchable
}

func NewSequence(elements ...Matchable) *Sequence {
	return &Sequence{Elements: elements}
}

func (s *Sequence) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	currIdx := idx
	var accumulated []BaseSegment

	for _, elem := range s.Elements {
		// Skip whitespace/newlines between code tokens in grammar
		for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
			accumulated = append(accumulated, segments[currIdx])
			currIdx++
		}

		res := elem.Match(segments, currIdx, ctx)
		if !res.Success {
			return MatchFailure()
		}
		accumulated = append(accumulated, res.MatchedSegments...)
		currIdx += res.MatchedLength
	}

	return MatchSuccess(accumulated, currIdx-idx)
}

// OneOf matches the first successful branch among options
type OneOf struct {
	Options []Matchable
}

func NewOneOf(options ...Matchable) *OneOf {
	return &OneOf{Options: options}
}

func (o *OneOf) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	for _, opt := range o.Options {
		res := opt.Match(segments, idx, ctx)
		if res.Success {
			return res
		}
	}
	return MatchFailure()
}

// AnyNumberOf matches 0 or more occurrences of grammar elements
type AnyNumberOf struct {
	Element Matchable
	Min     int
	Max     int // 0 means unbounded
}

func NewAnyNumberOf(elem Matchable) *AnyNumberOf {
	return &AnyNumberOf{Element: elem, Min: 0, Max: 0}
}

func (a *AnyNumberOf) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	currIdx := idx
	var accumulated []BaseSegment
	matches := 0

	for currIdx < len(segments) {
		// skip whitespace between repetitions
		origIdx := currIdx
		for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
			accumulated = append(accumulated, segments[currIdx])
			currIdx++
		}

		res := a.Element.Match(segments, currIdx, ctx)
		if !res.Success || res.MatchedLength == 0 {
			// Backtrack whitespace if grammar didn't match next
			currIdx = origIdx
			break
		}
		accumulated = append(accumulated, res.MatchedSegments...)
		currIdx += res.MatchedLength
		matches++
		if a.Max > 0 && matches >= a.Max {
			break
		}
	}

	if matches < a.Min {
		return MatchFailure()
	}
	return MatchSuccess(accumulated, currIdx-idx)
}

// Delimited matches elements separated by a delimiter (e.g. comma)
type Delimited struct {
	Element   Matchable
	Delimiter Matchable
}

func NewDelimited(elem Matchable, delimiter Matchable) *Delimited {
	return &Delimited{Element: elem, Delimiter: delimiter}
}

func (d *Delimited) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	currIdx := idx
	var accumulated []BaseSegment

	// Match first
	res := d.Element.Match(segments, currIdx, ctx)
	if !res.Success {
		return MatchFailure()
	}
	accumulated = append(accumulated, res.MatchedSegments...)
	currIdx += res.MatchedLength

	for currIdx < len(segments) {
		mark := currIdx
		var interim []BaseSegment
		for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
			interim = append(interim, segments[currIdx])
			currIdx++
		}

		delimRes := d.Delimiter.Match(segments, currIdx, ctx)
		if !delimRes.Success {
			currIdx = mark
			break
		}
		interim = append(interim, delimRes.MatchedSegments...)
		currIdx += delimRes.MatchedLength

		for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
			interim = append(interim, segments[currIdx])
			currIdx++
		}

		nextRes := d.Element.Match(segments, currIdx, ctx)
		if !nextRes.Success {
			currIdx = mark
			break
		}
		accumulated = append(accumulated, interim...)
		accumulated = append(accumulated, nextRes.MatchedSegments...)
		currIdx += nextRes.MatchedLength
	}

	return MatchSuccess(accumulated, currIdx-idx)
}

// GrammarSegment is a CompositeSegment definition backed by a Matchable grammar
type GrammarSegment struct {
	Name    string
	Grammar Matchable
}

func NewGrammarSegment(name string, grammar Matchable) *GrammarSegment {
	return &GrammarSegment{Name: name, Grammar: grammar}
}

func (g *GrammarSegment) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	res := g.Grammar.Match(segments, idx, ctx)
	if !res.Success {
		return MatchFailure()
	}
	comp := NewCompositeSegment(g.Name, res.MatchedSegments)
	return MatchSuccess([]BaseSegment{comp}, res.MatchedLength)
}
