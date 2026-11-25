// Package reportingrules applies shared reporting rules to strictgo analyzers.
package reportingrules

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// IgnoreGenerated suppresses generated-file diagnostics and unsafe fixes. It
// returns the original analyzer pointers, preserving dependency/fact identity.
// Install the policy once at registration, before any analysis starts.
func IgnoreGenerated(analyzers ...*analysis.Analyzer) []*analysis.Analyzer {
	visited := make(map[*analysis.Analyzer]bool)
	for _, analyzer := range analyzers {
		installGeneratedFilter(analyzer, visited)
	}
	return analyzers
}

func installGeneratedFilter(analyzer *analysis.Analyzer, visited map[*analysis.Analyzer]bool) {
	if analyzer == nil || visited[analyzer] {
		return
	}
	visited[analyzer] = true
	for _, dependency := range analyzer.Requires {
		installGeneratedFilter(dependency, visited)
	}
	if analyzer.Run == nil {
		return
	}
	execute := analyzer.Run
	analyzer.Run = func(pass *analysis.Pass) (any, error) {
		files := generatedFiles(pass)
		filtered := *pass
		filtered.Report = func(diagnostic analysis.Diagnostic) {
			// Physical token files cannot be disguised by //line directives.
			if files[pass.Fset.File(diagnostic.Pos)] {
				return
			}
			if len(diagnostic.SuggestedFixes) > 0 {
				var fixes []analysis.SuggestedFix
				for _, fix := range diagnostic.SuggestedFixes {
					if safeFix(pass.Fset, files, fix) {
						fixes = append(fixes, fix)
					}
				}
				diagnostic.SuggestedFixes = fixes
			}
			pass.Report(diagnostic)
		}
		// Keep Files, TypesInfo, ResultOf, facts, results, and errors intact:
		// generated code can still affect analysis of handwritten declarations.
		return execute(&filtered)
	}
}

func generatedFiles(pass *analysis.Pass) map[*token.File]bool {
	files := make(map[*token.File]bool)
	for _, file := range pass.Files {
		files[pass.Fset.File(file.Pos())] = ast.IsGenerated(file)
	}
	return files
}

func safeFix(files *token.FileSet, generated map[*token.File]bool, fix analysis.SuggestedFix) bool {
	for _, edit := range fix.TextEdits {
		file := files.File(edit.Pos)
		isGenerated, known := generated[file]
		// Unknown files cannot be verified, so withhold that fix as well. Never
		// partially apply a fix that also needs changes to generated code.
		if !known || isGenerated || (edit.End.IsValid() && files.File(edit.End) != file) {
			return false
		}
	}
	return true
}
