package dialects

import (
	"testing"
)

func TestSupportedDialects(t *testing.T) {
	if len(SupportedDialects) != 28 {
		t.Fatalf("expected 28 dialects, got %d", len(SupportedDialects))
	}

	for _, name := range SupportedDialects {
		d, err := GetDialect(name)
		if err != nil {
			t.Errorf("failed to get dialect %s: %v", name, err)
		}
		if d == nil {
			t.Errorf("dialect %s is nil", name)
		}
		if d.GetName() != name {
			t.Errorf("expected dialect name %s, got %s", name, d.GetName())
		}
	}
}

func TestUnknownDialect(t *testing.T) {
	_, err := GetDialect("nonexistent_dialect_xyz")
	if err == nil {
		t.Errorf("expected error for nonexistent dialect")
	}
}

func TestPostgresSpecialSegments(t *testing.T) {
	pg, err := GetDialect("postgres")
	if err != nil {
		t.Fatalf("unexpected error getting postgres dialect: %v", err)
	}
	if pg.GetSegment("CastOperator") == nil {
		t.Errorf("expected CastOperator in postgres dialect")
	}
	if pg.GetSegment("JsonOperator") == nil {
		t.Errorf("expected JsonOperator in postgres dialect")
	}
}
