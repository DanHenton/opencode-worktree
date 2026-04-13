package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunNoSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree"}

	err := run()
	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent, got %v", err)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "unknown-cmd"}

	r, w, _ := os.Pipe()
	origStderr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = origStderr; w.Close() }()

	err := run()

	w.Close()
	os.Stderr = origStderr

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if !errors.Is(err, errSilent) {
		t.Errorf("expected errSilent for unknown command, got %v", err)
	}
	if !strings.Contains(buf.String(), "Unknown command") {
		t.Errorf("expected stderr to contain 'Unknown command', got: %q", buf.String())
	}
}

func TestRunHelpShortFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "-h"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout; w.Close() }()

	err := run()

	w.Close()
	os.Stdout = origStdout
	io.Copy(io.Discard, r)

	if err != nil {
		t.Errorf("expected nil for -h, got %v", err)
	}
}

func TestRunHelpLongFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "--help"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout; w.Close() }()

	err := run()

	w.Close()
	os.Stdout = origStdout
	io.Copy(io.Discard, r)

	if err != nil {
		t.Errorf("expected nil for --help, got %v", err)
	}
}

func TestRunHelpSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "help"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout; w.Close() }()

	err := run()

	w.Close()
	os.Stdout = origStdout
	io.Copy(io.Discard, r)

	if err != nil {
		t.Errorf("expected nil for help subcommand, got %v", err)
	}
}

func TestRunVersionSubcommand(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "version"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout; w.Close() }()

	err := run()

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Errorf("expected nil for version subcommand, got %v", err)
	}
	if !strings.Contains(buf.String(), "opencode-worktree dev\n") {
		t.Errorf("expected version output 'opencode-worktree dev\\n', got: %q", buf.String())
	}
}

func TestRunVersionFlag(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"opencode-worktree", "--version"}

	r, w, _ := os.Pipe()
	origStdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = origStdout; w.Close() }()

	err := run()

	w.Close()
	os.Stdout = origStdout
	io.Copy(io.Discard, r)

	if err != nil {
		t.Errorf("expected nil for --version, got %v", err)
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
