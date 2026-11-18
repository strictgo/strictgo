package declorder

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	return strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go")
}

func checkTestOrder(pass *analysis.Pass, file *ast.File) {
	var lastTest token.Pos
	for _, declaration := range file.Decls {
		if function, ok := declaration.(*ast.FuncDecl); ok && testEntrypoint(pass.TypesInfo, function) {
			lastTest = function.Pos()
		}
	}
	// Files containing only shared test helpers do not need a test function.
	for _, declaration := range file.Decls {
		if declaration.Pos() >= lastTest {
			break
		}
		switch declaration := declaration.(type) {
		case *ast.GenDecl:
			if declaration.Tok == token.TYPE {
				pass.Reportf(declaration.Pos(), "helper type declarations must appear after test entrypoints")
			}
		case *ast.FuncDecl:
			if !testEntrypoint(pass.TypesInfo, declaration) {
				pass.Reportf(declaration.Pos(), "helper functions and methods must appear after test entrypoints")
			}
		}
	}
}

func testEntrypoint(info *types.Info, function *ast.FuncDecl) bool {
	if function.Recv != nil {
		return false
	}
	object, ok := info.Defs[function.Name].(*types.Func)
	if !ok {
		return false
	}
	signature := object.Type().(*types.Signature)
	if signature.Results().Len() != 0 || signature.Variadic() || signature.TypeParams().Len() != 0 {
		return false
	}
	name := function.Name.Name
	if signature.Params().Len() == 0 {
		return testingName(name, "Example")
	}
	if signature.Params().Len() != 1 {
		return false
	}
	pointer, ok := types.Unalias(signature.Params().At(0).Type()).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(pointer.Elem()).(*types.Named)
	if !ok || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != "testing" {
		return false
	}
	switch named.Obj().Name() {
	case "T":
		return testingName(name, "Test")
	case "B":
		return testingName(name, "Benchmark")
	case "F":
		return testingName(name, "Fuzz")
	case "M":
		return name == "TestMain"
	}
	return false
}

// Match the Go testing convention: the prefix must not be followed by a
// lowercase letter. TestHelper is a test name; Testhelper is not.
func testingName(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if name == prefix {
		return true
	}
	letter, _ := utf8.DecodeRuneInString(strings.TrimPrefix(name, prefix))
	return !unicode.IsLower(letter)
}
