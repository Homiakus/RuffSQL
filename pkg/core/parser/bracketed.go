package parser

import (
	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

// Bracketed matches expressions enclosed by open and close brackets: ( ... )
type Bracketed struct {
	Element Matchable
}

func NewBracketed(elem Matchable) *Bracketed {
	return &Bracketed{Element: elem}
}

func (b *Bracketed) Match(segments []BaseSegment, idx int, ctx *ParseContext) MatchResult {
	currIdx := idx
	// skip whitespace
	for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
		currIdx++
	}
	if currIdx >= len(segments) {
		return MatchFailure()
	}

	// Must start with '('
	if segments[currIdx].GetType() != string(lexer.TokenOpenBracket) && segments[currIdx].GetRaw() != "(" {
		return MatchFailure()
	}

	var accumulated []BaseSegment
	accumulated = append(accumulated, segments[currIdx])
	currIdx++

	// Match inner element
	innerRes := b.Element.Match(segments, currIdx, ctx)
	if !innerRes.Success {
		return MatchFailure()
	}
	accumulated = append(accumulated, innerRes.MatchedSegments...)
	currIdx += innerRes.MatchedLength

	// Skip whitespace before closing ')'
	for currIdx < len(segments) && (segments[currIdx].GetType() == string(lexer.TokenWhitespace) || segments[currIdx].GetType() == string(lexer.TokenNewline)) {
		accumulated = append(accumulated, segments[currIdx])
		currIdx++
	}

	if currIdx >= len(segments) || (segments[currIdx].GetType() != string(lexer.TokenCloseBracket) && segments[currIdx].GetRaw() != ")") {
		return MatchFailure()
	}
	accumulated = append(accumulated, segments[currIdx])
	currIdx++

	bracketedSegment := NewCompositeSegment("bracketed", accumulated)
	return MatchSuccess([]BaseSegment{bracketedSegment}, currIdx-idx)
}
