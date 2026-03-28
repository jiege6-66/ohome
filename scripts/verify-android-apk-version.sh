#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 3 ]]; then
  echo "Usage: $0 <apk-path> <expected-version-name> <expected-version-code>" >&2
  exit 1
fi

APK_PATH="$1"
EXPECTED_VERSION_NAME="$2"
EXPECTED_VERSION_CODE="$3"

if [[ ! -f "$APK_PATH" ]]; then
  echo "[ERROR] APK not found: $APK_PATH" >&2
  exit 1
fi

find_aapt() {
  local sdk_roots=()
  if [[ -n "${ANDROID_HOME:-}" ]]; then
    sdk_roots+=("${ANDROID_HOME}")
  fi
  if [[ -n "${ANDROID_SDK_ROOT:-}" ]]; then
    sdk_roots+=("${ANDROID_SDK_ROOT}")
  fi

  local root
  for root in "${sdk_roots[@]}"; do
    [[ -d "$root/build-tools" ]] || continue

    mapfile -t candidates < <(find "$root/build-tools" -type f -name aapt 2>/dev/null | sort -V)
    if [[ ${#candidates[@]} -gt 0 ]]; then
      echo "${candidates[-1]}"
      return 0
    fi
  done

  return 1
}

AAPT_BIN="$(find_aapt || true)"
if [[ -z "$AAPT_BIN" ]]; then
  echo "[ERROR] Unable to find Android build-tools aapt under ANDROID_HOME/ANDROID_SDK_ROOT." >&2
  exit 1
fi

BADGING_OUTPUT="$("$AAPT_BIN" dump badging "$APK_PATH")"
PACKAGE_LINE="$(printf '%s\n' "$BADGING_OUTPUT" | grep -m1 "^package:")"

if [[ -z "$PACKAGE_LINE" ]]; then
  echo "[ERROR] Failed to read APK package metadata from: $APK_PATH" >&2
  exit 1
fi

ACTUAL_VERSION_CODE="$(printf '%s\n' "$PACKAGE_LINE" | sed -n "s/.*versionCode='\([^']*\)'.*/\1/p")"
ACTUAL_VERSION_NAME="$(printf '%s\n' "$PACKAGE_LINE" | sed -n "s/.*versionName='\([^']*\)'.*/\1/p")"

if [[ -z "$ACTUAL_VERSION_CODE" || -z "$ACTUAL_VERSION_NAME" ]]; then
  echo "[ERROR] Failed to parse versionCode/versionName from APK: $APK_PATH" >&2
  echo "$PACKAGE_LINE" >&2
  exit 1
fi

echo "APK versionName : $ACTUAL_VERSION_NAME"
echo "APK versionCode : $ACTUAL_VERSION_CODE"
echo "Expected name   : $EXPECTED_VERSION_NAME"
echo "Expected code   : $EXPECTED_VERSION_CODE"

if [[ "$ACTUAL_VERSION_NAME" != "$EXPECTED_VERSION_NAME" ]]; then
  echo "[ERROR] APK versionName mismatch: expected $EXPECTED_VERSION_NAME, got $ACTUAL_VERSION_NAME" >&2
  exit 1
fi

if [[ "$ACTUAL_VERSION_CODE" != "$EXPECTED_VERSION_CODE" ]]; then
  echo "[ERROR] APK versionCode mismatch: expected $EXPECTED_VERSION_CODE, got $ACTUAL_VERSION_CODE" >&2
  exit 1
fi

echo "[OK] APK internal version matches expected release metadata."
