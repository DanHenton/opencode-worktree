# `opencode-worktree list`

## Summary

Shows active managed agent worktrees for the current repository.

## Usage

```bash
opencode-worktree list
```

## Arguments

This command does not take positional arguments.

## Options

This command has no flags.

## What it does

1. Confirms you are inside a git repository.
2. Reads `git worktree list` output for the repository.
3. Filters the result down to worktrees whose branch name starts with `agent/`.
4. Marks entries with `(uncommitted changes)` when the worktree contains changes outside the tool's marker/config files.

## When to use it

Use `list` to see which agent sessions are still active and to find task names you can pass to `attach`.

## Examples

```bash
opencode-worktree list
```

## Important behavior

- If no managed agent worktrees exist, the command prints `(none)`.
- The output is informational and does not modify anything.
- Dirty-state detection ignores `.agent-parent-branch`, `.agent-context`, `opencode.json`, and `.opencode/`.

## Common failure cases

- Running outside a git repository.

## Related commands

- [`attach`](attach.md) to reopen a listed session.
- [`cleanup`](cleanup.md) to remove stale sessions that are no longer active.
