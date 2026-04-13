# `opencode-worktree merge`

## Summary

Merges an agent branch back into its recorded parent branch and optionally removes the agent worktree and branch afterward.

## Usage

```bash
opencode-worktree merge [path] [--no-cleanup]
```

## Arguments

- `[path]` — optional path to an agent worktree. If omitted, the command tries to detect the current directory as a managed agent worktree.

## Options

- `--no-cleanup` — keep the worktree and branch after a successful merge.

## What it does

1. Resolves the target worktree path, either from the argument or by auto-detecting the current directory.
2. Reads `.agent-parent-branch` to find the branch the agent work should merge back into.
3. Confirms the current worktree branch starts with `agent/`.
4. Counts commits that exist on the agent branch but not the parent branch.
5. Checks whether the worktree has uncommitted changes, excluding marker and copied config files.
6. Takes a file lock before switching the main repo to the parent branch and merging, so concurrent merges do not race.
7. On success, optionally removes the worktree, deletes the agent branch, and prunes stale worktree metadata.

## When to use it

Use `merge` when you skipped auto-merge, when you want to merge from a specific worktree manually, or when you need to retry after resolving issues.

## Examples

```bash
# Merge from the current agent worktree
opencode-worktree merge

# Merge a specific worktree from anywhere
opencode-worktree merge /path/to/worktree

# Merge but keep the worktree and branch around
opencode-worktree merge --no-cleanup
```

## Important behavior

- If the worktree has no new commits, the command skips the merge itself.
- If there are no new commits and no uncommitted changes, cleanup still happens by default.
- If the worktree has uncommitted changes, the worktree is preserved even when cleanup is enabled.
- Merge conflicts are detected, the merge is aborted automatically, and conflicting files are listed.
- The command only works for managed agent worktrees that contain `.agent-parent-branch` and are on an `agent/*` branch.

## Common failure cases

- Running without a path from a directory that is not an agent worktree.
- Missing `.agent-parent-branch` marker.
- Worktree is on a non-agent branch.
- Merge conflicts between the parent branch and the agent branch.
- Cleanup failing after a successful merge.

## Related commands

- [`task`](task.md) and [`attach`](attach.md) for auto-merge on exit.
- [`sync`](sync.md) to rebase the agent branch onto the latest parent branch before merging.
- [`cleanup`](cleanup.md) to remove leftover worktrees or branches later.
