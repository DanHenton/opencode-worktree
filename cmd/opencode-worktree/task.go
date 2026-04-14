package main

import (
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
	"github.com/spf13/cobra"
)

func newTaskCmd() *cobra.Command {
	var noMerge bool

	cmd := &cobra.Command{
		Use:   "task <name> [message]",
		Short: "Create agent worktree and launch opencode",
		Long: `Create an agent worktree and launch opencode in it.

Examples:
  opencode-worktree task fix-auth-bug
  opencode-worktree task fix-auth-bug "Fix the JWT token expiry bug"
  opencode-worktree task add-feature --no-merge`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("task name is required\n\nUsage: opencode-worktree task <name> [message] [--no-merge]")
			}
			if len(args) > 2 {
				return fmt.Errorf("unexpected extra argument: %s", args[2])
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			taskName := args[0]
			var initialPrompt string
			if len(args) > 1 {
				initialPrompt = args[1]
			}

			if err := worktree.ValidateTaskName(taskName); err != nil {
				return err
			}

			repoRoot, err := git.RepoRoot(".")
			if err != nil {
				return fmt.Errorf("not inside a git repository")
			}

			parentBranch, err := git.CurrentBranch(repoRoot)
			if err != nil || parentBranch == "" {
				return fmt.Errorf("not on a named branch (detached HEAD) — checkout a branch first")
			}

			exists, err := worktree.AlreadyExists(repoRoot, taskName)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("a worktree for '%s%s' already exists — use 'opencode-worktree list' to see active sessions", worktree.BranchPrefix, taskName)
			}

			worktreeDir := worktree.WorktreeDir(repoRoot, taskName)
			branchName := worktree.BranchName(taskName)

			fmt.Printf("%sCreating worktree for task: %s\n", emoji("🌿 ", ""), taskName)
			fmt.Printf("   Branch:   %s\n", branchName)
			fmt.Printf("   From:     %s\n", parentBranch)
			fmt.Printf("   Path:     %s\n\n", worktreeDir)

			createdDir, err := worktree.Create(repoRoot, taskName, parentBranch)
			if err != nil {
				return err
			}

			fmt.Printf("%sAgent session '%s' starting.\n", emoji("✅ ", ""), taskName)
			fmt.Printf("   Worktree: %s\n", createdDir)
			if noMerge {
				fmt.Fprintf(os.Stderr, "   %s--no-merge is set. Run 'opencode-worktree merge' manually when done.\n", emoji("⚠️  ", "Note: "))
			}
			fmt.Println()

			return launchAndMaybeMerge(createdDir, initialPrompt, noMerge)
		},
	}

	cmd.Flags().BoolVarP(&noMerge, "no-merge", "n", false, "Skip auto-merge after opencode exits")

	return cmd
}
