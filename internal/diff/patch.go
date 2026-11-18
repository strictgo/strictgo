package diff

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var hunkHeader = regexp.MustCompile(`^@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)

type addedLines map[string]map[int]bool

// Consume hunk bodies rather than assuming every line in a hunk was added.
// This also distinguishes +++ source lines from actual file headers.
func parsePatch(patch []byte) (addedLines, error) {
	added := make(addedLines)
	var filename string
	var haveHeader bool
	var oldRemaining, newRemaining, nextLine int
	for _, line := range bytes.Split(patch, []byte{'\n'}) {
		if bytes.HasPrefix(line, []byte(`\ No newline at end of file`)) {
			continue
		}
		if oldRemaining > 0 || newRemaining > 0 {
			if len(line) == 0 {
				return nil, fmt.Errorf("truncated diff hunk for %q", filename)
			}
			switch line[0] {
			case '+':
				if strings.HasSuffix(filename, ".go") {
					if added[filename] == nil {
						added[filename] = make(map[int]bool)
					}
					added[filename][nextLine] = true
				}
				newRemaining--
				nextLine++
			case '-':
				oldRemaining--
			case ' ':
				oldRemaining--
				newRemaining--
				nextLine++
			default:
				return nil, fmt.Errorf("invalid diff hunk for %q", filename)
			}
			if oldRemaining < 0 || newRemaining < 0 {
				return nil, fmt.Errorf("invalid diff line counts for %q", filename)
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("diff --git ")) {
			filename = ""
			haveHeader = false
			continue
		}
		if bytes.HasPrefix(line, []byte("+++ ")) {
			var err error
			filename, err = patchFilename(string(line[4:]))
			if err != nil {
				return nil, err
			}
			haveHeader = true
			continue
		}
		if bytes.HasPrefix(line, []byte("@@")) {
			if !haveHeader {
				return nil, fmt.Errorf("diff hunk has no file header")
			}
			matches := hunkHeader.FindSubmatch(line)
			if matches == nil {
				return nil, fmt.Errorf("invalid diff hunk header: %s", line)
			}
			var err error
			oldRemaining, err = hunkCount(matches[2])
			if err != nil {
				return nil, err
			}
			newRemaining, err = hunkCount(matches[4])
			if err != nil {
				return nil, err
			}
			nextLine, err = strconv.Atoi(string(matches[3]))
			if err != nil {
				return nil, fmt.Errorf("invalid diff line number: %w", err)
			}
			if newRemaining > 0 && (filename == "" || nextLine == 0) {
				return nil, fmt.Errorf("invalid added-line location in diff")
			}
		}
	}
	if oldRemaining != 0 || newRemaining != 0 {
		return nil, fmt.Errorf("truncated diff hunk for %q", filename)
	}
	return added, nil
}

func hunkCount(text []byte) (int, error) {
	if len(text) == 0 {
		return 1, nil
	}
	count, err := strconv.Atoi(string(text))
	if err != nil {
		return 0, fmt.Errorf("invalid diff line count: %w", err)
	}
	return count, nil
}

func patchFilename(header string) (string, error) {
	// Git adds a tab separator to some unquoted headers containing spaces.
	name := strings.TrimSuffix(header, "\t")
	if strings.HasPrefix(name, `"`) {
		var err error
		name, err = strconv.Unquote(name)
		if err != nil {
			return "", fmt.Errorf("invalid quoted diff path: %w", err)
		}
	}
	if name == "/dev/null" {
		return "", nil
	}
	if !strings.HasPrefix(name, "b/") {
		return "", fmt.Errorf("unexpected diff path %q", name)
	}
	name = strings.TrimPrefix(name, "b/")
	if path.IsAbs(name) || path.Clean(name) != name || name == "." || name == ".." || strings.HasPrefix(name, "../") {
		return "", fmt.Errorf("invalid repository-relative diff path %q", name)
	}
	return name, nil
}
