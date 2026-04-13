package merge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/gofrs/flock"
)

type Result struct {
	Merged        bool
	ConflictFiles []string
	AgentBranch   string
	ParentBranch  string
	NoNewCommits  bool
}

func Run(worktreePath string, cleanup bool) (*Result, error) {
	worktreePath, err := filepath.Abs(worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve worktree path: %w", err)
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("worktree path does not exist: %s", worktreePath)
	}

	parentBranch, err := readParentBranch(worktreePath)
	if err != nil {
		return nil, err
	}

	agentBranch, err := git.CurrentBranch(worktreePath)
	if err != nil {
		return nil, fmt.Errorf("could not determine agent branch for worktree: %s", worktreePath)
	}
	if agentBranch == "" {
		return nil, fmt.Errorf("could not determine agent branch for worktree: %s (detached HEAD?)", worktreePath)
	}

	repoRoot, err := resolveRepoRoot(worktreePath)
	if err != nil {
		return nil, err
	}

	commitCount, err := git.CommitCountBetween(repoRoot, parentBranch, agentBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to count commits: %w", err)
	}

	result := &Result{
		AgentBranch:  agentBranch,
		ParentBranch: parentBranch,
	}

	if commitCount == 0 {
		result.NoNewCommits = true
		if cleanup {
			return result, cleanupWorktree(repoRoot, worktreePath, agentBranch)
		}
		return result, nil
	}

	lockPath := filepath.Join(os.TempDir(), filepath.Base(repoRoot)+"-merge.lock")
	fileLock := flock.New(lockPath)

	if err := fileLock.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire merge lock: %w", err)
	}
	defer fileLock.Unlock()

	currentBranch, err := git.CurrentBranch(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	if currentBranch != parentBranch {
		if err := git.Checkout(repoRoot, parentBranch); err != nil {
			return nil, fmt.Errorf("failed to checkout parent branch %s: %w", parentBranch, err)
		}
	}

	if err := git.Merge(repoRoot, agentBranch); err != nil {
		conflicts, _ := git.ConflictingFiles(repoRoot)
		_ = git.MergeAbort(repoRoot)
		result.ConflictFiles = conflicts
		return result, fmt.Errorf("merge conflict detected while merging %s into %s", agentBranch, parentBranch)
	}

	result.Merged = true

	if cleanup {
		if err := cleanupWorktree(repoRoot, worktreePath, agentBranch); err != nil {
			return result, fmt.Errorf("merge succeeded but cleanup failed: %w", err)
		}
	}

	return result, nil
}

func DetectWorktree() (string, error) {
	dir, err := git.RepoRoot(".")
	if err != nil {
		return "", fmt.Errorf("not inside a git repository")
	}

	markerPath := filepath.Join(dir, ".agent-parent-branch")
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		return "", fmt.Errorf("current directory is not an agent worktree (missing .agent-parent-branch)")
	}

	return dir, nil
}

func readParentBranch(worktreePath string) (string, error) {
	markerPath := filepath.Join(worktreePath, ".agent-parent-branch")
	data, err := os.ReadFile(markerPath)
	if err != nil {
		return "", fmt.Errorf("missing parent branch marker: %s", markerPath)
	}
	branch := strings.TrimSpace(string(data))
	if branch == "" {
		return "", fmt.Errorf("parent branch marker is empty: %s", markerPath)
	}
	return branch, nil
}

func resolveRepoRoot(worktreePath string) (string, error) {
	commonDir, err := git.GitCommonDir(worktreePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve git common dir: %w", err)
	}

	absCommonDir, err := filepath.Abs(commonDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute common dir: %w", err)
	}

	repoRoot := filepath.Dir(absCommonDir)
	return repoRoot, nil
}

func cleanupWorktree(repoRoot, worktreePath, agentBranch string) error {
	if err := git.WorktreeRemove(repoRoot, worktreePath); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}
	if _, err := git.BranchDelete(repoRoot, agentBranch); err != nil {
		return fmt.Errorf("failed to delete branch %s: %w", agentBranch, err)
	}
	return git.WorktreePrune(repoRoot)
}
