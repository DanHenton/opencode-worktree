package main

import (
	"fmt"
	"os"

	"strings"

	"github.com/danhenton/opencode-worktree/internal/git"
	"github.com/danhenton/opencode-worktree/internal/merge"
	"github.com/danhenton/opencode-worktree/internal/worktree"
)

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
		runCleanup()
	case "--completions":
		runCompletions(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "❌ Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Usage: opencode-worktree <command> [options]

Commands:
  task <name> [message]   Create agent worktree and launch opencode
  attach <name>           Reattach to an existing agent worktree session
  merge [path]            Merge agent branch back into parent
  list                    Show active agent worktrees
  cleanup                 Remove orphaned worktrees and branches

Task Options:
  --no-merge              Skip auto-merge after opencode exits

Attach Options:
  --no-merge              Skip auto-merge after opencode exits

Merge Options:
  --no-cleanup            Merge but keep worktree and branch

General:
  -h, --help              Show this help message

Alias:
  The installer adds 'ocwt' as a shell alias for opencode-worktree.
`)
}

func runTask(args []string) {
	var taskName, initialPrompt string
	noMerge := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-merge":
			noMerge = true
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		default:
			if len(args[i]) > 0 && args[i][0] == '-' {
				exitError("unknown option: %s", args[i])
			}
			if taskName == "" {
				taskName = args[i]
			} else if initialPrompt == "" {
				initialPrompt = args[i]
			} else {
				exitError("unexpected extra argument: %s", args[i])
			}
		}
	}

	if taskName == "" {
		exitError("task name is required\n\nUsage: opencode-worktree task <name> [message] [--no-merge]")
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

	if worktree.AlreadyExists(repoRoot, taskName) {
		exitError("a worktree for '%s%s' already exists — use 'opencode-worktree list' to see active sessions", worktree.BranchPrefix, taskName)
	}

	worktreeDir := worktree.WorktreeDir(repoRoot, taskName)
	branchName := worktree.BranchName(taskName)

	fmt.Printf("🌿 Creating worktree for task: %s\n", taskName)
	fmt.Printf("   Branch:   %s\n", branchName)
	fmt.Printf("   From:     %s\n", parentBranch)
	fmt.Printf("   Path:     %s\n\n", worktreeDir)

	createdDir, err := worktree.Create(repoRoot, taskName, parentBranch)
	if err != nil {
		exitError("%v", err)
	}

	fmt.Printf("✅ Agent session '%s' starting.\n", taskName)
	fmt.Printf("   Worktree: %s\n", createdDir)
	if noMerge {
		fmt.Println("   ⚠️  --no-merge is set. Run 'opencode-worktree merge' manually when done.")
	}
	fmt.Println()

	_ = worktree.LaunchOpenCode(createdDir, initialPrompt)

	if noMerge {
		return
	}

	fmt.Println()
	result, err := merge.Run(createdDir, true)
	if err != nil {
		if result != nil && len(result.ConflictFiles) > 0 {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			fmt.Fprintln(os.Stderr, "Conflicting files:")
			for _, f := range result.ConflictFiles {
				fmt.Fprintf(os.Stderr, "  %s\n", f)
			}
			os.Exit(1)
		}
		exitError("%v", err)
	}
	printMergeResult(result)
}

func runAttach(args []string) {
	var taskName string
	noMerge := false

	for _, arg := range args {
		switch arg {
		case "--no-merge":
			noMerge = true
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		default:
			if len(arg) > 0 && arg[0] == '-' {
				exitError("unknown option: %s", arg)
			}
			if taskName == "" {
				taskName = arg
			} else {
				exitError("unexpected extra argument: %s", arg)
			}
		}
	}

	if taskName == "" {
		exitError("task name is required\n\nUsage: opencode-worktree attach <name> [--no-merge]")
	}

	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	worktreeDir, err := worktree.ResolveWorktreeDir(repoRoot, taskName)
	if err != nil {
		exitError("%v", err)
	}

	fmt.Printf("🔗 Attaching to agent session: %s\n", taskName)
	fmt.Printf("   Path: %s\n\n", worktreeDir)

	_ = worktree.LaunchOpenCode(worktreeDir, "")

	if noMerge {
		return
	}

	fmt.Println()
	result, err := merge.Run(worktreeDir, true)
	if err != nil {
		if result != nil && len(result.ConflictFiles) > 0 {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			fmt.Fprintln(os.Stderr, "Conflicting files:")
			for _, f := range result.ConflictFiles {
				fmt.Fprintf(os.Stderr, "  %s\n", f)
			}
			os.Exit(1)
		}
		exitError("%v", err)
	}
	printMergeResult(result)
}

func runCompletions(args []string) {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		os.Exit(1)
	}

	if len(args) == 0 {
		fmt.Println(strings.Join([]string{"task", "attach", "merge", "list", "cleanup"}, "\n"))
		return
	}

	switch args[0] {
	case "attach":
		names, err := worktree.ActiveTaskNames(repoRoot)
		if err != nil {
			os.Exit(1)
		}
		if len(names) > 0 {
			fmt.Println(strings.Join(names, "\n"))
		}
	}
}

func runMerge(args []string) {
	var worktreePath string
	noCleanup := false

	for _, arg := range args {
		switch arg {
		case "--no-cleanup":
			noCleanup = true
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		default:
			if len(arg) > 0 && arg[0] == '-' {
				exitError("unknown option: %s", arg)
			}
			if worktreePath == "" {
				worktreePath = arg
			} else {
				exitError("unexpected extra argument: %s", arg)
			}
		}
	}

	if worktreePath == "" {
		detected, err := merge.DetectWorktree()
		if err != nil {
			exitError("%v\n\nUsage: opencode-worktree merge [worktree-path] [--no-cleanup]")
		}
		worktreePath = detected
	}

	cleanup := !noCleanup
	result, err := merge.Run(worktreePath, cleanup)
	if err != nil {
		if result != nil && len(result.ConflictFiles) > 0 {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			fmt.Fprintln(os.Stderr, "Conflicting files:")
			for _, f := range result.ConflictFiles {
				fmt.Fprintf(os.Stderr, "  %s\n", f)
			}
			os.Exit(1)
		}
		exitError("%v", err)
	}
	printMergeResult(result)
}

func runList() {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	fmt.Println("🗂️  Active agent worktrees:")
	out, err := worktree.List(repoRoot)
	if err != nil {
		exitError("%v", err)
	}
	fmt.Println(out)
}

func runCleanup() {
	repoRoot, err := git.RepoRoot(".")
	if err != nil {
		exitError("not inside a git repository")
	}

	fmt.Println("🧹 Cleaning up orphaned agent worktrees and branches...")
	if err := worktree.Cleanup(repoRoot); err != nil {
		exitError("%v", err)
	}
	fmt.Println("✅ Cleanup complete.")
}

func printMergeResult(result *merge.Result) {
	if result.DirtyWorktree {
		fmt.Printf("⚠️  Worktree has uncommitted changes — preserved at: %s\n", result.WorktreePath)
		fmt.Println("   Commit or discard your changes, then run 'opencode-worktree merge' to finish.")
	}
	if result.NoNewCommits && !result.DirtyWorktree {
		fmt.Printf("⚠️  No new commits found on %s. Cleaned up worktree only.\n", result.AgentBranch)
		return
	}
	if result.Merged {
		if result.DirtyWorktree {
			fmt.Printf("🚀 Merged %s into %s (worktree kept due to uncommitted changes).\n", result.AgentBranch, result.ParentBranch)
		} else {
			fmt.Printf("🚀 Merged %s into %s and cleaned up.\n", result.AgentBranch, result.ParentBranch)
		}
	}
}

func exitError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}
