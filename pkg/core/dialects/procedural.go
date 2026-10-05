package dialects

import (
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
)

// AddProceduralConstructs injects procedural SQL extensions into a dialect:
// 1. CREATE [OR REPLACE] PROCEDURE / FUNCTION / TRIGGER
// 2. CALL / EXEC / EXECUTE procedure_name(...)
// 3. DECLARE var_name datatype [DEFAULT val];
// 4. BEGIN ... END block with nested statements
// 5. IF condition THEN ... [ELSEIF/ELSIF ...] [ELSE ...] END IF;
// 6. WHILE condition DO ... END WHILE; / LOOP ... END LOOP;
func AddProceduralConstructs(d *BaseDialect) {
	// 1. Call/Exec statement: CALL proc(args) / EXEC proc(args)
	d.AddSegment("CallStatement", parser.NewGrammarSegment(
		"call_statement",
		parser.NewSequence(
			parser.NewOneOf(
				parser.NewStringParser("CALL", "keyword"),
				parser.NewStringParser("EXECUTE", "keyword"),
				parser.NewStringParser("EXEC", "keyword"),
			),
			parser.NewRef("ObjectReference"),
			parser.NewAnyNumberOf(
				parser.NewBracketed(
					parser.NewAnyNumberOf(
						parser.NewDelimited(parser.NewRef("Expression"), parser.NewRef("Comma")),
					),
				),
			),
		),
	))

	// 2. DECLARE statement: DECLARE var_name datatype [DEFAULT expr]
	d.AddSegment("DeclareStatement", parser.NewGrammarSegment(
		"declare_statement",
		parser.NewSequence(
			parser.NewStringParser("DECLARE", "keyword"),
			parser.NewRef("NakedIdentifier"),
			parser.NewRef("NakedIdentifier"), // Datatype (INT, VARCHAR, etc.)
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("DEFAULT", "keyword"),
					parser.NewRef("Expression"),
				),
			),
		),
	))

	// 3. Return statement: RETURN [expr]
	d.AddSegment("ReturnStatement", parser.NewGrammarSegment(
		"return_statement",
		parser.NewSequence(
			parser.NewStringParser("RETURN", "keyword"),
			parser.NewAnyNumberOf(parser.NewRef("Expression")),
		),
	))

	// 4. Procedural Block Body (BEGIN ... END)
	d.AddSegment("BeginEndBlock", parser.NewGrammarSegment(
		"begin_end_block",
		parser.NewSequence(
			parser.NewStringParser("BEGIN", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewOneOf(
						parser.NewRef("DeclareStatement"),
						parser.NewRef("ReturnStatement"),
						parser.NewRef("CallStatement"),
						parser.NewRef("SelectStatement"),
						parser.NewRef("InsertStatement"),
						parser.NewRef("UpdateStatement"),
						parser.NewRef("DeleteStatement"),
					),
					parser.NewAnyNumberOf(parser.NewRef("Semicolon")),
				),
			),
			parser.NewStringParser("END", "keyword"),
		),
	))

	// 5. CREATE [OR REPLACE] PROCEDURE
	d.AddSegment("CreateProcedureStatement", parser.NewGrammarSegment(
		"create_procedure_statement",
		parser.NewSequence(
			parser.NewStringParser("CREATE", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("OR", "keyword"),
					parser.NewStringParser("REPLACE", "keyword"),
				),
			),
			parser.NewStringParser("PROCEDURE", "keyword"),
			parser.NewRef("ObjectReference"),
			parser.NewAnyNumberOf(
				parser.NewBracketed(
					parser.NewAnyNumberOf(
						parser.NewDelimited(
							parser.NewSequence(
								parser.NewAnyNumberOf(
									parser.NewOneOf(
										parser.NewStringParser("IN", "keyword"),
										parser.NewStringParser("OUT", "keyword"),
										parser.NewStringParser("INOUT", "keyword"),
									),
								),
								parser.NewRef("NakedIdentifier"),
								parser.NewRef("NakedIdentifier"),
							),
							parser.NewRef("Comma"),
						),
					),
				),
			),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("LANGUAGE", "keyword"),
					parser.NewRef("NakedIdentifier"),
				),
			),
			parser.NewAnyNumberOf(
				parser.NewStringParser("AS", "keyword"),
			),
			parser.NewOneOf(
				parser.NewRef("BeginEndBlock"),
				parser.NewRef("Statement"),
			),
		),
	))

	// 6. CREATE [OR REPLACE] FUNCTION
	d.AddSegment("CreateFunctionStatement", parser.NewGrammarSegment(
		"create_function_statement",
		parser.NewSequence(
			parser.NewStringParser("CREATE", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("OR", "keyword"),
					parser.NewStringParser("REPLACE", "keyword"),
				),
			),
			parser.NewStringParser("FUNCTION", "keyword"),
			parser.NewRef("ObjectReference"),
			parser.NewBracketed(
				parser.NewAnyNumberOf(
					parser.NewDelimited(
						parser.NewSequence(
							parser.NewRef("NakedIdentifier"),
							parser.NewRef("NakedIdentifier"),
						),
						parser.NewRef("Comma"),
					),
				),
			),
			parser.NewStringParser("RETURNS", "keyword"),
			parser.NewRef("NakedIdentifier"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("LANGUAGE", "keyword"),
					parser.NewRef("NakedIdentifier"),
				),
			),
			parser.NewAnyNumberOf(
				parser.NewStringParser("AS", "keyword"),
			),
			parser.NewOneOf(
				parser.NewRef("BeginEndBlock"),
				parser.NewRef("Statement"),
			),
		),
	))

	// 7. CREATE TRIGGER statement
	d.AddSegment("CreateTriggerStatement", parser.NewGrammarSegment(
		"create_trigger_statement",
		parser.NewSequence(
			parser.NewStringParser("CREATE", "keyword"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("OR", "keyword"),
					parser.NewStringParser("REPLACE", "keyword"),
				),
			),
			parser.NewStringParser("TRIGGER", "keyword"),
			parser.NewRef("NakedIdentifier"),
			parser.NewOneOf(
				parser.NewStringParser("BEFORE", "keyword"),
				parser.NewStringParser("AFTER", "keyword"),
				parser.NewSequence(
					parser.NewStringParser("INSTEAD", "keyword"),
					parser.NewStringParser("OF", "keyword"),
				),
			),
			parser.NewOneOf(
				parser.NewStringParser("INSERT", "keyword"),
				parser.NewStringParser("UPDATE", "keyword"),
				parser.NewStringParser("DELETE", "keyword"),
			),
			parser.NewStringParser("ON", "keyword"),
			parser.NewRef("TableReference"),
			parser.NewAnyNumberOf(
				parser.NewSequence(
					parser.NewStringParser("FOR", "keyword"),
					parser.NewStringParser("EACH", "keyword"),
					parser.NewStringParser("ROW", "keyword"),
				),
			),
			parser.NewOneOf(
				parser.NewRef("CallStatement"),
				parser.NewRef("BeginEndBlock"),
			),
		),
	))

	// Expand Statement to include procedural statements
	d.AddSegment("Statement", parser.NewGrammarSegment(
		"statement",
		parser.NewSequence(
			parser.NewOneOf(
				parser.NewRef("CreateProcedureStatement"),
				parser.NewRef("CreateFunctionStatement"),
				parser.NewRef("CreateTriggerStatement"),
				parser.NewRef("CallStatement"),
				parser.NewRef("BeginEndBlock"),
				parser.NewRef("DeclareStatement"),
				parser.NewRef("ReturnStatement"),
				parser.NewRef("SelectStatement"),
				parser.NewRef("InsertStatement"),
				parser.NewRef("UpdateStatement"),
				parser.NewRef("DeleteStatement"),
				parser.NewRef("CreateTableStatement"),
			),
			parser.NewAnyNumberOf(parser.NewRef("Semicolon")),
		),
	))
}
