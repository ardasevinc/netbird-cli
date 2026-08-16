#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
config="$repo_root/.goreleaser.yaml"
test -f "$config"

required_patterns=(
  '^version: 2$'
  '^project_name: nb$'
  '^checksum:$'
  '^sboms:$'
  '^signs:$'
  '^brews:$'
  '^release:$'
)
for pattern in "${required_patterns[@]}"; do
  rg -q -- "$pattern" "$config" || { echo "release config missing: $pattern" >&2; exit 1; }
done

for literal in \
  'internal/version.Version={{ .Version }}' \
  'internal/version.Commit={{ .FullCommit }}' \
  'internal/version.Date={{ .Date }}'; do
  rg -F -q -- "$literal" "$config" || { echo "release config missing: $literal" >&2; exit 1; }
done

test -f "$repo_root/.github/workflows/release.yml"
rg -q 'push:' "$repo_root/.github/workflows/release.yml"
rg -q 'tags: \["v\*"\]' "$repo_root/.github/workflows/release.yml"
rg -q 'goreleaser/goreleaser-action@v7' "$repo_root/.github/workflows/release.yml"

echo "release configuration is present and provenance-capable"
