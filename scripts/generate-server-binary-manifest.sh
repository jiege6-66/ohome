#!/usr/bin/env bash
set -euo pipefail

OUTPUT_PATH="${1:?missing output path}"

VERSION="${VERSION:-}"
CHANNEL="${CHANNEL:-stable}"
BUILD_TIME="${BUILD_TIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"
RELEASE_NOTES="${RELEASE_NOTES:-Release ${VERSION}}"
MIN_RUNTIME_VERSION="${MIN_RUNTIME_VERSION:-}"
RECOMMENDED_RUNTIME_VERSION="${RECOMMENDED_RUNTIME_VERSION:-}"
LINUX_AMD64_URL="${LINUX_AMD64_URL:-}"
LINUX_AMD64_SHA256="${LINUX_AMD64_SHA256:-}"
LINUX_AMD64_FORMAT="${LINUX_AMD64_FORMAT:-tar.gz}"
LINUX_ARM64_URL="${LINUX_ARM64_URL:-}"
LINUX_ARM64_SHA256="${LINUX_ARM64_SHA256:-}"
LINUX_ARM64_FORMAT="${LINUX_ARM64_FORMAT:-tar.gz}"
WINDOWS_AMD64_URL="${WINDOWS_AMD64_URL:-}"
WINDOWS_AMD64_SHA256="${WINDOWS_AMD64_SHA256:-}"
WINDOWS_AMD64_FORMAT="${WINDOWS_AMD64_FORMAT:-zip}"
DARWIN_AMD64_URL="${DARWIN_AMD64_URL:-}"
DARWIN_AMD64_SHA256="${DARWIN_AMD64_SHA256:-}"
DARWIN_AMD64_FORMAT="${DARWIN_AMD64_FORMAT:-tar.gz}"
DARWIN_ARM64_URL="${DARWIN_ARM64_URL:-}"
DARWIN_ARM64_SHA256="${DARWIN_ARM64_SHA256:-}"
DARWIN_ARM64_FORMAT="${DARWIN_ARM64_FORMAT:-tar.gz}"

if [[ -z "$VERSION" ]]; then
  echo "VERSION is required" >&2
  exit 1
fi

if [[ -z "$LINUX_AMD64_URL" || -z "$LINUX_AMD64_SHA256" ]]; then
  echo "LINUX_AMD64_URL and LINUX_AMD64_SHA256 are required" >&2
  exit 1
fi

if [[ -z "$LINUX_ARM64_URL" || -z "$LINUX_ARM64_SHA256" ]]; then
  echo "LINUX_ARM64_URL and LINUX_ARM64_SHA256 are required" >&2
  exit 1
fi

if [[ -z "$WINDOWS_AMD64_URL" || -z "$WINDOWS_AMD64_SHA256" ]]; then
  echo "WINDOWS_AMD64_URL and WINDOWS_AMD64_SHA256 are required" >&2
  exit 1
fi

if [[ -z "$DARWIN_AMD64_URL" || -z "$DARWIN_AMD64_SHA256" ]]; then
  echo "DARWIN_AMD64_URL and DARWIN_AMD64_SHA256 are required" >&2
  exit 1
fi

if [[ -z "$DARWIN_ARM64_URL" || -z "$DARWIN_ARM64_SHA256" ]]; then
  echo "DARWIN_ARM64_URL and DARWIN_ARM64_SHA256 are required" >&2
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
      "url": "${LINUX_AMD64_URL}",
      "sha256": "${LINUX_AMD64_SHA256}",
      "format": "${LINUX_AMD64_FORMAT}"
    },
    "linux-arm64": {
      "url": "${LINUX_ARM64_URL}",
      "sha256": "${LINUX_ARM64_SHA256}",
      "format": "${LINUX_ARM64_FORMAT}"
    },
    "windows-amd64": {
      "url": "${WINDOWS_AMD64_URL}",
      "sha256": "${WINDOWS_AMD64_SHA256}",
      "format": "${WINDOWS_AMD64_FORMAT}"
    },
    "darwin-amd64": {
      "url": "${DARWIN_AMD64_URL}",
      "sha256": "${DARWIN_AMD64_SHA256}",
      "format": "${DARWIN_AMD64_FORMAT}"
    },
    "darwin-arm64": {
      "url": "${DARWIN_ARM64_URL}",
      "sha256": "${DARWIN_ARM64_SHA256}",
      "format": "${DARWIN_ARM64_FORMAT}"
    }
  }
}
EOF
