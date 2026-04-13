package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func run(repoDir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if repoDir != "" {
		cmd.Dir = repoDir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), strings.TrimSpace(stderr.String()), err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func RepoRoot(dir string) (string, error) {
	return run(dir, "rev-parse", "--show-toplevel")
}

func CurrentBranch(dir string) (string, error) {
	return run(dir, "branch", "--show-current")
}

func WorktreeList(dir string) (string, error) {
	return run(dir, "worktree", "list")
}

func WorktreeListPorcelain(dir string) (string, error) {
	return run(dir, "worktree", "list", "--porcelain")
}

func WorktreeAdd(dir, worktreePath, branchName, startPoint string) error {
	_, err := run(dir, "worktree", "add", worktreePath, "-b", branchName, startPoint)
	return err
}

func WorktreeRemove(dir, worktreePath string) error {
	_, err := run(dir, "worktree", "remove", worktreePath, "--force")
	return err
}

func WorktreePrune(dir string) error {
	_, err := run(dir, "worktree", "prune")
	return err
}

func BranchDelete(dir, branch string) (string, error) {
	return run(dir, "branch", "-d", branch)
}

func Checkout(dir, branch string) error {
	_, err := run(dir, "checkout", branch)
	return err
}

func Merge(dir, branch string) error {
	_, err := run(dir, "merge", branch, "--no-edit")
	return err
}

func MergeAbort(dir string) error {
	_, err := run(dir, "merge", "--abort")
	return err
}

func ConflictingFiles(dir string) ([]string, error) {
	out, err := run(dir, "diff", "--name-only", "--diff-filter=U")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

func CommitCountBetween(dir, base, head string) (int, error) {
	out, err := run(dir, "log", base+".."+head, "--oneline")
	if err != nil {
		return 0, err
	}
	if out == "" {
		return 0, nil
	}
	return len(strings.Split(out, "\n")), nil
}

func GitCommonDir(dir string) (string, error) {
	out, err := run(dir, "rev-parse", "--git-common-dir")
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(out) {
		return out, nil
	}
	return filepath.Join(dir, out), nil
}

func HasUncommittedChanges(dir string, excludePaths []string) (bool, error) {
	out, err := run(dir, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	if out == "" {
		return false, nil
	}
	if len(excludePaths) == 0 {
		return true, nil
	}

	excluded := make(map[string]bool, len(excludePaths))
	for _, p := range excludePaths {
		excluded[p] = true
	}

	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		filePath := strings.TrimSpace(line[3:])
		if !excluded[filePath] {
			return true, nil
		}
	}
	return false, nil
}

func BranchList(dir string) (string, error) {
	return run(dir, "branch")
}

func Rebase(dir, onto string) error {
	_, err := run(dir, "rebase", onto)
	return err
}

func RebaseAbort(dir string) error {
	_, err := run(dir, "rebase", "--abort")
	return err
}

func MergeBase(dir, a, b string) (string, error) {
	return run(dir, "merge-base", a, b)
}
