package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/merge"
)

var (
	errSilent = errors.New("")
	useEmoji  = detectTerminal()
)

func detectTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func emoji(e, fallback string) string {
	if useEmoji {
		return e
	}
	return fallback
}

func printMergeResult(result *merge.Result) {
	if result.DirtyWorktree {
		fmt.Fprintf(os.Stderr, "%sWorktree has uncommitted changes — preserved at: %s\n", emoji("⚠️  ", "warning: "), result.WorktreePath)
		fmt.Fprintf(os.Stderr, "   Next: commit or discard changes, then run 'opencode-worktree merge %s'\n", result.WorktreePath)
	}
	if result.NoNewCommits && !result.DirtyWorktree {
		fmt.Fprintf(os.Stderr, "%sNo new commits found on %s. Cleaned up worktree only.\n", emoji("⚠️  ", "warning: "), result.AgentBranch)
		return
	}
	if result.Merged {
		if result.DirtyWorktree {
			fmt.Printf("%sMerged %s into %s (worktree kept due to uncommitted changes).\n", emoji("🚀 ", ""), result.AgentBranch, result.ParentBranch)
		} else {
			fmt.Printf("%sMerged %s into %s and cleaned up.\n", emoji("🚀 ", ""), result.AgentBranch, result.ParentBranch)
		}
	}
}

func handleMergeError(result *merge.Result, err error) error {
	if result != nil && len(result.ConflictFiles) > 0 {
		fmt.Fprintf(os.Stderr, "%sMerge conflict: %s into %s\n", emoji("❌ ", "error: "), result.AgentBranch, result.ParentBranch)
		fmt.Fprintln(os.Stderr, "Conflicting files:")
		for _, f := range result.ConflictFiles {
			fmt.Fprintf(os.Stderr, "  %s\n", f)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "To resolve:")
		fmt.Fprintf(os.Stderr, "  cd %s\n", result.RepoRoot)
		fmt.Fprintln(os.Stderr, "  git status")
		fmt.Fprintln(os.Stderr, "  # Fix conflicts in the listed files")
		fmt.Fprintln(os.Stderr, "  git add <resolved-files>")
		fmt.Fprintln(os.Stderr, "  git commit")
		fmt.Fprintln(os.Stderr, "  opencode-worktree cleanup")
		return errSilent
	}
	return fmt.Errorf("%v", err)
}
