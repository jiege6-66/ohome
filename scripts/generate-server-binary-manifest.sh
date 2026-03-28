#!/usr/bin/env bash
set -euo pipefail

OUTPUT_PATH="${1:?missing output path}"

VERSION="${VERSION:-}"
CHANNEL="${CHANNEL:-stable}"
BUILD_TIME="${BUILD_TIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"
RELEASE_NOTES="${RELEASE_NOTES:-Release ${VERSION}}"
MIN_RUNTIME_VERSION="${MIN_RUNTIME_VERSION:-}"
RECOMMENDED_RUNTIME_VERSION="${RECOMMENDED_RUNTIME_VERSION:-}"
AMD64_URL="${AMD64_URL:-}"
AMD64_SHA256="${AMD64_SHA256:-}"
ARM64_URL="${ARM64_URL:-}"
ARM64_SHA256="${ARM64_SHA256:-}"

if [[ -z "$VERSION" ]]; then
  echo "VERSION is required" >&2
  exit 1
fi

if [[ -z "$AMD64_URL" || -z "$AMD64_SHA256" ]]; then
  echo "AMD64_URL and AMD64_SHA256 are required" >&2
  exit 1
fi

if [[ -z "$ARM64_URL" || -z "$ARM64_SHA256" ]]; then
  echo "ARM64_URL and ARM64_SHA256 are required" >&2
  exit 1
fi

mkdir -p "$(dirname "$OUTPUT_PATH")"

cat > "$OUTPUT_PATH" <<EOF
{
  "channel": "${CHANNEL}",
  "version": "${VERSION}",
  "releaseNotes": "${RELEASE_NOTES}",
  "publishedAt": "${BUILD_TIME}",
  "minRuntimeVersion": "${MIN_RUNTIME_VERSION}",
  "recommendedRuntimeVersion": "${RECOMMENDED_RUNTIME_VERSION}",
  "artifacts": {
    "linux-amd64": {
      "url": "${AMD64_URL}",
      "sha256": "${AMD64_SHA256}"
    },
    "linux-arm64": {
      "url": "${ARM64_URL}",
      "sha256": "${ARM64_SHA256}"
    }
  }
}
EOF
