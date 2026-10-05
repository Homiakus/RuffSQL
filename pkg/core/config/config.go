package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// FluffConfig stores loaded configuration options
type FluffConfig struct {
	Dialect      string
	Rules        map[string]bool
	ExcludeRules map[string]bool
	IgnorePaths  []string
}

func DefaultConfig() *FluffConfig {
	return &FluffConfig{
		Dialect:      "ansi",
		Rules:        make(map[string]bool),
		ExcludeRules: make(map[string]bool),
		IgnorePaths:  []string{},
	}
}

// LoadConfig attempts to locate and load .sqlfluff or .sqlfluffignore in dir or parents
func LoadConfig(dir string) (*FluffConfig, error) {
	cfg := DefaultConfig()

	// Search upwards for .sqlfluff
	curr := dir
	for {
		configFile := filepath.Join(curr, ".sqlfluff")
		if _, err := os.Stat(configFile); err == nil {
			_ = parseConfigFile(configFile, cfg)
			break
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	// Search upwards for .sqlfluffignore
	curr = dir
	for {
		ignoreFile := filepath.Join(curr, ".sqlfluffignore")
		if _, err := os.Stat(ignoreFile); err == nil {
			_ = parseIgnoreFile(ignoreFile, cfg)
			break
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return cfg, nil
}

func parseConfigFile(path string, cfg *FluffConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentSection := ""

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			if currentSection == "sqlfluff" {
				if key == "dialect" && val != "None" && val != "" {
					cfg.Dialect = val
				} else if key == "rules" && val != "all" && val != "None" {
					for _, r := range strings.Split(val, ",") {
						cfg.Rules[strings.TrimSpace(r)] = true
					}
				} else if key == "exclude_rules" && val != "None" {
					for _, r := range strings.Split(val, ",") {
						cfg.ExcludeRules[strings.TrimSpace(r)] = true
					}
				}
			}
		}
	}

	return scanner.Err()
}

func parseIgnoreFile(path string, cfg *FluffConfig) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cfg.IgnorePaths = append(cfg.IgnorePaths, line)
	}

	return scanner.Err()
}

// IsRuleEnabled checks if a rule code should run given current config
func (c *FluffConfig) IsRuleEnabled(code string) bool {
	if c.ExcludeRules[code] {
		return false
	}
	if len(c.Rules) > 0 {
		return c.Rules[code]
	}
	return true
}

// IsPathIgnored checks if relative path matches any ignore pattern
func (c *FluffConfig) IsPathIgnored(targetPath string) bool {
	normalized := filepath.ToSlash(targetPath)
	for _, pattern := range c.IgnorePaths {
		pat := filepath.ToSlash(pattern)
		// Direct base name match
		if match, _ := filepath.Match(pat, filepath.Base(targetPath)); match {
			return true
		}
		// Full path match
		if match, _ := filepath.Match(pat, normalized); match {
			return true
		}
		// Trim trailing /* for directory matching
		trimmed := strings.TrimSuffix(pat, "/*")
		trimmed = strings.TrimSuffix(trimmed, "*")
		if trimmed != "" && strings.Contains(normalized, trimmed) {
			return true
		}
	}
	return false
}
