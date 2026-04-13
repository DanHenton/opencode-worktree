package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runCompletions(args []string) error {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		return errSilent
	}

	if len(args) == 0 {
		for _, cmd := range []string{"task", "attach", "merge", "sync", "list", "cleanup"} {
			fmt.Println(cmd)
		}
		return nil
	}

	switch args[0] {
	case "attach":
		names, err := worktree.ActiveTaskNames(repoRoot)
		if err != nil {
			return errSilent
		}
		for _, name := range names {
			fmt.Println(name)
		}
	}

	return nil
}
