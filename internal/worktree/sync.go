package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danhenton/opencode-worktree/internal/git"
)

type SyncResult struct {
	AgentBranch    string
	ParentBranch   string
	WorktreePath   string
	AlreadyCurrent bool
	Rebased        bool
	ConflictFiles  []string
}

func Sync(worktreePath string) (*SyncResult, error) {
	parentBranch, err := ReadParentBranch(worktreePath)
	if err != nil {
		return nil, err
	}

	agentBranch, err := git.CurrentBranch(worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine agent branch: %w", err)
	}
	if agentBranch == "" {
		return nil, fmt.Errorf("detached HEAD in worktree — cannot sync")
	}
	if !strings.HasPrefix(agentBranch, BranchPrefix) {
		return nil, fmt.Errorf("not a managed agent worktree: branch %q does not have %s prefix", agentBranch, BranchPrefix)
	}

	result := &SyncResult{
		AgentBranch:  agentBranch,
		ParentBranch: parentBranch,
		WorktreePath: worktreePath,
	}

	dirty, err := git.HasUncommittedChanges(worktreePath, MarkerFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to check worktree status: %w", err)
	}
	if dirty {
		return nil, fmt.Errorf("worktree has uncommitted changes — commit or stash before syncing")
	}

	mergeBase, err := git.MergeBase(worktreePath, agentBranch, parentBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to find merge base: %w", err)
	}

	parentTip, err := git.CommitCountBetween(worktreePath, mergeBase, parentBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to check parent branch: %w", err)
	}
	if parentTip == 0 {
		result.AlreadyCurrent = true
		return result, nil
	}

	if err := git.Rebase(worktreePath, parentBranch); err != nil {
		conflicts, _ := git.ConflictingFiles(worktreePath)
		_ = git.RebaseAbort(worktreePath)
		result.ConflictFiles = conflicts
		return result, fmt.Errorf("rebase conflict while syncing %s onto %s", agentBranch, parentBranch)
	}

	result.Rebased = true
	return result, nil
}

func ReadParentBranch(worktreePath string) (string, error) {
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
