package dialects

import (
	
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// BuildPostgresDialect constructs the PostgreSQL specific grammar overrides
func BuildPostgresDialect(ansi *BaseDialect) *BaseDialect {
	pg := ansi.Copy("postgres")

	// 1. JSON operators: ->, ->>, #>, #>>, ?, ?&, ?|
	pg.AddSegment("JsonOperator", parser.NewOneOf(
		parser.NewStringParser("->>", "json_operator"),
		parser.NewStringParser("->", "json_operator"),
		parser.NewStringParser("#>>", "json_operator"),
		parser.NewStringParser("#>", "json_operator"),
	))

	// 2. Postgres Cast operator: ::
	pg.AddSegment("CastOperator", parser.NewStringParser("::", "cast_operator"))

	// 3. Postgres Datatypes
	pg.AddSegment("PostgresDataType", parser.NewOneOf(
		parser.NewStringParser("JSONB", "data_type"),
		parser.NewStringParser("JSON", "data_type"),
		parser.NewStringParser("UUID", "data_type"),
		parser.NewStringParser("TIMESTAMPTZ", "data_type"),
		parser.NewStringParser("TIMESTAMP", "data_type"),
		parser.NewStringParser("TEXT", "data_type"),
		parser.NewStringParser("BIGINT", "data_type"),
		parser.NewStringParser("INT", "data_type"),
		parser.NewStringParser("BOOLEAN", "data_type"),
	))

	// 4. Overrides for Expressions with :: Casting and JSON operators
	pg.AddSegment("PrimaryExpression", parser.NewSequence(
		ansi.GetSegment("PrimaryExpression"),
		parser.NewAnyNumberOf(
			parser.NewOneOf(
				parser.NewSequence(
					parser.NewRef("CastOperator"),
					parser.NewRef("PostgresDataType"),
				),
				parser.NewSequence(
					parser.NewRef("JsonOperator"),
					parser.NewOneOf(parser.NewRef("SingleQuote"), parser.NewRef("NumericLiteral"), parser.NewRef("NakedIdentifier")),
				),
			),
		),
	))

	// 5. ON CONFLICT clause for INSERT statements
	pg.AddSegment("OnConflictClause", parser.NewGrammarSegment(
		"on_conflict_clause",
		parser.NewSequence(
			parser.NewStringParser("ON", "keyword"),
			parser.NewStringParser("CONFLICT", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewBracketed(parser.NewRef("NakedIdentifier")),
			),
			parser.NewStringParser("DO", "keyword"),
			parser.NewOneOf(
				parser.NewStringParser("NOTHING", "keyword"),
				parser.NewSequence(
					parser.NewStringParser("UPDATE", "keyword"),
					parser.NewStringParser("SET", "keyword"),
					parser.NewDelimited(parser.NewSequence(parser.NewRef("NakedIdentifier"), parser.NewRef("ComparisonOperator"), parser.NewRef("Expression")), parser.NewRef("Comma")),
				),
			),
		),
	))

	// 6. Postgres Extended INSERT statement (with ON CONFLICT)
	pg.AddSegment("InsertStatement", parser.NewGrammarSegment(
		"insert_statement",
		parser.NewSequence(
			parser.NewStringParser("INSERT", "keyword"),
			parser.NewStringParser("INTO", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewStringParser("VALUES", "keyword"),
			parser.NewBracketed(parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma"))),
			parser.NewRef("OnConflictClause"),
		),
	))

	// 7. Statement: Select | Insert | Update | Delete | CreateTable
	pg.AddSegment("Statement", parser.NewGrammarSegment(
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

	return pg
}
