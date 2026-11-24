package diff

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFilterDiagnostics(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "sub")
	primary := []driverDiagnostic{
		{Position: filepath.Join(root, "a:b.go") + ":3:2", Message: "new"},
		{Position: filepath.Join(root, "a:b.go") + ":2:2", Message: "old"},
		{Position: "../relative.go:4:5", Message: "relative"},
		{Position: filepath.Join(filepath.Dir(root), "outside.go") + ":3:2", Message: "outside"},
	}
	tree := map[string]map[string]any{
		"p":          {"badname": primary},
		"p [p.test]": {"badname": primary},
		"related": {"badname": []map[string]any{{
			"posn": filepath.Join(root, "a:b.go") + ":1:1", "message": "unchanged primary",
			"related": []map[string]string{{"posn": filepath.Join(root, "a:b.go") + ":3:2"}},
		}}},
	}
	added := addedLines{"a:b.go": {3: true}, "relative.go": {4: true}}
	got, err := filterDiagnostics(jsonBytes(t, tree), added, root, cwd)
	if err != nil {
		t.Fatal(err)
	}
	want := []finding{
		{filename: "a:b.go", line: 3, column: 2, analyzer: "badname", message: "new"},
		{filename: "relative.go", line: 4, column: 5, analyzer: "badname", message: "relative"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findings = %#v, want %#v", got, want)
	}
}

func TestAnalysisFailuresAreNotFiltered(t *testing.T) {
	for _, text := range []string{
		`{"p":{"badname":{"error":"analysis skipped due to errors"}}}`,
		`{"p":{"badname":[{"posn":"-","message":"invalid primary"}]}}`,
		`{"p":{"badname":{"unexpected":"result"}}}`,
		`{"p":null}`,
		`null`,
		`[]`,
		`not JSON`,
	} {
		if _, err := filterDiagnostics([]byte(text), addedLines{}, t.TempDir(), t.TempDir()); err != nil {
			continue
		}
		t.Errorf("accepted failed or invalid analysis: %s", text)
	}
	if got, err := filterDiagnostics([]byte(`{}`), addedLines{}, t.TempDir(), t.TempDir()); err != nil || len(got) != 0 {
		t.Errorf("empty analysis = %v, %v", got, err)
	}
}

func TestParsePosition(t *testing.T) {
	got, err := parsePosition(`C:\repo\a.go:12:4`)
	if err != nil {
		t.Fatal(err)
	}
	want := location{filename: `C:\repo\a.go`, line: 12, column: 4}
	if got != want {
		t.Errorf("position = %v, want %v", got, want)
	}
	for _, text := range []string{"file.go", "file.go:1", "file.go:0:1", "file.go:1:0", "file.go:x:1"} {
		if _, err := parsePosition(text); err != nil {
			continue
		}
		t.Errorf("accepted invalid position %q", text)
	}
}

func jsonBytes(tb testing.TB, value any) []byte {
	tb.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		tb.Fatal(err)
	}
	return data
}
