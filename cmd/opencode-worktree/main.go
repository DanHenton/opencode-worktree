package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

// version is set at build time via ldflags. Defaults to "dev" for local builds.
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "task":
		runTask(os.Args[2:])
	case "attach":
		runAttach(os.Args[2:])
	case "merge":
		runMerge(os.Args[2:])
	case "list":
		runList()
	case "cleanup":
		runCleanup(os.Args[2:])
	case "sync":
		runSync(os.Args[2:])
	case "--completions":
		runCompletions(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
	case "version", "--version":
		fmt.Printf("opencode-worktree %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "%sUnknown command: %s\n\n", emoji("❌ ", "error: "), os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`opencode-worktree %s

Usage: opencode-worktree <command> [options]

Commands:
  task <name> [message]   Create agent worktree and launch opencode
  attach <name>           Reattach to an existing agent worktree session
  merge [path]            Merge agent branch back into parent
  sync [path]             Rebase agent branch onto latest parent
  list                    Show active agent worktrees
  cleanup                 Remove orphaned worktrees and branches

Run 'opencode-worktree <command> --help' for command-specific help.

General:
  -h, --help              Show this help message
  version, --version      Show version

Alias:
  The installer adds 'ocwt' as a shell alias for opencode-worktree.
`, version)
}

func runTask(args []string) {
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

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) == 0 {
		exitError("task name is required\n\nUsage: opencode-worktree task <name> [message] [--no-merge]")
	}
	if len(positional) > 2 {
		exitError("unexpected extra argument: %s", positional[2])
	}

	taskName := positional[0]
	var initialPrompt string
	if len(positional) > 1 {
		initialPrompt = positional[1]
	}

	if err := worktree.ValidateTaskName(taskName); err != nil {
		exitError("%v", err)
	}

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	parentBranch, err := git.CurrentBranch(repoRoot)
	if err != nil || parentBranch == "" {
		exitError("not on a named branch (detached HEAD) — checkout a branch first")
	}

	exists, err := worktree.AlreadyExists(repoRoot, taskName)
	if err != nil {
		exitError("%v", err)
	}
	if exists {
		exitError("a worktree for '%s%s' already exists — use 'opencode-worktree list' to see active sessions", worktree.BranchPrefix, taskName)
	}

	worktreeDir := worktree.WorktreeDir(repoRoot, taskName)
	branchName := worktree.BranchName(taskName)

	fmt.Printf("%sCreating worktree for task: %s\n", emoji("🌿 ", ""), taskName)
	fmt.Printf("   Branch:   %s\n", branchName)
	fmt.Printf("   From:     %s\n", parentBranch)
	fmt.Printf("   Path:     %s\n\n", worktreeDir)

	createdDir, err := worktree.Create(repoRoot, taskName, parentBranch)
	if err != nil {
		exitError("%v", err)
	}

	fmt.Printf("%sAgent session '%s' starting.\n", emoji("✅ ", ""), taskName)
	fmt.Printf("   Worktree: %s\n", createdDir)
	if *noMerge {
		fmt.Fprintf(os.Stderr, "   %s--no-merge is set. Run 'opencode-worktree merge' manually when done.\n", emoji("⚠️  ", "Note: "))
	}
	fmt.Println()

	_ = worktree.LaunchOpenCode(createdDir, initialPrompt)

	if *noMerge {
		return
	}

	fmt.Println()
	result, err := merge.Run(createdDir, true)
	if err != nil {
		handleMergeError(result, err)
	}
	printMergeResult(result)
}

func runAttach(args []string) {
	fs := flag.NewFlagSet("attach", flag.ContinueOnError)
	noMerge := fs.Bool("no-merge", false, "Skip auto-merge after opencode exits")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree attach <name> [--no-merge]

Reattach to an existing agent worktree session.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree attach fix-auth-bug
  opencode-worktree attach fix-auth-bug --no-merge
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) == 0 {
		exitError("task name is required\n\nUsage: opencode-worktree attach <name> [--no-merge]")
	}
	if len(positional) > 1 {
		exitError("unexpected extra argument: %s", positional[1])
	}

	taskName := positional[0]

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	worktreeDir, err := worktree.ResolveWorktreeDir(repoRoot, taskName)
	if err != nil {
		exitError("%v", err)
	}

	fmt.Printf("%sAttaching to agent session: %s\n", emoji("🔗 ", ""), taskName)
	fmt.Printf("   Path: %s\n\n", worktreeDir)

	_ = worktree.LaunchOpenCode(worktreeDir, "")

	if *noMerge {
		return
	}

	fmt.Println()
	result, err := merge.Run(worktreeDir, true)
	if err != nil {
		handleMergeError(result, err)
	}
	printMergeResult(result)
}

func runCompletions(args []string) {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		os.Exit(1)
	}

	if len(args) == 0 {
		for _, cmd := range []string{"task", "attach", "merge", "sync", "list", "cleanup"} {
			fmt.Println(cmd)
		}
		return
	}

	switch args[0] {
	case "attach":
		names, err := worktree.ActiveTaskNames(repoRoot)
		if err != nil {
			os.Exit(1)
		}
		for _, name := range names {
			fmt.Println(name)
		}
	}
}

func runMerge(args []string) {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	noCleanup := fs.Bool("no-cleanup", false, "Merge but keep worktree and branch")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree merge [path] [--no-cleanup]

Merge agent branch back into parent. If no path is given,
auto-detects the current directory as an agent worktree.

Options:
`)
		fs.PrintDefaults()
		fmt.Fprint(os.Stderr, `
Examples:
  opencode-worktree merge
  opencode-worktree merge /path/to/worktree
  opencode-worktree merge --no-cleanup
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) > 1 {
		exitError("unexpected extra argument: %s", positional[1])
	}

	var worktreePath string
	if len(positional) == 1 {
		worktreePath = positional[0]
	}

	if worktreePath == "" {
		detected, err := merge.DetectWorktree()
		if err != nil {
			exitError("%v\n\nUsage: opencode-worktree merge [worktree-path] [--no-cleanup]")
		}
		worktreePath = detected
	}

	cleanup := !*noCleanup
	result, err := merge.Run(worktreePath, cleanup)
	if err != nil {
		handleMergeError(result, err)
	}
	printMergeResult(result)
}

func runSync(args []string) {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: opencode-worktree sync [path]

Rebase the agent branch onto the latest parent branch, pulling in
upstream changes. If no path is given, auto-detects the current
directory as an agent worktree.

Examples:
  opencode-worktree sync
  opencode-worktree sync /path/to/worktree
`)
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	positional := fs.Args()
	if len(positional) > 1 {
		exitError("unexpected extra argument: %s", positional[1])
	}

	var worktreePath string
	if len(positional) == 1 {
		worktreePath = positional[0]
	}

	if worktreePath == "" {
		detected, err := merge.DetectWorktree()
		if err != nil {
			exitError("%v\n\nUsage: opencode-worktree sync [worktree-path]")
		}
		worktreePath = detected
	}

	result, err := worktree.Sync(worktreePath)
	if err != nil {
		handleSyncError(result, err)
	}

	if result.AlreadyCurrent {
		fmt.Printf("%sAlready up to date with %s.\n", emoji("✅ ", ""), result.ParentBranch)
		return
	}

	fmt.Printf("%sRebased %s onto %s.\n", emoji("✅ ", ""), result.AgentBranch, result.ParentBranch)
}

