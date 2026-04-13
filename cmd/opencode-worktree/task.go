package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runTask(args []string) error {
	fs := flag.NewFlagSet("task", flag.ContinueOnError)
	noMerge := fs.Bool("no-merge", false, "Skip auto-merge after opencode exits")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree task <name> [message] [--no-merge]

Create an agent worktree and launch opencode in it.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree task fix-auth-bug
  opencode-worktree task fix-auth-bug "Fix the JWT token expiry bug"
  opencode-worktree task add-feature --no-merge
`)
	}

	if err := fs.Parse(reorderKnownBoolFlags(args, "--no-merge")); err != nil {
		return errSilent
	}

	positional := fs.Args()
	if len(positional) == 0 {
		return fmt.Errorf("task name is required\n\nUsage: opencode-worktree task <name> [message] [--no-merge]")
	}
	if len(positional) > 2 {
		return fmt.Errorf("unexpected extra argument: %s", positional[2])
	}

	taskName := positional[0]
	var initialPrompt string
	if len(positional) > 1 {
		initialPrompt = positional[1]
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
	if *noMerge {
		fmt.Fprintf(os.Stderr, "   %s--no-merge is set. Run 'opencode-worktree merge' manually when done.\n", emoji("⚠️  ", "Note: "))
	}
	fmt.Println()

	return launchAndMaybeMerge(createdDir, initialPrompt, *noMerge)
}
