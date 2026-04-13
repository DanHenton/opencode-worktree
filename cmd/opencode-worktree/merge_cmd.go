package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/merge"
)

func runMerge(args []string) {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	noCleanup := fs.Bool("no-cleanup", false, "Merge but keep worktree and branch")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree merge [path] [--no-cleanup]

Merge agent branch back into parent. If no path is given,
auto-detects the current directory as an agent worktree.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree merge
  opencode-worktree merge /path/to/worktree
  opencode-worktree merge --no-cleanup
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
			exitError("%v\n\nUsage: opencode-worktree merge [worktree-path] [--no-cleanup]", err)
		}
		worktreePath = detected
	}

	cleanup := !*noCleanup
	result, err := merge.Run(worktreePath, cleanup)
	if err != nil {
		handleMergeError(result, err)
	}
	printMergeResult(result)
}
