# `opencode-worktree sync`

## Summary

Rebases an agent branch onto the latest version of its parent branch so the worktree picks up upstream changes before merge time.

## Usage

```bash
opencode-worktree sync [path]
```

## Arguments

- `[path]` — optional path to an agent worktree. If omitted, the command tries to detect the current directory as a managed agent worktree.

## Options

This command has no flags.

## What it does

1. Resolves the target worktree path, either from the argument or by auto-detecting the current directory.
2. Reads `.agent-parent-branch` to determine which parent branch to sync against.
3. Confirms the worktree is on an `agent/*` branch.
4. Refuses to continue if the worktree has uncommitted changes.
5. Checks whether the parent branch has moved ahead of the branch point.
6. If needed, rebases the agent branch onto the parent branch.
7. If a conflict occurs, aborts the rebase and reports the conflicting files.

## When to use it

Use `sync` when the parent branch has new commits and you want to integrate them into the agent worktree before finishing or merging.

## Examples

```bash
# Sync from inside an agent worktree
opencode-worktree sync

# Sync a specific worktree from anywhere
opencode-worktree sync /path/to/worktree
```

## Important behavior

- If the parent branch has no new commits, the command reports that the worktree is already up to date.
- `sync` is stricter than `merge`: it refuses to run if the worktree has uncommitted changes.
- Rebase conflicts are aborted automatically so the worktree is not left mid-rebase.
- The command only works for managed agent worktrees with a valid `.agent-parent-branch` marker.

## Common failure cases

- Running without a path from a directory that is not an agent worktree.
- Missing or empty `.agent-parent-branch` marker.
- Worktree is dirty.
- Worktree is on a non-agent branch.
- Rebase conflicts while replaying agent commits onto the parent branch.

## Related commands

- [`merge`](merge.md) to finish the session after syncing.
- [`task`](task.md) to start a new synced-off-current-branch session.
