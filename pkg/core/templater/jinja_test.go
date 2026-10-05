package templater

import (
	"strings"
	"testing"
)

func TestIsTemplated(t *testing.T) {
	if !IsTemplated("SELECT * FROM {{ ref('my_table') }}") {
		t.Errorf("expected true for jinja expression")
	}
	if !IsTemplated("{% set x = 1 %} SELECT 1") {
		t.Errorf("expected true for jinja block")
	}
	if !IsTemplated("SELECT 1 {# comment #}") {
		t.Errorf("expected true for jinja comment")
	}
	if IsTemplated("SELECT * FROM my_table") {
		t.Errorf("expected false for plain sql")
	}
}

func TestSanitizeJinjaPreservesLengthAndLines(t *testing.T) {
	input := `SELECT
    {{ config(materialized='table') }},
    col_a
FROM
    {% if is_incremental() %}
        {{ ref('source_table') }}
    {% else %}
        base_table
    {% endif %}
`
	res := SanitizeJinja(input)
	if len(res.Slices) == 0 {
		t.Fatalf("expected slices to be extracted")
	}

	// Line counts must match exactly
	inputLines := strings.Count(input, "\n")
	sanitizedLines := strings.Count(res.SanitizedSQL, "\n")
	if inputLines != sanitizedLines {
		t.Errorf("expected %d lines, got %d", inputLines, sanitizedLines)
	}

	// Byte length must match exactly
	if len(input) != len(res.SanitizedSQL) {
		t.Errorf("expected length %d, got %d", len(input), len(res.SanitizedSQL))
	}
}
