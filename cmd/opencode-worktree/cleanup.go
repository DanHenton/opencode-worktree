package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

func newCleanupCmd() *cobra.Command {
	var dryRun bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Remove orphaned worktrees and branches",
		Long: `Remove orphaned agent worktrees and branches.

Examples:
  opencode-worktree cleanup
  opencode-worktree cleanup --dry-run
  opencode-worktree cleanup --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRoot, err := git.RepoRoot(".")
			if err != nil {
				return fmt.Errorf("not inside a git repository")
			}

			fmt.Printf("%sCleaning up orphaned agent worktrees and branches...\n", emoji("🧹 ", ""))
			opts := worktree.CleanupOptions{DryRun: dryRun, Yes: yes}
			if err := worktree.Cleanup(repoRoot, opts); err != nil {
				return err
			}
			if !dryRun {
				fmt.Printf("%sCleanup complete.\n", emoji("✅ ", ""))
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be removed without removing anything")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
