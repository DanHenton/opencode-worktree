# Installation

## For Humans

Paste this into your LLM agent session:

```
Install opencode-worktree by following the instructions here:
https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/docs/install.md
```

**Or install it yourself:**

```bash
curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/install.sh | sh
```

## For LLM Agents

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

Run the installer script. This auto-detects OS (macOS/Linux) and architecture (x86_64/arm64/armv7/i386), downloads the correct release from GitHub, and places it in `~/.local/bin/`.

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

This should display the available subcommands: `task`, `merge`, `list`, `cleanup`.

### Step 5: Confirm with user

Tell the user installation is complete. Show them a quick usage example:

```bash
# Start an isolated agent session
opencode-worktree task my-task-name

# Start with an initial prompt
opencode-worktree task my-task-name "Description of what to work on"
```

### Requirements

- **OS**: macOS or Linux (Windows: download manually from [Releases](https://github.com/DanHenton/opencode-worktree/releases))
- **Dependencies**: `git`, `opencode` (the OpenCode TUI must be installed separately)
- **Optional**: Go 1.24+ (only for `go install` method)
