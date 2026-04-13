# Command Reference

Command-level reference for `opencode-worktree`. Each page follows the same structure so a human or an agent can understand behavior, inputs, and edge cases without reading the Go code.

## Commands

- [`task`](task.md) — create a new agent worktree and launch OpenCode
- [`attach`](attach.md) — reopen an existing agent worktree session
- [`merge`](merge.md) — merge an agent branch back into its parent branch
- [`sync`](sync.md) — rebase an agent branch onto the latest parent branch
- [`list`](list.md) — show active agent worktrees
- [`cleanup`](cleanup.md) — remove orphaned worktrees and stale agent branches

## Shared conventions

- Agent branches are named `agent/<task-name>`.
- Agent worktrees are created as sibling directories named `<repo>-agent-<task-name>`.
- Managed worktrees contain `.agent-parent-branch` and `.agent-context` marker files.
- The installer adds `ocwt` as a shell alias for `opencode-worktree`.
- Run `opencode-worktree <command> --help` for built-in CLI help.
