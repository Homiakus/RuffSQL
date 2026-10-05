package dialects

import (
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// BuildBigQueryDialect adds BigQuery-specific syntax (backticks, SAFE_CAST, STRUCT, QUALIFY)
func BuildBigQueryDialect(ansi *BaseDialect) *BaseDialect {
	bq := ansi.Copy("bigquery")

	// 1. Safe Cast & BigQuery functions
	bq.AddSegment("BigQueryFunctions", parser.NewOneOf(
		parser.NewStringParser("SAFE_CAST", "function_name"),
		parser.NewStringParser("STRUCT", "function_name"),
		parser.NewStringParser("ARRAY", "function_name"),
		parser.NewStringParser("PARSE_DATE", "function_name"),
		parser.NewStringParser("FORMAT_DATE", "function_name"),
	))

	// 2. QUALIFY clause: QUALIFY expression
	bq.AddSegment("QualifyClause", parser.NewGrammarSegment(
		"qualify_clause",
		parser.NewSequence(
			parser.NewStringParser("QUALIFY", "keyword"),
			parser.NewRef("Expression"),
		),
	))

	// 3. Override SelectStatement to include QUALIFY
	bq.AddSegment("SelectStatement", parser.NewGrammarSegment(
		"select_statement",
		parser.NewSequence(
			parser.NewRef("SelectClause"),
			parser.NewAnyNumberOf(parser.NewRef("FromClause")),
			parser.NewAnyNumberOf(parser.NewRef("WhereClause")),
			parser.NewAnyNumberOf(parser.NewRef("GroupByClause")),
			parser.NewAnyNumberOf(parser.NewRef("HavingClause")),
			parser.NewAnyNumberOf(parser.NewRef("QualifyClause")),
			parser.NewAnyNumberOf(parser.NewRef("OrderByClause")),
			parser.NewAnyNumberOf(parser.NewRef("LimitClause")),
		),
	))

	return bq
}

// BuildSnowflakeDialect adds Snowflake-specific syntax (QUALIFY, ILIKE, semi-structured accessor)
func BuildSnowflakeDialect(ansi *BaseDialect) *BaseDialect {
	sf := ansi.Copy("snowflake")

	// 1. QUALIFY clause
	sf.AddSegment("QualifyClause", parser.NewGrammarSegment(
		"qualify_clause",
		parser.NewSequence(
			parser.NewStringParser("QUALIFY", "keyword"),
			parser.NewRef("Expression"),
		),
	))

	// 2. ILIKE operator
	sf.AddSegment("IlikeOperator", parser.NewStringParser("ILIKE", "comparison_operator"))

	// 3. Override SelectStatement with QUALIFY
	sf.AddSegment("SelectStatement", parser.NewGrammarSegment(
		"select_statement",
		parser.NewSequence(
			parser.NewRef("SelectClause"),
			parser.NewAnyNumberOf(parser.NewRef("FromClause")),
			parser.NewAnyNumberOf(parser.NewRef("WhereClause")),
			parser.NewAnyNumberOf(parser.NewRef("GroupByClause")),
			parser.NewAnyNumberOf(parser.NewRef("HavingClause")),
			parser.NewAnyNumberOf(parser.NewRef("QualifyClause")),
			parser.NewAnyNumberOf(parser.NewRef("OrderByClause")),
			parser.NewAnyNumberOf(parser.NewRef("LimitClause")),
		),
	))

	return sf
}

// BuildMySQLDialect adds MySQL-specific syntax (ON DUPLICATE KEY UPDATE)
func BuildMySQLDialect(ansi *BaseDialect) *BaseDialect {
	my := ansi.Copy("mysql")

	// ON DUPLICATE KEY UPDATE col = val, ...
	my.AddSegment("OnDuplicateKeyUpdate", parser.NewGrammarSegment(
		"on_duplicate_key_update",
		parser.NewSequence(
			parser.NewStringParser("ON", "keyword"),
			parser.NewStringParser("DUPLICATE", "keyword"),
			parser.NewStringParser("KEY", "keyword"),
			parser.NewStringParser("UPDATE", "keyword"),
			parser.NewDelimited(
				parser.NewSequence(
					parser.NewRef("NakedIdentifier"),
					parser.NewRef("ComparisonOperator"),
					parser.NewRef("Expression"),
				),
				parser.NewRef("Comma"),
			),
		),
	))

	my.AddSegment("InsertStatement", parser.NewGrammarSegment(
		"insert_statement",
		parser.NewSequence(
			parser.NewStringParser("INSERT", "keyword"),
			parser.NewStringParser("INTO", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewAnyNumberOf(parser.NewBracketed(parser.NewDelimited(parser.NewRef("NakedIdentifier"), parser.NewRef("Comma")))),
			parser.NewStringParser("VALUES", "keyword"),
			parser.NewBracketed(parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma"))),
			parser.NewAnyNumberOf(parser.NewRef("OnDuplicateKeyUpdate")),
		),
	))

	return my
}

// BuildClickHouseDialect adds ClickHouse-specific syntax (ENGINE = MergeTree(), FINAL, ARRAY JOIN)
func BuildClickHouseDialect(ansi *BaseDialect) *BaseDialect {
	ch := ansi.Copy("clickhouse")

	ch.AddSegment("ClickHouseKeywords", parser.NewOneOf(
		parser.NewStringParser("FINAL", "keyword"),
		parser.NewStringParser("ARRAY", "keyword"),
		parser.NewStringParser("GLOBAL", "keyword"),
	))

	ch.AddSegment("EngineClause", parser.NewGrammarSegment(
		"engine_clause",
		parser.NewSequence(
			parser.NewStringParser("ENGINE", "keyword"),
			parser.NewRef("Equals"),
			parser.NewRef("NakedIdentifier"),
			parser.NewAnyNumberOf(parser.NewBracketed(parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma")))),
		),
	))

	return ch
}

// BuildDuckDBDialect adds DuckDB-specific syntax (COLUMNS(*), EXCLUDE, REPLACE)
func BuildDuckDBDialect(ansi *BaseDialect) *BaseDialect {
	dd := ansi.Copy("duckdb")

	dd.AddSegment("DuckDBKeywords", parser.NewOneOf(
		parser.NewStringParser("EXCLUDE", "keyword"),
		parser.NewStringParser("REPLACE", "keyword"),
		parser.NewStringParser("COLUMNS", "keyword"),
	))

	return dd
}
