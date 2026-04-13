package main

import (
	"errors"
	"fmt"
	"os"
)

// version is set at build time via ldflags. Defaults to "dev" for local builds.
var version = "dev"

func main() {
	if err := run(); err != nil {
		if !errors.Is(err, errSilent) {
			fmt.Fprint(os.Stderr, emoji("❌ ", "error: ")+err.Error()+"\n")
		}
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return errSilent
	}

	switch os.Args[1] {
	case "task":
		return runTask(os.Args[2:])
	case "attach":
		return runAttach(os.Args[2:])
	case "merge":
		return runMerge(os.Args[2:])
	case "list":
		return runList(os.Args[2:])
	case "cleanup":
		return runCleanup(os.Args[2:])
	case "sync":
		return runSync(os.Args[2:])
	case "--completions":
		return runCompletions(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
		return nil
	case "version", "--version":
		fmt.Printf("opencode-worktree %s\n", version)
		return nil
	default:
		fmt.Fprintf(os.Stderr, "%sUnknown command: %s\n\n", emoji("❌ ", "error: "), os.Args[1])
		printUsage()
		return errSilent
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
