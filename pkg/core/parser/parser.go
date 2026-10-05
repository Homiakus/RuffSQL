package parser

import (
	"fmt"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
)

// Parser manages tokenization and tree generation
type Parser struct {
	Lexer   *lexer.Lexer
	Dialect Dialect
}

func NewParser(d Dialect) *Parser {
	return &Parser{
		Lexer:   lexer.NewLexer(),
		Dialect: d,
	}
}

// ParseSQL parses raw SQL string into an AST root ("file")
func (p *Parser) ParseSQL(sql string) (BaseSegment, error) {
	tokens, err := p.Lexer.Tokenize(sql)
	if err != nil {
		return nil, fmt.Errorf("lexing error: %w", err)
	}

	var rawSegments []BaseSegment
	for _, tok := range tokens {
		rawSegments = append(rawSegments, NewRawSegment(string(tok.Type), tok.Raw, tok.Position))
	}

	ctx := &ParseContext{
		Dialect: p.Dialect,
	}

	stmtGrammar := p.Dialect.GetSegment("Statement")
	if stmtGrammar == nil {
		return nil, fmt.Errorf("root Statement segment not defined in dialect")
	}

	currIdx := 0
	var fileChildren []BaseSegment

	for currIdx < len(rawSegments) {
		// Skip whitespace/newlines at file level
		if rawSegments[currIdx].GetType() == string(lexer.TokenWhitespace) || rawSegments[currIdx].GetType() == string(lexer.TokenNewline) {
			fileChildren = append(fileChildren, rawSegments[currIdx])
			currIdx++
			continue
		}

		res := stmtGrammar.Match(rawSegments, currIdx, ctx)
		if !res.Success {
			// Unparsable segment
			unparsable := NewCompositeSegment("unparsable", []BaseSegment{rawSegments[currIdx]})
			fileChildren = append(fileChildren, unparsable)
			currIdx++
			continue
		}

		fileChildren = append(fileChildren, res.MatchedSegments...)
		currIdx += res.MatchedLength
	}

	pos := lexer.Position{Line: 1, Pos: 1}
	if len(rawSegments) > 0 {
		pos = rawSegments[0].GetPosition()
	}

	fileSeg := &CompositeSegment{
		Type:     "file",
		Segments: fileChildren,
		Position: pos,
	}

	return fileSeg, nil
}
