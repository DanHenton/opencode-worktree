package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/testutil"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func TestRunTaskMissingName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runTask([]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "task name is required") {
		t.Errorf("expected 'task name is required', got: %v", err)
	}
}

func TestRunTaskInvalidName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runTask([]string{"bad name with spaces"})
	if err == nil {
		t.Fatal("expected error for invalid task name, got nil")
	}
}

func TestRunTaskExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runTask([]string{"valid-name", "msg", "extra"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestRunTaskNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runTask([]string{"some-task"})
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestRunTaskAlreadyExists(t *testing.T) {
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

	err = runTask([]string{taskName})
	if err == nil {
		t.Fatal("expected error for already-existing worktree, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected 'already exists' in error, got: %v", err)
	}
}

func TestRunTaskUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runTask([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown flag, got: %v", err)
	}
}

func TestRunAttachMissingName(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runAttach([]string{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "task name is required") {
		t.Errorf("expected 'task name is required', got: %v", err)
	}
}

func TestRunAttachExtraArg(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runAttach([]string{"name", "extra"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestRunAttachNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runAttach([]string{"some-task"})
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestRunAttachWorktreeNotFound(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runAttach([]string{"nonexistent-task"})
	if err == nil {
		t.Fatal("expected error for nonexistent worktree, got nil")
	}
}

func TestRunAttachUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runAttach([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown flag, got: %v", err)
	}
}

func TestRunMergeNotInAgentWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runMerge([]string{})
	if err == nil {
		t.Fatal("expected error when not in agent worktree, got nil")
	}
}

func TestRunMergeExtraArgs(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runMerge([]string{"arg1", "arg2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestRunMergeUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runMerge([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown flag, got: %v", err)
	}
}

func TestRunSyncNotInAgentWorktree(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runSync([]string{})
	if err == nil {
		t.Fatal("expected error when not in agent worktree, got nil")
	}
}

func TestRunSyncExtraArgs(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runSync([]string{"arg1", "arg2"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected extra argument") {
		t.Errorf("expected 'unexpected extra argument', got: %v", err)
	}
}

func TestRunSyncUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runSync([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown flag, got: %v", err)
	}
}

func TestRunListNoWorktrees(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runList([]string{})
	if err != nil {
		t.Errorf("expected nil error for list in empty repo, got: %v", err)
	}
}

func TestRunListNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runList([]string{})
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestRunCleanupInGitRepo(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runCleanup([]string{"--yes"})
	if err != nil {
		t.Errorf("expected nil error for cleanup in git repo, got: %v", err)
	}
}

func TestRunCleanupDryRun(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runCleanup([]string{"--dry-run"})
	if err != nil {
		t.Errorf("expected nil error for --dry-run cleanup, got: %v", err)
	}
}

func TestRunCleanupNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runCleanup([]string{"--dry-run"})
	if err == nil {
		t.Fatal("expected error when not in git repo, got nil")
	}
	if !strings.Contains(err.Error(), "not inside a git repository") {
		t.Errorf("expected 'not inside a git repository', got: %v", err)
	}
}

func TestRunCleanupUnknownFlag(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runCleanup([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown flag, got: %v", err)
	}
}

func TestRunCompletionsNoArgs(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runCompletions([]string{})
	if err != nil {
		t.Errorf("expected nil error for completions with no args, got: %v", err)
	}
}

func TestRunCompletionsAttachSubcommand(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	t.Chdir(repoDir)

	err := runCompletions([]string{"attach"})
	if err != nil {
		t.Errorf("expected nil error for completions attach in repo, got: %v", err)
	}
}

func TestRunCompletionsNotInGitRepo(t *testing.T) {
	t.Chdir(t.TempDir())

	err := runCompletions([]string{})
	if err != nil && !errors.Is(err, errSilent) {
		t.Errorf("expected nil or errSilent when not in git repo, got: %v", err)
	}
}
