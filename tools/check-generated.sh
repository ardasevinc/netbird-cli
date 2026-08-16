#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
asset_root="$repo_root/internal/catalog/assets/schemas/nb/v1"
public_root="$repo_root/schemas/nb/v1"

for source in "$asset_root"/*.json; do
  name="$(basename "$source")"
  public="$public_root/$name"
  test -f "$public" || { echo "missing generated schema: $public" >&2; exit 1; }
  cmp -s "$source" "$public" || { echo "stale generated schema: $public" >&2; exit 1; }
done

echo "generated schemas are current"
