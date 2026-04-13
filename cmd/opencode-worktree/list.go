package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runList(args []string) error {
	_ = args

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		return fmt.Errorf("not inside a git repository")
	}

	fmt.Printf("%sActive agent worktrees:\n", emoji("🗂️  ", ""))
	out, err := worktree.List(repoRoot)
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}
