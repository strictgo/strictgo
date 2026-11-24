package diff

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func TestChooseDefaultBase(t *testing.T) {
	cases := []struct {
		name string
		refs map[string]string
		want string
	}{
		{"origin HEAD", map[string]string{
			"refs/remotes/origin/HEAD":   "refs/remotes/origin/master",
			"refs/remotes/origin/master": "", "refs/remotes/origin/main": "", "refs/heads/main": "",
		}, "refs/remotes/origin/master"},
		{"nonstandard default", map[string]string{
			"refs/remotes/origin/HEAD":           "refs/remotes/origin/release/stable",
			"refs/remotes/origin/release/stable": "", "refs/heads/main": "",
		}, "refs/remotes/origin/release/stable"},
		{"origin preferred over upstream", map[string]string{
			"refs/remotes/origin/HEAD": "refs/remotes/origin/trunk", "refs/remotes/origin/trunk": "",
			"refs/remotes/upstream/HEAD": "refs/remotes/upstream/main", "refs/remotes/upstream/main": "",
		}, "refs/remotes/origin/trunk"},
		{"upstream HEAD before conventions", map[string]string{
			"refs/remotes/upstream/HEAD": "refs/remotes/upstream/develop", "refs/remotes/upstream/develop": "",
			"refs/remotes/origin/main": "",
		}, "refs/remotes/upstream/develop"},
		{"other remote HEAD", map[string]string{
			"refs/remotes/company/HEAD": "refs/remotes/company/default", "refs/remotes/company/default": "",
		}, "refs/remotes/company/default"},
		{"CI remote fallback", map[string]string{"refs/remotes/origin/main": ""}, "refs/remotes/origin/main"},
		{"remote before local", map[string]string{"refs/remotes/origin/master": "", "refs/heads/main": ""}, "refs/remotes/origin/master"},
		{"upstream fallback", map[string]string{"refs/remotes/upstream/master": ""}, "refs/remotes/upstream/master"},
		{"local main", map[string]string{"refs/heads/main": "", "refs/heads/topic": ""}, "refs/heads/main"},
		{"local master", map[string]string{"refs/heads/master": "", "refs/heads/topic": ""}, "refs/heads/master"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got, err := chooseDefaultBase(test.refs)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("default base = %q, want %q", got, test.want)
			}
		})
	}
}

func TestAmbiguousOrMissingBase(t *testing.T) {
	cases := []map[string]string{
		{},
		{"refs/heads/topic": ""},
		{"refs/heads/main": "", "refs/heads/master": ""},
		{"refs/remotes/origin/main": "", "refs/remotes/origin/master": "", "refs/heads/main": ""},
		{"refs/remotes/company/HEAD": "refs/remotes/company/main", "refs/remotes/company/main": "", "refs/remotes/fork/HEAD": "refs/remotes/fork/master", "refs/remotes/fork/master": ""},
		{"refs/remotes/origin/HEAD": "refs/remotes/origin/missing", "refs/heads/main": ""},
	}
	for _, refs := range cases {
		_, err := chooseDefaultBase(refs)
		if err != nil {
			if !strings.Contains(err.Error(), "--base") {
				t.Errorf("error must explain the override: %v", err)
			}
			continue
		}
		t.Errorf("guessed a default base from %v", refs)
	}
}

func TestDefaultBaseRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git is not installed")
	}
	repo := t.TempDir()
	gitTest(t, repo, "init", "-b", "main")
	gitTest(t, repo, "config", "user.name", "Strictgo Test")
	gitTest(t, repo, "config", "user.email", "strictgo@example.test")
	gitTest(t, repo, "config", "commit.gpgsign", "false")
	gitTest(t, repo, "commit", "--allow-empty", "-qm", "base")
	gitTest(t, repo, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	gitTest(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	got, err := defaultBase(context.Background(), repo)
	if err != nil {
		t.Fatal(err)
	}
	if got != "refs/remotes/origin/trunk" {
		t.Errorf("default base = %q", got)
	}
	gitTest(t, repo, "update-ref", "-d", "refs/remotes/origin/trunk")
	if _, err := defaultBase(context.Background(), repo); err != nil {
		return
	}
	t.Fatal("a dangling default HEAD must not silently fall back to local main")
}
