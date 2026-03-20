#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

STATE_FILE=".build_number"
PUBSPEC_FILE="pubspec.yaml"
APP_ENV="prod"
OUTPUT_DIR="build/app/outputs/flutter-apk"
DRY_RUN=""
REQUESTED_BUILD_NUMBER=""

# ── Parse arguments ──────────────────────────────────────────────
if [[ "${1:-}" == "--dry-run" ]]; then
  DRY_RUN=1
  REQUESTED_BUILD_NUMBER="${2:-}"
else
  REQUESTED_BUILD_NUMBER="${1:-}"
fi

# ── Helper ───────────────────────────────────────────────────────
is_number() {
  [[ "$1" =~ ^[0-9]+$ ]]
}

# ── Read version from pubspec.yaml ───────────────────────────────
if [[ ! -f "$PUBSPEC_FILE" ]]; then
  echo "[ERROR] \"$PUBSPEC_FILE\" not found."
  exit 1
fi

VERSION_LINE=$(grep -m1 '^version:' "$PUBSPEC_FILE" | cut -d':' -f2 | tr -d ' ')

if [[ -z "$VERSION_LINE" ]]; then
  echo "[ERROR] version was not found in \"$PUBSPEC_FILE\"."
  exit 1
fi

BUILD_NAME="${VERSION_LINE%%+*}"
PUBSPEC_BUILD_NUMBER="${VERSION_LINE#*+}"

if [[ -z "$BUILD_NAME" ]]; then
  echo "[ERROR] build name could not be parsed from \"$VERSION_LINE\"."
  exit 1
fi

# If there was no '+' in the version string, default to 1
if [[ "$PUBSPEC_BUILD_NUMBER" == "$VERSION_LINE" ]]; then
  PUBSPEC_BUILD_NUMBER=1
fi

# ── Determine build number ───────────────────────────────────────
if [[ -n "$REQUESTED_BUILD_NUMBER" ]]; then
  if ! is_number "$REQUESTED_BUILD_NUMBER"; then
    echo "[ERROR] Invalid build number \"$REQUESTED_BUILD_NUMBER\"."
    exit 1
  fi
  BUILD_NUMBER="$REQUESTED_BUILD_NUMBER"
  BUILD_NUMBER_SOURCE="argument"
else
  LAST_BUILD_NUMBER=""
  if [[ -f "$STATE_FILE" ]]; then
    LAST_BUILD_NUMBER=$(tr -d '[:space:]' < "$STATE_FILE")
    if ! is_number "$LAST_BUILD_NUMBER"; then
      echo "[ERROR] Invalid value in \"$STATE_FILE\": \"$LAST_BUILD_NUMBER\"."
      exit 1
    fi
  fi

  if [[ -n "$LAST_BUILD_NUMBER" ]]; then
    BUILD_NUMBER=$((LAST_BUILD_NUMBER + 1))
    BUILD_NUMBER_SOURCE="local state"
  else
    GIT_BUILD_NUMBER=$(git rev-list --count HEAD 2>/dev/null || true)

    if [[ -n "$GIT_BUILD_NUMBER" ]] && is_number "$GIT_BUILD_NUMBER"; then
      BUILD_NUMBER=$((GIT_BUILD_NUMBER + 1))
      BUILD_NUMBER_SOURCE="git commit count"
    else
      if ! is_number "$PUBSPEC_BUILD_NUMBER"; then
        echo "[ERROR] Invalid pubspec build number \"$PUBSPEC_BUILD_NUMBER\"."
        exit 1
      fi
      BUILD_NUMBER=$((PUBSPEC_BUILD_NUMBER + 1))
      BUILD_NUMBER_SOURCE="pubspec seed"
    fi
  fi
fi

if [[ "$BUILD_NUMBER" -le 0 ]]; then
  echo "[ERROR] build number must be greater than 0."
  exit 1
fi

# ── Version codes per ABI ────────────────────────────────────────
ARM32_VERSION_CODE=$((1000 + BUILD_NUMBER))
ARM64_VERSION_CODE=$((2000 + BUILD_NUMBER))
X64_VERSION_CODE=$((4000 + BUILD_NUMBER))

echo ""
echo "Build name   : $BUILD_NAME"
echo "Build number : $BUILD_NUMBER ($BUILD_NUMBER_SOURCE)"
echo "Version code : armeabi-v7a=$ARM32_VERSION_CODE, arm64-v8a=$ARM64_VERSION_CODE, x86_64=$X64_VERSION_CODE"
echo "Output dir   : $OUTPUT_DIR"
echo ""

# ── Dry run ──────────────────────────────────────────────────────
if [[ -n "$DRY_RUN" ]]; then
  echo "[DRY RUN] flutter pub get"
  echo "[DRY RUN] flutter build apk --release --split-per-abi --dart-define=APP_ENV=$APP_ENV --build-name=$BUILD_NAME --build-number=$BUILD_NUMBER"
  exit 0
fi

# ── Build ────────────────────────────────────────────────────────
flutter pub get
if [[ $? -ne 0 ]]; then
  echo "[ERROR] flutter pub get failed."
  exit 1
fi

flutter build apk --release --split-per-abi \
  --dart-define=APP_ENV="$APP_ENV" \
  --build-name="$BUILD_NAME" \
  --build-number="$BUILD_NUMBER"
if [[ $? -ne 0 ]]; then
  echo "[ERROR] flutter build apk failed."
  exit 1
fi

# ── Persist build number ─────────────────────────────────────────
echo "$BUILD_NUMBER" > "$STATE_FILE"

echo ""
echo "Build finished successfully."
echo "State file updated: $STATE_FILE"
