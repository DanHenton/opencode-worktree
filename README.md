# opencode-worktree

Git worktree manager for isolated [OpenCode](https://opencode.ai) agent sessions.

Run multiple AI agents in parallel — each gets its own branch and directory, fully isolated from your working tree. When they're done, changes merge back automatically. No stashing, no branch juggling, no conflicts between agents.

## Why?

AI coding agents work best with full autonomy over their environment. But you don't want them trampling your working directory — or each other. `opencode-worktree` gives each agent an isolated workspace and handles the entire lifecycle:

1. **Create** an `agent/<task>` branch in a sibling worktree directory
2. **Launch** OpenCode in that worktree with your config copied over
3. **Merge** the agent branch back into your parent branch on exit
4. **Clean up** the worktree and branch automatically

No manual git commands. No orphaned branches. Multiple agents can run simultaneously with safe, serialized merges.

## Install

### Quick Install (macOS / Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/install.sh | sh
```

### LLM Agent Install

Paste this into your LLM agent session:

```
Install opencode-worktree by following the instructions here:
https://raw.githubusercontent.com/DanHenton/opencode-worktree/main/docs/install.md
```

For other install methods (Go, manual download, PATH troubleshooting), see the [full install guide](docs/install.md).

## Usage

### Start an agent task

```bash
opencode-worktree task fix-auth-bug
```

This will:

- Create worktree at `../your-repo-agent-fix-auth-bug/`
- Create branch `agent/fix-auth-bug` from your current branch
- Copy `opencode.json` and `.opencode/` into the worktree
- Detect dependency manifests (`package.json`, `go.mod`, etc.) and tell the agent to install them
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

When `opencode-worktree` launches an agent session, it automatically injects worktree context into the initial message — including branch info, detected dependency manifests, and install commands. No configuration needed.

Optionally, add this section to your project's `AGENTS.md` for additional guidance:

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
