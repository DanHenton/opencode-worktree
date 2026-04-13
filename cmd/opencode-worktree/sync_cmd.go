package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runSync(args []string) error {
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
		return errSilent
	}

	positional := fs.Args()
	if len(positional) > 1 {
		return fmt.Errorf("unexpected extra argument: %s", positional[1])
	}

	var worktreePath string
	if len(positional) == 1 {
		worktreePath = positional[0]
	}

	if worktreePath == "" {
		detected, err := merge.DetectWorktree()
		if err != nil {
			return fmt.Errorf("%v\n\nUsage: opencode-worktree sync [worktree-path]", err)
		}
		worktreePath = detected
	}

	result, err := worktree.Sync(worktreePath)
	if err != nil {
		if err := handleSyncError(result, err); err != nil {
			return err
		}
	}

	if result.AlreadyCurrent {
		fmt.Printf("%sAlready up to date with %s.\n", emoji("✅ ", ""), result.ParentBranch)
		return nil
	}

	fmt.Printf("%sRebased %s onto %s.\n", emoji("✅ ", ""), result.AgentBranch, result.ParentBranch)
	return nil
}
