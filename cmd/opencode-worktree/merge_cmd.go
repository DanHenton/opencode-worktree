package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/spf13/cobra"
)

func newMergeCmd() *cobra.Command {
	var noCleanup bool

	cmd := &cobra.Command{
		Use:   "merge [worktree-path]",
		Short: "Merge agent branch back into parent",
		Long: `Merge agent branch back into parent. If no path is given,
auto-detects the current directory as an agent worktree.

Examples:
  opencode-worktree merge
  opencode-worktree merge /path/to/worktree
  opencode-worktree merge --no-cleanup`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var worktreePath string
			if len(args) == 1 {
				worktreePath = args[0]
			}

			if worktreePath == "" {
				detected, err := merge.DetectWorktree()
				if err != nil {
					return fmt.Errorf("%v\n\nUsage: opencode-worktree merge [worktree-path] [--no-cleanup]", err)
				}
				worktreePath = detected
			}

			cleanup := !noCleanup
			result, err := merge.Run(worktreePath, cleanup)
			if err != nil {
				if err := handleMergeError(result, err); err != nil {
					return err
				}
			}
			printMergeResult(result)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&noCleanup, "no-cleanup", "c", false, "Merge but keep worktree and branch")

	return cmd
}
