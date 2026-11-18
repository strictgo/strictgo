package diff

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type driverDiagnostic struct {
	Position string `json:"posn"`
	Message  string `json:"message"`
}

type location struct {
	filename string
	line     int
	column   int
}

type finding struct {
	location
	analyzer string
	message  string
}

func filterDiagnostics(data []byte, added addedLines, root, cwd string) ([]finding, error) {
	var tree map[string]map[string]json.RawMessage
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("invalid analyzer JSON: %w", err)
	}
	if tree == nil {
		return nil, errors.New("missing analyzer JSON results")
	}
	seen := make(map[finding]bool)
	failures := make(map[string]bool)
	var findings []finding
	for pkg, actions := range tree {
		if actions == nil {
			return nil, fmt.Errorf("missing analyzer results for %s", pkg)
		}
		for analyzer, raw := range actions {
			if !bytes.HasPrefix(bytes.TrimSpace(raw), []byte("[")) {
				var failure struct {
					Message string `json:"error"`
				}
				if err := json.Unmarshal(raw, &failure); err != nil || failure.Message == "" {
					return nil, fmt.Errorf("invalid result for %s: %s", pkg, analyzer)
				}
				failures[fmt.Sprintf("%s: %s: %s", pkg, analyzer, failure.Message)] = true
				continue
			}
			var diagnostics []driverDiagnostic
			if err := json.Unmarshal(raw, &diagnostics); err != nil {
				return nil, fmt.Errorf("invalid diagnostics for %s: %s: %w", pkg, analyzer, err)
			}
			for _, diagnostic := range diagnostics {
				position, err := parsePosition(diagnostic.Position)
				if err != nil {
					return nil, err
				}
				filename, ok := repositoryPath(position.filename, root, cwd)
				if !ok || !added[filename][position.line] {
					continue
				}
				position.filename = filename
				item := finding{location: position, analyzer: analyzer, message: diagnostic.Message}
				if !seen[item] {
					seen[item] = true
					findings = append(findings, item)
				}
			}
		}
	}
	if len(failures) > 0 {
		var messages []string
		for message := range failures {
			messages = append(messages, message)
		}
		slices.Sort(messages)
		return nil, errors.New("analysis failed:\n" + strings.Join(messages, "\n"))
	}
	slices.SortFunc(findings, func(left, right finding) int {
		return cmp.Or(strings.Compare(left.filename, right.filename), cmp.Compare(left.line, right.line), cmp.Compare(left.column, right.column), strings.Compare(left.analyzer, right.analyzer), strings.Compare(left.message, right.message))
	})
	return findings, nil
}

func parsePosition(text string) (location, error) {
	columnSeparator := strings.LastIndexByte(text, ':')
	if columnSeparator < 0 {
		return location{}, fmt.Errorf("invalid diagnostic position %q", text)
	}
	lineSeparator := strings.LastIndexByte(text[:columnSeparator], ':')
	if lineSeparator <= 0 {
		return location{}, fmt.Errorf("invalid diagnostic position %q", text)
	}
	line, err := strconv.Atoi(text[lineSeparator+1 : columnSeparator])
	if err != nil || line <= 0 {
		return location{}, fmt.Errorf("invalid diagnostic line in %q", text)
	}
	column, err := strconv.Atoi(text[columnSeparator+1:])
	if err != nil || column <= 0 {
		return location{}, fmt.Errorf("invalid diagnostic column in %q", text)
	}
	return location{filename: text[:lineSeparator], line: line, column: column}, nil
}

func repositoryPath(filename, root, cwd string) (string, bool) {
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(cwd, filename)
	}
	relative, err := filepath.Rel(root, filename)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.ToSlash(relative), true
}

func writeFindings(output io.Writer, findings []finding) error {
	for _, item := range findings {
		if _, err := fmt.Fprintf(output, "%s:%d:%d: [%s] %s\n", item.filename, item.line, item.column, item.analyzer, item.message); err != nil {
			return err
		}
	}
	return nil
}
