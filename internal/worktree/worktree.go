package worktree

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/danhenton/opencode-worktree/internal/git"
)

const BranchPrefix = "agent/"
const DirSuffix = "-agent-"

var validTaskName = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

func ValidateTaskName(name string) error {
	if !validTaskName.MatchString(name) {
		return fmt.Errorf("invalid task name: '%s' — only alphanumeric characters and hyphens allowed", name)
	}
	return nil
}

func BranchName(taskName string) string {
	return BranchPrefix + taskName
}

func WorktreeDir(repoRoot, taskName string) string {
	return filepath.Join(filepath.Dir(repoRoot), filepath.Base(repoRoot)+DirSuffix+taskName)
}

func AlreadyExists(repoRoot, taskName string) bool {
	out, err := git.WorktreeList(repoRoot)
	if err != nil {
		return false
	}
	return strings.Contains(out, BranchPrefix+taskName)
}

func Create(repoRoot, taskName, parentBranch string) (string, error) {
	worktreeDir := WorktreeDir(repoRoot, taskName)
	branch := BranchName(taskName)

	if err := git.WorktreeAdd(repoRoot, worktreeDir, branch, parentBranch); err != nil {
		return "", fmt.Errorf("failed to create worktree: %w", err)
	}

	if err := writeParentBranchMarker(worktreeDir, parentBranch); err != nil {
		return "", err
	}

	if err := writeAgentContext(worktreeDir, parentBranch, branch, repoRoot); err != nil {
		return "", err
	}

	if err := copyOpenCodeConfig(repoRoot, worktreeDir); err != nil {
		return "", err
	}

	return worktreeDir, nil
}

func writeParentBranchMarker(worktreeDir, parentBranch string) error {
	path := filepath.Join(worktreeDir, ".agent-parent-branch")
	return os.WriteFile(path, []byte(parentBranch+"\n"), 0644)
}

func writeAgentContext(worktreeDir, parentBranch, agentBranch, sourceRepo string) error {
	path := filepath.Join(worktreeDir, ".agent-context")
	content := fmt.Sprintf("WORKTREE_DIR=%s\nPARENT_BRANCH=%s\nAGENT_BRANCH=%s\nSOURCE_REPO=%s\n",
		worktreeDir, parentBranch, agentBranch, sourceRepo)
	return os.WriteFile(path, []byte(content), 0644)
}

func copyOpenCodeConfig(repoRoot, worktreeDir string) error {
	configFile := filepath.Join(repoRoot, "opencode.json")
	if info, err := os.Stat(configFile); err == nil && !info.IsDir() {
		dest := filepath.Join(worktreeDir, "opencode.json")
		cmd := exec.Command("cp", configFile, dest)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy opencode.json: %w", err)
		}
	}

	configDir := filepath.Join(repoRoot, ".opencode")
	if info, err := os.Stat(configDir); err == nil && info.IsDir() {
		dest := filepath.Join(worktreeDir, ".opencode")
		cmd := exec.Command("cp", "-r", configDir, dest)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy .opencode/: %w", err)
		}
	}

	return nil
}

func LaunchOpenCode(worktreeDir, initialPrompt string) error {
	var cmd *exec.Cmd
	if initialPrompt != "" {
		cmd = exec.Command("opencode", "--message", initialPrompt)
	} else {
		cmd = exec.Command("opencode")
	}
	cmd.Dir = worktreeDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func List(repoRoot string) (string, error) {
	out, err := git.WorktreeList(repoRoot)
	if err != nil {
		return "", err
	}

	var agentLines []string
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, BranchPrefix) {
			agentLines = append(agentLines, line)
		}
	}

	if len(agentLines) == 0 {
		return "  (none)", nil
	}
	return strings.Join(agentLines, "\n"), nil
}

func Cleanup(repoRoot string) error {
	if err := git.WorktreePrune(repoRoot); err != nil {
		return fmt.Errorf("failed to prune worktree entries: %w", err)
	}

	if err := cleanupOrphanedDirectories(repoRoot); err != nil {
		return err
	}

	return cleanupOrphanedBranches(repoRoot)
}

func cleanupOrphanedDirectories(repoRoot string) error {
	porcelain, err := git.WorktreeListPorcelain(repoRoot)
	if err != nil {
		return err
	}

	activeWorktrees := make(map[string]bool)
	for _, line := range strings.Split(porcelain, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			activeWorktrees[strings.TrimPrefix(line, "worktree ")] = true
		}
	}

	siblingPrefix := WorktreeDir(repoRoot, "")
	parentDir := filepath.Dir(repoRoot)
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return fmt.Errorf("failed to read parent directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(parentDir, entry.Name())
		if !strings.HasPrefix(fullPath, siblingPrefix) {
			continue
		}
		if activeWorktrees[fullPath] {
			fmt.Printf("⚠️  Skipping active worktree directory: %s\n", fullPath)
			continue
		}
		fmt.Printf("⚠️  Removing stale worktree directory: %s\n", fullPath)
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to remove %s: %w", fullPath, err)
		}
		fmt.Printf("✅ Removed: %s\n", fullPath)
	}

	return nil
}

func cleanupOrphanedBranches(repoRoot string) error {
	branchOutput, err := git.BranchList(repoRoot)
	if err != nil {
		return err
	}

	porcelain, err := git.WorktreeListPorcelain(repoRoot)
	if err != nil {
		return err
	}

	activeBranches := make(map[string]bool)
	for _, line := range strings.Split(porcelain, "\n") {
		if strings.HasPrefix(line, "branch refs/heads/") {
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			activeBranches[branch] = true
		}
	}

	for _, line := range strings.Split(branchOutput, "\n") {
		branch := strings.TrimSpace(strings.TrimPrefix(line, "* "))
		if !strings.HasPrefix(branch, BranchPrefix) {
			continue
		}
		if activeBranches[branch] {
			fmt.Printf("⚠️  Skipping active worktree branch: %s\n", branch)
			continue
		}
		if _, err := git.BranchDelete(repoRoot, branch); err != nil {
			fmt.Printf("⚠️  Could not delete branch (unmerged?): %s — use 'git branch -D %s' to force\n", branch, branch)
		} else {
			fmt.Printf("✅ Deleted branch: %s\n", branch)
		}
	}

	return nil
}
