package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/linter"
)

func TestRunnerParallel(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "query1.sql")
	_ = os.WriteFile(file1, []byte("SELECT a, b FROM t;\n"), 0644)

	file2 := filepath.Join(tmpDir, "query2.sql")
	_ = os.WriteFile(file2, []byte("select a from t;\n"), 0644) // lowercase select -> CP01

	l, err := linter.NewLinter("ansi")
	if err != nil {
		t.Fatalf("failed to create linter: %v", err)
	}

	r := NewRunner(l, FormatHuman, 2)
	reports, hasFailures := r.RunParallel([]string{tmpDir})

	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}

	if !hasFailures {
		t.Errorf("expected failures due to lowercase select in query2.sql")
	}
}

func TestRunnerRenderFormats(t *testing.T) {
	l, _ := linter.NewLinter("ansi")

	rJSON := NewRunner(l, FormatJSON, 1)
	rJSON.RenderOutput([]FileReport{
		{FilePath: "test.sql", Passed: true},
	})

	rGH := NewRunner(l, FormatGitHub, 1)
	rGH.RenderOutput([]FileReport{
		{FilePath: "test.sql", Passed: true},
	})
}
