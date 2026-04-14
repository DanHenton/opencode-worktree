package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/testutil"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func newTestCmd() (*bytes.Buffer, *bytes.Buffer, func(args ...string) error) {
	var outBuf, errBuf bytes.Buffer
	return &outBuf, &errBuf, func(args ...string) error {
		root := newRootCmd()
		root.SetOut(&outBuf)
		root.SetErr(&errBuf)
		root.SetArgs(args)
		return root.Execute()
	}
}

func TestTaskMissingName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, exec := newTestCmd()
	err := exec("task")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "task name is required") {
		t.Errorf("expected 'task name is required', got: %v", err)
	}
}

func TestTaskInvalidName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("task", "bad name with spaces")
	if err == nil {
		t.Fatal("expected error for invalid task name, got nil")
	}
}

func TestTaskExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("task", "valid-name", "msg", "extra")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestTaskNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	_, _, run := newTestCmd()
	err := run("task", "some-task")
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestTaskAlreadyExists(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	parentBranch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}

	taskName := "my-existing-task"
	if _, err := worktree.Create(repoDir, taskName, parentBranch); err != nil {
		t.Fatalf("failed to pre-create worktree: %v", err)
	}

	_, _, run := newTestCmd()
	err = run("task", taskName)
	if err == nil {
		t.Fatal("expected error for already-existing worktree, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' in error, got: %v", err)
	}
}

func TestTaskUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("task", "--unknown-flag")
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestTaskReturnsLaunchError(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("failed to find git: %v", err)
	}
	t.Setenv("PATH", filepath.Dir(gitPath))

	_, _, run := newTestCmd()
	err = run("task", "launch-fails", "--no-merge")
	if err == nil {
		t.Fatal("expected error when opencode is unavailable")
	}
	if !strings.Contains(err.Error(), "opencode not found in PATH") {
		t.Errorf("expected missing opencode error, got: %v", err)
	}
}

func TestTaskBranchExistsWithoutWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	parentBranch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}

	worktreeDir, err := worktree.Create(repoDir, "existing-branch", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}
	if err := git.WorktreeRemove(repoDir, worktreeDir); err != nil {
		t.Fatalf("failed to remove worktree: %v", err)
	}
	if err := git.WorktreePrune(repoDir); err != nil {
		t.Fatalf("failed to prune worktrees: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(worktreeDir, ".agent-parent-branch")); !os.IsNotExist(statErr) {
		t.Fatalf("expected worktree to be removed")
	}

	_, _, run := newTestCmd()
	err = run("task", "existing-branch", "--no-merge")
	if err == nil {
		t.Fatal("expected error when branch exists without worktree")
	}
	if !strings.Contains(err.Error(), "branch named 'agent/existing-branch' already exists") {
		t.Errorf("expected existing branch error, got: %v", err)
	}
}

func TestAttachMissingName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("attach")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "task name is required") {
		t.Errorf("expected 'task name is required', got: %v", err)
	}
}

func TestAttachExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("attach", "name", "extra")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestAttachNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	_, _, run := newTestCmd()
	err := run("attach", "some-task")
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestAttachWorktreeNotFound(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("attach", "nonexistent-task")
	if err == nil {
		t.Fatal("expected error for nonexistent worktree, got nil")
	}
}

func TestAttachUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("attach", "--unknown-flag")
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestAttachReturnsLaunchError(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)
	gitPath, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("failed to find git: %v", err)
	}
	t.Setenv("PATH", filepath.Dir(gitPath))

	parentBranch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}
	if _, err := worktree.Create(repoDir, "attach-launch-fails", parentBranch); err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	_, _, run := newTestCmd()
	err = run("attach", "attach-launch-fails", "--no-merge")
	if err == nil {
		t.Fatal("expected error when opencode is unavailable")
	}
	if !strings.Contains(err.Error(), "opencode not found in PATH") {
		t.Errorf("expected missing opencode error, got: %v", err)
	}
}

func TestMergeAcceptsTrailingNoCleanupFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, err := git.CurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("failed to get current branch: %v", err)
	}

	worktreeDir, err := worktree.Create(repoDir, "merge-trailing-flag", parentBranch)
	if err != nil {
		t.Fatalf("failed to create worktree: %v", err)
	}

	testutil.CommitFile(t, worktreeDir, "merged.txt", "content", "Agent commit")

	_, _, run := newTestCmd()
	err = run("merge", worktreeDir, "--no-cleanup")
	if err != nil {
		t.Fatalf("expected trailing --no-cleanup flag to parse, got: %v", err)
	}

	if _, statErr := os.Stat(worktreeDir); os.IsNotExist(statErr) {
		t.Fatalf("expected worktree to be preserved when --no-cleanup is set")
	}
}

func TestMergeNotInAgentWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("merge")
	if err == nil {
		t.Fatal("expected error when not in agent worktree, got nil")
	}
}

func TestMergeExtraArgs(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("merge", "arg1", "arg2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMergeUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("merge", "--unknown-flag")
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestSyncNotInAgentWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("sync")
	if err == nil {
		t.Fatal("expected error when not in agent worktree, got nil")
	}
}

func TestSyncExtraArgs(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("sync", "arg1", "arg2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSyncUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("sync", "--unknown-flag")
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestListNoWorktrees(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("list")
	if err != nil {
		t.Errorf("expected nil error for list in empty repo, got: %v", err)
	}
}

func TestListNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	_, _, run := newTestCmd()
	err := run("list")
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestListExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("list", "extra-arg")
	if err == nil {
		t.Fatal("expected error for extra arg to list, got nil")
	}
}

func TestCleanupInGitRepo(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("cleanup", "--yes")
	if err != nil {
		t.Errorf("expected nil error for cleanup in git repo, got: %v", err)
	}
}

func TestCleanupDryRun(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("cleanup", "--dry-run")
	if err != nil {
		t.Errorf("expected nil error for --dry-run cleanup, got: %v", err)
	}
}

func TestCleanupNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	_, _, run := newTestCmd()
	err := run("cleanup", "--dry-run")
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestCleanupUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("cleanup", "--unknown-flag")
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestCleanupExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	_, _, run := newTestCmd()
	err := run("cleanup", "extra-arg")
	if err == nil {
		t.Fatal("expected error for extra arg to cleanup, got nil")
	}
}
