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

func TestValidateTaskName(t *testing.T) {
	validNames := []string{"task", "task-123", "TASK-name", "123"}
	for _, name := range validNames {
		if err := worktree.ValidateTaskName(name); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", name, err)
		}
	}

	invalidNames := []string{"task 1", "task_1", "task@1", "", "task/1"}
	for _, name := range invalidNames {
		if err := worktree.ValidateTaskName(name); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", name)
		}
	}
}

func TestBranchName(t *testing.T) {
	expected := "agent/task-1"
	got := worktree.BranchName("task-1")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestWorktreeDir(t *testing.T) {
	repoRoot := "/path/to/my-repo"
	expected := "/path/to/my-repo-agent-task-1"
	got := filepath.ToSlash(worktree.WorktreeDir(repoRoot, "task-1"))
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestCreate(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	taskName := "test-task"
	parentBranch, _ := git.CurrentBranch(repoDir)

	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("unexpected error creating worktree: %v", err)
	}

	if _, err := os.Stat(worktreeDir); os.IsNotExist(err) {
		t.Errorf("expected worktree dir %q to exist", worktreeDir)
	}

	parentBranchData, err := os.ReadFile(filepath.Join(worktreeDir, ".agent-parent-branch"))
	if err != nil {
		t.Fatalf("unexpected error reading parent branch file: %v", err)
	}
	if strings.TrimSpace(string(parentBranchData)) != parentBranch {
		t.Errorf("expected parent branch %q, got %q", parentBranch, string(parentBranchData))
	}

	contextData, err := os.ReadFile(filepath.Join(worktreeDir, ".agent-context"))
	if err != nil {
		t.Fatalf("unexpected error reading agent context file: %v", err)
	}
	contextStr := string(contextData)
	if !strings.Contains(contextStr, "WORKTREE_DIR="+worktreeDir) {
		t.Errorf("agent context missing correct WORKTREE_DIR: %s", contextStr)
	}
	if !strings.Contains(contextStr, "AGENT_BRANCH=agent/"+taskName) {
		t.Errorf("agent context missing correct AGENT_BRANCH: %s", contextStr)
	}
	if !strings.Contains(contextStr, "PARENT_BRANCH="+parentBranch) {
		t.Errorf("agent context missing correct PARENT_BRANCH: %s", contextStr)
	}

	out, _ := git.BranchList(repoDir)
	if !strings.Contains(out, "agent/"+taskName) {
		t.Errorf("expected branch agent/%s to exist", taskName)
	}
}

func TestCreateCopiesOpenCodeConfig(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	taskName := "test-config"
	parentBranch, _ := git.CurrentBranch(repoDir)

	configContent := `{"test":true}`
	if err := os.WriteFile(filepath.Join(repoDir, "opencode.json"), []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write opencode.json: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repoDir, ".opencode"), 0755); err != nil {
		t.Fatalf("failed to create .opencode dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoDir, ".opencode", "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write .opencode file: %v", err)
	}

	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("unexpected error creating worktree: %v", err)
	}

	copiedConfig, err := os.ReadFile(filepath.Join(worktreeDir, "opencode.json"))
	if err != nil {
		t.Errorf("expected opencode.json to be copied: %v", err)
	} else if string(copiedConfig) != configContent {
		t.Errorf("expected opencode.json content %q, got %q", configContent, string(copiedConfig))
	}

	if _, err := os.Stat(filepath.Join(worktreeDir, ".opencode", "test.txt")); os.IsNotExist(err) {
		t.Errorf("expected .opencode dir and contents to be copied")
	}
}

