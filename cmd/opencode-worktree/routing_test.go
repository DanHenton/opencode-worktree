package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRouteNoSubcommand(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil error for no subcommand, got %v", err)
	}
}

func TestRouteUnknownCommand(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"unknown-cmd"})
	err := root.Execute()
	if err == nil {
		t.Errorf("expected non-nil error for unknown command, got nil")
	}
}

func TestRouteHelpShortFlag(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"-h"})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil for -h, got %v", err)
	}
}

func TestRouteHelpLongFlag(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"--help"})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil for --help, got %v", err)
	}
}

func TestRouteHelpSubcommand(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"help"})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil for help subcommand, got %v", err)
	}
}

func TestRouteVersionFlag(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"--version"})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil for --version, got %v", err)
	}
	if !strings.Contains(outBuf.String(), "opencode-worktree") {
		t.Errorf("expected version output to contain 'opencode-worktree', got: %q", outBuf.String())
	}
}

func TestRouteVersionShortFlag(t *testing.T) {
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs([]string{"-v"})
	err := root.Execute()
	if err != nil {
		t.Errorf("expected nil for -v, got %v", err)
	}
	if !strings.Contains(outBuf.String(), "opencode-worktree") {
		t.Errorf("expected version output to contain 'opencode-worktree', got: %q", outBuf.String())
	}
}

func TestEmojiDisabled(t *testing.T) {
	origEmoji := useEmoji
	defer func() { useEmoji = origEmoji }()
	useEmoji = false

	if got := emoji("🌿 ", ""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
	if got := emoji("❌ ", "error: "); got != "error: " {
		t.Errorf("expected 'error: ', got %q", got)
	}
}

func TestEmojiEnabled(t *testing.T) {
	origEmoji := useEmoji
	defer func() { useEmoji = origEmoji }()
	useEmoji = true

	if got := emoji("🌿 ", ""); got != "🌿 " {
		t.Errorf("expected '🌿 ', got %q", got)
	}
}

func TestDetectTerminalInTestContext(t *testing.T) {
	if detectTerminal() {
		t.Errorf("expected detectTerminal() to return false in test context (stdout is a pipe)")
	}
}

func TestErrSilentSentinel(t *testing.T) {
	if !errors.Is(errSilent, errSilent) {
		t.Errorf("expected errors.Is(errSilent, errSilent) to be true")
	}
	if errSilent.Error() != "" {
		t.Errorf("expected errSilent.Error() to be empty string, got %q", errSilent.Error())
	}
}
