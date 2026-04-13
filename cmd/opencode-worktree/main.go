package main

import (
	"fmt"
	"os"
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
