# `opencode-worktree attach`

## Summary

Reopens an existing agent worktree session by task name, launches OpenCode inside that worktree again, and by default merges the branch when the session exits.

## Usage

```bash
opencode-worktree attach <name> [--no-merge]
```

## Arguments

- `<name>` — required task name for an already-existing agent worktree.

## Options

- `--no-merge` — skip the automatic merge step after the OpenCode session ends.

## What it does

1. Confirms you are inside the source git repository.
2. Finds the active worktree that matches `agent/<name>`.
3. Launches `opencode` inside that worktree.
4. Unless `--no-merge` is set, runs the same merge-and-cleanup flow used by `task`.

## When to use it

Use `attach` when a task already has a live worktree and you want to continue that session instead of creating a second worktree.

## Examples

```bash
# Reopen an existing session
opencode-worktree attach fix-auth-bug

# Reopen without auto-merging on exit
opencode-worktree attach fix-auth-bug --no-merge
```

## Important behavior

- `attach` does not create a branch or worktree. It only reconnects to one that already exists.
- Task lookup is based on active git worktree metadata, not directory name guessing.
- If the named task does not exist, the command exits with an error.
- The same merge rules as `task` apply after OpenCode exits.

## Common failure cases

- Not inside the source git repository.
- No active worktree exists for the requested task name.
- Missing `opencode` binary.
- Merge conflicts during auto-merge.

## Related commands

- [`task`](task.md) to create a new session.
- [`list`](list.md) to discover active task names.
- [`merge`](merge.md) to merge manually later.
