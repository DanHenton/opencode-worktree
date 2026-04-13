package merge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/testutil"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func TestMergeWithCommits(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	taskName := "feature-1"
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, worktreeDir, "new-file.txt", "content", "Agent commit")

	result, err := merge.Run(worktreeDir, true)
	if err != nil {
		t.Fatalf("unexpected error during merge: %v", err)
	}

	if !result.Merged {
		t.Errorf("expected merge to succeed")
	}

	if result.AgentBranch != worktree.BranchName(taskName) {
		t.Errorf("expected agent branch %q, got %q", worktree.BranchName(taskName), result.AgentBranch)
	}

	if result.ParentBranch != parentBranch {
		t.Errorf("expected parent branch %q, got %q", parentBranch, result.ParentBranch)
	}

	if result.NoNewCommits {
		t.Errorf("expected NoNewCommits to be false")
	}

	mergedFilePath := filepath.Join(repoDir, "new-file.txt")
	if _, err := os.Stat(mergedFilePath); os.IsNotExist(err) {
		t.Errorf("expected merged file %q to exist in parent branch", mergedFilePath)
	}

	if _, err := os.Stat(worktreeDir); !os.IsNotExist(err) {
		t.Errorf("expected worktree to be cleaned up")
	}
}

func TestMergeNoNewCommits(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	taskName := "feature-empty"
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	result, err := merge.Run(worktreeDir, true)
	if err != nil {
		t.Fatalf("unexpected error during merge: %v", err)
	}

	if result.Merged {
		t.Errorf("expected merge to be skipped (no commits)")
	}

	if !result.NoNewCommits {
		t.Errorf("expected NoNewCommits to be true")
	}

	if _, err := os.Stat(worktreeDir); !os.IsNotExist(err) {
		t.Errorf("expected worktree to be cleaned up")
	}
}

func TestMergeConflict(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	taskName := "feature-conflict"
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, repoDir, "conflict.txt", "parent content", "Parent commit")
	testutil.CommitFile(t, worktreeDir, "conflict.txt", "agent content", "Agent commit")

	result, err := merge.Run(worktreeDir, true)
	if err == nil {
		t.Fatalf("expected merge conflict error, but got none")
	}

	if len(result.ConflictFiles) == 0 || result.ConflictFiles[0] != "conflict.txt" {
		t.Errorf("expected conflict in conflict.txt, got %v", result.ConflictFiles)
	}

	if _, err := os.Stat(worktreeDir); os.IsNotExist(err) {
		t.Errorf("expected worktree to still exist after conflict")
	}
}

func TestMergeNoCleanup(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	taskName := "feature-keep"
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, worktreeDir, "new-file.txt", "content", "Agent commit")

	result, err := merge.Run(worktreeDir, false)
	if err != nil {
		t.Fatalf("unexpected error during merge: %v", err)
	}

	if !result.Merged {
		t.Errorf("expected merge to succeed")
	}

	if _, err := os.Stat(worktreeDir); os.IsNotExist(err) {
		t.Errorf("expected worktree to be preserved since cleanup=false")
	}
}

func TestDetectWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	_, err = merge.DetectWorktree()
	if err == nil {
		t.Errorf("expected DetectWorktree to fail in standard repo dir, but it succeeded")
	}

	taskName := "feature-detect"
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	if err := os.Chdir(worktreeDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	detectedRoot, err := merge.DetectWorktree()
	if err != nil {
		t.Errorf("expected DetectWorktree to succeed in worktree dir, but failed: %v", err)
	}

	absWorktreeDir, _ := filepath.Abs(worktreeDir)
	absDetectedRoot, _ := filepath.Abs(detectedRoot)

	if absDetectedRoot != absWorktreeDir {
		t.Errorf("expected detected root %q, got %q", absWorktreeDir, absDetectedRoot)
	}
}
