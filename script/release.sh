#!/usr/bin/env bash

set -euo pipefail

day="$(date '+%Y-%m-%d')"
echo "preparing assets for release $day"

targets=(
  darwin-amd64
  darwin-arm64
  linux-amd64
  windows-amd64
)

# Create zipped binaries.
# chdir into dist so zip files sit at root of zip instead of in 'dist/pggen'.
pushd dist >/dev/null
for target in "${targets[@]}"; do
  binary="pggen-${target}"
  if [[ "$binary" == *windows* ]]; then
    binary+='.exe'
  fi
  echo -n "zipping ${binary} ... "
  zip --quiet -9 "pggen-${target}.zip" "${binary}"
  echo "done"
done
popd >/dev/null

# Download github-release if necessary.
GH_REL_BIN='github-release'
if ! command -v "$GH_REL_BIN"; then
  echo 'downloading github-release'
  GH_REL_BIN="$(mktemp)"
  url=https://github.com/github-release/github-release/releases/download/v0.10.0/darwin-amd64-github-release.bz2
  curl -L --fail --silent "${url}" | bzip2 -dc >"$GH_REL_BIN"
  chmod +x "$GH_REL_BIN"
fi

# Delete the remote tag since we're creating a new release tagged today.
git push origin ":refs/tags/$day" 2>/dev/null
# Create or move the day tag to the latest commit.
git tag -f "$day"
# Delete any existing releases. We only support 1 release per day.
# Ignore errors if we try to delete a release that doesn't exist.
"${GH_REL_BIN}" delete --user djsavvy --repo pggen --tag "$day" 2>/dev/null || true

echo
echo "creating release $day"
"${GH_REL_BIN}" release --user djsavvy --repo pggen --tag "$day" --name "$day"

# Upload each of the zipped binaries.
for target in "${targets[@]}"; do
  echo -n "uploading pggen-${target}.zip ... "
  "${GH_REL_BIN}" upload \
    --user djsavvy \
    --repo pggen \
    --tag "$day" \
    --name "pggen-${target}.zip" \
    --file "dist/pggen-${target}.zip"
  echo "done"
done
