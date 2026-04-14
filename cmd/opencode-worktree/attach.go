package main

import (
	"fmt"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

func newAttachCmd() *cobra.Command {
	var noMerge bool

	cmd := &cobra.Command{
		Use:   "attach <name>",
		Short: "Reattach to an existing agent worktree session",
		Long: `Reattach to an existing agent worktree session.

Examples:
  opencode-worktree attach fix-auth-bug
  opencode-worktree attach fix-auth-bug --no-merge`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("task name is required\n\nUsage: opencode-worktree attach <name> [--no-merge]")
			}
			if len(args) > 1 {
				return fmt.Errorf("unexpected extra argument: %s", args[1])
			}
			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			repoRoot, err := git.RepoRoot(".")
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			names, err := worktree.ActiveTaskNames(repoRoot)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return names, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]

			repoRoot, err := git.RepoRoot(".")
			if err != nil {
				return fmt.Errorf("not inside a git repository")
			}

			worktreeDir, err := worktree.ResolveWorktreeDir(repoRoot, taskName)
			if err != nil {
				return err
			}

			fmt.Printf("%sAttaching to agent session: %s\n", emoji("🔗 ", ""), taskName)
			fmt.Printf("   Path: %s\n\n", worktreeDir)

			return launchAndMaybeMerge(worktreeDir, "", noMerge)
		},
	}

	cmd.Flags().BoolVarP(&noMerge, "no-merge", "n", false, "Skip auto-merge after opencode exits")

	return cmd
}
