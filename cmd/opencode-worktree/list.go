package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runList() {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	fmt.Printf("%sActive agent worktrees:\n", emoji("🗂️  ", ""))
	out, err := worktree.List(repoRoot)
	if err != nil {
		exitError("%v", err)
	}
	fmt.Println(out)
}
