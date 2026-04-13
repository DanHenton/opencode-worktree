# opencode-worktree Agent Guide

## Build & Test

```bash
go build ./cmd/opencode-worktree
go test ./... -v
go vet ./...
go fmt ./...
```

## Project Structure

- `cmd/opencode-worktree/main.go` — Entry point, routing switch (`run()`), usage output, `errSilent` dispatch
- `cmd/opencode-worktree/task.go` — `runTask` handler
- `cmd/opencode-worktree/attach.go` — `runAttach` handler
- `cmd/opencode-worktree/merge_cmd.go` — `runMerge` handler
- `cmd/opencode-worktree/sync_cmd.go` — `runSync` handler
- `cmd/opencode-worktree/list.go` — `runList` handler
- `cmd/opencode-worktree/cleanup.go` — `runCleanup` handler
- `cmd/opencode-worktree/completions.go` — `runCompletions` handler
- `cmd/opencode-worktree/output.go` — Shared output helpers: `emoji`, `printMergeResult`, `handleMergeError`, `handleSyncError`, `errSilent`
- `internal/git/` — Thin wrappers around `exec.Command("git", ...)`, no abstractions
- `internal/worktree/` — Create, list, cleanup worktrees; launch opencode; copy config
- `internal/merge/` — Merge agent branch into parent with flock serialization

## Code Style

- Go 1.24, standard library `flag` package (no Cobra)
- No abstraction layers — package-level functions calling git directly
- No config files — all behavior from flags and git context
- Error handling: `fmt.Errorf("failed to X: %w", err)` pattern throughout
- Tests use real git repos in temp dirs (no mocks)

## Key Behaviors

- Worktree created as sibling: `<repo-dir>-agent-<task-name>`
- Branch naming: `agent/<task-name>`
- Marker files: `.agent-parent-branch`, `.agent-context`
- Merge lock: `/tmp/<repo-basename>-merge.lock`
- `git checkout` of parent branch happens INSIDE the flock (race condition fix from original bash scripts)
