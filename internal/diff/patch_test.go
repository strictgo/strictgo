package diff

import (
	"reflect"
	"testing"
)

func TestParsePatch(t *testing.T) {
	cases := []struct {
		name  string
		patch string
		want  addedLines
	}{
		{"empty", "", addedLines{}},
		{"addition", "+++ b/a.go\n@@ -0,0 +1,2 @@\n+package a\n+var value = 1\n", addedLines{"a.go": {1: true, 2: true}}},
		{"replacement", "+++ b/a.go\n@@ -3,2 +3,1 @@\n-old\n-old\n+new\n", addedLines{"a.go": {3: true}}},
		{"deletion", "+++ /dev/null\n@@ -1 +0,0 @@\n-old\n", addedLines{}},
		{"removed lines", "+++ b/a.go\n@@ -3,2 +2,0 @@\n-old\n-old\n", addedLines{}},
		{"context", "+++ b/a.go\n@@ -2,3 +2,3 @@\n context\n-old\n+new\n context\n", addedLines{"a.go": {3: true}}},
		{"multiple hunks", "+++ b/a.go\n@@ -1 +1 @@\n-old\n+new\n@@ -9,0 +10 @@\n+new\n", addedLines{"a.go": {1: true, 10: true}}},
		{"header-like source", "+++ b/a.go\n@@ -0,0 +1,2 @@\n+++ b/not-a-header.go\n+@@ not-a-hunk\n", addedLines{"a.go": {1: true, 2: true}}},
		{"spaces", "+++ b/a file.go\t\n@@ -0,0 +1 @@\n+new\n", addedLines{"a file.go": {1: true}}},
		{"quoted path", "+++ \"b/a\\tb.go\"\n@@ -0,0 +1 @@\n+new\n", addedLines{"a\tb.go": {1: true}}},
		{"unicode path", "+++ \"b/\\303\\251.go\"\n@@ -0,0 +1 @@\n+new\n", addedLines{"é.go": {1: true}}},
		{"no newline", "+++ b/a.go\n@@ -1 +1 @@\n-old\n\\ No newline at end of file\n+new\n\\ No newline at end of file\n", addedLines{"a.go": {1: true}}},
		{"non-Go file", "+++ b/README.md\n@@ -0,0 +1 @@\n+text\n", addedLines{}},
		{"rename", "diff --git a/old.go b/new.go\nsimilarity index 90%\nrename from old.go\nrename to new.go\n--- a/old.go\n+++ b/new.go\n@@ -3 +3 @@\n-old\n+new\n", addedLines{"new.go": {3: true}}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := parsePatch([]byte(test.patch))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("added lines = %v, want %v", got, test.want)
			}
		})
	}
}

func TestInvalidPatch(t *testing.T) {
	for _, patch := range []string{
		"+++ b/a.go\n@@ -0,0 +1,2 @@\n+one\n",
		"+++ b/a.go\n@@ invalid\n",
		"+++ b/a.go\n@@ -0,0 +1 @@\n context\n",
		"+++ b/../outside.go\n",
		"+++ \"b/bad\\q.go\"\n",
		"+++ unexpected.go\n",
		"@@ -0,0 +1 @@\n+new\n",
		"+++ /dev/null\n@@ -0,0 +1 @@\n+new\n",
		"+++ b/a.go\n@@ -0,0 +0 @@\n+new\n",
	} {
		if _, err := parsePatch([]byte(patch)); err != nil {
			continue
		}
		t.Errorf("accepted malformed patch %q", patch)
	}
}
