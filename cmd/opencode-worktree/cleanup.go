package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

func runCleanup(args []string) error {
	fs := flag.NewFlagSet("cleanup", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "Show what would be removed without removing anything")
	yes := fs.Bool("yes", false, "Skip confirmation prompt")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree cleanup [--dry-run] [--yes]

Remove orphaned agent worktrees and branches.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree cleanup
  opencode-worktree cleanup --dry-run
  opencode-worktree cleanup --yes
`)
	}

	if err := fs.Parse(args); err != nil {
		return errSilent
	}

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		return fmt.Errorf("not inside a git repository")
	}

	fmt.Printf("%sCleaning up orphaned agent worktrees and branches...\n", emoji("🧹 ", ""))
	opts := worktree.CleanupOptions{DryRun: *dryRun, Yes: *yes}
	if err := worktree.Cleanup(repoRoot, opts); err != nil {
		return fmt.Errorf("%v", err)
	}
	if !*dryRun {
		fmt.Printf("%sCleanup complete.\n", emoji("✅ ", ""))
	}
	return nil
}
