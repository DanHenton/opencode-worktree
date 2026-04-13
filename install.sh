#!/bin/sh
set -eu

REPO="DanHenton/opencode-worktree"
BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"

detect_os() {
  case "$(uname -s)" in
    Linux)  echo "Linux" ;;
    Darwin) echo "Darwin" ;;
    *) printf 'Unsupported OS: %s\n' "$(uname -s)" >&2; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64)  echo "x86_64" ;;
    arm64|aarch64) echo "arm64" ;;
    i386|i686)     echo "i386" ;;
    armv7l|armv7)  echo "armv7" ;;
    *) printf 'Unsupported architecture: %s\n' "$(uname -m)" >&2; exit 1 ;;
  esac
}

os="$(detect_os)"
arch="$(detect_arch)"
archive="opencode-worktree_${os}_${arch}.tar.gz"
url="https://github.com/${REPO}/releases/latest/download/${archive}"

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT INT TERM

printf 'Downloading %s...\n' "$archive"
curl -fsSL "$url" -o "$tmp/$archive"
tar -xzf "$tmp/$archive" -C "$tmp"

mkdir -p "$BIN_DIR"
install -m 0755 "$tmp/opencode-worktree" "$BIN_DIR/opencode-worktree"

printf 'Installed opencode-worktree to %s/opencode-worktree\n' "$BIN_DIR"

case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    printf '\nAdd this to your shell profile to make it available:\n'
    printf '  export PATH="%s:$PATH"\n\n' "$BIN_DIR"
    printf 'Then open a new terminal, or run:\n'
    printf '  export PATH="%s:$PATH"\n' "$BIN_DIR"
    ;;
esac
