#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
upstream_repo="${NB_E2E_NETBIRD_REPO:-/Users/arda/.agents/repo-ref/netbird}"
upstream_tag="v0.77.0"
image="${NB_E2E_NETBIRD_IMAGE:-nb-netbird-cli-e2e:v0.77.0}"
container="nb-cli-e2e-${RANDOM}-${RANDOM}"
work_dir="$(mktemp -d "${TMPDIR:-/tmp}/nb-cli-e2e.XXXXXX")"
fixture_id=""
pat=""
base_url=""

cleanup() {
	set +e
	if [[ -n "$fixture_id" && -n "$base_url" && -n "$pat" ]]; then
		curl -fsS -X DELETE -H "Authorization: Bearer $pat" "$base_url/api/groups/$fixture_id" >/dev/null 2>&1
	fi
	docker rm -f "$container" >/dev/null 2>&1
	find "$work_dir" -type f -delete >/dev/null 2>&1
	find "$work_dir" -depth -type d -empty -delete >/dev/null 2>&1
}
trap cleanup EXIT

command -v docker >/dev/null || { echo "docker is required for self-hosted e2e" >&2; exit 2; }
command -v curl >/dev/null || { echo "curl is required for self-hosted e2e" >&2; exit 2; }
command -v jq >/dev/null || { echo "jq is required for self-hosted e2e" >&2; exit 2; }
test -d "$upstream_repo" || { echo "set NB_E2E_NETBIRD_REPO to a NetBird source checkout" >&2; exit 2; }
git -C "$upstream_repo" cat-file -e "$upstream_tag^{commit}" || { echo "missing pinned NetBird tag $upstream_tag" >&2; exit 2; }

expected_openapi_sha="3470cede60e3263e715b42c37f0d4534efa5ac54ab81b0e87081ade26e6dbe83"
if command -v sha256sum >/dev/null; then
	actual_openapi_sha="$(git -C "$upstream_repo" show "$upstream_tag:shared/management/http/api/openapi.yml" | sha256sum | awk '{print $1}')"
else
	actual_openapi_sha="$(git -C "$upstream_repo" show "$upstream_tag:shared/management/http/api/openapi.yml" | shasum -a 256 | awk '{print $1}')"
fi
test "$actual_openapi_sha" = "$expected_openapi_sha" || { echo "pinned OpenAPI hash mismatch" >&2; exit 2; }

if [[ -z "${NB_E2E_NETBIRD_IMAGE:-}" ]]; then
	upstream_context="$work_dir/netbird"
	mkdir -p "$upstream_context"
	git -C "$upstream_repo" archive "$upstream_tag" | tar -xf - -C "$upstream_context"
	docker build --quiet -f "$upstream_context/combined/Dockerfile.multistage" -t "$image" "$upstream_context" >/dev/null
fi

mkdir -p "$work_dir/data"
sed "s|__EXPOSED_ADDRESS__|http://127.0.0.1:0|" "$repo_root/tests/e2e/netbird-config.yaml" > "$work_dir/config.yaml"
docker run -d --name "$container" -p 127.0.0.1::8080 -v "$work_dir:/nb" \
	-e NB_SETUP_PAT_ENABLED=true -e NB_DISABLE_GEOLOCATION=true \
	"$image" --config /nb/config.yaml >/dev/null

host_port=""
for _ in $(seq 1 60); do
	host_port="$(docker port "$container" 8080/tcp 2>/dev/null | sed 's/.*://')"
	if [[ -n "$host_port" ]] && curl -fsS "http://127.0.0.1:$host_port/api/instance" >/dev/null 2>&1; then
		break
	fi
	host_port=""
	sleep 2
done
test -n "$host_port" || { docker logs "$container" >&2; exit 1; }
base_url="http://127.0.0.1:$host_port"

setup_json="$(curl -fsS -X POST "$base_url/api/setup" -H 'content-type: application/json' --data '{"email":"admin@netbird.test","password":"Netbird-e2e-Passw0rd!","name":"NB E2E Admin","create_pat":true,"pat_expire_in":1}')"
pat="$(jq -er '.personal_access_token // empty' <<<"$setup_json")"
account_id="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/accounts" | jq -er '.[0].id')"

profile="$work_dir/profile.toml"
state="$work_dir/ledger.db"
printf '%s\n' \
	"[profiles.default]" \
	"url = \"$base_url\"" \
	'credential_ref = "env:NB_E2E_TOKEN"' \
	"server_identity = \"$base_url\"" \
	"account_id = \"$account_id\"" > "$profile"

nb() {
	NB_E2E_TOKEN="$pat" "$repo_root/bin/nb" --config "$profile" --state "$state" "$@"
}

accounts_json="$(nb --json accounts list)"
test "$(jq -r '.operation' <<<"$accounts_json")" = "accounts.list"
test "$(jq -r '.data.accounts[0].id' <<<"$accounts_json")" = "$account_id"
users_json="$(nb --json users list)"
test "$(jq -r '.operation' <<<"$users_json")" = "users.list"
groups_json="$(nb --json groups list)"
test "$(jq -r '.data.groups | length' <<<"$groups_json")" -ge 1
routes_json="$(nb --json routes list)"
test "$(jq -r '.operation' <<<"$routes_json")" = "routes.list"
networks_json="$(nb --json networks list)"
test "$(jq -r '.operation' <<<"$networks_json")" = "networks.list"

# The built-in All group is intentionally immutable. Create a disposable
# fixture through the server API, then exercise the consequential update only
# through nb stage/apply and verify its read-back.
fixture_id="$(curl -fsS -X POST "$base_url/api/groups" -H "Authorization: Bearer $pat" -H 'content-type: application/json' --data '{"name":"nb-e2e-fixture"}' | jq -er '.id')"
before="$(curl -fsS -H "Authorization: Bearer $pat" "$base_url/api/groups/$fixture_id")"
after="$(jq '.name = "nb-e2e-fixture-renamed"' <<<"$before")"
plan="$(jq -n --arg id "$fixture_id" --argjson before "$before" --argjson after "$after" '{operation:"groups.update",request:{id:$id,name:"nb-e2e-fixture-renamed"},before:$before,intended_after:$after,findings:[]}')"
stage_json="$(printf '%s' "$plan" | nb --json stage create --from-json)"
stage_id="$(jq -er '.data.stage_id' <<<"$stage_json")"
revision="$(jq -er '.data.revision' <<<"$stage_json")"
apply_json="$(nb --json apply "$stage_id@$revision")"
test "$(jq -r '.data.state' <<<"$apply_json")" = "confirmed_success"
readback_json="$(nb --json groups get "$fixture_id")"
test "$(jq -r '.data.group.name' <<<"$readback_json")" = "nb-e2e-fixture-renamed"

printf 'self-hosted e2e passed: pinned %s, account/user/group/route/network reads, staged group update, read-back confirmed\n' "$upstream_tag"
