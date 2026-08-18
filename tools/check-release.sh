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
	grep -Eq -- "$pattern" "$config" || { echo "release config missing: $pattern" >&2; exit 1; }
done

for literal in \
  'internal/version.Version={{ .Version }}' \
  'internal/version.Commit={{ .FullCommit }}' \
  'internal/version.Date={{ .Date }}'; do
	grep -Fq -- "$literal" "$config" || { echo "release config missing: $literal" >&2; exit 1; }
done

test -f "$repo_root/.github/workflows/release.yml"
grep -Eq 'push:' "$repo_root/.github/workflows/release.yml"
grep -Eq 'tags: \["v\*"\]' "$repo_root/.github/workflows/release.yml"
grep -Fq 'goreleaser/goreleaser-action@v7' "$repo_root/.github/workflows/release.yml"
grep -Fq 'sigstore/cosign-installer@ba7bc0a3fef59531c69a25acd34668d6d3fe6f22' "$repo_root/.github/workflows/release.yml"
grep -Fq 'anchore/sbom-action/download-syft@52542941bb775adef77f7a51ff5ea984e600fab8' "$repo_root/.github/workflows/release.yml"
! grep -Fq 'go install github.com/sigstore/cosign' "$repo_root/.github/workflows/release.yml"
! grep -Fq 'go install github.com/anchore/syft' "$repo_root/.github/workflows/release.yml"

echo "release configuration is present and provenance-capable"
