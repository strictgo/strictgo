// Package diff runs strictgo on committed HEAD and reports added-line diagnostics.
package diff

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const usage = `Usage: strictgo diff [--base origin/main] [packages...]

Analyze committed HEAD and report diagnostics on added Go lines.
Packages default to ./.... Uncommitted changes are excluded.
`

type options struct {
	base     string
	packages []string
}

// Run implements the diff subcommand. Exit codes match the analysis driver:
// 0 for no relevant findings, 3 for findings, 1 for failures, 2 for bad arguments.
func Run(args []string, stdout, stderr io.Writer) int {
	settings, err := parseOptions(args, stderr)
	if errors.Is(err, flag.ErrHelp) {
		return 0
	}
	if err != nil {
		return reportError(stderr, err, 2)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return reportError(stderr, err, 1)
	}
	executable, err := os.Executable()
	if err != nil {
		return reportError(stderr, err, 1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	code, err := check(ctx, cwd, executable, settings, stdout, stderr)
	if err != nil {
		return reportError(stderr, err, 1)
	}
	return code
}

func reportError(output io.Writer, err error, code int) int {
	if _, err := fmt.Fprintln(output, "strictgo diff:", err); err != nil {
		return 1
	}
	return code
}

func parseOptions(args []string, output io.Writer) (options, error) {
	var settings options
	flags := flag.NewFlagSet("strictgo diff", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var usageFailed bool
	flags.Usage = func() {
		if _, err := io.WriteString(output, usage); err != nil {
			usageFailed = true
		}
	}
	flags.StringVar(&settings.base, "base", "origin/main", "base revision")
	err := flags.Parse(args)
	if usageFailed {
		return options{}, errors.New("could not write diff usage")
	}
	if err != nil {
		return options{}, err
	}
	if settings.base == "" {
		return options{}, errors.New("--base must not be empty")
	}
	settings.packages = flags.Args()
	if len(settings.packages) == 0 {
		settings.packages = []string{"./..."}
	}
	for _, pattern := range settings.packages {
		if strings.HasPrefix(pattern, "-") || strings.HasSuffix(pattern, ".go") {
			return options{}, fmt.Errorf("expected a package pattern, got %q; diff mode supports only --base and package patterns", pattern)
		}
	}
	return settings, nil
}

func check(ctx context.Context, cwd, executable string, settings options, stdout, stderr io.Writer) (code int, err error) {
	rootOutput, err := git(ctx, cwd, "rev-parse", "--show-toplevel")
	if err != nil {
		return 0, err
	}
	root, err := filepath.EvalSymlinks(strings.TrimSuffix(string(rootOutput), "\n"))
	if err != nil {
		return 0, err
	}
	cwd, err = filepath.EvalSymlinks(cwd)
	if err != nil {
		return 0, err
	}
	head, err := revision(ctx, root, "HEAD")
	if err != nil {
		return 0, err
	}
	base, err := revision(ctx, root, settings.base)
	if err != nil {
		return 0, err
	}
	mergeBase, err := git(ctx, root, "merge-base", base, head)
	if err != nil {
		return 0, err
	}
	patch, err := git(ctx, root, "-c", "core.quotePath=true", "diff", "--no-ext-diff", "--no-textconv", "--no-color", "--find-renames", "--no-relative", "--src-prefix=a/", "--dst-prefix=b/", "--unified=0", "--inter-hunk-context=0", strings.TrimSpace(string(mergeBase)), head, "--")
	if err != nil {
		return 0, err
	}
	added, err := parsePatch(patch)
	if err != nil {
		return 0, err
	}
	environment, err := snapshotEnvironment()
	if err != nil {
		return 0, err
	}
	temporary, err := os.MkdirTemp("", "strictgo-diff-*")
	if err != nil {
		return 0, err
	}
	tree := filepath.Join(temporary, "source")
	var attached bool
	defer func() { err = errors.Join(err, cleanup(root, tree, temporary, attached)) }()
	if _, err := git(ctx, root, "worktree", "add", "--detach", "--quiet", tree, head); err != nil {
		return 0, err
	}
	attached = true
	patterns, err := snapshotPatterns(settings.packages, root, cwd, tree)
	if err != nil {
		return 0, err
	}
	relative, err := filepath.Rel(root, cwd)
	if err != nil {
		return 0, err
	}
	directory := filepath.Join(tree, relative)
	output, messages, err := command(ctx, directory, environment, executable, append([]string{"-json"}, patterns...)...)
	if len(messages) > 0 {
		messages = bytes.ReplaceAll(messages, []byte(tree), []byte(root))
		if _, err := stderr.Write(messages); err != nil {
			return 0, err
		}
	}
	if err != nil {
		return 0, fmt.Errorf("analyzer subprocess failed: %w", err)
	}
	status, err := git(ctx, tree, "status", "--porcelain=v1", "--untracked-files=all", "-z")
	if err != nil {
		return 0, err
	}
	if len(status) != 0 {
		return 0, fmt.Errorf("analysis modified the committed snapshot: %q", status)
	}
	findings, err := filterDiagnostics(output, added, tree, directory)
	if err != nil {
		return 0, errors.New(strings.ReplaceAll(err.Error(), tree, root))
	}
	if err := writeFindings(stdout, findings); err != nil {
		return 0, err
	}
	if len(findings) > 0 {
		return 3, nil
	}
	return 0, nil
}

func revision(ctx context.Context, root, name string) (string, error) {
	output, err := git(ctx, root, "rev-parse", "--verify", "--end-of-options", name+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("cannot resolve revision %q: %w", name, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func git(ctx context.Context, directory string, args ...string) ([]byte, error) {
	args = append([]string{"--no-pager", "-c", "core.hooksPath=" + os.DevNull}, args...)
	output, messages, err := command(ctx, directory, nil, "git", args...)
	if err != nil {
		return nil, fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, messages)
	}
	return output, nil
}

func command(ctx context.Context, directory string, environment []string, executable string, args ...string) ([]byte, []byte, error) {
	process := exec.CommandContext(ctx, executable, args...)
	process.Dir = directory
	if environment != nil {
		process.Env = environment
	}
	var output, messages bytes.Buffer
	process.Stdout = &output
	process.Stderr = &messages
	err := process.Run()
	return output.Bytes(), messages.Bytes(), err
}

func cleanup(root, tree, temporary string, attached bool) (err error) {
	if attached {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err = git(ctx, root, "worktree", "remove", "--force", tree)
	}
	return errors.Join(err, os.RemoveAll(temporary))
}

// Ignore host workspaces and custom package drivers that could load the dirty
// checkout instead of the snapshot. Source overlays and alternate modfiles are
// rejected; other GOFLAGS (for example build tags) retain their usual behavior.
func snapshotEnvironment() ([]string, error) {
	var environment []string
	for _, entry := range os.Environ() {
		name, value, ok := strings.Cut(entry, "=")
		if !ok || name == "GOWORK" || name == "GOPACKAGESDRIVER" {
			continue
		}
		if name == "GOFLAGS" {
			for _, field := range strings.Fields(value) {
				field = strings.Trim(field, "\"'")
				for _, unsafe := range []string{"-overlay", "-modfile", "-C"} {
					if field == unsafe || strings.HasPrefix(field, unsafe+"=") {
						return nil, fmt.Errorf("diff mode does not support %s in GOFLAGS", unsafe)
					}
				}
			}
		}
		environment = append(environment, entry)
	}
	return append(environment, "GOWORK=off", "GOPACKAGESDRIVER=off"), nil
}

func snapshotPatterns(patterns []string, root, cwd, tree string) ([]string, error) {
	var mapped []string
	for _, pattern := range patterns {
		if filepath.IsAbs(pattern) || strings.HasPrefix(pattern, ".") {
			absolute := pattern
			if !filepath.IsAbs(absolute) {
				absolute = filepath.Join(cwd, absolute)
			}
			relative, ok := repositoryPath(absolute, root, cwd)
			if !ok {
				return nil, fmt.Errorf("package pattern %q is outside the Git repository", pattern)
			}
			pattern = filepath.Join(tree, filepath.FromSlash(relative))
		}
		mapped = append(mapped, pattern)
	}
	return mapped, nil
}
