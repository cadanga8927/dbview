#!/usr/bin/env sh

set -eu

REPO="${DBVIEW_REPO:-pageton/dbview}"
INSTALL_DIR="${DBVIEW_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${DBVIEW_VERSION:-latest}"
APP_NAME="dbview"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf 'Error: required command not found: %s\n' "$1" >&2
    exit 1
  fi
}

fetch() {
  url="$1"
  out="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out"
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -qO "$out" "$url"
    return
  fi

  printf 'Error: curl or wget is required\n' >&2
  exit 1
}

resolve_version() {
  if [ "$VERSION" != "latest" ]; then
    case "$VERSION" in
      v*) printf '%s\n' "$VERSION" ;;
      *) printf 'v%s\n' "$VERSION" ;;
    esac
    return
  fi

  api_url="https://api.github.com/repos/$REPO/releases/latest"
  release_json="$(mktemp)"
  trap 'rm -f "$release_json" "$tmpfile"' EXIT INT TERM
  fetch "$api_url" "$release_json"

  resolved_version="$(grep -m1 '"tag_name"' "$release_json" | sed -E 's/.*"tag_name"[[:space:]]*:[[:space:]]*"([^"]+)".*/\1/')"
  if [ -z "$resolved_version" ]; then
    printf 'Error: unable to resolve latest release version from %s\n' "$api_url" >&2
    exit 1
  fi

  printf '%s\n' "$resolved_version"
}

detect_asset() {
  os_name="$(uname -s)"
  arch_name="$(uname -m)"

  case "$os_name" in
    Linux) platform_os="linux" ;;
    Darwin) platform_os="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) platform_os="windows" ;;
    *)
      printf 'Error: unsupported operating system: %s\n' "$os_name" >&2
      exit 1
      ;;
  esac

  case "$arch_name" in
    x86_64|amd64) platform_arch="amd64" ;;
    aarch64|arm64) platform_arch="arm64" ;;
    *)
      printf 'Error: unsupported architecture: %s\n' "$arch_name" >&2
      exit 1
      ;;
  esac

  if [ "$platform_os" = "windows" ] && [ "$platform_arch" != "amd64" ]; then
    printf 'Error: Windows installer currently supports amd64 only\n' >&2
    exit 1
  fi

  if [ "$platform_os" = "windows" ]; then
    asset_name="${APP_NAME}_${resolved_version}_${platform_os}_${platform_arch}.exe"
    install_name="${APP_NAME}.exe"
  else
    asset_name="${APP_NAME}_${resolved_version}_${platform_os}_${platform_arch}"
    install_name="$APP_NAME"
  fi
}

need_cmd uname
need_cmd mktemp
need_cmd grep
need_cmd sed
need_cmd chmod
need_cmd mkdir

tmpfile="$(mktemp)"
trap 'rm -f "$tmpfile" "$release_json"' EXIT INT TERM

resolved_version="$(resolve_version)"
detect_asset

download_url="https://github.com/$REPO/releases/download/$resolved_version/$asset_name"
mkdir -p "$INSTALL_DIR"
fetch "$download_url" "$tmpfile"

target_path="$INSTALL_DIR/$install_name"
mv "$tmpfile" "$target_path"
chmod +x "$target_path"

printf 'Installed %s %s to %s\n' "$APP_NAME" "$resolved_version" "$target_path"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *) printf 'Note: add %s to your PATH to run `%s` directly.\n' "$INSTALL_DIR" "$install_name" ;;
esac
