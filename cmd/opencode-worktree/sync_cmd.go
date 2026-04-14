package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

func newSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync [worktree-path]",
		Short: "Rebase agent branch onto latest parent",
		Long: `Rebase the agent branch onto the latest parent branch, pulling in
upstream changes. If no path is given, auto-detects the current
directory as an agent worktree.

Examples:
  opencode-worktree sync
  opencode-worktree sync /path/to/worktree`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var worktreePath string
			if len(args) == 1 {
				worktreePath = args[0]
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
		},
	}
	return cmd
}
