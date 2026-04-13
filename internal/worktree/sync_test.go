package worktree_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/testutil"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func TestSyncAlreadyCurrent(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	worktreeDir, err := worktree.Create(repoDir, "sync-noop", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	result, err := worktree.Sync(worktreeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.AlreadyCurrent {
		t.Error("expected AlreadyCurrent to be true")
	}
	if result.Rebased {
		t.Error("expected Rebased to be false")
	}
}

func TestSyncRebasesParentChanges(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	worktreeDir, err := worktree.Create(repoDir, "sync-rebase", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, worktreeDir, "agent-work.txt", "agent changes\n", "Agent commit")

	testutil.CommitFile(t, repoDir, "parent-update.txt", "new parent work\n", "Parent commit")

	result, err := worktree.Sync(worktreeDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.AlreadyCurrent {
		t.Error("expected AlreadyCurrent to be false")
	}
	if !result.Rebased {
		t.Error("expected Rebased to be true")
	}

	if _, err := os.Stat(filepath.Join(worktreeDir, "parent-update.txt")); os.IsNotExist(err) {
		t.Error("expected parent-update.txt to be present after sync")
	}
	if _, err := os.Stat(filepath.Join(worktreeDir, "agent-work.txt")); os.IsNotExist(err) {
		t.Error("expected agent-work.txt to be preserved after sync")
	}
}

func TestSyncConflictAbortsRebase(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	worktreeDir, err := worktree.Create(repoDir, "sync-conflict", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, worktreeDir, "shared.txt", "agent version\n", "Agent edits shared")

	testutil.CommitFile(t, repoDir, "shared.txt", "parent version\n", "Parent edits shared")

	result, err := worktree.Sync(worktreeDir)
	if err == nil {
		t.Fatal("expected error on conflict")
	}
	if !strings.Contains(err.Error(), "rebase conflict") {
		t.Errorf("expected rebase conflict error, got: %v", err)
	}
	if len(result.ConflictFiles) == 0 {
		t.Error("expected conflict files to be reported")
	}

	branch, _ := git.CurrentBranch(worktreeDir)
	if branch == "" {
		t.Error("expected worktree to be on a branch after abort (not in rebase state)")
	}
}

func TestSyncRejectsDirtyWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	worktreeDir, err := worktree.Create(repoDir, "sync-dirty", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	if err := os.WriteFile(filepath.Join(worktreeDir, "uncommitted.txt"), []byte("wip"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	_, err = worktree.Sync(worktreeDir)
	if err == nil {
		t.Fatal("expected error for dirty worktree")
	}
	if !strings.Contains(err.Error(), "uncommitted changes") {
		t.Errorf("expected uncommitted changes error, got: %v", err)
	}
}

func TestSyncRejectsNonAgentBranch(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	if err := os.WriteFile(filepath.Join(repoDir, ".agent-parent-branch"), []byte("main\n"), 0644); err != nil {
		t.Fatalf("failed to write marker: %v", err)
	}

	_, err := worktree.Sync(repoDir)
	if err == nil {
		t.Fatal("expected error for non-agent branch")
	}
	if !strings.Contains(err.Error(), "not a managed agent worktree") {
		t.Errorf("expected non-agent worktree error, got: %v", err)
	}
}
