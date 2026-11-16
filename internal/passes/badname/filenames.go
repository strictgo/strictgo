package badname

import (
	"go/ast"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func checkPackageAndFileNames(pass *analysis.Pass) {
	var checkedPackage bool
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}
		if !checkedPackage {
			name := strings.TrimSuffix(file.Name.Name, "_test")
			if vagueName(name) {
				pass.Reportf(file.Name.Pos(), "package name %q is too generic; use a name that describes its purpose", file.Name.Name)
			}
			checkedPackage = true
		}
		// Use the physical filename, not a filename supplied by a //line directive.
		filename := filepath.Base(pass.Fset.File(file.Pos()).Name())
		stem := strings.TrimSuffix(filename, ".go")
		for _, part := range strings.Split(stem, "_") {
			if vagueName(part) {
				pass.Reportf(file.Name.Pos(), "Go file name %q is too generic; use a name that describes its purpose", filename)
				break
			}
		}
	}
}

func vagueName(name string) bool {
	name = strings.ToLower(name)
	for _, suffix := range []string{"util", "utils", "helper", "helpers", "misc"} {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}
