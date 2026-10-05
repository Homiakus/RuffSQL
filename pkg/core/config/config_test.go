package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Dialect != "ansi" {
		t.Errorf("expected ansi dialect, got %s", cfg.Dialect)
	}
	if !cfg.IsRuleEnabled("CP01") {
		t.Errorf("expected rule CP01 to be enabled by default")
	}
}

func TestParseConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, ".sqlfluff")
	content := `
[sqlfluff]
dialect = postgres
rules = CP01, LT01, CV01
exclude_rules = LT01
`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Dialect != "postgres" {
		t.Errorf("expected dialect postgres, got %s", cfg.Dialect)
	}
	if !cfg.IsRuleEnabled("CP01") {
		t.Errorf("expected CP01 to be enabled")
	}
	if cfg.IsRuleEnabled("LT01") {
		t.Errorf("expected LT01 to be disabled via exclude_rules")
	}
	if cfg.IsRuleEnabled("CV02") {
		t.Errorf("expected CV02 to be disabled because not in rules list")
	}
}

func TestParseIgnoreFile(t *testing.T) {
	tmpDir := t.TempDir()
	ignorePath := filepath.Join(tmpDir, ".sqlfluffignore")
	content := `
# ignore migrations
migrations/*
ignored.sql
`
	if err := os.WriteFile(ignorePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write ignore: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if !cfg.IsPathIgnored("ignored.sql") {
		t.Errorf("expected ignored.sql to be ignored")
	}
	if !cfg.IsPathIgnored("/some/path/migrations/001_init.sql") {
		t.Errorf("expected migration file to be ignored")
	}
	if cfg.IsPathIgnored("valid.sql") {
		t.Errorf("expected valid.sql NOT to be ignored")
	}
}
