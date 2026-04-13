# opencode-worktree

Git worktree manager for isolated [OpenCode](https://opencode.ai) agent sessions. Create an isolated worktree, launch OpenCode in it, and auto-merge back when done.

## Why?

When using AI coding agents on complex tasks, you want isolation — the agent works on a branch in a separate directory, and you merge its changes back only when ready. This tool automates the entire lifecycle:

1. **Create** an `agent/<task>` branch in a sibling worktree directory
2. **Launch** OpenCode TUI in that worktree
3. **Merge** the agent branch back into your parent branch on exit
4. **Clean up** the worktree and branch automatically

## Install

### Pre-compiled Binaries

Download the latest release for your OS from the [Releases Page](https://github.com/DanHenton/opencode-worktree/releases).

```bash
tar -xzf opencode-worktree_Linux_x86_64.tar.gz
sudo mv opencode-worktree /usr/local/bin/
```

### Build from Source

```bash
go install github.com/danhenton/opencode-worktree/cmd/opencode-worktree@latest
```

## Usage

### Start an agent task

```bash
opencode-worktree task fix-auth-bug
```

This will:

- Create worktree at `../your-repo-agent-fix-auth-bug/`
- Create branch `agent/fix-auth-bug` from your current branch
- Copy `opencode.json` and `.opencode/` into the worktree
- Launch `opencode` in the worktree
- Auto-merge back into your branch when you exit OpenCode

### Send an initial prompt

```bash
opencode-worktree task fix-auth-bug "Fix the JWT token expiry bug in the auth middleware"
```

### Skip auto-merge

```bash
opencode-worktree task add-dark-mode --no-merge
```

### Merge manually

From inside the agent worktree:

```bash
opencode-worktree merge
```

Or from anywhere:

```bash
opencode-worktree merge /path/to/worktree
```

Keep the worktree after merge:

```bash
opencode-worktree merge --no-cleanup
```

### List active sessions

```bash
opencode-worktree list
```

### Clean up orphaned worktrees

```bash
opencode-worktree cleanup
```

Removes stale worktree directories and agent branches that no longer have an active worktree.

## How It Works

### Directory Layout

```
~/Development/
├── your-project/                        # Main repo (you're here)
├── your-project-agent-fix-auth/         # Agent worktree (created by tool)
└── your-project-agent-add-feature/      # Another agent worktree
```

### Marker Files

Each worktree gets two marker files:

- `.agent-parent-branch` — the branch to merge back into
- `.agent-context` — metadata for AI agents (worktree dir, parent branch, agent branch, source repo)

### Merge Safety

- Merges are serialized with a file lock (`/tmp/<repo-name>-merge.lock`) to prevent races when multiple agents merge simultaneously
- On conflict, the merge is aborted and conflicting files are listed
- Empty branches (no new commits) skip the merge and clean up directly

## AGENTS.md Integration

Add this section to your project's `AGENTS.md` to give AI agents context about their worktree environment:

```markdown
## Worktree Agent Sessions

When running in an agent worktree session, you are working on an isolated
`agent/<task-name>` branch.

### Path Discipline

- Read `.agent-context` to confirm your WORKTREE_DIR
- NEVER edit files in the SOURCE_REPO directory
- Make commits normally — they're isolated to your agent branch

### On-demand Merge

Run `opencode-worktree merge` from within the worktree to merge back.
```

## License

MIT
