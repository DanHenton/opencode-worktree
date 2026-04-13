# `opencode-worktree task`

## Summary

Creates a new agent worktree from your current branch, launches OpenCode inside it, and by default merges the finished work back into the parent branch when OpenCode exits.

## Usage

```bash
opencode-worktree task <name> [message] [--no-merge]
```

## Arguments

- `<name>` — required task name used for both the branch and sibling worktree directory. Only letters, numbers, and hyphens are allowed.
- `[message]` — optional initial prompt passed to `opencode --prompt`.

## Options

- `--no-merge` — skip the automatic merge step after the OpenCode session ends.

## What it does

1. Confirms you are inside a git repository and currently on a named branch.
2. Validates the task name and rejects duplicates if a matching agent worktree already exists.
3. Creates a new branch named `agent/<name>`.
4. Creates a sibling worktree at `<repo-dir>-agent-<name>`.
5. Writes `.agent-parent-branch` and `.agent-context` marker files into the worktree.
6. Copies `opencode.json` and `.opencode/` into the worktree if they exist in the source repo.
7. Launches `opencode` in the new worktree.
8. Unless `--no-merge` is set, attempts to merge the agent branch back into the parent branch and clean up.

## When to use it

Use `task` when you want to start a brand-new isolated agent session for a piece of work.

## Examples

```bash
# Start a new isolated session
opencode-worktree task fix-auth-bug

# Start a session with an initial prompt
opencode-worktree task fix-auth-bug "Fix the JWT token expiry bug"

# Keep the worktree and merge manually later
opencode-worktree task add-dark-mode --no-merge
```

## Important behavior

- If `opencode` is not installed or not on your `PATH`, launch fails.
- Task names cannot contain spaces, underscores, slashes, or special characters.
- If a worktree already exists for the same task name, the command stops instead of reusing it.
- Automatic merge happens only after the OpenCode process exits.
- If there are no new commits, the merge step skips the merge itself and only cleans up.
- If the worktree still has uncommitted changes, the merge result is preserved and the worktree is kept.

## Common failure cases

- Not inside a git repository.
- Running from detached `HEAD` instead of a named branch.
- Reusing an existing task name.
- Missing `opencode` binary.
- Merge conflicts during auto-merge.

## Related commands

- [`attach`](attach.md) to reopen an existing session.
- [`merge`](merge.md) to merge later when `--no-merge` was used.
- [`list`](list.md) to see active sessions.