func handleSyncError(result *worktree.SyncResult, err error) {
	if result != nil && len(result.ConflictFiles) > 0 {
		fmt.Fprintf(os.Stderr, "%sRebase conflict: %s onto %s\n", emoji("❌ ", "error: "), result.AgentBranch, result.ParentBranch)
		fmt.Fprintln(os.Stderr, "Conflicting files:")
		for _, f := range result.ConflictFiles {
			fmt.Fprintf(os.Stderr, "  %s\n", f)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "The rebase was aborted. To resolve manually:")
		fmt.Fprintf(os.Stderr, "  cd %s\n", result.WorktreePath)
		fmt.Fprintf(os.Stderr, "  git rebase %s\n", result.ParentBranch)
		fmt.Fprintln(os.Stderr, "  # Fix conflicts in the listed files")
		fmt.Fprintln(os.Stderr, "  git add <resolved-files>")
		fmt.Fprintln(os.Stderr, "  git rebase --continue")
		os.Exit(1)
	}
	exitError("%v", err)
}

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

func runCleanup(args []string) {
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
		os.Exit(1)
	}

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	fmt.Printf("%sCleaning up orphaned agent worktrees and branches...\n", emoji("🧹 ", ""))
	opts := worktree.CleanupOptions{DryRun: *dryRun, Yes: *yes}
	if err := worktree.Cleanup(repoRoot, opts); err != nil {
		exitError("%v", err)
	}
	if !*dryRun {
		fmt.Printf("%sCleanup complete.\n", emoji("✅ ", ""))
	}
}

func printMergeResult(result *merge.Result) {
	if result.DirtyWorktree {
		fmt.Fprintf(os.Stderr, "%sWorktree has uncommitted changes — preserved at: %s\n", emoji("⚠️  ", "warning: "), result.WorktreePath)
		fmt.Fprintf(os.Stderr, "   Next: commit or discard changes, then run 'opencode-worktree merge %s'\n", result.WorktreePath)
	}
	if result.NoNewCommits && !result.DirtyWorktree {
		fmt.Fprintf(os.Stderr, "%sNo new commits found on %s. Cleaned up worktree only.\n", emoji("⚠️  ", "warning: "), result.AgentBranch)
		return
	}
	if result.Merged {
		if result.DirtyWorktree {
			fmt.Printf("%sMerged %s into %s (worktree kept due to uncommitted changes).\n", emoji("🚀 ", ""), result.AgentBranch, result.ParentBranch)
		} else {
			fmt.Printf("%sMerged %s into %s and cleaned up.\n", emoji("🚀 ", ""), result.AgentBranch, result.ParentBranch)
		}
	}
}

func handleMergeError(result *merge.Result, err error) {
	if result != nil && len(result.ConflictFiles) > 0 {
		fmt.Fprintf(os.Stderr, "%sMerge conflict: %s into %s\n", emoji("❌ ", "error: "), result.AgentBranch, result.ParentBranch)
		fmt.Fprintln(os.Stderr, "Conflicting files:")
		for _, f := range result.ConflictFiles {
			fmt.Fprintf(os.Stderr, "  %s\n", f)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "To resolve:")
		fmt.Fprintf(os.Stderr, "  cd %s\n", result.RepoRoot)
		fmt.Fprintln(os.Stderr, "  git status")
		fmt.Fprintln(os.Stderr, "  # Fix conflicts in the listed files")
		fmt.Fprintln(os.Stderr, "  git add <resolved-files>")
		fmt.Fprintln(os.Stderr, "  git commit")
		fmt.Fprintln(os.Stderr, "  opencode-worktree cleanup")
		os.Exit(1)
	}
	exitError("%v", err)
}

var useEmoji = detectTerminal()

func detectTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func emoji(e, fallback string) string {
	if useEmoji {
		return e
	}
	return fallback
}

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, emoji("❌ ", "error: ")+format+"\n", args...)
	os.Exit(1)
}
