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

install_completions() {
  completion_marker="# opencode-worktree completions"

  zsh_snippet='# opencode-worktree completions
_opencode_worktree() {
  if (( CURRENT == 2 )); then
    compadd $(opencode-worktree --completions 2>/dev/null)
  elif (( CURRENT == 3 )); then
    compadd $(opencode-worktree --completions ${words[2]} 2>/dev/null)
  fi
}
compdef _opencode_worktree opencode-worktree'

  bash_snippet='# opencode-worktree completions
_opencode_worktree() {
  if [ "${#COMP_WORDS[@]}" -eq 2 ]; then
    COMPREPLY=($(compgen -W "$(opencode-worktree --completions 2>/dev/null)" -- "${COMP_WORDS[1]}"))
  elif [ "${#COMP_WORDS[@]}" -eq 3 ]; then
    COMPREPLY=($(compgen -W "$(opencode-worktree --completions "${COMP_WORDS[1]}" 2>/dev/null)" -- "${COMP_WORDS[2]}"))
  fi
}
complete -F _opencode_worktree opencode-worktree'

  shell="$(basename "${SHELL:-}")"
  rc_file=""
  snippet=""

  case "$shell" in
    zsh)
      snippet="$zsh_snippet"
      if [ -f "$HOME/.zshrc" ]; then
        rc_file="$HOME/.zshrc"
      fi
      ;;
    bash)
      snippet="$bash_snippet"
      if [ -f "$HOME/.bashrc" ]; then
        rc_file="$HOME/.bashrc"
      elif [ -f "$HOME/.bash_profile" ]; then
        rc_file="$HOME/.bash_profile"
      fi
      ;;
  esac

  if [ -z "$rc_file" ] || [ -z "$snippet" ]; then
    printf '\nShell completions: could not detect shell rc file. Add manually:\n'
    printf '  See: https://github.com/%s#shell-completions\n' "$REPO"
    return
  fi

  if grep -qF "$completion_marker" "$rc_file" 2>/dev/null; then
    printf '\nShell completions already installed in %s\n' "$rc_file"
    return
  fi

  printf '\n%s\n' "$snippet" >> "$rc_file"
  printf '\nShell completions installed in %s\n' "$rc_file"
  printf 'Restart your shell or run: source %s\n' "$rc_file"
}

install_alias() {
  alias_marker="# opencode-worktree alias"
  alias_line="alias ocwt='opencode-worktree'"
  snippet="$alias_marker
$alias_line"

  shell="$(basename "${SHELL:-}")"
  rc_file=""

  case "$shell" in
    zsh)
      if [ -f "$HOME/.zshrc" ]; then
        rc_file="$HOME/.zshrc"
      fi
      ;;
    bash)
      if [ -f "$HOME/.bashrc" ]; then
        rc_file="$HOME/.bashrc"
      elif [ -f "$HOME/.bash_profile" ]; then
        rc_file="$HOME/.bash_profile"
      fi
      ;;
  esac

  if [ -z "$rc_file" ]; then
    printf '\nShortcut: you can add this alias to your shell profile:\n'
    printf '  %s\n' "$alias_line"
    return
  fi

  if grep -qF "$alias_marker" "$rc_file" 2>/dev/null; then
    printf '\nShell alias (ocwt) already installed in %s\n' "$rc_file"
    return
  fi

  printf '\n%s\n' "$snippet" >> "$rc_file"
  printf '\nShell alias installed: ocwt → opencode-worktree\n'
  printf 'Restart your shell or run: source %s\n' "$rc_file"
}

install_completions
install_alias
