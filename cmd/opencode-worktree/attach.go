package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runAttach(args []string) {
	fs := flag.NewFlagSet("attach", flag.ContinueOnError)
	noMerge := fs.Bool("no-merge", false, "Skip auto-merge after opencode exits")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree attach <name> [--no-merge]

Reattach to an existing agent worktree session.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree attach fix-auth-bug
  opencode-worktree attach fix-auth-bug --no-merge
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) == 0 {
		exitError("task name is required\n\nUsage: opencode-worktree attach <name> [--no-merge]")
	}
	if len(positional) > 1 {
		exitError("unexpected extra argument: %s", positional[1])
	}

	taskName := positional[0]

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	worktreeDir, err := worktree.ResolveWorktreeDir(repoRoot, taskName)
	if err != nil {
		exitError("%v", err)
	}

	fmt.Printf("%sAttaching to agent session: %s\n", emoji("🔗 ", ""), taskName)
	fmt.Printf("   Path: %s\n\n", worktreeDir)

	_ = worktree.LaunchOpenCode(worktreeDir, "")

	if *noMerge {
		return
	}

	fmt.Println()
	result, err := merge.Run(worktreeDir, true)
	if err != nil {
		handleMergeError(result, err)
	}
	printMergeResult(result)
}
