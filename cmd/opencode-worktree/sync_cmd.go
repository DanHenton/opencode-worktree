package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runSync(args []string) {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree sync [path]

Rebase the agent branch onto the latest parent branch, pulling in
upstream changes. If no path is given, auto-detects the current
directory as an agent worktree.

Examples:
  opencode-worktree sync
  opencode-worktree sync /path/to/worktree
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) > 1 {
		exitError("unexpected extra argument: %s", positional[1])
	}

	var worktreePath string
	if len(positional) == 1 {
		worktreePath = positional[0]
	}

	if worktreePath == "" {
		detected, err := merge.DetectWorktree()
		if err != nil {
			exitError("%v\n\nUsage: opencode-worktree sync [worktree-path]", err)
		}
		worktreePath = detected
	}

	result, err := worktree.Sync(worktreePath)
	if err != nil {
		handleSyncError(result, err)
	}

	if result.AlreadyCurrent {
		fmt.Printf("%sAlready up to date with %s.\n", emoji("✅ ", ""), result.ParentBranch)
		return
	}

	fmt.Printf("%sRebased %s onto %s.\n", emoji("✅ ", ""), result.AgentBranch, result.ParentBranch)
}

func handleSyncError(result *worktree.SyncResult, err error) {
	if result != nil && len(result.ConflictFiles) > 0 {
		fmt.Fprintf(os.Stderr, "%sRebase conflict: %s onto %s\n", emoji("❌ ", "error: "), result.AgentBranch, result.ParentBranch)
		fmt.Fprintln(os.Stderr, "Conflicting files:")
		for _, f := range result.ConflictFiles {
			fmt.Fprintf(os.Stderr, "  %s\n", f)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "The rebase was aborted. To resolve manually:")
		fmt.Fprintf(os.Stderr, "  cd %s\n", result.WorktreePath)
		fmt.Fprintf(os.Stderr, "  git rebase %s\n", result.ParentBranch)
		fmt.Fprintln(os.Stderr, "  # Fix conflicts in the listed files")
		fmt.Fprintln(os.Stderr, "  git add <resolved-files>")
		fmt.Fprintln(os.Stderr, "  git rebase --continue")
		os.Exit(1)
	}
	exitError("%v", err)
}
