# Install Guide

All methods for installing `opencode-worktree`. For a quick overview, see the [README](../README.md).

## Quick Install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/install.sh | sh
```

Auto-detects OS and architecture, downloads the correct release from GitHub, and places it in `~/.local/bin/`.

Override the install directory:

```bash
curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/install.sh | BIN_DIR=/usr/local/bin sh
```

## Install with Go

Requires Go 1.24+.

```bash
go install github.com/danhenton/opencode-worktree/cmd/opencode-worktree@latest
```

## Manual Download

Download a release directly from the [Releases Page](https://github.com/DanHenton/opencode-worktree/releases):

```bash
# macOS (Apple Silicon)
curl -fsSLO https://github.com/DanHenton/opencode-worktree/releases/latest/download/opencode-worktree_Darwin_arm64.tar.gz

# macOS (Intel)
curl -fsSLO https://github.com/DanHenton/opencode-worktree/releases/latest/download/opencode-worktree_Darwin_x86_64.tar.gz

# Linux (x86_64)
curl -fsSLO https://github.com/DanHenton/opencode-worktree/releases/latest/download/opencode-worktree_Linux_x86_64.tar.gz

# Linux (ARM64)
curl -fsSLO https://github.com/DanHenton/opencode-worktree/releases/latest/download/opencode-worktree_Linux_arm64.tar.gz
```

Extract and install:

```bash
tar -xzf opencode-worktree_*.tar.gz
mkdir -p ~/.local/bin
install -m 0755 opencode-worktree ~/.local/bin/opencode-worktree
```

## PATH Troubleshooting

If `opencode-worktree` isn't found after install, the install directory isn't in your PATH.

**If installed via install.sh or manual download** (default: `~/.local/bin`):

```bash
# zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc

# bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

**If installed via `go install`:**

```bash
# zsh
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc

# bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
```

Then open a new terminal, or run `export PATH=...` directly to apply immediately.

## Requirements

- **OS**: macOS or Linux (Windows: download manually from [Releases](https://github.com/DanHenton/opencode-worktree/releases))
- **Dependencies**: `git`, `opencode` (the [OpenCode TUI](https://opencode.ai) must be installed separately)
- **Optional**: Go 1.24+ (only for `go install` method)

---

## LLM Agent Install

> **IMPORTANT: Use `curl` to fetch this file, NOT WebFetch.** WebFetch may summarize content and lose critical details. Always use:
>
> ```bash
> curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/docs/install.md
> ```

You are helping a user install `opencode-worktree`, a Go CLI tool for managing git worktrees with isolated OpenCode agent sessions.

### Step 1: Check if already installed

```bash
if command -v opencode-worktree &> /dev/null; then
    echo "opencode-worktree $(opencode-worktree --help 2>&1 | head -1) is already installed"
else
    echo "opencode-worktree is not installed"
fi
```

If already installed, skip to Step 4.

### Step 2: Install the binary

Run the installer script:

```bash
curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/install.sh | sh
```

If the user prefers `go install` and has Go 1.24+ available:

```bash
go install github.com/danhenton/opencode-worktree/cmd/opencode-worktree@latest
```

### Step 3: Verify PATH

After install, verify the binary is accessible:

```bash
command -v opencode-worktree
```

If not found, the install directory is not in the user's PATH. Fix it:

**If installed via install.sh** (default: `~/.local/bin`):

```bash
# Detect shell and add to profile
if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "*/zsh" ]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "*/bash" ]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
fi
export PATH="$HOME/.local/bin:$PATH"
```

**If installed via `go install`:**

```bash
GO_BIN="$(go env GOPATH)/bin"
if [ -n "$ZSH_VERSION" ] || [ "$SHELL" = "*/zsh" ]; then
    printf 'export PATH="%s:$PATH"\n' "$GO_BIN" >> ~/.zshrc
elif [ -n "$BASH_VERSION" ] || [ "$SHELL" = "*/bash" ]; then
    printf 'export PATH="%s:$PATH"\n' "$GO_BIN" >> ~/.bashrc
fi
export PATH="$GO_BIN:$PATH"
```

Verify again:

```bash
command -v opencode-worktree
```

### Step 4: Verify it works

```bash
opencode-worktree --help
```

This should display the available subcommands: `task`, `attach`, `merge`, `sync`, `list`, `cleanup`.

### Step 5: Confirm with user

Tell the user installation is complete. Mention the `ocwt` shortcut alias that was installed automatically:

```bash
# Start an isolated agent session
opencode-worktree task my-task-name

# Start with an initial prompt
opencode-worktree task my-task-name "Description of what to work on"
```

> **Tip:** The installer adds `ocwt` as a shell alias for `opencode-worktree` — use either interchangeably. If you installed manually or via `go install`, add it yourself:
> ```bash
> echo "alias ocwt='opencode-worktree'" >> ~/.zshrc  # or ~/.bashrc
> ```
