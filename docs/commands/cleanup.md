# `opencode-worktree cleanup`

## Summary

Removes orphaned agent worktree directories and stale `agent/*` branches that no longer correspond to an active git worktree.

## Usage

```bash
opencode-worktree cleanup [--dry-run] [--yes]
```

## Arguments

This command does not take positional arguments.

## Options

- `--dry-run` — show what would be removed without deleting anything.
- `--yes` — skip the confirmation prompt.

## What it does

1. Prunes stale git worktree metadata.
2. Scans sibling directories for names matching the tool's worktree naming pattern.
3. Compares those directories against active git worktrees and marks unmatched ones as orphaned.
4. Scans local branches for `agent/*` branches that are no longer attached to an active worktree.
5. In normal mode, prompts for confirmation when running interactively unless `--yes` is set.
6. Removes orphaned directories and attempts to delete orphaned branches.

## When to use it

Use `cleanup` when a worktree directory or agent branch is left behind after manual interruption, failed cleanup, or other unexpected session endings.

## Examples

```bash
# Review what would be removed
opencode-worktree cleanup --dry-run

# Remove orphaned items without an interactive prompt
opencode-worktree cleanup --yes
```

## Important behavior

- Active worktree directories and active `agent/*` branches are preserved.
- In `--dry-run` mode, nothing is removed.
- Without `--yes`, interactive terminals get a confirmation prompt.
- If a stale branch cannot be deleted because it is unmerged, the command leaves it in place and prints the `git branch -D` command to force deletion manually.

## Common failure cases

- Running outside a git repository.
- File-system permission failures while removing directories.
- Branch deletion blocked because the branch is not fully merged.

## Related commands

- [`list`](list.md) to inspect active sessions before cleaning.
- [`merge`](merge.md) if a worktree is still valid and should be merged instead of deleted.
