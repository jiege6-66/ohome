#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

if ! command -v flutter >/dev/null 2>&1; then
  echo "[ERROR] flutter not found in PATH." >&2
  exit 1
fi

flutter config --enable-macos-desktop >/dev/null

if [[ -d macos ]]; then
  echo "[INFO] app/macos already exists, skip generation."
  exit 0
fi

echo "[INFO] Generating macOS host project with flutter create..."
flutter create \
  --platforms=macos \
  --project-name ohome \
  --org iosjk.xyz \
  .
