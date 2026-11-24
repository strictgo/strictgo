package diff

import (
	"bytes"
	"errors"
	"flag"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestOptions(t *testing.T) {
	var output bytes.Buffer
	settings, err := parseOptions(nil, &output)
	if err != nil {
		t.Fatal(err)
	}
	if settings.base != "" || !reflect.DeepEqual(settings.packages, []string{"./..."}) {
		t.Errorf("defaults = %+v", settings)
	}
	settings, err = parseOptions([]string{"--base", "main", "./cmd/..."}, &output)
	if err != nil || settings.base != "main" || !reflect.DeepEqual(settings.packages, []string{"./cmd/..."}) {
		t.Errorf("explicit options = %+v, %v", settings, err)
	}
	for _, args := range [][]string{{"--base="}, {"--base"}, {"-fix"}, {"-json=false"}, {"main.go"}, {"./...", "-fix"}} {
		if _, err := parseOptions(args, &output); err != nil {
			continue
		}
		t.Errorf("accepted invalid arguments %v", args)
	}
	if _, err := parseOptions([]string{"-h"}, &output); !errors.Is(err, flag.ErrHelp) || !strings.Contains(output.String(), "Usage:") {
		t.Errorf("help = %v, %q", err, output.String())
	}
}

func TestSnapshotEnvironment(t *testing.T) {
	t.Setenv("GOWORK", "/original/go.work")
	t.Setenv("GOPACKAGESDRIVER", "/original/driver")
	t.Setenv("GOFLAGS", "-tags=example")
	environment, err := snapshotEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(environment, "\n")
	if strings.Contains(joined, "/original/") || !strings.Contains(joined, "GOWORK=off") || !strings.Contains(joined, "GOPACKAGESDRIVER=off") || !strings.Contains(joined, "GOFLAGS=-tags=example") {
		t.Errorf("unsafe or incomplete environment: %s", joined)
	}
	for _, value := range []string{"-overlay=/original/overlay.json", "-modfile other.mod", `"-overlay=/a path/overlay.json"`, "-C=/original"} {
		t.Setenv("GOFLAGS", value)
		if _, err := snapshotEnvironment(); err != nil {
			continue
		}
		t.Errorf("accepted unsafe GOFLAGS %q", value)
	}
}

func TestSnapshotPatterns(t *testing.T) {
	root := t.TempDir()
	cwd := filepath.Join(root, "module")
	tree := filepath.Join(t.TempDir(), "source")
	got, err := snapshotPatterns([]string{"./...", "../other", filepath.Join(root, "module", "pkg"), "example.com/module/..."}, root, cwd, tree)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{filepath.Join(tree, "module", "..."), filepath.Join(tree, "other"), filepath.Join(tree, "module", "pkg"), "example.com/module/..."}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("patterns = %v, want %v", got, want)
	}
	if _, err := snapshotPatterns([]string{"../../outside"}, root, cwd, tree); err != nil {
		return
	}
	t.Error("accepted package outside the snapshot")
}
