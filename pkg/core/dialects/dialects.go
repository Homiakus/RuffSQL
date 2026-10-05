package dialects

import (
	"fmt"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/lexer"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// BaseDialect represents an extensible dialect dictionary
type BaseDialect struct {
	Name     string
	Segments map[string]parser.Matchable
}

func NewBaseDialect(name string) *BaseDialect {
	return &BaseDialect{
		Name:     name,
		Segments: make(map[string]parser.Matchable),
	}
}

func (d *BaseDialect) GetName() string {
	return d.Name
}

func (d *BaseDialect) GetSegment(name string) parser.Matchable {
	return d.Segments[name]
}

func (d *BaseDialect) AddSegment(name string, m parser.Matchable) {
	d.Segments[name] = m
}

// Copy creates an inherited clone for dialect inheritance
func (d *BaseDialect) Copy(newName string) *BaseDialect {
	cloned := NewBaseDialect(newName)
	for k, v := range d.Segments {
		cloned.Segments[k] = v
	}
	return cloned
}

// TokenMatcher matches tokens by lexer TokenType
type TokenMatcher struct {
	TokenType lexer.TokenType
	NodeName  string
}

func (t *TokenMatcher) Match(segments []parser.BaseSegment, idx int, ctx *parser.ParseContext) parser.MatchResult {
	if idx >= len(segments) {
		return parser.MatchFailure()
	}
	seg := segments[idx]
	if seg.GetType() == string(t.TokenType) {
		name := t.NodeName
		if name == "" {
			name = string(t.TokenType)
		}
		res := parser.NewRawSegment(name, seg.GetRaw(), seg.GetPosition())
		return parser.MatchSuccess([]parser.BaseSegment{res}, 1)
	}
	return parser.MatchFailure()
}

// BuildANSIDialect defines the standard SQL grammar mimicking SQLFluff ANSI dialect
func BuildANSIDialect() *BaseDialect {
	ansi := NewBaseDialect("ansi")

	// Primitive segments
	ansi.AddSegment("NakedIdentifier", &TokenMatcher{TokenType: lexer.TokenIdentifier, NodeName: "naked_identifier"})
	ansi.AddSegment("QuotedIdentifier", &TokenMatcher{TokenType: lexer.TokenQuotedIdentifier, NodeName: "quoted_identifier"})
	ansi.AddSegment("AnyIdentifier", parser.NewOneOf(parser.NewRef("NakedIdentifier"), parser.NewRef("QuotedIdentifier")))
	ansi.AddSegment("NumericLiteral", &TokenMatcher{TokenType: lexer.TokenNumericLiteral, NodeName: "numeric_literal"})
	ansi.AddSegment("SingleQuote", &TokenMatcher{TokenType: lexer.TokenStringLiteral, NodeName: "single_quote"})
	ansi.AddSegment("Comma", &TokenMatcher{TokenType: lexer.TokenComma, NodeName: "comma"})
	ansi.AddSegment("Dot", &TokenMatcher{TokenType: lexer.TokenDot, NodeName: "dot"})
	ansi.AddSegment("ComparisonOperator", &TokenMatcher{TokenType: lexer.TokenComparisonOperator, NodeName: "comparison_operator"})
	ansi.AddSegment("BinaryOperator", &TokenMatcher{TokenType: lexer.TokenBinaryOperator, NodeName: "binary_operator"})
	ansi.AddSegment("Semicolon", &TokenMatcher{TokenType: lexer.TokenSemicolon, NodeName: "statement_terminator"})
	ansi.AddSegment("Star", parser.NewStringParser("*", "star"))

	// ObjectReference (supports schema.table or table.column: e.g. u.id, public.users)
	ansi.AddSegment("ObjectReference", parser.NewGrammarSegment(
		"column_reference",
		parser.NewDelimited(parser.NewRef("AnyIdentifier"), parser.NewRef("Dot")),
	))

	ansi.AddSegment("ColumnReference", parser.NewRef("ObjectReference"))
	ansi.AddSegment("TableReference", parser.NewRef("ObjectReference"))

	
	// WindowSpecification: OVER ( [PARTITION BY expr] [ORDER BY expr] )
	ansi.AddSegment("WindowSpecification", parser.NewSequence(
		parser.NewStringParser("OVER", "keyword"),
		parser.NewBracketed(
			parser.NewSequence(
				parser.NewAnyNumberOf(
					parser.NewSequence(
						parser.NewStringParser("PARTITION", "keyword"),
						parser.NewStringParser("BY", "keyword"),
						parser.NewDelimited(parser.NewOneOf(parser.NewRef("CastExpression"), parser.NewRef("Expression")), parser.NewRef("Comma")),
					),
				),
				parser.NewAnyNumberOf(
					parser.NewSequence(
						parser.NewStringParser("ORDER", "keyword"),
						parser.NewStringParser("BY", "keyword"),
						parser.NewDelimited(
							parser.NewSequence(
								parser.NewRef("Expression"),
								parser.NewAnyNumberOf(
									parser.NewOneOf(
										parser.NewStringParser("ASC", "keyword"),
										parser.NewStringParser("DESC", "keyword"),
									),
								),
							),
							parser.NewRef("Comma"),
						),
					),
				),
			),
		),
	))

	
	// CastExpression: expression AS datatype
	ansi.AddSegment("CastExpression", parser.NewSequence(
		parser.NewRef("Expression"),
		parser.NewStringParser("AS", "keyword"),
		parser.NewRef("NakedIdentifier"),
	))

	// FunctionCall: identifier ( expression [, ...] | star )
	ansi.AddSegment("FunctionCall", parser.NewGrammarSegment(
		"function",
		parser.NewSequence(
			parser.NewRef("NakedIdentifier"),
			parser.NewBracketed(
				parser.NewOneOf(
					parser.NewRef("Star"),
					parser.NewDelimited(parser.NewOneOf(parser.NewRef("CastExpression"), parser.NewRef("Expression")), parser.NewRef("Comma")),
					parser.NewSequence(),
				),
			),
			parser.NewAnyNumberOf(parser.NewRef("WindowSpecification")),
		),
	))

	ansi.AddSegment("NullLiteral", parser.NewStringParser("NULL", "null_literal"))

	// CaseExpression: CASE [WHEN expr THEN expr ...] [ELSE expr] END
	ansi.AddSegment("CaseExpression", parser.NewGrammarSegment(
		"case_expression",
		parser.NewSequence(
			parser.NewStringParser("CASE", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("WHEN", "keyword"),
					parser.NewRef("Expression"),
					parser.NewStringParser("THEN", "keyword"),
					parser.NewRef("Expression"),
				),
			),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("ELSE", "keyword"),
					parser.NewOneOf(parser.NewRef("NullLiteral"), parser.NewRef("Expression")),
				),
			),
			parser.NewStringParser("END", "keyword"),
		),
	))

	// PrimaryExpression
	ansi.AddSegment("PrimaryExpression", parser.NewOneOf(
		parser.NewRef("NullLiteral"),
		parser.NewRef("CaseExpression"),
		parser.NewRef("FunctionCall"),
		parser.NewRef("ColumnReference"),
		parser.NewRef("NumericLiteral"),
		parser.NewRef("SingleQuote"),
		parser.NewRef("Star"),
		parser.NewBracketed(parser.NewRef("Expression")),
	))

	// Expression
	ansi.AddSegment("Expression", parser.NewGrammarSegment(
		"expression",
		parser.NewSequence(
			parser.NewRef("PrimaryExpression"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewOneOf(parser.NewRef("ComparisonOperator"), parser.NewRef("BinaryOperator")),
					parser.NewRef("PrimaryExpression"),
				),
			),
			parser.NewAnyNumberOf(parser.NewRef("WindowSpecification")),
		),
	))

	// SelectClauseElement: Expression [AS alias]
	ansi.AddSegment("SelectClauseElement", parser.NewGrammarSegment(
		"select_clause_element",
		parser.NewSequence(
			parser.NewRef("Expression"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("AS", "keyword"),
					parser.NewRef("NakedIdentifier"),
				),
			),
			parser.NewAnyNumberOf(parser.NewRef("WindowSpecification")),
		),
	))

	// SelectClause: SELECT [DISTINCT [ON (...) ]] [select_clause_element , ...]
	ansi.AddSegment("SelectClause", parser.NewGrammarSegment(
		"select_clause",
		parser.NewSequence(
			parser.NewStringParser("SELECT", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewOneOf(
					parser.NewStringParser("DISTINCT", "keyword"),
					parser.NewStringParser("ALL", "keyword"),
					parser.NewSequence(
						parser.NewStringParser("DISTINCT", "keyword"),
						parser.NewBracketed(parser.NewRef("Expression")),
					),
				),
			),
			parser.NewDelimited(parser.NewRef("SelectClauseElement"), parser.NewRef("Comma")),
		),
	))

	// FromExpressionElement: table_reference [AS alias]
	ansi.AddSegment("FromExpressionElement", parser.NewGrammarSegment(
		"from_expression_element",
		parser.NewSequence(
			parser.NewGrammarSegment("table_expression", parser.NewRef("TableReference")),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("AS", "keyword"),
					parser.NewRef("NakedIdentifier"),
				),
			),
			parser.NewAnyNumberOf(parser.NewRef("WindowSpecification")),
		),
	))

	// JoinClause: [LEFT|RIGHT|INNER|FULL|CROSS] JOIN from_expression_element [ON expression]
	ansi.AddSegment("JoinClause", parser.NewGrammarSegment(
		"join_clause",
		parser.NewSequence(
			parser.NewAnyNumberOf(
				parser.NewOneOf(
					parser.NewStringParser("LEFT", "keyword"),
					parser.NewStringParser("RIGHT", "keyword"),
					parser.NewStringParser("INNER", "keyword"),
					parser.NewStringParser("FULL", "keyword"),
					parser.NewStringParser("CROSS", "keyword"),
				),
			),
			parser.NewStringParser("JOIN", "keyword"),
			parser.NewRef("FromExpressionElement"),
			parser.NewAnyNumberOf(
				parser.NewOneOf(
					parser.NewSequence(
						parser.NewStringParser("ON", "keyword"),
						parser.NewRef("Expression"),
					),
					parser.NewSequence(
						parser.NewStringParser("USING", "keyword"),
						parser.NewBracketed(parser.NewDelimited(parser.NewRef("NakedIdentifier"), parser.NewRef("Comma"))),
					),
				),
			),
			parser.NewAnyNumberOf(parser.NewRef("WindowSpecification")),
		),
	))

	// FromClause: FROM from_expression_element [join_clause ...]
	ansi.AddSegment("FromClause", parser.NewGrammarSegment(
		"from_clause",
		parser.NewSequence(
			parser.NewStringParser("FROM", "keyword"),
			parser.NewGrammarSegment("from_expression", parser.NewRef("FromExpressionElement")),
			parser.NewAnyNumberOf(parser.NewRef("JoinClause")),
		),
	))

	// WhereClause: WHERE expression
	ansi.AddSegment("WhereClause", parser.NewGrammarSegment(
		"where_clause",
		parser.NewSequence(
			parser.NewStringParser("WHERE", "keyword"),
			parser.NewRef("Expression"),
		),
	))

	// GroupByClause: GROUP BY expression [, ...]
	ansi.AddSegment("GroupByClause", parser.NewGrammarSegment(
		"groupby_clause",
		parser.NewSequence(
			parser.NewStringParser("GROUP", "keyword"),
			parser.NewStringParser("BY", "keyword"),
			parser.NewDelimited(parser.NewOneOf(parser.NewRef("CastExpression"), parser.NewRef("Expression")), parser.NewRef("Comma")),
					parser.NewSequence(),
		),
	))

	// HavingClause: HAVING expression
	ansi.AddSegment("HavingClause", parser.NewGrammarSegment(
		"having_clause",
		parser.NewSequence(
			parser.NewStringParser("HAVING", "keyword"),
			parser.NewRef("Expression"),
		),
	))

	// OrderByClause: ORDER BY expression [ASC|DESC] [, ...]
	ansi.AddSegment("OrderByClause", parser.NewGrammarSegment(
		"orderby_clause",
		parser.NewSequence(
			parser.NewStringParser("ORDER", "keyword"),
			parser.NewStringParser("BY", "keyword"),
			parser.NewDelimited(
				parser.NewSequence(
					parser.NewRef("Expression"),
					parser.NewAnyNumberOf(
						parser.NewOneOf(
							parser.NewStringParser("ASC", "keyword"),
							parser.NewStringParser("DESC", "keyword"),
						),
					),
				),
				parser.NewRef("Comma"),
			),
		),
	))

	// LimitClause: LIMIT numeric_literal
	ansi.AddSegment("LimitClause", parser.NewGrammarSegment(
		"limit_clause",
		parser.NewSequence(
			parser.NewStringParser("LIMIT", "keyword"),
			parser.NewRef("NumericLiteral"),
		),
	))

	// SelectStatement: SelectClause [FromClause] [WhereClause] [GroupByClause] [HavingClause] [OrderByClause] [LimitClause]
	ansi.AddSegment("SelectStatement", parser.NewGrammarSegment(
		"select_statement",
		parser.NewSequence(
			parser.NewRef("SelectClause"),
			parser.NewAnyNumberOf(parser.NewRef("FromClause")),
			parser.NewAnyNumberOf(parser.NewRef("WhereClause")),
			parser.NewAnyNumberOf(parser.NewRef("GroupByClause")),
			parser.NewAnyNumberOf(parser.NewRef("HavingClause")),
			parser.NewAnyNumberOf(parser.NewRef("OrderByClause")),
			parser.NewAnyNumberOf(parser.NewRef("LimitClause")),
		),
	))

	ansi.AddSegment("OnConflictClause", parser.NewOneOf())

	// InsertStatement: INSERT INTO table_reference VALUES ( ... )
	
	ansi.AddSegment("InsertStatement", parser.NewGrammarSegment(
		"insert_statement",
		parser.NewSequence(
			parser.NewStringParser("INSERT", "keyword"),
			parser.NewStringParser("INTO", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewAnyNumberOf(parser.NewBracketed(parser.NewDelimited(parser.NewRef("NakedIdentifier"), parser.NewRef("Comma")))),
			parser.NewStringParser("VALUES", "keyword"),
			parser.NewBracketed(parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma"))),
			parser.NewAnyNumberOf(parser.NewRef("OnConflictClause")),
		),
	))

	// UpdateStatement: UPDATE table_reference SET expression WHERE ...
	ansi.AddSegment("UpdateStatement", parser.NewGrammarSegment(
		"update_statement",
		parser.NewSequence(
			parser.NewStringParser("UPDATE", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewStringParser("SET", "keyword"),
			parser.NewDelimited(parser.NewOneOf(parser.NewRef("CastExpression"), parser.NewRef("Expression")), parser.NewRef("Comma")),
					parser.NewSequence(),
			parser.NewAnyNumberOf(parser.NewRef("WhereClause")),
		),
	))

	// DeleteStatement: DELETE FROM table_reference WHERE ...
	ansi.AddSegment("DeleteStatement", parser.NewGrammarSegment(
		"delete_statement",
		parser.NewSequence(
			parser.NewStringParser("DELETE", "keyword"),
			parser.NewStringParser("FROM", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewAnyNumberOf(parser.NewRef("WhereClause")),
		),
	))

	// CreateTableStatement: CREATE TABLE table_reference ( ... )
	ansi.AddSegment("CreateTableStatement", parser.NewGrammarSegment(
		"create_table_statement",
		parser.NewSequence(
			parser.NewStringParser("CREATE", "keyword"),
			parser.NewStringParser("TABLE", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewBracketed(parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma"))),
			parser.NewAnyNumberOf(parser.NewRef("OnConflictClause")),
		),
	))

	// Statement: Select | Insert | Update | Delete | CreateTable
	ansi.AddSegment("Statement", parser.NewGrammarSegment(
		"statement",
		parser.NewSequence(
			parser.NewOneOf(
				parser.NewRef("SelectStatement"),
				parser.NewRef("InsertStatement"),
				parser.NewRef("UpdateStatement"),
				parser.NewRef("DeleteStatement"),
				parser.NewRef("CreateTableStatement"),
			),
			parser.NewAnyNumberOf(parser.NewRef("Semicolon")),
		),
	))

	return ansi
}

// Registry stores all registered dialects
var registry = make(map[string]*BaseDialect)

// All 28 SQLFluff Dialect names
var SupportedDialects = []string{
	"ansi", "athena", "bigquery", "clickhouse", "databricks", "db2", "doris",
	"duckdb", "exasol", "flink", "greenplum", "hive", "impala", "mariadb",
	"materialize", "mysql", "oracle", "postgres", "redshift", "snowflake",
	"soql", "sparksql", "sqlite", "starrocks", "teradata", "trino", "tsql", "vertica",
}

func init() {
	ansi := BuildANSIDialect()
	AddProceduralConstructs(ansi)
	registry["ansi"] = ansi

	for _, name := range SupportedDialects {
		if name == "ansi" {
			continue
		}
		if name == "postgres" {
			registry["postgres"] = BuildPostgresDialect(ansi)
		} else if name == "bigquery" {
			registry["bigquery"] = BuildBigQueryDialect(ansi)
		} else if name == "snowflake" {
			registry["snowflake"] = BuildSnowflakeDialect(ansi)
		} else if name == "mysql" {
			registry["mysql"] = BuildMySQLDialect(ansi)
		} else if name == "clickhouse" {
			registry["clickhouse"] = BuildClickHouseDialect(ansi)
		} else if name == "duckdb" {
			registry["duckdb"] = BuildDuckDBDialect(ansi)
		} else {
			registry[name] = ansi.Copy(name)
		}
		AddProceduralConstructs(registry[name])
	}
}

func GetDialect(name string) (*BaseDialect, error) {
	d, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("dialect '%s' not supported (supported: %v)", name, SupportedDialects)
	}
	return d, nil
}
