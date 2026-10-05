package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/config"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/dialects"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/linter"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/parser"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/runner"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Automatically load .sqlfluff and .sqlfluffignore from working directory
	cwd, _ := os.Getwd()
	cfg, _ := config.LoadConfig(cwd)

	switch command {
	case "parse":
		parseCmd := flag.NewFlagSet("parse", flag.ExitOnError)
		dialectName := parseCmd.String("dialect", cfg.Dialect, "Dialect to use")
		_ = parseCmd.Parse(os.Args[2:])

		args := parseCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: No SQL file or string provided")
			os.Exit(1)
		}

		filePath := args[0]
		data, err := os.ReadFile(filePath)
		var sql string
		if err != nil {
			sql = filePath
		} else {
			sql = string(data)
		}

		d, err := dialects.GetDialect(*dialectName)
		if err != nil {
			fmt.Printf("Dialect error: %v\n", err)
			os.Exit(1)
		}

		p := parser.NewParser(d)
		tree, err := p.ParseSQL(sql)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(tree.String(0))

	case "lint":
		lintCmd := flag.NewFlagSet("lint", flag.ExitOnError)
		dialectName := lintCmd.String("dialect", cfg.Dialect, "Dialect to use")
		format := lintCmd.String("format", "human", "Output format: human, json, github-annotation, sarif")
		workers := lintCmd.Int("workers", 8, "Number of concurrent worker threads")
		_ = lintCmd.Parse(os.Args[2:])

		args := lintCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: No SQL file or path provided")
			os.Exit(1)
		}

		// Update dialect in config if overridden by CLI flag
		if *dialectName != "" {
			cfg.Dialect = *dialectName
		}

		lintEngine, err := linter.NewLinterWithConfig(cfg)
		if err != nil {
			fmt.Printf("Linter initialization error: %v\n", err)
			os.Exit(1)
		}

		// Support reading from STDIN if '-' is passed
		if len(args) == 1 && args[0] == "-" {
			stdInData, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
			res, err := lintEngine.LintString(string(stdInData), "stdin.sql")
			if err != nil {
				fmt.Printf("Lint error: %v\n", err)
				os.Exit(1)
			}
			r := runner.NewRunner(lintEngine, runner.OutputFormat(*format), 1)
			r.RenderOutput([]runner.FileReport{{
				FilePath:   "stdin.sql",
				Violations: res.Violations,
				Passed:     len(res.Violations) == 0,
			}})
			if len(res.Violations) > 0 {
				os.Exit(1)
			}
			return
		}

		r := runner.NewRunner(lintEngine, runner.OutputFormat(*format), *workers)
		_, hasFailures := r.RunParallel(args)

		if hasFailures {
			if *format == "human" {
				fmt.Println("\nAll Finished 📜 💥!")
			}
			os.Exit(1)
		} else {
			if *format == "human" {
				fmt.Println("\nAll Finished 📜 🎉!")
			}
		}

	case "fix":
		fixCmd := flag.NewFlagSet("fix", flag.ExitOnError)
		dialectName := fixCmd.String("dialect", cfg.Dialect, "Dialect to use")
		inPlace := fixCmd.Bool("write", false, "Write fixed output back to file")
		showDiff := fixCmd.Bool("diff", false, "Show unified diff between original and fixed SQL")
		_ = fixCmd.Parse(os.Args[2:])

		args := fixCmd.Args()
		if len(args) == 0 {
			fmt.Println("Error: No SQL file or path provided")
			os.Exit(1)
		}

		if *dialectName != "" {
			cfg.Dialect = *dialectName
		}

		lintEngine, err := linter.NewLinterWithConfig(cfg)
		if err != nil {
			fmt.Printf("Linter initialization error: %v\n", err)
			os.Exit(1)
		}

		for _, target := range args {
			data, err := os.ReadFile(target)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", target, err)
				continue
			}
			origSQL := string(data)
			fixedSQL, count, err := lintEngine.FixString(origSQL, 5)
			if err != nil {
				fmt.Printf("Error fixing %s: %v\n", target, err)
				continue
			}

			if count > 0 {
				fmt.Printf("Fixed %d violations in %s\n", count, target)
				if *showDiff {
					fmt.Printf("--- %s (original)\n+++ %s (fixed)\n", target, target)
					origLines := strings.Split(origSQL, "\n")
					fixedLines := strings.Split(fixedSQL, "\n")
					for i := 0; i < len(origLines) || i < len(fixedLines); i++ {
						if i < len(origLines) && i < len(fixedLines) && origLines[i] != fixedLines[i] {
							fmt.Printf("- %s\n+ %s\n", origLines[i], fixedLines[i])
						} else if i >= len(origLines) && i < len(fixedLines) {
							fmt.Printf("+ %s\n", fixedLines[i])
						}
					}
				} else if *inPlace {
					_ = os.WriteFile(target, []byte(fixedSQL), 0644)
					fmt.Printf("Changes written to %s\n", target)
				} else {
					fmt.Println("--- Result ---")
					fmt.Print(fixedSQL)
				}
			} else {
				fmt.Printf("No fixable violations in %s\n", target)
			}
		}

	case "version":
		fmt.Println("ruffsql, version 0.5.0 (High Performance SQL Linter + Multi-Dialect + Diff & Stdin)")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("ruffsql - High Performance SQL Linter and Fixer in Go")
	fmt.Println("\nUsage:")
	fmt.Println("  ruffsql parse [--dialect <name>] <file.sql>")
	fmt.Println("  ruffsql lint [--dialect <name>] [--format human|json|github-annotation|sarif] [--workers 8] <path... | ->")
	fmt.Println("  ruffsql fix [--dialect <name>] [--write] [--diff] <path...>")
	fmt.Println("  ruffsql version")
}
