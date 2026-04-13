package worktree

import (
	"bufio"
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

func AlreadyExists(repoRoot, taskName string) (bool, error) {
	out, err := git.WorktreeList(repoRoot)
	if err != nil {
		return false, fmt.Errorf("failed to check existing worktrees: %w", err)
	}
	return strings.Contains(out, BranchPrefix+taskName), nil
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
		if err := copyFile(configFile, dest, info.Mode()); err != nil {
			return fmt.Errorf("failed to copy opencode.json: %w", err)
		}
	}

	configDir := filepath.Join(repoRoot, ".opencode")
	if info, err := os.Stat(configDir); err == nil && info.IsDir() {
		dest := filepath.Join(worktreeDir, ".opencode")
		if err := copyDir(configDir, dest); err != nil {
			return fmt.Errorf("failed to copy .opencode/: %w", err)
		}
	}

	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, mode)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		return copyFile(path, target, info.Mode())
	})
}

func LaunchOpenCode(worktreeDir, initialPrompt string) error {
	if _, err := exec.LookPath("opencode"); err != nil {
		return fmt.Errorf("opencode not found in PATH — install it from https://opencode.ai")
	}

	var cmd *exec.Cmd
	if initialPrompt != "" {
		cmd = exec.Command("opencode", "--prompt", initialPrompt)
	} else {
		cmd = exec.Command("opencode")
	}
	cmd.Dir = worktreeDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

var MarkerFiles = []string{
	".agent-parent-branch",
	".agent-context",
	"opencode.json",
	".opencode/",
	".sisyphus/",
}

func List(repoRoot string) (string, error) {
	out, err := git.WorktreeList(repoRoot)
	if err != nil {
		return "", err
	}

	var agentLines []string
	for line := range strings.SplitSeq(out, "\n") {
		if !strings.Contains(line, BranchPrefix) {
			continue
		}

		worktreePath := strings.Fields(line)[0]
		dirty, _ := git.HasUncommittedChanges(worktreePath, MarkerFiles)
		if dirty {
			line += " (uncommitted changes)"
		}

		agentLines = append(agentLines, line)
	}

	if len(agentLines) == 0 {
		return "  (none)", nil
	}
	return strings.Join(agentLines, "\n"), nil
}

func ActiveTaskNames(repoRoot string) ([]string, error) {
	out, err := git.WorktreeList(repoRoot)
	if err != nil {
		return nil, err
	}

	var names []string
	for line := range strings.SplitSeq(out, "\n") {
		start := strings.Index(line, "["+BranchPrefix)
		if start == -1 {
			continue
		}
		end := strings.Index(line[start:], "]")
		if end == -1 {
			continue
		}
		branch := line[start+1 : start+end]
		name := strings.TrimPrefix(branch, BranchPrefix)
		if name != "" {
			names = append(names, name)
		}
	}
	return names, nil
}

func ResolveWorktreeDir(repoRoot, taskName string) (string, error) {
	porcelain, err := git.WorktreeListPorcelain(repoRoot)
	if err != nil {
		return "", err
	}

	targetBranch := "branch refs/heads/" + BranchPrefix + taskName
	var currentWorktree string

	for line := range strings.SplitSeq(porcelain, "\n") {
		if worktreePath, ok := strings.CutPrefix(line, "worktree "); ok {
			currentWorktree = worktreePath
		}
		if strings.TrimSpace(line) == targetBranch && currentWorktree != "" {
			return currentWorktree, nil
		}
	}

	return "", fmt.Errorf("no worktree found for task '%s'", taskName)
}

type CleanupOptions struct {
	DryRun bool
	Yes    bool
}

func Cleanup(repoRoot string, opts CleanupOptions) error {
	if err := git.WorktreePrune(repoRoot); err != nil {
		return fmt.Errorf("failed to prune worktree entries: %w", err)
	}

	staleDirs, err := findOrphanedDirectories(repoRoot)
	if err != nil {
		return err
	}

	staleBranches, err := findOrphanedBranches(repoRoot)
	if err != nil {
		return err
	}

	if len(staleDirs) == 0 && len(staleBranches) == 0 {
		return nil
	}

	if opts.DryRun {
		for _, dir := range staleDirs {
			fmt.Fprintf(os.Stderr, "Would remove directory: %s\n", dir)
		}
		for _, branch := range staleBranches {
			fmt.Fprintf(os.Stderr, "Would delete branch: %s\n", branch)
		}
		return nil
	}

	if !opts.Yes {
		if fi, _ := os.Stdin.Stat(); fi != nil && fi.Mode()&os.ModeCharDevice != 0 {
			fmt.Fprintf(os.Stderr, "Remove %d worktree(s) and %d branch(es)? [y/N] ", len(staleDirs), len(staleBranches))
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() || (scanner.Text() != "y" && scanner.Text() != "Y") {
				return fmt.Errorf("cleanup aborted")
			}
		}
	}

	for _, dir := range staleDirs {
		fmt.Fprintf(os.Stderr, "Removing stale worktree directory: %s\n", dir)
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove %s: %w", dir, err)
		}
		fmt.Printf("Removed: %s\n", dir)
	}

	for _, branch := range staleBranches {
		if _, err := git.BranchDelete(repoRoot, branch); err != nil {
			fmt.Fprintf(os.Stderr, "Could not delete branch (unmerged?): %s — use 'git branch -D %s' to force\n", branch, branch)
		} else {
			fmt.Printf("Deleted branch: %s\n", branch)
		}
	}

	return nil
}

func findOrphanedDirectories(repoRoot string) ([]string, error) {
	porcelain, err := git.WorktreeListPorcelain(repoRoot)
	if err != nil {
		return nil, err
	}

	activeWorktrees := make(map[string]bool)
	for line := range strings.SplitSeq(porcelain, "\n") {
		if worktreePath, ok := strings.CutPrefix(line, "worktree "); ok {
			activeWorktrees[worktreePath] = true
		}
	}

	siblingPrefix := WorktreeDir(repoRoot, "")
	parentDir := filepath.Dir(repoRoot)
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read parent directory: %w", err)
	}

	var stale []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		fullPath := filepath.Join(parentDir, entry.Name())
		if !strings.HasPrefix(fullPath, siblingPrefix) {
			continue
		}
		if activeWorktrees[fullPath] {
			fmt.Fprintf(os.Stderr, "Skipping active worktree directory: %s\n", fullPath)
			continue
		}
		stale = append(stale, fullPath)
	}

	return stale, nil
}

func findOrphanedBranches(repoRoot string) ([]string, error) {
	branchOutput, err := git.BranchList(repoRoot)
	if err != nil {
		return nil, err
	}

	porcelain, err := git.WorktreeListPorcelain(repoRoot)
	if err != nil {
		return nil, err
	}

	activeBranches := make(map[string]bool)
	for line := range strings.SplitSeq(porcelain, "\n") {
		if branch, ok := strings.CutPrefix(line, "branch refs/heads/"); ok {
			activeBranches[branch] = true
		}
	}

	var stale []string
	for line := range strings.SplitSeq(branchOutput, "\n") {
		branch := strings.TrimSpace(strings.TrimPrefix(line, "* "))
		if !strings.HasPrefix(branch, BranchPrefix) {
			continue
		}
		if activeBranches[branch] {
			fmt.Fprintf(os.Stderr, "Skipping active worktree branch: %s\n", branch)
			continue
		}
		stale = append(stale, branch)
	}

	return stale, nil
}
