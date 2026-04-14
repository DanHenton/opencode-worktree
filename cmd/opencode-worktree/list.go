package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show active agent worktrees",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
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
		},
	}
	return cmd
}
