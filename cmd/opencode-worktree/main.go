package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via ldflags. Defaults to "dev" for local builds.
var version = "dev"

func main() {
	if err := newRootCmd().Execute(); err != nil {
		if !errors.Is(err, errSilent) {
			fmt.Fprint(os.Stderr, emoji("❌ ", "error: ")+err.Error()+"\n")
		}
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "opencode-worktree",
		Short:         "Git worktree manager for isolated OpenCode agent sessions",
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	cmd.AddCommand(
		newTaskCmd(),
		newAttachCmd(),
		newMergeCmd(),
		newSyncCmd(),
		newListCmd(),
		newCleanupCmd(),
	)

	return cmd
}
