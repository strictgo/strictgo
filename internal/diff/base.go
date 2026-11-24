package diff

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// Read only local refs: choosing a base must not trigger network access or
// accidentally use a feature branch's tracking upstream as the default branch.
func defaultBase(ctx context.Context, root string) (string, error) {
	output, err := git(ctx, root, "for-each-ref", "--format=%(refname)%09%(symref)", "refs/remotes", "refs/heads")
	if err != nil {
		return "", err
	}
	refs := make(map[string]string)
	for line := range strings.SplitSeq(strings.TrimSuffix(string(output), "\n"), "\n") {
		if line == "" {
			continue
		}
		name, target, ok := strings.Cut(line, "\t")
		if !ok {
			return "", fmt.Errorf("invalid Git reference metadata: %q", line)
		}
		refs[name] = target
	}
	if err := readRemoteHeads(ctx, root, refs); err != nil {
		return "", err
	}
	return chooseDefaultBase(refs)
}

// for-each-ref omits dangling symbolic refs. Read configured remote HEADs
// explicitly as well, so a missing default branch fails instead of guessing.
func readRemoteHeads(ctx context.Context, root string, refs map[string]string) error {
	output, err := git(ctx, root, "remote")
	if err != nil {
		return err
	}
	remotes := []string{"origin", "upstream"}
	for remote := range strings.FieldsSeq(string(output)) {
		if !slices.Contains(remotes, remote) {
			remotes = append(remotes, remote)
		}
	}
	for _, remote := range remotes {
		name := "refs/remotes/" + remote + "/HEAD"
		target, err := git(ctx, root, "symbolic-ref", "--quiet", name)
		if err != nil {
			var exit *exec.ExitError
			if errors.As(err, &exit) && exit.ExitCode() == 1 {
				continue
			}
			return err
		}
		refs[name] = strings.TrimSpace(string(target))
	}
	return nil
}

func chooseDefaultBase(refs map[string]string) (string, error) {
	for _, remote := range []string{"origin", "upstream"} {
		if target := refs["refs/remotes/"+remote+"/HEAD"]; target != "" {
			return availableDefault(refs, target)
		}
	}
	var defaults []string
	for name, target := range refs {
		if strings.HasPrefix(name, "refs/remotes/") && strings.HasSuffix(name, "/HEAD") && target != "" && !slices.Contains(defaults, target) {
			defaults = append(defaults, target)
		}
	}
	if len(defaults) > 0 {
		target, err := singleBase(defaults)
		if err != nil {
			return "", err
		}
		return availableDefault(refs, target)
	}
	// CI checkouts and local-only repositories may not have remote HEAD refs.
	// Prefer a remote tracking branch, but never guess between main and master
	// at the same priority or fall back to the current feature branch.
	for _, prefix := range []string{"refs/remotes/origin/", "refs/remotes/upstream/", "refs/heads/"} {
		var candidates []string
		for _, name := range []string{"main", "master"} {
			candidate := prefix + name
			if _, ok := refs[candidate]; ok {
				candidates = append(candidates, candidate)
			}
		}
		if len(candidates) > 0 {
			return singleBase(candidates)
		}
	}
	return "", errors.New("cannot determine the repository's default branch; set a remote HEAD (git remote set-head origin --auto) or use --base")
}

func availableDefault(refs map[string]string, target string) (string, error) {
	if _, ok := refs[target]; !ok {
		return "", fmt.Errorf("default branch %q is unavailable locally; fetch it or use --base", target)
	}
	return target, nil
}

func singleBase(candidates []string) (string, error) {
	if len(candidates) != 1 {
		slices.Sort(candidates)
		return "", fmt.Errorf("ambiguous default branch (%s); use --base explicitly", strings.Join(candidates, ", "))
	}
	return candidates[0], nil
}
