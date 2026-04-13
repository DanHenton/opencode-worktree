package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/testutil"
)

func TestRepoRoot(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	root, err := git.RepoRoot(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root != filepath.Clean(repoDir) {
		t.Errorf("expected %q, got %q", filepath.Clean(repoDir), root)
	}
}

func TestCurrentBranch(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	testutil.RunGit(t, repoDir, "checkout", "-b", "test-branch")

	branch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if branch != "test-branch" {
		t.Errorf("expected test-branch, got %q", branch)
	}
}

func TestCommitCountBetween(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	baseBranch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testutil.RunGit(t, repoDir, "checkout", "-b", "feature-branch")

	count, err := git.CommitCountBetween(repoDir, baseBranch, "feature-branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 commits, got %d", count)
	}

	testutil.CommitFile(t, repoDir, "test1.txt", "content1", "First commit")

	count, err = git.CommitCountBetween(repoDir, baseBranch, "feature-branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 commit, got %d", count)
	}

	testutil.CommitFile(t, repoDir, "test2.txt", "content2", "Second commit")
	testutil.CommitFile(t, repoDir, "test3.txt", "content3", "Third commit")

	count, err = git.CommitCountBetween(repoDir, baseBranch, "feature-branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 commits, got %d", count)
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	dirty, err := git.HasUncommittedChanges(repoDir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dirty {
		t.Errorf("expected clean repo, got dirty")
	}

	if err := os.WriteFile(filepath.Join(repoDir, "untracked.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	dirty, err = git.HasUncommittedChanges(repoDir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dirty {
		t.Errorf("expected dirty repo after adding untracked file")
	}

	dirty, err = git.HasUncommittedChanges(repoDir, []string{"untracked.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dirty {
		t.Errorf("expected clean repo when excluding the untracked file")
	}

	testutil.CommitFile(t, repoDir, "untracked.txt", "hello", "Commit untracked")

	dirty, err = git.HasUncommittedChanges(repoDir, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dirty {
		t.Errorf("expected clean repo after committing")
	}
}

func TestHasUncommittedChangesExcludesDirectories(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	if err := os.MkdirAll(filepath.Join(repoDir, ".sisyphus", "plans"), 0755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoDir, ".sisyphus", "plans", "plan.md"), []byte("draft"), 0644); err != nil {
		t.Fatalf("failed to write nested file: %v", err)
	}

	dirty, err := git.HasUncommittedChanges(repoDir, []string{".sisyphus/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dirty {
		t.Errorf("expected clean repo when excluding .sisyphus/ contents")
	}

	if err := os.WriteFile(filepath.Join(repoDir, "real-change.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write regular file: %v", err)
	}

	dirty, err = git.HasUncommittedChanges(repoDir, []string{".sisyphus/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dirty {
		t.Errorf("expected dirty repo when non-excluded files exist")
	}
}
