package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func NewTestRepo(t *testing.T) string {
	t.Helper()

	parentDir := t.TempDir()
	repoDir := filepath.Join(parentDir, "repo")
	if err := os.Mkdir(repoDir, 0755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}

	RunGit(t, repoDir, "init")
	RunGit(t, repoDir, "config", "user.name", "Test User")
	RunGit(t, repoDir, "config", "user.email", "test@example.com")

	CommitFile(t, repoDir, "README.md", "# Test Repo\n", "Initial commit")

	return repoDir
}

func CommitFile(t *testing.T, repoDir, filename, content, message string) {
	t.Helper()

	path := filepath.Join(repoDir, filename)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", filename, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", filename, err)
	}

	RunGit(t, repoDir, "add", filename)
	RunGit(t, repoDir, "commit", "-m", message)
}

func RunGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test User",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test User",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
	}
}
