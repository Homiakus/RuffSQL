package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sqlfluff/sqlfluff-go/pkg/core/linter"
	"github.com/sqlfluff/sqlfluff-go/pkg/core/rules"
)

type OutputFormat string

const (
	FormatHuman  OutputFormat = "human"
	FormatJSON   OutputFormat = "json"
	FormatGitHub OutputFormat = "github-annotation"
	FormatSARIF  OutputFormat = "sarif"
)

type FileReport struct {
	FilePath   string                `json:"filepath"`
	Violations []rules.LintViolation `json:"violations"`
	Passed     bool                  `json:"passed"`
}

type Runner struct {
	Linter     *linter.Linter
	Format     OutputFormat
	NumWorkers int
}

func NewRunner(l *linter.Linter, format OutputFormat, numWorkers int) *Runner {
	if numWorkers <= 0 {
		numWorkers = 8
	}
	return &Runner{
		Linter:     l,
		Format:     format,
		NumWorkers: numWorkers,
	}
}

// RunParallel lints multiple files concurrently using a worker pool
func (r *Runner) RunParallel(paths []string) ([]FileReport, bool) {
	// Discover all sql files if directory is passed
	var allFiles []string
	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			allFiles = append(allFiles, p)
			continue
		}
		if info.IsDir() {
			_ = filepath.Walk(p, func(path string, fi os.FileInfo, err error) error {
				if err == nil && !fi.IsDir() && filepath.Ext(path) == ".sql" {
					allFiles = append(allFiles, path)
				}
				return nil
			})
		} else {
			allFiles = append(allFiles, p)
		}
	}

	filesChan := make(chan string, len(allFiles))
	reportsChan := make(chan FileReport, len(allFiles))

	var wg sync.WaitGroup
	for w := 0; w < r.NumWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				res, err := r.Linter.LintFile(file)
				if err != nil {
					reportsChan <- FileReport{FilePath: file, Passed: false}
					continue
				}
				reportsChan <- FileReport{
					FilePath:   file,
					Violations: res.Violations,
					Passed:     len(res.Violations) == 0,
				}
			}
		}()
	}

	for _, f := range allFiles {
		filesChan <- f
	}
	close(filesChan)

	wg.Wait()
	close(reportsChan)

	var reports []FileReport
	hasFailures := false
	for rep := range reportsChan {
		reports = append(reports, rep)
		if !rep.Passed {
			hasFailures = true
		}
	}

	r.RenderOutput(reports)
	return reports, hasFailures
}

func (r *Runner) RenderOutput(reports []FileReport) {
	switch r.Format {
	case FormatJSON:
		data, _ := json.MarshalIndent(reports, "", "  ")
		fmt.Println(string(data))

	case FormatGitHub:
		// GitHub Actions format: ::warning file={name},line={line},col={col}::{message}
		for _, rep := range reports {
			for _, v := range rep.Violations {
				fmt.Printf("::error file=%s,line=%d,col=%d::[%s] %s\n", rep.FilePath, v.Line, v.Pos, v.RuleCode, v.Description)
			}
		}

	case FormatSARIF:
		sarifOutput := map[string]interface{}{
			"version": "2.1.0",
			"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
			"runs": []map[string]interface{}{
				{
					"tool": map[string]interface{}{
						"driver": map[string]interface{}{
							"name":           "ruffsql",
							"version":        "0.5.0",
							"informationUri": "https://github.com/sqlfluff/sqlfluff-go",
						},
					},
					"results": func() []map[string]interface{} {
						var res []map[string]interface{}
						for _, rep := range reports {
							for _, v := range rep.Violations {
								res = append(res, map[string]interface{}{
									"ruleId": v.RuleCode,
									"message": map[string]interface{}{
										"text": v.Description,
									},
									"locations": []map[string]interface{}{
										{
											"physicalLocation": map[string]interface{}{
												"artifactLocation": map[string]interface{}{
													"uri": rep.FilePath,
												},
												"region": map[string]interface{}{
													"startLine":   v.Line,
													"startColumn": v.Pos,
												},
											},
										},
									},
								})
							}
						}
						return res
					}(),
				},
			},
		}
		data, _ := json.MarshalIndent(sarifOutput, "", "  ")
		fmt.Println(string(data))

	default: // Human
		for _, rep := range reports {
			if len(rep.Violations) > 0 {
				fmt.Printf("== [%s] FAIL\n", rep.FilePath)
				for _, v := range rep.Violations {
					fmt.Printf("L: %3d | P: %3d | %-4s | %s\n", v.Line, v.Pos, v.RuleCode, v.Description)
				}
			} else {
				fmt.Printf("== [%s] PASS\n", rep.FilePath)
			}
		}
	}
}
