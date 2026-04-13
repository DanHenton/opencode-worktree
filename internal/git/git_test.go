package git_test

import (
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
