#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

test -s LICENSE
test -s NOTICE
test -s go.mod

packages="$(go list -deps ./... )"
if printf '%s\n' "$packages" | grep -E '^github.com/netbirdio/netbird/(management|signal|relay|combined|proxy)(/|$)' >/dev/null; then
  echo "disallowed AGPL NetBird package family detected" >&2
  exit 1
fi

echo "license boundary passed"
