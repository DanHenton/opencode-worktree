package main

import (
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runCompletions(args []string) {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		os.Exit(1)
	}

	if len(args) == 0 {
		for _, cmd := range []string{"task", "attach", "merge", "sync", "list", "cleanup"} {
			fmt.Println(cmd)
		}
		return
	}

	switch args[0] {
	case "attach":
		names, err := worktree.ActiveTaskNames(repoRoot)
		if err != nil {
			os.Exit(1)
		}
		for _, name := range names {
			fmt.Println(name)
		}
	}
}