func TestAlreadyExists(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	taskName := "test-exists"
	parentBranch, _ := git.CurrentBranch(repoDir)

	exists, err := worktree.AlreadyExists(repoDir, taskName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Errorf("expected worktree not to exist yet")
	}

	if _, err := worktree.Create(repoDir, taskName, parentBranch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err = worktree.AlreadyExists(repoDir, taskName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("expected worktree to exist")
	}
}

func TestList(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	list, err := worktree.List(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(list) != "(none)" {
		t.Errorf("expected (none), got %q", list)
	}

	taskName := "test-list"
	parentBranch, _ := git.CurrentBranch(repoDir)
	worktreeDir, err := worktree.Create(repoDir, taskName, parentBranch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, err = worktree.List(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(list, "agent/"+taskName) {
		t.Errorf("expected list to contain branch %q, got %q", "agent/"+taskName, list)
	}
	if strings.Contains(list, "(uncommitted changes)") {
		t.Errorf("expected no uncommitted changes indicator for clean worktree, got %q", list)
	}

	if err := os.WriteFile(filepath.Join(worktreeDir, "dirty.txt"), []byte("wip"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	list, err = worktree.List(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(list, "(uncommitted changes)") {
		t.Errorf("expected uncommitted changes indicator for dirty worktree, got %q", list)
	}
}

func TestActiveTaskNames(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	names, err := worktree.ActiveTaskNames(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected no task names, got %v", names)
	}

	if _, err := worktree.Create(repoDir, "alpha", parentBranch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := worktree.Create(repoDir, "beta", parentBranch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names, err = worktree.ActiveTaskNames(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 task names, got %v", names)
	}

	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["alpha"] || !found["beta"] {
		t.Errorf("expected alpha and beta, got %v", names)
	}
}

func TestResolveWorktreeDir(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	_, err := worktree.ResolveWorktreeDir(repoDir, "nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent task")
	}

	createdDir, err := worktree.Create(repoDir, "resolve-me", parentBranch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resolved, err := worktree.ResolveWorktreeDir(repoDir, "resolve-me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != createdDir {
		t.Errorf("expected %q, got %q", createdDir, resolved)
	}
}

func TestCleanup(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)
	parentBranch, _ := git.CurrentBranch(repoDir)

	orphanedTask := "orphaned"
	orphanedDir := worktree.WorktreeDir(repoDir, orphanedTask)
	if err := os.MkdirAll(orphanedDir, 0755); err != nil {
		t.Fatalf("failed to create orphaned dir: %v", err)
	}

	activeTask := "active"
	activeDir, err := worktree.Create(repoDir, activeTask, parentBranch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := worktree.Cleanup(repoDir, worktree.CleanupOptions{Yes: true}); err != nil {
		t.Fatalf("unexpected error during cleanup: %v", err)
	}

	if _, err := os.Stat(orphanedDir); !os.IsNotExist(err) {
		t.Errorf("expected orphaned dir to be removed")
	}

	if _, err := os.Stat(activeDir); os.IsNotExist(err) {
		t.Errorf("expected active worktree dir to be preserved")
	}
}

func TestCleanupDryRun(t *testing.T) {
	repoDir := testutil.NewTestRepo(t)

	orphanedTask := "orphaned-dry"
	orphanedDir := worktree.WorktreeDir(repoDir, orphanedTask)
	if err := os.MkdirAll(orphanedDir, 0755); err != nil {
		t.Fatalf("failed to create orphaned dir: %v", err)
	}

	if err := worktree.Cleanup(repoDir, worktree.CleanupOptions{DryRun: true}); err != nil {
		t.Fatalf("unexpected error during dry-run cleanup: %v", err)
	}

	if _, err := os.Stat(orphanedDir); os.IsNotExist(err) {
		t.Errorf("expected orphaned dir to still exist after dry-run")
	}
}

func TestLaunchOpenCodeMissingBinary(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	err := worktree.LaunchOpenCode(t.TempDir(), "")
	if err == nil {
		t.Fatal("expected error when opencode is not in PATH")
	}
	if !strings.Contains(err.Error(), "opencode not found in PATH") {
		t.Errorf("expected error about missing opencode, got: %v", err)
	}
}
