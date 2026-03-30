#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

PUBSPEC_FILE="pubspec.yaml"
APP_ENV="prod"
OUTPUT_DIR="build/macos/Build/Products/Release"
REQUESTED_BUILD_NUMBER="${1:-}"
BUILD_METADATA_FILE="${BUILD_METADATA_FILE:-}"
RELEASE_TAG="${RELEASE_TAG:-}"
CI_BUILD_NUMBER="${CI_BUILD_NUMBER:-}"

is_number() {
  [[ "$1" =~ ^[0-9]+$ ]]
}

TAG_BUILD_NAME=""
parse_release_tag() {
  local normalized_tag="$1"

  if [[ "$normalized_tag" =~ ^v?([0-9]+\.[0-9]+\.[0-9]+(-rc[0-9]+)?)$ ]]; then
    TAG_BUILD_NAME="${BASH_REMATCH[1]}"
    return 0
  fi

  TAG_BUILD_NAME=""
  return 1
}

if [[ ! -f "$PUBSPEC_FILE" ]]; then
  echo "[ERROR] \"$PUBSPEC_FILE\" not found."
  exit 1
fi

VERSION_LINE=$(grep -m1 '^version:' "$PUBSPEC_FILE" | cut -d':' -f2 | tr -d ' \r')
if [[ -z "$VERSION_LINE" ]]; then
  echo "[ERROR] version was not found in \"$PUBSPEC_FILE\"."
  exit 1
fi

BUILD_NAME="${VERSION_LINE%%+*}"
PUBSPEC_BUILD_NUMBER="${VERSION_LINE#*+}"
if [[ "$PUBSPEC_BUILD_NUMBER" == "$VERSION_LINE" ]]; then
  PUBSPEC_BUILD_NUMBER=1
fi

if [[ -n "$RELEASE_TAG" ]]; then
  NORMALIZED_TAG="${RELEASE_TAG#refs/tags/}"
  parse_release_tag "$NORMALIZED_TAG" || true
fi

if [[ -n "$TAG_BUILD_NAME" ]]; then
  BUILD_NAME="$TAG_BUILD_NAME"
fi

if [[ -n "$REQUESTED_BUILD_NUMBER" ]]; then
  if ! is_number "$REQUESTED_BUILD_NUMBER"; then
    echo "[ERROR] Invalid build number \"$REQUESTED_BUILD_NUMBER\"."
    exit 1
  fi
  BUILD_NUMBER="$REQUESTED_BUILD_NUMBER"
elif [[ -n "$CI_BUILD_NUMBER" ]]; then
  if ! is_number "$CI_BUILD_NUMBER"; then
    echo "[ERROR] Invalid CI build number \"$CI_BUILD_NUMBER\"."
    exit 1
  fi
  BUILD_NUMBER="$CI_BUILD_NUMBER"
else
  if ! is_number "$PUBSPEC_BUILD_NUMBER"; then
    echo "[ERROR] Invalid pubspec build number \"$PUBSPEC_BUILD_NUMBER\"."
    exit 1
  fi
  BUILD_NUMBER="$PUBSPEC_BUILD_NUMBER"
fi

if [[ "$BUILD_NUMBER" -le 0 ]]; then
  echo "[ERROR] build number must be greater than 0."
  exit 1
fi

bash ensure_macos_platform.sh
flutter pub get
flutter build macos \
  --release \
  --dart-define=APP_ENV="$APP_ENV" \
  --build-name="$BUILD_NAME" \
  --build-number="$BUILD_NUMBER"

APP_BUNDLE_PATH=$(find "$OUTPUT_DIR" -maxdepth 1 -type d -name '*.app' | sort | head -n 1)
if [[ -z "$APP_BUNDLE_PATH" ]]; then
  echo "[ERROR] Failed to locate built .app bundle under $OUTPUT_DIR"
  exit 1
fi

RAW_ARCH=$(uname -m)
case "$RAW_ARCH" in
  x86_64) ARCH_LABEL="amd64" ;;
  arm64) ARCH_LABEL="arm64" ;;
  *) ARCH_LABEL="$RAW_ARCH" ;;
esac

ZIP_FILENAME="ohome-macos-${ARCH_LABEL}.zip"
ZIP_PATH="$OUTPUT_DIR/$ZIP_FILENAME"
CHECKSUM_PATH="$OUTPUT_DIR/checksums-macos-${ARCH_LABEL}.txt"

rm -f "$ZIP_PATH"
ditto -c -k --sequesterRsrc --keepParent "$APP_BUNDLE_PATH" "$ZIP_PATH"

ZIP_BASENAME=$(basename "$ZIP_PATH")
APP_BUNDLE_BASENAME=$(basename "$APP_BUNDLE_PATH")
SHA256=$(shasum -a 256 "$ZIP_PATH" | awk '{print $1}')
printf '%s  %s\n' "$SHA256" "$ZIP_BASENAME" > "$CHECKSUM_PATH"

if [[ -n "$BUILD_METADATA_FILE" ]]; then
  mkdir -p "$(dirname "$BUILD_METADATA_FILE")"
  cat > "$BUILD_METADATA_FILE" <<EOF
BUILD_NAME=$BUILD_NAME
BUILD_NUMBER=$BUILD_NUMBER
MACOS_ARCH_LABEL=$ARCH_LABEL
MACOS_OUTPUT_DIR=$OUTPUT_DIR
MACOS_APP_BUNDLE_PATH=$APP_BUNDLE_PATH
MACOS_APP_BUNDLE_NAME=$APP_BUNDLE_BASENAME
MACOS_ZIP_PATH=$ZIP_PATH
MACOS_ZIP_FILENAME=$ZIP_BASENAME
MACOS_CHECKSUM_PATH=$CHECKSUM_PATH
MACOS_ZIP_SHA256=$SHA256
EOF
fi

echo "[INFO] Built macOS app bundle: $APP_BUNDLE_PATH"
echo "[INFO] Packaged macOS archive: $ZIP_PATH"
