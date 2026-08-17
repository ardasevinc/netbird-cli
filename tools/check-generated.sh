#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
if ! cmp -s "$repo_root/coverage/manifest.json" "$repo_root/internal/catalog/assets/coverage/manifest.json"; then
  echo "stale generated coverage manifest: internal/catalog/assets/coverage/manifest.json" >&2
  exit 1
fi

asset_root="$repo_root/internal/catalog/assets/schemas/nb/v1"
public_root="$repo_root/schemas/nb/v1"

for source in "$asset_root"/*.json; do
  name="$(basename "$source")"
  public="$public_root/$name"
  test -f "$public" || { echo "missing generated schema: $public" >&2; exit 1; }
  cmp -s "$source" "$public" || { echo "stale generated schema: $public" >&2; exit 1; }
done

skill_asset_root="$repo_root/internal/catalog/assets/skills"
skill_public_root="$repo_root/skills"
while IFS= read -r -d '' source; do
  relative="${source#"$skill_asset_root/"}"
  public="$skill_public_root/$relative"
  test -f "$public" || { echo "missing generated skill: $public" >&2; exit 1; }
  cmp -s "$source" "$public" || { echo "stale generated skill: $public" >&2; exit 1; }
done < <(find "$skill_asset_root" -type f -name 'SKILL.md' -print0)

echo "generated schemas and skills are current"
